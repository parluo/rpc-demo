package client

import (
	"net"

	"github.com/simple-rpc-golang/common"
	"github.com/simple-rpc-golang/stream"
	"github.com/simple-rpc-golang/transport"
)

type Client interface {
	Invoke(method string, in interface{}, out interface{}) error
}

type ClientImp struct {
	conn        net.Conn
	transporter transport.Transporter
	encoder     stream.Encoder
}

func NewClientImp(conn net.Conn, encoder stream.Encoder) *ClientImp {
	return &ClientImp{
		conn:        conn,
		transporter: transport.NewTransporterImp(encoder),
		encoder:     encoder,
	}
}

func handleError(err interface{}) *common.RPCError {
	if e, ok := err.(interface {
		GetErrorCode() int
		Error() string
	}); ok {
		return &common.RPCError{
			Code:    e.GetErrorCode(),
			Message: e.Error(),
		}
	}
	if e, ok := err.(interface {
		Error() string
	}); ok {
		return common.WrapRPCError(e)
	}
	return common.UnknownRPCError
}

// in: 入参列表
// out: 出参列表
func (c *ClientImp) Invoke(method string, in interface{}, out interface{}) (err error) {
	inBytes, err := c.encoder.Encode(in)
	if err != nil {
		return err
	}
	req := &common.RPCData{
		Method: method,
		Params: inBytes,
	}

	err = c.transporter.Send(c.conn, req)
	if err != nil {
		return err
	}
	// todo 序列化后的字节数据如何返回为对应的数据
	resp := &common.RPCData{Method: method, Status: nil}
	err = c.transporter.Recv(c.conn, &resp)
	if err != nil {
		return err
	}
	if resp.Status != nil && resp.Status.Code != 0 {
		return resp.Status
	}
	err = c.encoder.Decode(resp.Params, out)
	if err != nil {
		return err
	}
	return nil
}
