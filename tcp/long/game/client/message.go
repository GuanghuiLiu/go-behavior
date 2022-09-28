package client

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/GuanghuiLiu/behavior/tcp/long/game"
	"google.golang.org/protobuf/proto"
	"io"
	"net"
)

type Message struct {
	Name string
}

func NewMessage(name string) *Message {
	return &Message{Name: name}
}
func (m *Message) UnpackHead(conn net.Conn) (msg *game.Message, err error) {
	head := make([]byte, game.HeadLen)

	if _, err = io.ReadFull(conn, head); err != nil {
		fmt.Println(m.Name, "read msg head error ", err)
		return msg, err
	}
	headBuf := bytes.NewReader(head)
	msg = &game.Message{}
	if err = binary.Read(headBuf, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	if err = binary.Read(headBuf, binary.LittleEndian, &msg.ProtoID); err != nil {
		return nil, err
	}
	return msg, nil
}

// todo client 拆包，包过大时，分段拆包，参见server拆包
func (m *Message) UnpackData(conn net.Conn, msg *game.Message) (*game.Message, error) {
	var data []byte

	if msg != nil && msg.DataLen > 0 {
		data = make([]byte, msg.DataLen)
		if _, err := io.ReadFull(conn, data); err != nil {
			fmt.Println(m.Name, "read msg data error ", err)
			return msg, err
		}
		msg.Data = data
	}
	return msg, nil
}
func (m *Message) Pack(protoID uint64, dataBytes []byte) (out []byte, err error) {

	byteLen := uint32(len(dataBytes))

	// if byteLen > game.MaxDataSize {
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

func (m *Message) SendMsg(conn net.Conn, protoID uint64, data proto.Message) {

	// 进行编码
	binaryData, err := proto.Marshal(data)
	if err != nil {
		fmt.Println(fmt.Sprintf("marshaling error:  %s", err))
		return
	}

	sendData, err := m.Pack(protoID, binaryData)

	// fmt.Println(m.Name, "sent", protoID, data)

	if err == nil {
		if _, e := conn.Write(sendData); e != nil {
			fmt.Println(e)
		}
	} else {
		fmt.Println(err)
	}
	return
}
