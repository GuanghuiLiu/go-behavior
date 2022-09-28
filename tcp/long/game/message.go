package game

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"google.golang.org/protobuf/proto"
)

type Packer struct {
	Name   string
	buffer buffer
	// data *CircleBinary
}

func NewPack(name string, conn *net.TCPConn) *Packer {
	return &Packer{Name: name,
		buffer: newBuffer(conn, MaxPackSize)} // buffer: newBuffer(conn, MaxPackSize)

}

func (a *Packer) CircleUnpack() (msg *Message, err error) {
	_, err = a.buffer.readFromReader()
	if err != nil {
		return
	}
	return a.doUnpack()
}
func (a *Packer) doUnpack() (msg *Message, err error) {
	msg = &Message{}
	dataLenBuf, err := a.buffer.readData(LengthDateLen)
	if err != nil {
		return
	}
	msg.DataLen = binary.LittleEndian.Uint32(dataLenBuf)
	ProtoIDBuf, err := a.buffer.readData(LengthProtoID)
	if err != nil {
		return
	}
	msg.ProtoID = binary.LittleEndian.Uint64(ProtoIDBuf)
	msg.Data, err = a.buffer.readData(msg.DataLen)
	return msg, err
}
func (a *Packer) Unpack(conn *net.TCPConn) (msg *Message, err error) {
	msg, err = a.UnpackHead(conn)
	if err != nil {
		return msg, err
	}
	return a.UnpackData(conn, msg)
}

func (a *Packer) UnpackHead(conn *net.TCPConn) (msg *Message, err error) {

	head := make([]byte, HeadLen)

	if _, err = io.ReadFull(conn, head); err != nil {
		fmt.Println(a.Name, "read msg head error:", err)
		return msg, err
	}

	headBuf := bytes.NewReader(head)

	msg = &Message{}

	if err = binary.Read(headBuf, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	// fmt.Println(r.Name, "read tail data len", msg.DataLen)

	if err = binary.Read(headBuf, binary.LittleEndian, &msg.ProtoID); err != nil {
		return nil, err
	}

	headBuf = nil

	return msg, nil
}
func (a *Packer) UnpackData(conn *net.TCPConn, msg *Message) (*Message, error) {

	var data []byte

label:
	if msg != nil && msg.DataLen > MaxDataSize {
		data = make([]byte, MaxDataSize)
		a.DoUnpackData(conn, msg, data)
		data = nil
		msg.DataLen -= MaxDataSize
		goto label
	}

	if msg != nil && msg.DataLen > 0 {
		data = make([]byte, msg.DataLen)
		a.DoUnpackData(conn, msg, data)
		msg.DataLen = 0
	}
	data = nil
	return msg, nil
}
func (a *Packer) DoUnpackData(conn *net.TCPConn, msg *Message, data []byte) (*Message, error) {

	if _, err := io.ReadFull(conn, data); err != nil {
		fmt.Println(a.Name, "read msg data error ", err)
		panic(err.(any))
		return msg, err
	}
	msg.Data = append(msg.Data, data...)
	return msg, nil
}

func (a *Packer) Pack(protoID uint64, dataBytes []byte) (out []byte, err error) {

	byteLen := uint32(len(dataBytes))

	// if byteLen > MaxDataSize {
	// 	err = errors.New("packet over size")
	// 	return
	// }
	outbuff := bytes.NewBuffer([]byte{})
	// 写Len
	if err = binary.Write(outbuff, binary.LittleEndian, byteLen); err != nil {
		return
	}
	// 写MsgID
	if err = binary.Write(outbuff, binary.LittleEndian, protoID); err != nil {
		return
	}

	//all pkg data
	if err = binary.Write(outbuff, binary.LittleEndian, dataBytes); err != nil {
		return
	}

	out = outbuff.Bytes()

	return
}
func (a *Packer) Send2Client(conn *net.TCPConn, protoID uint64, data []byte) error {

	sendData, err := a.Pack(protoID, data)
	if err != nil {
		return err
	}

	if _, e := conn.Write(sendData); e != nil {
		return e
	}

	return nil
}
func (a *Packer) Send2ClientPB(conn *net.TCPConn, protoID uint64, data proto.Message) error {

	binaryData, err := proto.Marshal(data)
	if err != nil {
		fmt.Println(fmt.Sprintf(a.Name, "marshaling error:  %s", err))
		return err
	}

	sendData, err := a.Pack(protoID, binaryData)

	if err == nil {
		conn.Write(sendData)
	} else {
		fmt.Println(err)
	}

	return err
}
