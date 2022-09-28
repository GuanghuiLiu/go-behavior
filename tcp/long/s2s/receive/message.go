package receive

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/GuanghuiLiu/behavior/tcp/long/s2s"
	"io"
	"net"
)

// 不用面向对象，原因是：短链接，不需要process管理链接；数量大，不用每次链接都重复创建数据

func Unpack(conn *net.TCPConn) (*name_server.Message, error) {
	msg, err := UnpackHead(conn)
	if err != nil {
		return msg, err
	}
	return UnpackData(conn, msg)

}

func UnpackData(conn *net.TCPConn, msg *name_server.Message) (*name_server.Message, error) {

	var data []byte

	if msg.DataLen > 0 && msg.DataLen < name_server.MaxPacketSize {
		data = make([]byte, msg.DataLen)
		if _, err := io.ReadFull(conn, data); err != nil {
			fmt.Println("read msg data error ", err)
			return msg, err
		}
		msg.Data = data
	}
	return msg, nil
}

func UnpackHead(conn *net.TCPConn) (msg *name_server.Message, err error) {

	head := make([]byte, name_server.HeadLen)

	if _, err = io.ReadFull(conn, head); err != nil {
		fmt.Println("read msg head error ", err)
		return msg, err
	}

	headBuf := bytes.NewReader(head)

	msg = &name_server.Message{}

	if err = binary.Read(headBuf, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	if err = binary.Read(headBuf, binary.LittleEndian, &msg.ProtoID); err != nil {
		return nil, err
	}

	return msg, nil
}

func Pack(protoID uint16, dataBytes []byte) (out []byte, err error) {

	byteLen := uint16(len(dataBytes))

	// if byteLen > MaxPacketSize {
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
