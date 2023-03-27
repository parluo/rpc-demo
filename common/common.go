package common

import "reflect"

// 传输过程结构对象
type RPCData struct {
	Method string
	Params []byte
	Status *RPCError
}

func RPCDataWithError(err error) *RPCData {
	return &RPCData{
		Status: WrapRPCError(err),
	}
}

func RPCDataWithData(params []byte) *RPCData {
	return &RPCData{
		Params: params,
		Status: &RPCError{},
	}
}

var UnknownRPCError = &RPCError{Code: -2, Message: "unknown error"}

type RPCError struct {
	Message string
	Code    int
}

func (e *RPCError) Error() string {
	return e.Message
}
func (e *RPCError) GetErrorCode() int {
	return e.Code
}

func WrapRPCError(err error) *RPCError {
	return &RPCError{
		Message: err.Error(),
		Code:    -1,
	}
}

func IsValidFunc(method string, f interface{}) {
	fType := reflect.TypeOf(f)
	fPointer := fType
	if fType.Kind() == reflect.Interface || fType.Kind() == reflect.Ptr {
		fPointer = fType.Elem()
	}

	// 入参只能有一个参数，出参只有两个参数
	if fPointer.NumIn() != 1 {
		panic("only one input parameters is allowed")
	}
	// 出参必须有两个参数
	if fPointer.NumOut() != 2 {
		panic("only two output parameters are allowed")
	}
	// 注册的方法最后一个必须是error，且只能有一个error
	for id := 0; id < fPointer.NumOut(); id++ {
		if _, ok := reflect.New(fPointer.Out(id)).Interface().(*error); ok {
			if id != fPointer.NumOut()-1 {
				panic("only lasted output parameter allowed to be an error")
			}
		} else {
			if id == fPointer.NumOut()-1 {
				panic("lasted output parameter must be an error")
			}
		}

	}

}
