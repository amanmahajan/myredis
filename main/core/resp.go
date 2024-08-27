package core

import (
	"bytes"
	"errors"
	"fmt"
)

// Read the length value stored in the RESP byte stream
// Keep reading the integer value till we /r/n
// Example  10/r/n. In this case length = 10 and /r/n is the ending.
func readLength(b []byte) (int, int) {

	idx, length := 0, 0

	for idx = range b {
		val := b[idx]
		if !(val >= '0' && val <= '9') {
			return length, idx + 2
		}
		length = length*10 + int(val-'0')
	}

	return 0, 0
}

// Simple string is in the form of +abc\r\n
// So read the bytes till we reach \r
func readSimpleString(b []byte) (string, int, error) {
	idx := 1

	for ; b[idx] != '\r'; idx++ {

	}
	return string(b[1:idx]), idx + 2, nil

}

// Simple int64 is in the form of :123\r\n
// So read the bytes till we reach \r
func readInt64(b []byte) (int64, int, error) {
	idx := 1

	var val int64 = 0

	for ; b[idx] != '\r'; idx++ {

		val = val*10 + int64(b[idx]-'0')

	}
	return val, idx + 2, nil

}

func readBulkString(b []byte) (string, int, error) {
	idx := 1
	// Read length value first
	length, delta := readLength(b[idx:])

	idx += delta
	return string(b[idx : idx+length]), idx + length + 2, nil
}

// reads a RESP encoded array from data and returns
// the array, the delta, and the error
func readArray(data []byte) (interface{}, int, error) {
	// first character *
	pos := 1

	// reading the length
	count, delta := readLength(data[pos:])
	pos += delta

	var elems []interface{} = make([]interface{}, count)
	for i := range elems {
		elem, delta, err := DecodeOne(data[pos:])
		if err != nil {
			return nil, 0, err
		}
		elems[i] = elem
		pos += delta
	}
	return elems, pos, nil
}

func DecodeOne(data []byte) (interface{}, int, error) {
	if len(data) == 0 {
		return nil, 0, errors.New("no data")
	}
	switch data[0] {
	case '+':
		return readSimpleString(data)
	case '-':
		//return readError(data)
	case ':':
		return readInt64(data)
	case '$':
		return readBulkString(data)
	case '*':
		return readArray(data)
	}
	return nil, 0, nil
}

func Decode(data []byte) ([]interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("no data")
	}

	intfs := make([]interface{}, 0)
	currIdx := 0

	for currIdx < len(data) {
		value, idx, err := DecodeOne(data[currIdx:])
		if err != nil {
			return intfs, err
		}
		intfs = append(intfs, value)
		currIdx += idx
	}

	return intfs, nil
}

func Encode(value interface{}, isSimpleStr bool) []byte {
	switch t := value.(type) {
	case string:
		if isSimpleStr {
			return []byte(fmt.Sprintf("+%s\r\n", t))
		}
		return encodeString(t) // bulk string

	case int64, int32, int16, int8:
		return []byte(fmt.Sprintf(":%d\r\n", t))
	case []string:
		b := make([]byte, 0)
		buffer := bytes.NewBuffer(b)
		for _, b := range value.([]string) {
			buffer.Write(encodeString(b))
		}
		return []byte(fmt.Sprintf("*%d\r\n%s", len(t), buffer.Bytes()))
	case error:
		return []byte(fmt.Sprintf("-%s\r\n", t))
	default:
		return NilResp
	}

}

func encodeString(t string) []byte {
	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(t), t))
}
