package pbf

import (
	"math"
)

type Reader struct {
	Pbf    []byte
	Pos    int
	Length int
}

func (pbf *Reader) ReadTag() (TagType, WireType) {
	var key byte
	var val byte
	if pbf.Pos > pbf.Length-1 {
		key, val = 100, 100
	} else {
		key, val = Key(pbf.Pbf[pbf.Pos])
		pbf.Pos += 1
	}
	return TagType(key), WireType(val)
}

func (pbf *Reader) Reset() {
	pbf.Pos = 0
}

func (pbf *Reader) ReadNext() (TagType, WireType) {
	_, w := pbf.ReadTag()
	pbf.skip(w)
	return pbf.ReadTag()
}

func (pbf *Reader) skip(val WireType) {
	if val == Varint {
		for pbf.Pbf[pbf.Pos] > 0x7f {
			pbf.Pos++
		}
	} else if val == Bytes {
		pbf.Pos = pbf.ReadVarint() + pbf.Pos
	} else if val == Fixed32 {
		pbf.Pos += 4
	} else if val == Fixed64 {
		pbf.Pos += 8
	} else {
		panic("error")
	}
}

func (pbf *Reader) ReadFields(readField func(tag TagType, tp WireType, result interface{}, pbf *Reader), result interface{}, endpos int) interface{} {
	if endpos <= 0 {
		endpos = pbf.Length
	}
	for pbf.Pos < endpos {
		tag, val := pbf.ReadTag()
		startPos := pbf.Pos

		readField(tag, val, result, pbf)

		if pbf.Pos == startPos {
			pbf.skip(val)
		}
	}
	return result
}

func (pbf *Reader) ReadMessage(readField func(tag TagType, tp WireType, result interface{}, pbf *Reader), result interface{}) interface{} {
	return pbf.ReadFields(readField, result, pbf.ReadVarint()+pbf.Pos)
}

func (pbf *Reader) ReadVarint2() int {
	if pbf.Pos+1 >= pbf.Length {
		if pbf.Pos+1 == pbf.Length {
			pbf.Pos += 1
		}
		return 0
	}
	if pbf.Pbf[pbf.Pos] <= 127 {
		pbf.Pos += 1
		return int(pbf.Pbf[pbf.Pos-1])
	}

	startPos := pbf.Pos
	for pbf.Pbf[pbf.Pos] > 127 {
		pbf.Pos += 1
	}
	pbf.Pos += 1
	return int(DecodeVarint(pbf.Pbf[startPos:pbf.Pos]))
}

func (pbf *Reader) ReadSVarint() float64 {
	num := int(pbf.ReadVarint())
	if num%2 == 1 {
		return float64((num + 1) / -2)
	} else {
		return float64(num / 2)
	}
}

func (pbf *Reader) Varint() []byte {
	startPos := pbf.Pos
	for pbf.Pbf[pbf.Pos] > 127 {
		pbf.Pos += 1
	}
	pbf.Pos += 1
	return pbf.Pbf[startPos:pbf.Pos]
}

func (pbf *Reader) ReadFixed32() uint32 {
	val := ReadUInt32(pbf.Pbf[pbf.Pos : pbf.Pos+4])

	pbf.Pos += 4
	return val
}

func (pbf *Reader) ReadUInt322() uint32 {
	return uint32(pbf.Pbf[pbf.Pos+0]) | uint32(pbf.Pbf[pbf.Pos+1])<<8 | uint32(pbf.Pbf[pbf.Pos+2])<<16 | uint32(pbf.Pbf[pbf.Pos+3])<<24

}

func (pbf *Reader) ReadUInt32() uint32 {
	if pbf.Pbf[pbf.Pos] < 128 {
		pbf.Pos += 1
		return uint32(pbf.Pbf[pbf.Pos-1])
	} else if pbf.Pbf[pbf.Pos+1] < 128 {
		a := pbf.Pos
		pbf.Pos += 2
		return uint32(pbf.Pbf[a]) | uint32(pbf.Pbf[a+1])<<8
	} else if pbf.Pbf[pbf.Pos+2] < 128 {
		a := pbf.Pos
		pbf.Pos += 3
		return uint32(pbf.Pbf[a]) | uint32(pbf.Pbf[a+1])<<8 | uint32(pbf.Pbf[a+2])<<16
	} else {
		a := pbf.Pos
		pbf.Pos += 4
		return uint32(pbf.Pbf[a]) | uint32(pbf.Pbf[a+1])<<8 | uint32(pbf.Pbf[a+2])<<16 | uint32(pbf.Pbf[a+3])<<24
	}
}

