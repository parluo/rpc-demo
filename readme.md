# RPC-Demo
一个简易的RPC demo实现。
结构如下
├── client		客户端
├── common 		server和client公共文件
├── namespace 	服务发现
├── server 		服务端
├── stream		序列化
└── transport	传输层

## start
```
go run main.go 
```

```golang
package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"time"

	"github.com/simple-rpc-golang/client"
	"github.com/simple-rpc-golang/server"
	"github.com/simple-rpc-golang/stream"
)

type Person struct {
	Name string
	Age  int
}

func GetPerson(id int) (*Person, error) {
	return &Person{
		Name: "wang",
		Age:  id,
	}, nil
}

func GetPersonName(p *Person) (string, error) {
	if p == nil {
		return "", fmt.Errorf("fail to get person")
	}
	return p.Name, nil
}

func run() {
	gob.Register(&Person{})
	listen, err := net.Listen("tcp", ":9001")
	if err != nil {
		panic(err)
	}
	encoder := &stream.GobEncoder{} // 可以注入不同的的序列化实现
	serverObj := server.NewServerImp(listen, encoder)
	serverObj.Register("GetName", GetPerson)
	serverObj.Register("GetPersonName", GetPersonName)
	go serverObj.Run()
	time.Sleep(2 * time.Second)

	conn, err := net.Dial("tcp", ":9001")
	if err != nil {
		panic(err)
	}
	clientObj := client.NewClientImp(conn, encoder)
	out := &Person{}
	e := clientObj.Invoke("GetName", 13, out)
	fmt.Println("result:", out, e)

	var name string
	err = clientObj.Invoke("GetPersonName", out, &name)
	fmt.Println("result:", err, name)
}

func main() {
	run()
}

```

## statement
- 限制了注入方法的入参只有一个
- 限制了注入方法的出参只有两个，且第二位为error
- 不要在入参和出参中使用interface类型，必须是明确的对象
- 支持gob，json序列化
- 服务发现 todo
