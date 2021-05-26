package pbf

import (
	"math"
	"reflect"
)

type TagType uint32

type WireType uint32

type LengthType = uint32

const (
	Varint  WireType = 0  // varint: int32, int64, uint32, uint64, sint32, sint64, bool, enum
	Fixed64 WireType = 1  // 64-bit: double, fixed64, sfixed64
	Bytes   WireType = 2  // length-delimited: string, bytes, embedded messages, packed repeated fields
	Fixed32 WireType = 5  // 32-bit: float, fixed32, sfixed32
	Unknown WireType = 99 // used for default setting in this library
)

func tagAndType(t TagType, w WireType) int {
	return int((uint32(t) << 3) | uint32(w))
}

const maxVarintBytes = 10

func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

func EncodeVarint(x uint64) []byte {
	var buf [maxVarintBytes]byte
	var n int
	for n = 0; x > 127; n++ {
		buf[n] = 0x80 | uint8(x&0x7F)
		x >>= 7
	}
	buf[n] = uint8(x)
	n++
	return buf[0:n]
}

func DecodeVarint2(buf []byte) (x uint64, n int) {
	for shift := uint(0); shift < 64; shift += 7 {
		if n >= len(buf) {
			return 0, 0
		}
		b := uint64(buf[n])
		n++
		x |= (b & 0x7F) << shift
		if (b & 0x80) == 0 {
			return x, n
		}
	}

	return 0, 0
}

func DecodeVarint(buf []byte) (x uint64) {
	i := 0
	if buf[i] < 0x80 {
		return uint64(buf[i])
	}

	var b uint64

	x = uint64(buf[i]) - 0x80
	i++

	b = uint64(buf[i])
	i++
	x += b << 7
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 7

	b = uint64(buf[i])
	i++
	x += b << 14
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 14

	b = uint64(buf[i])
	i++
	x += b << 21
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 21

	b = uint64(buf[i])
	i++
	x += b << 28
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 28

	b = uint64(buf[i])
	i++
	x += b << 35
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 35

	b = uint64(buf[i])
	i++
	x += b << 42
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 42

	b = uint64(buf[i])
	i++
	x += b << 49
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 49

	b = uint64(buf[i])
	i++
	x += b << 56
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 56

	b = uint64(buf[i])
	i++
	x += b << 63
	if b&0x80 == 0 {
		goto done
	}

	return 0

done:
	return x
}

func Key(x byte) (byte, byte) {
	val := x >> 3

	if int(x) >= 8 {
		x = x & 0x07

	} else {
		return val, x
	}

	if int(x) >= 16 {
		x = x & 0x0f

	} else {
		return val, x
	}

	if int(x) >= 32 {
		x = x & 0x1f

	} else {
		return val, x
	}

	if int(x) >= 64 {
		x = x & 0x3f

	} else {
		return val, x
	}

	if int(x) >= 128 {
		x = x & 0x7f

	} else {
		return val, x
	}

	return val, x
}

func ReadUInt32(buf []byte) uint32 {
	val := len(buf)
	switch val {
	case 1:
		return uint32(buf[0])
	case 2:
		return uint32(buf[0]) | uint32(buf[1])<<8
	case 3:
		return uint32(buf[0]) | uint32(buf[1])<<8 | uint32(buf[2])<<16
	case 4:
		return uint32(buf[0]) | uint32(buf[1])<<8 | uint32(buf[2])<<16 | uint32(buf[3])<<24
	}

	return uint32(0)
}

func ReadInt32(buf []byte) int32 {
	val := len(buf)
	switch val {
	case 1:
		return int32(buf[0])
	case 2:
		return int32(buf[0]) | int32(buf[1])<<8
	case 3:
		return int32(buf[0]) | int32(buf[1])<<8 | int32(buf[2])<<16
	case 4:
		return int32(buf[0]) | int32(buf[1])<<8 | int32(buf[2])<<16 + int32(buf[3])<<24
	}
	return int32(0)
}

func ReadUint64(bytes []byte) uint64 {
	return DecodeVarint(bytes)
}

func ReadInt64(bytes []byte) int64 {
	return int64(DecodeVarint(bytes))
}

func EncodeUInt32(val uint32) []byte {
	var buf [4]byte
	buf[0] = byte(val)
	buf[1] = byte(val >> 8)
	buf[2] = byte(val >> 16)
	buf[3] = byte(val >> 24)
	return buf[:]
}

func EncodeInt32(val int32) []byte {
	var buf [4]byte
	buf[0] = byte(val)
	buf[1] = byte(val >> 8)
	buf[2] = byte(val >> 16)
	buf[3] = byte(val >> 24)
	return buf[:]
}

func EncodeInt64(val int64) []byte {
	var buf [8]byte
	buf[0] = byte(val)
	buf[1] = byte(val >> 8)
	buf[2] = byte(val >> 16)
	buf[3] = byte(val >> 24)
	buf[4] = byte(val >> 32)
	buf[5] = byte(val >> 40)
	buf[6] = byte(val >> 48)
	buf[7] = byte(val >> 56)
	return buf[:]
}

func EncodeUInt64(val uint64) []byte {
	var buf [8]byte
	buf[0] = byte(val)
	buf[1] = byte(val >> 8)
	buf[2] = byte(val >> 16)
	buf[3] = byte(val >> 24)
	buf[4] = byte(val >> 32)
	buf[5] = byte(val >> 40)
	buf[6] = byte(val >> 48)
	buf[7] = byte(val >> 56)
	return buf[:]
}

