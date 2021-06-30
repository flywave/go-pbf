package pbf

import (
	"math"
)

const (
	smallBufferSize = 64
)

type Writer struct {
	Pbf    []byte
	Pos    int
	Length int
}

func (pbf *Writer) WriteTag(tag TagType, w WireType) {
	pbf.writeValue(tagAndType(tag, w))
}

func (pbf *Writer) tryGrowByReslice(n int) (int, bool) {
	if l := pbf.len(); n <= pbf.cap()-l {
		pbf.Pbf = pbf.Pbf[:l+n]
		return l, true
	}
	return 0, false
}

func (pbf *Writer) len() int { return len(pbf.Pbf) }

func (pbf *Writer) cap() int { return cap(pbf.Pbf) }

func (pbf *Writer) Reset() {
	pbf.Pbf = pbf.Pbf[:0]
	pbf.Pos = 0
}

func (pbf *Writer) realloc(min int) int {
	length := pbf.Length
	if length <= 0 {
		length = smallBufferSize
	}

	for length < pbf.Pos+min {
		length *= 2
	}

	if i, ok := pbf.tryGrowByReslice(min); ok {
		pbf.Length += min
		return i
	}

	if pbf.Pbf == nil && min <= smallBufferSize {
		pbf.Pbf = make([]byte, min, smallBufferSize)
		pbf.Pos = 0
		pbf.Length = min
		return 0
	}

	buf := make([]byte, length)
	copy(buf, pbf.Pbf[:pbf.Pos])
	pbf.Pbf = buf[:pbf.Pos+min]
	pbf.Length = length
	return pbf.Pos
}

func (pbf *Writer) Finish() []byte {
	pbf.Length = pbf.Pos
	pbf.Pos = 0
	bytes := pbf.Pbf
	return bytes[:pbf.Length]
}

func (pbf *Writer) writeFixed32(v uint32) {
	i := pbf.realloc(4)
	buf := EncodeUInt32(v)
	copy(pbf.Pbf[i:], buf)
	pbf.Pos += 4
}

func (pbf *Writer) writeSFixed32(v int32) {
	i := pbf.realloc(4)
	buf := EncodeInt32(v)
	copy(pbf.Pbf[i:], buf)
	pbf.Pos += 4
}

func (pbf *Writer) writeFixed64(v uint64) {
	i := pbf.realloc(8)
	buf := EncodeUInt64(v)
	copy(pbf.Pbf[i:], buf)
	pbf.Pos += 8
}

func (pbf *Writer) writeSFixed64(v int64) {
	i := pbf.realloc(8)
	buf := EncodeInt64(v)
	copy(pbf.Pbf[i:], buf)
	pbf.Pos += 8
}

func (pbf *Writer) WriteFixed32(tag TagType, v uint32) {
	pbf.WriteTag(tag, Fixed32)
	pbf.writeFixed32(v)
}

func (pbf *Writer) WriteSFixed32(tag TagType, v int32) {
	pbf.WriteTag(tag, Fixed32)
	pbf.writeSFixed32(v)
}

func (pbf *Writer) WriteFixed64(tag TagType, v uint64) {
	pbf.WriteTag(tag, Fixed64)
	pbf.writeFixed64(v)
}

func (pbf *Writer) WriteSFixed64(tag TagType, v int64) {
	pbf.WriteTag(tag, Fixed64)
	pbf.writeSFixed64(v)
}

func (pbf *Writer) writeValue(v interface{}) {
	buf := WriteValue(v)
	n := len(buf)
	i := pbf.realloc(n)
	copy(pbf.Pbf[i:], buf)
	pbf.Pos += n
}

func (pbf *Writer) WriteValue(tag TagType, v interface{}) {
	pbf.WriteTag(tag, Varint)
	pbf.writeValue(v)
}

func (pbf *Writer) WriteUInt32(tag TagType, v uint32) {
	pbf.WriteTag(tag, Varint)
	pbf.writeValue(v)
}

func (pbf *Writer) WriteVarint(tag TagType, v int) {
	pbf.WriteTag(tag, Varint)
	pbf.writeValue(v)
}

func (pbf *Writer) WriteInt32(tag TagType, v int32) {
	pbf.WriteTag(tag, Varint)
	pbf.writeValue(v)
}

func (pbf *Writer) WriteUInt64(tag TagType, v uint64) {
	pbf.WriteTag(tag, Varint)
	pbf.writeValue(v)
}

func (pbf *Writer) WriteInt64(tag TagType, v int64) {
	pbf.WriteTag(tag, Varint)
	pbf.writeValue(v)
}

func (pbf *Writer) WriteDouble(tag TagType, v float64) {
	pbf.WriteTag(tag, Fixed64)
	pbf.writeValue(v)
}

