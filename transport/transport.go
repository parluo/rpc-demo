package transport

import (
	"encoding/binary"
	"io"
	"log"
	"net"

	"github.com/simple-rpc-golang/common"
	"github.com/simple-rpc-golang/stream"
)


type TransporterImp struct {
	parser stream.Encoder
}

func NewTransporterImp(parser stream.Encoder) *TransporterImp {
	return &TransporterImp{
		parser: parser,
	}
}

func (t *TransporterImp) Send(conn net.Conn, in interface{}) error {
	inBytes, err := t.parser.Encode(in)
	if err != nil {
		return err
	}
	data := &common.RPCData{
		Params: inBytes,
	}
	msg, err := t.parser.Encode(data)
	if err != nil {
		return err
	}

	lenByte := [4]byte{}
	binary.BigEndian.PutUint32(lenByte[:], uint32(len(msg)))
	_, err = conn.Write(lenByte[:])
	if err != nil {
		return err
	}
	_, err = conn.Write(msg)
	return err
}

func (t *TransporterImp) Recv(conn net.Conn, out interface{}) error {
	lenByte := [4]byte{}
	_, err := io.ReadFull(conn, lenByte[:])
	if err != nil {
		if err != io.EOF {
			log.Printf("transport recv error:%v\n", err)
			return err
		}
	}
	dataLen := binary.BigEndian.Uint32(lenByte[:])
	dataByte := make([]byte, dataLen)
	_, err = io.ReadFull(conn, dataByte)
	if err != nil {
		return err
	}
	resp := &common.RPCData{}
	err = t.parser.Decode(dataByte, resp)
	if err != nil {
		return err
	}
	if resp.Status != nil && resp.Status.Code != 0 {
		return resp.Status
	}

	err = t.parser.Decode(resp.Params, out)
	if err != nil {
		return err
	}
	return nil
}
