package main

import (
	"StK8s/util"
	"flag"
	"fmt"
)

//注意：port、address、apiPrefix均为指针类型
var (
	port           = flag.Uint("port", 8080, "The port to listen on.  Default 8080.")
	address        = flag.String("address", "127.0.0.1", "The address on the local server to listen to. Default 127.0.0.1")
	apiPrefix      = flag.String("api_prefix", "/api/v1beta1", "The prefix for API requests on the server. Default '/api/v1beta1'")
	etcdServerList util.StringList
)

func init() {
	//Var方法使用指定的名字、使用信息注册一个flag。该flag的类型和值由第一个参数表示，该参数应实现了Value接口。
	flag.Var(&etcdServerList, "etcd_servers", "Value type is string; Servers for the etcd (http://ip:port), comma separated")
}

func main() {
	//从os.Args[1:]中解析注册的flag。 os.Args []string 保管了命令行参数，第一个是程序名。
	flag.Parse()

	fmt.Printf("port type is %T \n", port)

	fmt.Println("port:", *port)
	fmt.Println("address:", *address)
	fmt.Println("apiPrefix:", *apiPrefix)
	fmt.Println("etcdServerList:", etcdServerList)
}