func (pbf *Reader) ReadVarint() int {
	left := pbf.Length - pbf.Pos
	if pbf.Pbf[pbf.Pos+0] < 128 && left >= 1 {
		a := pbf.Pos
		pbf.Pos += 1
		return int(DecodeVarint(pbf.Pbf[a:pbf.Pos]))
	} else if pbf.Pbf[pbf.Pos+1] < 128 && left >= 2 {
		a := pbf.Pos
		pbf.Pos += 2
		return int(DecodeVarint(pbf.Pbf[a:pbf.Pos]))
	} else if pbf.Pbf[pbf.Pos+2] < 128 && left >= 3 {
		a := pbf.Pos
		pbf.Pos += 3
		return int(DecodeVarint(pbf.Pbf[a:pbf.Pos]))

	} else if pbf.Pbf[pbf.Pos+3] < 128 && left >= 4 {
		a := pbf.Pos
		pbf.Pos += 4
		return int(DecodeVarint(pbf.Pbf[a:pbf.Pos]))
	} else if pbf.Pbf[pbf.Pos+4] < 128 && left >= 5 {
		a := pbf.Pos
		pbf.Pos += 5
		return int(DecodeVarint(pbf.Pbf[a:pbf.Pos]))
	} else if pbf.Pbf[pbf.Pos+5] < 128 && left >= 6 {
		a := pbf.Pos
		pbf.Pos += 6
		return int(DecodeVarint(pbf.Pbf[a:pbf.Pos]))
	} else if pbf.Pbf[pbf.Pos+6] < 128 && left >= 7 {
		a := pbf.Pos
		pbf.Pos += 7
		return int(DecodeVarint(pbf.Pbf[a:pbf.Pos]))
	} else if pbf.Pbf[pbf.Pos+7] < 128 && left >= 8 {
		a := pbf.Pos
		pbf.Pos += 8
		return int(DecodeVarint(pbf.Pbf[a:pbf.Pos]))

	}
	return int(0)
}

func (pbf *Reader) ReadInt32() int32 {
	return int32(pbf.Pbf[pbf.Pos+0]) | int32(pbf.Pbf[+1])<<8 | int32(pbf.Pbf[pbf.Pos+2])<<16 + int32(pbf.Pbf[pbf.Pos+3])<<24
}

func (pbf *Reader) ReadSFixed32() int32 {
	val := ReadInt32(pbf.Pbf[pbf.Pos : pbf.Pos+4])
	pbf.Pos += 4
	return val
}

func (pbf *Reader) ReadFixed64() uint64 {
	a := pbf.Pos
	val := uint64(pbf.Pbf[a]) | uint64(pbf.Pbf[a+1])<<8 | uint64(pbf.Pbf[a+2])<<16 | uint64(pbf.Pbf[a+3])<<24 | uint64(pbf.Pbf[a+4])<<32 | uint64(pbf.Pbf[a+5])<<40 | uint64(pbf.Pbf[a+6])<<48 | uint64(pbf.Pbf[a+7])<<56
	pbf.Pos += 8
	return val
}

func (pbf *Reader) ReadUInt64() uint64 {
	return ReadUint64(pbf.Varint())
}

func (pbf *Reader) ReadSFixed64() int64 {
	val := pbf.ReadFixed64()
	return int64(val)
}

func (pbf *Reader) ReadInt64() int64 {
	return ReadInt64(pbf.Varint())
}

func (pbf *Reader) ReadDouble() float64 {
	a := pbf.Pos
	pbf.Pos += 8
	return math.Float64frombits(uint64(pbf.Pbf[a]) | uint64(pbf.Pbf[a+1])<<8 | uint64(pbf.Pbf[a+2])<<16 | uint64(pbf.Pbf[a+3])<<24 | uint64(pbf.Pbf[a+4])<<32 | uint64(pbf.Pbf[a+5])<<40 | uint64(pbf.Pbf[a+6])<<48 | uint64(pbf.Pbf[a+7])<<56)
}

func (pbf *Reader) ReadFloat() float32 {
	a := pbf.Pos
	pbf.Pos += 4
	return math.Float32frombits(uint32(pbf.Pbf[a]) | uint32(pbf.Pbf[a+1])<<8 | uint32(pbf.Pbf[a+2])<<16 | uint32(pbf.Pbf[a+3])<<24)
}

func (pbf *Reader) ReadString() string {
	size := pbf.ReadVarint()
	stringval := string(pbf.Pbf[pbf.Pos : pbf.Pos+size])
	pbf.Pos += size
	return stringval
}

func (pbf *Reader) ReadBool() bool {
	if pbf.Pbf[pbf.Pos] == 1 {
		pbf.Pos += 1
		return true
	} else {
		pbf.Pos += 1
		return false
	}
}

