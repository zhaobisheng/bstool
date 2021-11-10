package Stream

import (
	"bytes"
	"encoding/binary"
	//"encoding/hex"
	//"hex"
)

type StreamStruct struct {
	Postion int
	Raw     []byte
}

func NewStream() *StreamStruct {
	return &StreamStruct{Postion: 0}
}

func SetStream(data []byte) *StreamStruct {
	return &StreamStruct{Postion: 0, Raw: data}
}

func (stream *StreamStruct) ReadByte() byte {
	rs := stream.Raw[stream.Postion]
	stream.Postion += 1
	return rs
}

func (stream *StreamStruct) ReadBytes(length int) []byte {
	rs := stream.Raw[stream.Postion : stream.Postion+length]
	stream.Postion += length
	return rs
}

func (stream *StreamStruct) ReadInt() int {
	buf := stream.ReadBytes(4)
	rs := int(binary.BigEndian.Uint32(buf))
	//stream.Postion += 4
	return rs
}

func (stream *StreamStruct) ReadInt16() int {
	buf := stream.ReadBytes(2)
	rs := int(binary.BigEndian.Uint16(buf))
	return rs
}

func (stream *StreamStruct) WriteByte(data byte) {
	stream.Raw = append(stream.Raw, data)
}

func (stream *StreamStruct) WriteInt(data int) {
	x := int16(data)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	stream.WriteBytes(stream.Raw, bytesBuffer.Bytes())
}

func (stream *StreamStruct) WriteBytes(data ...[]byte) {
	stream.Raw = bytes.Join(data, []byte(""))
}

func (stream *StreamStruct) Find(data []byte) int {
	sLen := len(data)
	rawLen := len(stream.Raw)
	tempByte := make([]byte, 0)
	for i := 0; i < rawLen-sLen; i++ {
		tempByte = stream.Raw[i : i+sLen]
		if bytes.Equal(data, tempByte) {
			return i
		}
	}
	return -1
}

/*
func ArrEq(a, b []byte) bool {
    // If one is nil, the other must also be nil.
    if (a == nil) != (b == nil) {
        return false;
    }

    if len(a) != len(b) {
        return false
    }

    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }
    return true
}
*/
