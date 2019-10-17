package resp3

import (
	"bytes"
	"testing"
)

func TestReader_ReadRaw(t *testing.T) {
	var responses = map[string]string{
		"blobstring":      "$11\r\nhello world\r\n",
		"simple string":   "+hello world\r\n",
		"verbatim string": "=15\r\ntxt:Some string\r\n",
		"err":             "-ERR this is the error description\r\n",
		"number":          ":1234\r\n",
		"double":          ",1.23\r\n",
		"inf":             ",inf\r\n",
		"null":            "_\r\n",
		"bool:true":       "#t\r\n",
		"bool:false":      "#f\r\n",
		"bignumber":       "(3492890328409238509324850943850943825024385\r\n",
		"array":           "*3\r\n:1\r\n:2\r\n:3\r\n",
		"map":             "%2\r\n+first\r\n:1\r\n+second\r\n:2\r\n",
		"set":             "~5\r\n+orange\r\n+apple\r\n#t\r\n:100\r\n:999\r\n",
		"push":            ">4\r\n+pubsub\r\n+message\r\n+somechannel\r\n+this is the message\r\n",
	}

	var responses2 = map[string]string{
		"blobstring":      "$11\r\nhello world\r\nabcd",
		"simple string":   "+hello world\r\nabcd",
		"verbatim string": "=15\r\ntxt:Some string\r\nabcd",
		"err":             "-ERR this is the error description\r\nabcd",
		"number":          ":1234\r\nabcd",
		"double":          ",1.23\r\nabcd",
		"inf":             ",inf\r\nabcd",
		"null":            "_\r\nabcd",
		"bool:true":       "#t\r\nabcd",
		"bool:false":      "#f\r\nabcd",
		"bignumber":       "(3492890328409238509324850943850943825024385\r\nabcd",
		"array":           "*3\r\n:1\r\n:2\r\n:3\r\nabcd",
		"map":             "%2\r\n+first\r\n:1\r\n+second\r\n:2\r\nabcd",
		"set":             "~5\r\n+orange\r\n+apple\r\n#t\r\n:100\r\n:999\r\nabcd",
		"push":            ">4\r\n+pubsub\r\n+message\r\n+somechannel\r\n+this is the message\r\nabcd",
	}

	for k, v := range responses {
		buf := bytes.NewBuffer(nil)
		reader := NewReader(buf)
		buf.WriteString(v)
		data, err := reader.ReadRaw()
		if err != nil {
			t.Errorf("%v failed. err: %v", k, err)
		}
		if v != string(data) {
			t.Errorf("%v failed. expect %s, but got %s", k, v, data)
		}
	}

	for k, v := range responses2 {
		buf := bytes.NewBuffer(nil)
		reader := NewReader(buf)
		buf.WriteString(v)
		data, err := reader.ReadRaw()
		if err != nil {
			t.Errorf("%v failed. err: %v", k, err)
		}
		if v[:len(v)-4] != string(data) {
			t.Errorf("%v failed. expect %s, but got %s", k, v, data)
		}
	}
}

func TestRawReader_Push(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	reader := NewReader(buf)
	buf.WriteString(">4\r\n+pubsub\r\n+message\r\n+somechannel\r\n+this is the message\r\n")
	v, marker, err := reader.ReadValue()
	if err = isError(err, marker); err != nil {
		t.Errorf("failed to read: %v", err)
	}
	if v.Type != TypePush || len(v.Elems) != 4 {
		t.Errorf("not expected, got %v", v.Elems)
	}
	if v.Elems[0].Str != "pubsub" {
		t.Errorf("not expected, got %v", v.Elems[0])
	}
	if v.Elems[1].Str != "message" {
		t.Errorf("not expected, got %v", v.Elems[1])
	}
	if v.Elems[2].Str != "somechannel" {
		t.Errorf("not expected, got %v", v.Elems[2])
	}
	if v.Elems[3].Str != "this is the message" {
		t.Errorf("not expected, got %v", v.Elems[3])
	}
}
