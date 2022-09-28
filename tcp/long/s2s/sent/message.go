package sent

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"

	s2s "github.com/GuanghuiLiu/behavior/tcp/long/s2s"
	"google.golang.org/protobuf/proto"
)

func Unpack(conn net.Conn) (*s2s.Message, error) {
	msg, err := UnpackHead(conn)
	if err != nil {
		return nil, err
	}
	return UnpackData(conn, msg)
}

func UnpackHead(conn net.Conn) (msg *s2s.Message, err error) {
	head := make([]byte, s2s.HeadLen)

	if _, err = io.ReadFull(conn, head); err != nil {
		fmt.Println("read msg head error ", err)
		return msg, err
	}
	headBuf := bytes.NewReader(head)
	msg = &s2s.Message{}
	if err = binary.Read(headBuf, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	if err = binary.Read(headBuf, binary.LittleEndian, &msg.ProtoID); err != nil {
		return nil, err
	}
	return msg, nil
}

// todo client 拆包，包过大时，分段拆包，参见server拆包
func UnpackData(conn net.Conn, msg *s2s.Message) (*s2s.Message, error) {
	var data []byte

	if msg != nil && msg.DataLen > 0 {
		data = make([]byte, msg.DataLen)
		if _, err := io.ReadFull(conn, data); err != nil {
			fmt.Println("read msg data error ", err)
			return msg, err
		}
		msg.Data = data
	}
	return msg, nil
}

func pack(protoID uint16, dataBytes []byte) (out []byte, err error) {

	byteLen := uint16(len(dataBytes))

	// if byteLen > s2s.MaxDataSize {
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

func SendMsg(conn net.Conn, protoID uint16, base proto.Message) error {

	// 进行编码
	binaryData, err := proto.Marshal(base)
	if err != nil {
		fmt.Println(fmt.Sprintf("marshaling error:  %s", err))
		return err
	}

	sendData, err := pack(protoID, binaryData)

	if err != nil {
		fmt.Println(fmt.Sprintf("marshaling error:  %s", err))
		return err
	}

	_, err = conn.Write(sendData)

	return err
}
