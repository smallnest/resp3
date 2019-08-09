package resp3

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/emirpasic/gods/maps/linkedhashmap"
)

// Errors
var (
	ErrInvalidSyntax = errors.New("resp: invalid syntax")
)

// Reader is reader to parse responses/requests from the underlying reader.
type Reader struct {
	*bufio.Reader
}

// NewReader returns a RESP3 reader.
func NewReader(reader io.Reader) *Reader {
	return NewReaderSize(reader, 32*1024)
}

// NewReaderSize returns a new Reader whose buffer has at least the specified size.
func NewReaderSize(reader io.Reader, size int) *Reader {
	return &Reader{
		Reader: bufio.NewReaderSize(reader, size),
	}
}

// ReadValue parses a RESP3 value.
func (r *Reader) ReadValue() (*Value, []byte, error) {
	line, err := r.readLine()
	if err != nil {
		return nil, nil, err
	}
	if len(line) < 3 {
		return nil, nil, ErrInvalidSyntax
	}

	var attrs *linkedhashmap.Map
	if line[0] == TypeAttribute {
		attrs, err = r.readAttr(line)
		if err != nil {
			return nil, nil, err
		}
		line, err = r.readLine()
	}

	// check stream. if it is stream, return the stream marker
	if line[0] == TypeBlobString && len(line) == 45 && bytes.Compare(line[:5], StreamMarkerPrefix) == 0 {
		return nil, line[5:], nil
	}

	v := &Value{
		Type:  line[0],
		Attrs: attrs,
	}

	switch v.Type {
	case TypeSimpleString:
		v.Str = string(line[1 : len(line)-2])
	case TypeBlobString:
		v.Str, err = r.readBlobString(line)
	case TypeVerbatimString:
		var s string
		s, err = r.readBlobString(line)
		if err == nil {
			if len(s) < 4 {
				err = ErrInvalidSyntax
			} else {
				v.Str = s[4:]
				v.StrFmt = s[:3]
			}
		}
	case TypeSimpleError:
		v.Err = string(line[1 : len(line)-2])
	case TypeBlobError:
		v.Err, err = r.readBlobString(line)
	case TypeNumber:
		v.Integer, err = r.readNumber(line)
	case TypeDouble:
		v.Double, err = r.readDouble(line)
	case TypeBigNumber:
		v.BigInt, err = r.readBigNumber(line)
	case TypeNull:
		if len(line) != 3 {
			err = ErrInvalidSyntax
		}
	case TypeBoolean:
		v.Boolean, err = r.readBoolean(line)
	case TypeArray, TypeSet, TypePush:
		v.Elems, err = r.readArray(line)
	case TypeMap:
		v.KV, err = r.readMap(line)
	}

	return v, nil, err
}

func (r *Reader) readLine() (line []byte, err error) {
	line, err = r.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	if len(line) > 1 && line[len(line)-2] == '\r' {
		return line, nil
	}
	return nil, ErrInvalidSyntax
}

func (r *Reader) getCount(line []byte) (int, error) {
	end := bytes.IndexByte(line, '\r')
	return strconv.Atoi(string(line[1:end]))
}

func (r *Reader) readBlobString(line []byte) (string, error) {
	count, err := r.getCount(line)
	if err != nil {
		return "", err
	}
	if count < 0 {
		return "", ErrInvalidSyntax
	}

	buf := make([]byte, count+2)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return "", err
	}
	return string(buf[:count]), nil
}

func (r *Reader) readNumber(line []byte) (int64, error) {
	v := string(line[1 : len(line)-2])
	return strconv.ParseInt(v, 10, 64)
}

func (r *Reader) readDouble(line []byte) (float64, error) {
	v := string(line[1 : len(line)-2])
	if v == "inf" {
		return math.NaN(), nil
	} else if v == "-inf" {
		return -math.NaN(), nil
	}
	return strconv.ParseFloat(v, 64)
}

func (r *Reader) readBigNumber(line []byte) (*big.Int, error) {
	v := string(line[1 : len(line)-2])
	i := new(big.Int)
	if i, ok := i.SetString(v, 10); ok {
		return i, nil
	}
	return nil, ErrInvalidSyntax
}

func (r *Reader) readBoolean(line []byte) (bool, error) {
	v := string(line[1 : len(line)-2])
	if v == "t" {
		return true, nil
	} else if v == "f" {
		return false, nil
	}

	return false, ErrInvalidSyntax
}

func (r *Reader) readArray(line []byte) ([]*Value, error) {
	count, err := r.getCount(line)
	if err != nil {
		return nil, err
	}

	var rt []*Value
	for i := 0; i < count; i++ {
		v, streamMarkerPrefix, err := r.ReadValue()
		if err = isError(err, streamMarkerPrefix); err != nil {
			return nil, err
		}
		rt = append(rt, v)
	}
	return rt, nil
}

func isError(err error, streamMarkerPrefix []byte) error {
	if err != nil {
		return err
	}
	if len(streamMarkerPrefix) > 0 {
		return ErrInvalidSyntax
	}

	return nil
}

func (r *Reader) readMap(line []byte) (*linkedhashmap.Map, error) {
	count, err := r.getCount(line)
	if err != nil {
		return nil, err
	}

	rt := linkedhashmap.New()
	for i := 0; i < count; i++ {
		k, streamMarkerPrefix, err := r.ReadValue()
		if err = isError(err, streamMarkerPrefix); err != nil {
			return nil, err
		}
		v, streamMarkerPrefix, err := r.ReadValue()
		if err = isError(err, streamMarkerPrefix); err != nil {
			return nil, err
		}
		rt.Put(k, v)
	}
	return rt, nil
}

func (r *Reader) readAttr(line []byte) (*linkedhashmap.Map, error) {
	count, err := r.getCount(line)
	if err != nil {
		return nil, err
	}

	rt := linkedhashmap.New()
	for i := 0; i < count; i++ {
		k, streamMarkerPrefix, err := r.ReadValue()
		if err = isError(err, streamMarkerPrefix); err != nil {
			return nil, err
		}
		v, streamMarkerPrefix, err := r.ReadValue()
		if err = isError(err, streamMarkerPrefix); err != nil {
			return nil, err
		}
		rt.Put(k, v)
	}
	return rt, nil
}

// FromString convert a string into a Value.
func FromString(data string) (*Value, error) {
	r := NewReader(strings.NewReader(data))
	v, _, err := r.ReadValue()
	return v, err
}