func (pbf *Reader) ReadPacked() []uint32 {
	endpos := pbf.Pos + pbf.ReadVarint()

	vals := make([]uint32, pbf.Length)
	currentpos := 0

	for pbf.Pos < endpos {
		startpos := pbf.Pos

		for pbf.Pbf[pbf.Pos] > 127 {
			pbf.Pos += 1
		}
		pbf.Pos += 1

		switch pbf.Pos - startpos {

		case 1:
			vals[currentpos] = uint32(pbf.Pbf[startpos])
			currentpos += 1
		case 2:
			vals[currentpos] = (uint32(pbf.Pbf[startpos])) | (uint32(pbf.Pbf[startpos+1]) << 8)
			currentpos += 1
		case 3:
			vals[currentpos] = (uint32(pbf.Pbf[startpos])) | (uint32(pbf.Pbf[startpos+1]) << 8) | (uint32(pbf.Pbf[startpos+2]) << 16)
			currentpos += 1
		case 4:
			vals[currentpos] = (uint32(pbf.Pbf[startpos])) | (uint32(pbf.Pbf[startpos+1]) << 8) | (uint32(pbf.Pbf[startpos+2]) << 16) + (uint32(pbf.Pbf[startpos+3]) * 0x1000000)
			currentpos += 1
		}
	}
	return vals[:currentpos]
}

func (pbf *Reader) ReadPackedInt32() []int32 {
	size := pbf.ReadVarint()
	arr := make([]int32, size)
	endpos := pbf.Pos + size

	i := 0
	for pbf.Pos < endpos {
		arr[i] = int32(pbf.ReadUInt32())
		i++
	}

	return arr[:i]
}

func (pbf *Reader) ReadPackedInt64() []int64 {
	size := pbf.ReadVarint()
	arr := make([]int64, size)

	endpos := pbf.Pos + size
	i := 0
	for pbf.Pos < endpos {
		arr[i] = int64(pbf.ReadVarint())
		i++
	}

	return arr[:i]
}

func (pbf *Reader) ReadPackedUInt64() []uint64 {
	size := pbf.ReadVarint()
	arr := make([]uint64, size)

	endpos := pbf.Pos + size
	i := 0
	for pbf.Pos < endpos {
		arr[i] = uint64(pbf.ReadVarint())
		i++
	}
	return arr[:i]
}

func (pbf *Reader) ReadPackedUInt32() []uint32 {
	size := pbf.ReadVarint()

	arr := make([]uint32, size)
	endpos := pbf.Pos + size
	i := 0
	for pbf.Pos < endpos {
		arr[i] = uint32(pbf.ReadVarint())
		i++
	}
	return arr[:i]
}

func (pbf *Reader) ReadPackedFloat() []float32 {
	size := pbf.ReadVarint()

	arr := make([]float32, size)
	endpos := pbf.Pos + size
	i := 0
	for pbf.Pos < endpos {
		arr[i] = pbf.ReadFloat()
		i++
	}
	return arr[:i]
}

func (pbf *Reader) ReadPackedDouble() []float64 {
	size := pbf.ReadVarint()

	arr := make([]float64, size)
	endpos := pbf.Pos + size
	i := 0
	for pbf.Pos < endpos {
		arr[i] = pbf.ReadDouble()
		i++
	}
	return arr[:i]
}

func (pbf *Reader) ReadPackedBool() []bool {
	size := pbf.ReadVarint()

	arr := make([]bool, size)
	endpos := pbf.Pos + size
	i := 0
	for pbf.Pos < endpos {
		arr[i] = pbf.ReadBool()
		i++
	}
	return arr[:i]
}

func (pbf *Reader) ReadPackedFixed32() []uint32 {
	size := pbf.ReadVarint()

	arr := make([]uint32, size)
	endpos := pbf.Pos + size
	i := 0
	for pbf.Pos < endpos {
		arr[i] = pbf.ReadFixed32()
		i++
	}
	return arr[:i]
}

func (pbf *Reader) ReadPackedSFixed32() []int32 {
	size := pbf.ReadVarint()

	arr := make([]int32, size)
	endpos := pbf.Pos + size
	i := 0
	for pbf.Pos < endpos {
		arr[i] = pbf.ReadSFixed32()
		i++
	}
	return arr[:i]
}

func (pbf *Reader) ReadPackedFixed64() []uint64 {
	size := pbf.ReadVarint()

	arr := make([]uint64, size)
	endpos := pbf.Pos + size
	i := 0
	for pbf.Pos < endpos {
		arr[i] = pbf.ReadFixed64()
		i++
	}
	return arr[:i]
}

func (pbf *Reader) ReadPackedSFixed64() []int64 {
	size := pbf.ReadVarint()

	arr := make([]int64, size)
	endpos := pbf.Pos + size
	i := 0
	for pbf.Pos < endpos {
		arr[i] = pbf.ReadSFixed64()
		i++
	}
	return arr[:i]
}

func (pbf *Reader) ReadPackedVarint() []int {
	size := pbf.ReadVarint()

	arr := make([]int, size)
	endpos := pbf.Pos + size
	i := 0
	for pbf.Pos < endpos {
		arr[i] = pbf.ReadVarint()
		i++
	}
	return arr[:i]
}

func NewReader(bytevals []byte) *Reader {
	return &Reader{Pbf: bytevals, Length: len(bytevals)}
}
