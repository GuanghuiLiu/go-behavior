package game

import (
	"net"
)

type buffer struct {
	reader *net.TCPConn
	buf    []byte
	start  uint32
	end    uint32
}

func newBuffer(reader *net.TCPConn, len uint32) buffer {
	buf := make([]byte, len)
	return buffer{reader, buf, 0, 0}
}

func (b *buffer) Len() uint32 {
	return b.end - b.start
}

//将有用的字节前移
func (b *buffer) grow() {
	if b.start == 0 {
		return
	}
	if b.start == b.end {
		b.start = 0
		b.end = 0
		return
	}
	copy(b.buf, b.buf[b.start:b.end])
	b.end -= b.start
	b.start = 0
}

//从reader里面读取数据，如果reader阻塞，会发生阻塞
func (b *buffer) readFromReader() (uint32, error) {
	b.grow()
	n0, err := b.reader.Read(b.buf[b.end:])
	n := uint32(n0)
	if err != nil {
		return n, err
	}
	b.end += n
	return n, nil
}

func (b *buffer) read(n uint32) (buf []byte) {
	buf = b.buf[b.start : b.start+n]
	b.start += n
	return buf

}
func (b *buffer) readData(n uint32) (buf []byte, err error) {

label:
	length := b.end - b.start
	if n > length {
		newBuff := b.read(length)
		buf = append(buf, newBuff...)
		n -= length
		_, err = b.readFromReader()
		if err != nil {
			return
		}
		goto label
	}
	newBuff := b.read(n)
	return append(buf, newBuff...), err
}