func (pbf *Writer) WriteFloat(tag TagType, v float32) {
	pbf.WriteTag(tag, Fixed32)
	pbf.writeValue(v)
}

func (pbf *Writer) WriteBool(tag TagType, b bool) {
	pbf.WriteTag(tag, Varint)
	pbf.writeValue(b)
}

func (pbf *Writer) makeRoomForExtraLength(startPos int, len int) {
	var extraLen int
	if len <= 0x3fff {
		extraLen = 1
	} else if len <= 0x1fffff {
		extraLen = 2
	} else if len <= 0xfffffff {
		extraLen = 3
	} else {
		extraLen = int(math.Floor(math.Log(float64(len)) / (math.Ln2 * 7)))
	}

	pbf.realloc(extraLen)
	for i := pbf.Pos - 1; i >= startPos; i-- {
		pbf.Pbf[i+extraLen] = pbf.Pbf[i]
	}
}

func (pbf *Writer) writeString(s string) {
	pbf.realloc(len(s) * 4)

	pbf.Pos++

	startPos := pbf.Pos

	pbf.Pos = writeUtf8(pbf.Pbf, []byte(s), pbf.Pos)
	l := pbf.Pos - startPos

	if l >= 0x80 {
		pbf.makeRoomForExtraLength(startPos, l)
	}

	buf := WriteValue(l)
	i := startPos - 1
	copy(pbf.Pbf[i:], buf)
}

func (pbf *Writer) WriteString(tag TagType, s string) {
	pbf.WriteTag(tag, Bytes)
	pbf.writeString(s)
}

func (pbf *Writer) writeRawMessage(fn func(w *Writer)) {
	pbf.realloc(1)
	pbf.Pos++

	startPos := pbf.Pos

	fn(pbf)
	l := pbf.Pos - startPos

	if l >= 0x80 {
		pbf.makeRoomForExtraLength(startPos, l)
	}

	buf := WriteValue(l)
	i := startPos - 1
	copy(pbf.Pbf[i:], buf)
}

func (pbf *Writer) writeMessage(tag TagType, fn func(w *Writer)) {
	pbf.WriteTag(tag, Bytes)
	pbf.writeRawMessage(fn)
}

func (pbf *Writer) WriteMessage(tag TagType, fn func(w *Writer)) {
	pbf.writeMessage(tag, fn)
}

func (pbf *Writer) WritePackedVarint(tag TagType, p []int) {
	pbf.WriteTag(tag, Bytes)
	pbf.writeRawMessage(func(w *Writer) {
		for i := range p {
			pbf.writeValue(p[i])
		}
	})
}

func (pbf *Writer) WritePackedBoolean(tag TagType, p []bool) {
	pbf.WriteTag(tag, Bytes)
	pbf.writeRawMessage(func(w *Writer) {
		for i := range p {
			pbf.writeValue(p[i])
		}
	})
}

func (pbf *Writer) WritePackedUInt32(tag TagType, p []uint32) {
	pbf.WriteTag(tag, Bytes)
	pbf.writeRawMessage(func(w *Writer) {
		for i := range p {
			pbf.writeValue(p[i])
		}
	})
}

func (pbf *Writer) WritePackedInt32(tag TagType, p []int32) {
	pbf.WriteTag(tag, Bytes)
	pbf.writeRawMessage(func(w *Writer) {
		for i := range p {
			pbf.writeValue(p[i])
		}
	})
}

func (pbf *Writer) WritePackedUInt64(tag TagType, p []uint64) {
	pbf.WriteTag(tag, Bytes)
	pbf.writeRawMessage(func(w *Writer) {
		for i := range p {
			pbf.writeValue(p[i])
		}
	})
}

func (pbf *Writer) WritePackedInt64(tag TagType, p []int64) {
	pbf.WriteTag(tag, Bytes)
	pbf.writeRawMessage(func(w *Writer) {
		for i := range p {
			pbf.writeValue(p[i])
		}
	})
}

func (pbf *Writer) WritePackedDouble(tag TagType, p []float64) {
	pbf.WriteTag(tag, Bytes)
	pbf.writeRawMessage(func(w *Writer) {
		for i := range p {
			pbf.writeValue(p[i])
		}
	})
}

func (pbf *Writer) WritePackedFloat(tag TagType, p []float32) {
	pbf.WriteTag(tag, Bytes)
	pbf.writeRawMessage(func(w *Writer) {
		for i := range p {
			pbf.writeValue(p[i])
		}
	})
}

func (pbf *Writer) WriteRaw(buf []byte) {
	n := len(buf)
	i := pbf.realloc(n)
	copy(pbf.Pbf[i:], buf)
	pbf.Pos += n
}

func NewWriter() *Writer {
	return &Writer{Pbf: nil, Length: 0}
}
