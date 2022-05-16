package apiserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// Status is a return value for calls that don't return other objects
type Status struct {
	success bool
}

// RESTStorage is a generic interface for RESTful storage services
type RESTStorage interface {
	List(*url.URL) (interface{}, error)
	Get(id string) (interface{}, error)
	Delete(id string) error
	Extract(body string) (interface{}, error)
	Create(interface{}) error
	Update(interface{}) error
}

type ApiServer struct {
	prefix  string
	storage map[string]RESTStorage
}

type RESTStorageData string

// New creates a new ApiServer object.
// 'storage' contains a map of handlers.
// 'prefix' is the hosting path prefix.
func New(storage map[string]RESTStorage, prefix string) *ApiServer {
	return &ApiServer{
		storage: storage,
		prefix:  prefix,
	}
}

// HTTP Handler interface
func (server *ApiServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Printf("%s %s", req.Method, req.RequestURI)
	// ParseRequestURI函数解析rawurl为一个  URL结构体
	/*
		type URL struct {
		    Scheme   string
		    Opaque   string    // 编码后的不透明数据
		    User     *Userinfo // 用户名和密码信息
		    Host     string    // host或host:port
		    Path     string
		    RawQuery string // 编码后的查询字符串，没有'?'
		    Fragment string // 引用的片段（文档位置），没有'#'
		}
	*/
	url, err := url.ParseRequestURI(req.RequestURI)
	log.Printf("%s", url)
	if err != nil {
		server.error(err, w)
		return
	}
	if url.Path == "/index.html" || url.Path == "/" || url.Path == "" {
		server.handleIndex(w)
		return
	}
	if !strings.HasPrefix(url.Path, server.prefix) {
		server.notFound(req, w)
		return
	}
	requestParts := strings.Split(url.Path[len(server.prefix):], "/")[1:]
	log.Println("requestParts:", requestParts)
	if len(requestParts) < 1 {
		server.notFound(req, w)
		return
	}
	storage := server.storage[requestParts[0]]
	if storage == nil {
		server.notFound(req, w)
		return
	} else {
		server.handleREST(requestParts, url, req, w, storage)
	}
}

func (server *ApiServer) error(err error, w http.ResponseWriter) {
	w.WriteHeader(500)
	fmt.Fprintf(w, "Internal Error: %#v", err)
}

func (server *ApiServer) handleIndex(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	// TODO: serve this out of a file?
	data := "<html><body>Welcome to Kubernetes</body></html>"
	fmt.Fprint(w, data)
}

func (server *ApiServer) notFound(req *http.Request, w http.ResponseWriter) {
	w.WriteHeader(404)
	fmt.Fprintf(w, "Not Found: %#v", req)
}

func (server *ApiServer) handleREST(parts []string, url *url.URL, req *http.Request, w http.ResponseWriter, storage RESTStorage) {
	switch req.Method {
	case "GET":
		switch len(parts) {
		case 1:
			controllers, err := storage.List(url)
			if err != nil {
				server.error(err, w)
				return
			}
			server.write(200, controllers, w)
		case 2:
			task, err := storage.Get(parts[1])
			if err != nil {
				server.error(err, w)
				return
			}
			if task == nil {
				server.notFound(req, w)
				return
			}
			server.write(200, task, w)
		default:
			server.notFound(req, w)
		}
		return
	case "POST":
		if len(parts) != 1 {
			server.notFound(req, w)
			return
		}
		body, err := server.readBody(req)
		if err != nil {
			server.error(err, w)
			return
		}
		obj, err := storage.Extract(body)
		if err != nil {
			server.error(err, w)
			return
		}
		storage.Create(obj)
		server.write(200, obj, w)
		return
	case "DELETE":
		if len(parts) != 2 {
			server.notFound(req, w)
			return
		}
		err := storage.Delete(parts[1])
		if err != nil {
			server.error(err, w)
			return
		}
		server.write(200, Status{success: true}, w)
		return
	case "PUT":
		if len(parts) != 2 {
			server.notFound(req, w)
			return
		}
		body, err := server.readBody(req)
		if err != nil {
			server.error(err, w)
		}
		obj, err := storage.Extract(body)
		if err != nil {
			server.error(err, w)
			return
		}
		err = storage.Update(obj)
		if err != nil {
			server.error(err, w)
			return
		}
		server.write(200, obj, w)
		return
	default:
		server.notFound(req, w)
	}
}

func (server *ApiServer) write(statusCode int, object interface{}, w http.ResponseWriter) {
	w.WriteHeader(statusCode)
	output, err := json.MarshalIndent(object, "", "    ")
	if err != nil {
		server.error(err, w)
		return
	}
	w.Write(output)
}

func (server *ApiServer) readBody(req *http.Request) (string, error) {
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	return string(body), err
}

func (sr *RESTStorageData) List(*url.URL) (interface{}, error) {
	return "RESTStorageData  List is executed", nil
}

func (sr *RESTStorageData) Get(id string) (interface{}, error) {
	return "RESTStorageData Get task1" + id, nil
}

func (sr *RESTStorageData) Delete(id string) error {
	return nil
}

func (sr *RESTStorageData) Extract(body string) (interface{}, error) {
	return nil, nil
}

func (sr *RESTStorageData) Create(obj interface{}) error {
	return nil
}

func (sr *RESTStorageData) Update(obj interface{}) error {
	return nil
}
