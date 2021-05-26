package pbf

import (
	"math"
)

type Reader struct {
	Pbf    []byte
	Pos    int
	Length int
}

var powerfactor = math.Pow(10.0, 7.0)

func (pbf *Reader) ReadKey() (TagType, WireType) {
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

func (pbf *Reader) ReadSVarintPower() float64 {
	num := int(pbf.ReadVarint())
	if num%2 == 1 {
		return float64((num+1)/-2) / powerfactor
	} else {
		return float64(num/2) / powerfactor
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

func (pbf *Reader) ReadPoint(endpos int) []float64 {
	for pbf.Pos < endpos {
		x := pbf.ReadSVarintPower()
		y := pbf.ReadSVarintPower()
		return []float64{Round(x, .5, 7), Round(y, .5, 7)}
	}
	return []float64{}
}

func (pbf *Reader) ReadLine(num int, endpos int) [][]float64 {
	var x, y float64
	if num == 0 {

		for startpos := pbf.Pos; startpos < endpos; startpos++ {
			if pbf.Pbf[startpos] <= 127 {
				num += 1
			}
		}
		newlist := make([][]float64, num/2)

		for i := 0; i < num/2; i++ {
			x += pbf.ReadSVarintPower()
			y += pbf.ReadSVarintPower()
			newlist[i] = []float64{Round(x, .5, 7), Round(y, .5, 7)}
		}

		return newlist
	} else {
		newlist := make([][]float64, num/2)

		for i := 0; i < num/2; i++ {
			x += pbf.ReadSVarintPower()
			y += pbf.ReadSVarintPower()

			newlist[i] = []float64{Round(x, .5, 7), Round(y, .5, 7)}

		}
		return newlist
	}
}

func (pbf *Reader) ReadPolygon(endpos int) [][][]float64 {
	polygon := [][][]float64{}
	for pbf.Pos < endpos {
		num := pbf.ReadVarint()
		polygon = append(polygon, pbf.ReadLine(num, endpos))
	}
	return polygon
}

func (pbf *Reader) ReadMultiPolygon(endpos int) [][][][]float64 {
	multipolygon := [][][][]float64{}
	for pbf.Pos < endpos {
		num_rings := pbf.ReadVarint()
		polygon := make([][][]float64, num_rings)
		for i := 0; i < num_rings; i++ {
			num := pbf.ReadVarint()
			polygon[i] = pbf.ReadLine(num, endpos)
		}
		multipolygon = append(multipolygon, polygon)
	}
	return multipolygon
}

func (pbf *Reader) ReadBoundingBox() []float64 {
	bb := make([]float64, 4)
	pbf.ReadVarint()
	bb[0] = float64(pbf.ReadSVarintPower())
	bb[1] = float64(pbf.ReadSVarintPower())
	bb[2] = float64(pbf.ReadSVarintPower())
	bb[3] = float64(pbf.ReadSVarintPower())
	return bb
}

func (pbf *Reader) ReadPackedInt32() []int32 {
	size := pbf.ReadVarint()
	arr := []int32{}
	endpos := pbf.Pos + size

	for pbf.Pos < endpos {
		arr = append(arr, int32(pbf.ReadUInt32()))
	}

	return arr
}

func (pbf *Reader) ReadPackedUInt64() []int {
	size := pbf.ReadVarint()
	arr := []int{}
	endpos := pbf.Pos + size

	for pbf.Pos < endpos {
		arr = append(arr, pbf.ReadVarint())
	}

	return arr
}

func NewReader(bytevals []byte) *Reader {
	return &Reader{Pbf: bytevals, Length: len(bytevals)}
}

func (pbf *Reader) ReadPackedUInt32_3() []uint32 {
	size := pbf.ReadVarint()
	endpos := pbf.Pos + size

	count := 0
	for startpos := pbf.Pos; startpos < endpos; startpos++ {
		if pbf.Pbf[startpos] <= 127 {
			count += 1
		}

	}

	arr := make([]uint32, count)

	for pos := 0; pbf.Pos < endpos; pos++ {
		arr[pos] = pbf.ReadUInt32()
	}

	return arr
}

func (pbf *Reader) ReadPackedUInt32_2() []uint32 {
	size := pbf.ReadVarint()

	arr := []uint32{}
	endpos := pbf.Pos + size

	for pbf.Pos < endpos {
		arr = append(arr, pbf.ReadUInt32())

	}

	return arr
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
