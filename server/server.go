package server

import (
	"fmt"
	"log"
	"net"
	"reflect"
	"sync"

	"github.com/simple-rpc-golang/common"
	"github.com/simple-rpc-golang/stream"
	"github.com/simple-rpc-golang/transport"
)



type ServerImp struct {
	transporter transport.Transporter
	listener    net.Listener
	allMethod   sync.Map
	encoder     stream.Encoder
}

func NewServerImp(listener net.Listener, encoder stream.Encoder) *ServerImp {
	return &ServerImp{
		transporter: transport.NewTransporterImp(encoder),
		listener:    listener,
		encoder:     encoder,
	}
}

func (s *ServerImp) Register(method string, f interface{}) {
	common.IsValidFunc(method, f)
	s.allMethod.Store(method, f)
}

func (s *ServerImp) GetMethod(method string) (interface{}, error) {
	if f, ok := s.allMethod.Load(method); ok {
		return f, nil
	}
	return nil, fmt.Errorf("not found method: %v", method)
}

func (s *ServerImp) handle(conn net.Conn) error {
	in := &common.RPCData{}
	err := s.transporter.Recv(conn, in)
	if err != nil {
		return err
	}
	log.Printf("server收到请求:%s\n", in.Method)

	out := s.execute(in)
	err = s.transporter.Send(conn, out)
	if err != nil {
		return err
	}
	log.Printf("server请求处理完成:%s\n", in.Method)
	return nil
}

func (s *ServerImp) execute(req *common.RPCData) *common.RPCData {
	handlerFunc, err := s.GetMethod(req.Method)
	if err != nil {
		return common.RPCDataWithError(err)
	}
	fPointer := reflect.TypeOf(handlerFunc)

	var inData interface{}
	for i := 0; i < fPointer.NumIn(); i++ {
		var newObj reflect.Value
		if fPointer.In(i).Kind() == reflect.Ptr {
			newObj = reflect.New(fPointer.In(i))
		} else {
			newObj = reflect.New(fPointer.In(i))
		}
		inData = newObj.Interface()
	}
	// 指针先接住
	err = s.encoder.Decode(req.Params, inData)
	if err != nil {
		return common.RPCDataWithError(err)
	}

	out := reflect.ValueOf(handlerFunc).Call([]reflect.Value{reflect.ValueOf(inData).Elem()})
	outData := make([]interface{}, fPointer.NumOut())
	for id := range out {
		outData[id] = out[id].Interface()
	}

	// 处理error情况
	if outData[len(outData)-1] != nil {
		if e, ok := outData[len(outData)-1].(error); ok {
			return common.RPCDataWithError(e)
		}
	}
	// 正常返回
	outBytes, err := s.encoder.Encode(outData[0])
	if err != nil {
		return common.RPCDataWithError(err)
	}
	return common.RPCDataWithData(outBytes)
}

func (s *ServerImp) Run() {
	log.Println("rpc server is running...")
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println(err)
		}
		go func() {
			for {
				err := s.handle(conn)
				if err != nil {
					log.Printf("server handle error:%v\n", err)
				}
			}
		}()
	}
}