func EncodeVarint32(x uint32) []byte {
	var buf [4]byte
	var n int
	for n = 0; x > 127; n++ {
		buf[n] = 0x80 | uint8(x&0x7F)
		x >>= 7
	}
	buf[n] = uint8(x)
	n++
	return buf[0:n]
}

func AppendAll(b ...[]byte) []byte {
	total := 0
	for _, i := range b {
		total += len(i)
	}
	pos := 0
	totalbytes := make([]byte, total)
	for _, i := range b {
		for _, byteval := range i {
			totalbytes[pos] = byteval
			pos += 1
		}
	}
	return totalbytes
}

func WritePackedUint32_2(geom []uint32) []byte {
	buf := make([]byte, len(geom)*4+8)
	pos := 8
	for _, x := range geom {

		for x > 127 {
			buf[pos] = 0x80 | uint8(x&0x7F)
			x >>= 7
			pos++
		}
		buf[pos] = uint8(x)
		pos++
	}
	beg := EncodeVarint(uint64(pos - 8))
	startpos := 8 - len(beg)
	currentpos := startpos
	i := 0
	for currentpos < 8 {
		buf[currentpos] = beg[i]
		i++
		currentpos++
	}

	return buf[startpos:pos]
}

func WritePackedUint32(geom []uint32) []byte {
	buf := make([]byte, len(geom)*4+8)
	pos := 8
	for _, x := range geom {
		if x < 128 {
			buf[pos] = uint8(x)
			x >>= 7
			pos++
		} else if x < 16384 {
			buf[pos] = 0x80 | uint8(x&0x7F)
			x >>= 7
			pos++
			buf[pos] = uint8(x)
			pos++
		} else if x < 2097152 {
			buf[pos] = 0x80 | uint8(x&0x7F)
			x >>= 7
			pos++
			buf[pos] = 0x80 | uint8(x&0x7F)
			x >>= 7
			pos++
			buf[pos] = uint8(x)
			pos++
		} else {
			buf[pos] = 0x80 | uint8(x&0x7F)
			x >>= 7
			pos++
			buf[pos] = 0x80 | uint8(x&0x7F)
			x >>= 7
			pos++
			buf[pos] = 0x80 | uint8(x&0x7F)
			x >>= 7
			pos++
			buf[pos] = uint8(x)
			pos++
		}
	}
	beg := EncodeVarint(uint64(pos - 8))
	startpos := 8 - len(beg)
	currentpos := startpos
	i := 0
	for currentpos < 8 {
		buf[currentpos] = beg[i]
		i++
		currentpos++
	}

	return buf[startpos:pos]
}

func WritePackedUint32_3(geom []uint32) []byte {
	count := 0
	for _, x := range geom {
		if x < 128 {
			count += 1
		} else if x < 16384 {
			count += 2
		} else if x < 2097152 {
			count += 3
		} else {
			count += 4
		}
	}

	beg := EncodeVarint(uint64(count))
	buf := make([]byte, len(beg)+count)
	for pos, i := range beg {
		buf[pos] = i
	}

	pos := len(beg)
	for _, x := range geom {

		for x > 127 {
			buf[pos] = 0x80 | uint8(x&0x7F)
			x >>= 7
			pos++
		}
		buf[pos] = uint8(x)
		pos++
	}
	return buf
}

func EncodeVarint_Value(x uint64, typeint int) []byte {
	var buf [maxVarintBytes]byte
	var n int
	for n = 0; x > 127; n++ {
		buf[n] = 0x80 | uint8(x&0x7F)
		x >>= 7
	}
	buf[n] = uint8(x)
	n++
	total := []byte{34, byte(n + 1), byte(typeint)}
	return append(total, buf[0:n]...)
}

func FloatVal32(f float32) []byte {
	buf := make([]byte, 4)
	n := math.Float32bits(f)
	buf[3] = byte(n >> 24)
	buf[2] = byte(n >> 16)
	buf[1] = byte(n >> 8)
	buf[0] = byte(n)
	return append([]byte{34, 5, 21}, buf...)
}

func FloatVal64(f float64) []byte {
	buf := make([]byte, 8)
	n := math.Float64bits(f)
	buf[7] = byte(n >> 56)
	buf[6] = byte(n >> 48)
	buf[5] = byte(n >> 40)
	buf[4] = byte(n >> 32)
	buf[3] = byte(n >> 24)
	buf[2] = byte(n >> 16)
	buf[1] = byte(n >> 8)
	buf[0] = byte(n)
	return append([]byte{34, 9, 25}, buf...)
}

func WriteValue(value interface{}) []byte {
	vv := reflect.ValueOf(value)
	kd := vv.Kind()

	switch kd {
	case reflect.String:
		if len(vv.String()) > 0 {
			size := uint64(len(vv.String()))
			size_bytes := EncodeVarint(size)
			bytevals := []byte{10}
			bytevals = append(bytevals, size_bytes...)
			bytevals = append(bytevals, []byte(vv.String())...)
			bytevals = append(EncodeVarint(uint64(len(bytevals))), bytevals...)
			return append([]byte{34}, bytevals...)
		}
	case reflect.Float32:
		return FloatVal32(float32(vv.Float()))
	case reflect.Float64:
		return FloatVal64(float64(vv.Float()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return EncodeVarint_Value(uint64(vv.Int()), 32)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return EncodeVarint_Value(uint64(vv.Uint()), 40)
	case reflect.Bool:
		if vv.Bool() == true {
			return []byte{34, 2, 56, 1}
		} else if vv.Bool() == false {
			return []byte{34, 2, 56, 0}
		}
	}

	return []byte{}
}
