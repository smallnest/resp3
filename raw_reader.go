package resp3

import (
	"bytes"
	"io"
)

// ReadRaw parses a RESP3 raw byte slice.
func (r *Reader) ReadRaw() ([]byte, error) {
	var buf bytes.Buffer
	err := r.readRaw(&buf)
	return buf.Bytes(), err
}

// ReadRaw parses a RESP3 raw byte slice.
func (r *Reader) readRaw(w io.Writer) error {
	line, err := r.readLine()
	if err != nil {
		return err
	}
	if len(line) < 3 {
		return ErrInvalidSyntax
	}

	var data []byte

	if line[0] == TypeAttribute {
		w.Write(line)
		err = r.readRawAttr(w, line)
		if err != nil {
			return err
		}
		w.Write(data)
		line, err = r.readLine()
	}

	// check stream. if it is stream, return the stream marker
	if line[0] == TypeBlobString && len(line) == 45 && bytes.Compare(line[:5], StreamMarkerPrefix) == 0 {
		return ErrStreamingUnsupport
	}

	w.Write(line)

	switch line[0] {
	case TypeSimpleString, TypeSimpleError, TypeBlobError:
		return nil
	case TypeNumber, TypeDouble, TypeBigNumber:
		return nil
	case TypeNull, TypeBoolean:
		return nil
	case TypeBlobString:
		err = r.readRawBlobString(w, line)
	case TypeVerbatimString:
		err = r.readRawBlobString(w, line)
	case TypeArray, TypeSet, TypePush:
		err = r.readRawArray(w, line)
	case TypeMap:
		err = r.readRawMap(w, line)
	}

	return err
}

func (r *Reader) readRawAttr(w io.Writer, line []byte) error {
	count, err := r.getCount(line)
	if err != nil {
		return err
	}

	for i := 0; i < count; i++ {
		err = r.readRaw(w)
		if err != nil {
			return err
		}
		err := r.readRaw(w)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Reader) readRawBlobString(w io.Writer, line []byte) error {
	count, err := r.getCount(line)
	if err != nil {
		return err
	}
	if count < 0 {
		return ErrInvalidSyntax
	}

	buf := make([]byte, count+2)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return err
	}
	w.Write(buf)
	return nil
}

func (r *Reader) readRawArray(w io.Writer, line []byte) error {
	count, err := r.getCount(line)
	if err != nil {
		return err
	}

	for i := 0; i < count; i++ {
		err = r.readRaw(w)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Reader) readRawMap(w io.Writer, line []byte) error {
	count, err := r.getCount(line)
	if err != nil {
		return err
	}

	for i := 0; i < count; i++ {
		err = r.readRaw(w)
		if err != nil {
			return err
		}
		err = r.readRaw(w)
		if err != nil {
			return err
		}
	}
	return nil
}
