package pbf

import (
	"bytes"
)

const (
	BufferSize = 16
)

type Writer struct {
	Pbf    *bytes.Buffer
	Pos    int
	Length int
}

func (pbf *Writer) WriteKey(tag TagType, w WireType) {
	pbf.WriteVarint2(tagAndType(tag, w))
}

func (pbf *Writer) WriteVarint2(v int) {
}

func (pbf *Writer) WriteSVarint(v float64) {
}

func (pbf *Writer) WriteSVarintPower(v float64) {
}

func (pbf *Writer) realloc(min int) {
	length := pbf.Length
	if length <= 0 {
		length = BufferSize
	}

	for length < pbf.Pos+min {
		length *= 2
	}
	grow := length - pbf.Length

	if length != pbf.Length && grow > 0 {
		pbf.Pbf.Grow(grow)
		pbf.Length = length
	}

}

func (pbf *Writer) Finish() []byte {
	pbf.Length = pbf.Pos
	pbf.Pos = 0
	bytes := pbf.Pbf.Bytes()
	return bytes[0:pbf.Length]
}

func (pbf *Writer) WriteFixed32(v uint32) {
	pbf.realloc(4)
	buf := EncodeUInt32(v)
	pbf.Pbf.Write(buf)
	pbf.Pos += 4
}

func (pbf *Writer) WriteSFixed32(v int32) {
	pbf.realloc(4)
	buf := EncodeInt32(v)
	pbf.Pbf.Write(buf)
	pbf.Pos += 4
}

func (pbf *Writer) WriteFixed64(v uint64) {
	pbf.realloc(8)
	buf := EncodeUInt64(v)
	pbf.Pbf.Write(buf)
	pbf.Pos += 8
}

func (pbf *Writer) WriteSFixed64(v int64) {
	pbf.realloc(8)
	buf := EncodeInt64(v)
	pbf.Pbf.Write(buf)
	pbf.Pos += 8
}

func (pbf *Writer) writeValue(v interface{}) {
	buf := WriteValue(v)
	n := len(buf)
	pbf.realloc(n)
	pbf.Pbf.Write(buf)
	pbf.Pos += n
}

func (pbf *Writer) WriteUInt32(v uint32) {
	pbf.writeValue(v)
}

func (pbf *Writer) WriteVarint(v int) {
	pbf.writeValue(v)
}

func (pbf *Writer) WriteInt32(v int32) {
	pbf.writeValue(v)
}

func (pbf *Writer) WriteUInt64(v uint64) {
	pbf.writeValue(v)
}

func (pbf *Writer) WriteInt64(v int64) {
	pbf.writeValue(v)
}

func (pbf *Writer) WriteDouble(v float64) {
	pbf.writeValue(v)
}

func (pbf *Writer) WriteFloat(v float32) {
	pbf.writeValue(v)
}

func (pbf *Writer) WriteString(s string) {
	pbf.writeValue(s)
}

func (pbf *Writer) WriteBool(b bool) {
	pbf.writeValue(b)
}

func (pbf *Writer) writeRawMessage(fn func(v interface{}) (int, []byte), v interface{}) {
	l, pack := fn(v)

	n := len(pack)
	pbf.WriteVarint(l)

	pbf.realloc(n)
	pbf.Pbf.Write(pack)
	pbf.Pos += n
}

func (pbf *Writer) writeMessage(tag TagType, fn func(v interface{}) (int, []byte), v interface{}) {
	pbf.WriteKey(tag, Bytes)
	pbf.writeRawMessage(fn, v)
}

func (pbf *Writer) WritePackedUint32(tag TagType, p []uint32) {

}

func (pbf *Writer) WritePackedUint32_2(tag TagType, p []uint32) {
}

func (pbf *Writer) WritePackedUint32_3(tag TagType, p []uint32) {
}

func (pbf *Writer) WritePackedInt32(tag TagType, p []int32) {
}

func (pbf *Writer) WritePackedUInt64(tag TagType, p []int) {
}

func NewWriter() *Writer {
	buf := make([]byte, BufferSize)
	return &Writer{Pbf: bytes.NewBuffer(buf), Length: len(buf)}
}
