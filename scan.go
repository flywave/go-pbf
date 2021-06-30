package pbf

import (
	"bufio"
	"io"
)

type Scanner struct {
	Scanner        *bufio.Scanner
	BoolVal        bool
	EndBool        bool
	TotalPosition  int
	BufferPosition int
	increment      int
	SizeBuffer     int
}

var SizeBuffer = 64 * 1028
var SizeBufferLarge = 64 * 1028 * 1028

func NewScanner(ioreader io.Reader) *Scanner {
	scanner := bufio.NewScanner(ioreader)
	buf := make([]byte, SizeBuffer)
	scanner.Buffer(buf, SizeBuffer)
	scannerval := &Scanner{Scanner: scanner, BoolVal: true, SizeBuffer: SizeBuffer}
	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if len(data) < scannerval.increment {
			token = make([]byte, scannerval.increment)
			copy(token, data[:scannerval.increment])
			advance = len(data)
		} else {
			token = make([]byte, scannerval.increment)
			copy(token, data)
			advance = scannerval.increment
		}
		if atEOF {
			scannerval.EndBool = true
		}
		return
	}

	scannerval.Scanner.Split(split)
	return scannerval
}

func NewScannerSize(ioreader io.Reader, size_buffer int) *Scanner {
	scanner := bufio.NewScanner(ioreader)
	buf := make([]byte, size_buffer)
	scanner.Buffer(buf, size_buffer)
	scannerval := &Scanner{Scanner: scanner, BoolVal: true, SizeBuffer: size_buffer}
	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if len(data) < scannerval.increment {
			token = make([]byte, scannerval.increment)
			copy(token, data[:scannerval.increment])
			advance = len(data)
		} else {
			token = make([]byte, scannerval.increment)
			copy(token, data)
			advance = scannerval.increment
		}
		if atEOF {
			scannerval.EndBool = true
		}
		return
	}

	scannerval.Scanner.Split(split)
	return scannerval
}

func (scanner *Scanner) Reset() {
	scanner.increment = 0
	scanner.BoolVal = true
	scanner.EndBool = false
	scanner.TotalPosition = 0
	scanner.BufferPosition = 0
}

func (scanner *Scanner) Scan() bool {
	scanner.GetIncrement(0)
	if scanner.EndBool {
		return false
	}
	return scanner.BoolVal
}

func (scanner *Scanner) GetIncrement(step int) []byte {
	scanner.TotalPosition += step

	buffer_left := scanner.SizeBuffer - scanner.BufferPosition

	if step > scanner.SizeBuffer {
		var newlist []byte
		if scanner.BufferPosition != 0 {
			scanner.increment = buffer_left
			scanner.BoolVal = scanner.Scanner.Scan()
			newlist = scanner.Scanner.Bytes()[:scanner.increment]
			scanner.BufferPosition = 0
		}

		for len(newlist) != step {
			remaining_bytes := step - len(newlist)
			if remaining_bytes > scanner.SizeBuffer {
				scanner.increment = scanner.SizeBuffer
				scanner.BoolVal = scanner.Scanner.Scan()
				tmpbytes := scanner.Scanner.Bytes()
				newlist = append(newlist, tmpbytes...)
			} else {
				scanner.increment = remaining_bytes
				scanner.BufferPosition = scanner.BufferPosition + scanner.increment
				scanner.BoolVal = scanner.Scanner.Scan()
				tmpbytes := scanner.Scanner.Bytes()[:scanner.increment]
				newlist = append(newlist, tmpbytes...)
			}
		}
		return newlist
	} else {
		var newlist []byte
		if buffer_left > step {
			scanner.increment = step
			scanner.BoolVal = scanner.Scanner.Scan()
			newlist = scanner.Scanner.Bytes()[:scanner.increment]
			scanner.BufferPosition = scanner.BufferPosition + scanner.increment
		} else {
			scanner.increment = buffer_left
			scanner.BoolVal = scanner.Scanner.Scan()
			newlist = scanner.Scanner.Bytes()[:scanner.increment]
			scanner.increment = step - buffer_left

			buffer_left = 0
			scanner.BufferPosition = scanner.increment
			scanner.BoolVal = scanner.Scanner.Scan()
			tmpbytes := scanner.Scanner.Bytes()[:scanner.increment]
			newlist = append(newlist, tmpbytes...)
		}
		return newlist
	}
}

func (scanner *Scanner) Protobuf() []byte {
	size := scanner.GetIncrement(1)
	size = scanner.GetIncrement(1)
	size_bytes := []byte{size[0]}
	for size[0] > 127 {
		size = scanner.GetIncrement(1)
		size_bytes = append(size_bytes, size[0])
	}
	size_protobuf := int(DecodeVarint(size_bytes))
	return scanner.GetIncrement(size_protobuf)
}

func (scanner *Scanner) ProtobufIndicies() ([]byte, [2]int) {
	size := scanner.GetIncrement(1)
	size = scanner.GetIncrement(1)
	size_bytes := []byte{size[0]}
	for size[0] > 127 {
		size = scanner.GetIncrement(1)
		size_bytes = append(size_bytes, size[0])
	}
	pos1 := scanner.TotalPosition
	size_protobuf := int(DecodeVarint(size_bytes))

	return scanner.GetIncrement(size_protobuf), [2]int{pos1, scanner.TotalPosition}
}
