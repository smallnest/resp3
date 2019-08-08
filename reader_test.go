package resp3

import (
	"bytes"
	"math"
	"testing"
)

func TestReader_String(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	reader := NewReader(buf)

	// blobstring
	buf.WriteString("$11\r\nhello world\r\n")
	v, marker, err := reader.ReadValue()
	if err = isError(err, marker); err != nil {
		t.Errorf("failed to read: %v", err)
	}
	if v.Str != "hello world" {
		t.Errorf("not expected, got %s", v.Str)
	}

	//simple string
	buf.Reset()
	buf.WriteString("+hello world\r\n")
	v, marker, err = reader.ReadValue()
	if err = isError(err, marker); err != nil {
		t.Errorf("failed to read: %v", err)
	}
	if v.Str != "hello world" {
		t.Errorf("not expected, got %s", v.Str)
	}

	//verbatim string
	buf.Reset()
	buf.WriteString("=15\r\ntxt:Some string\r\n")
	v, marker, err = reader.ReadValue()
	if err = isError(err, marker); err != nil {
		t.Errorf("failed to read: %v", err)
	}
	if v.Str != "Some string" || v.StrFmt != "txt" {
		t.Errorf("not expected, got %s, %s", v.StrFmt, v.Str)
	}
}

func TestReader_Error(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	reader := NewReader(buf)
	buf.WriteString("-ERR this is the error description\r\n")
	v, marker, err := reader.ReadValue()
	if err = isError(err, marker); err != nil {
		t.Errorf("failed to read: %v", err)
	}
	if v.Err != "ERR this is the error description" {
		t.Errorf("not expected, got %s", v.Err)
	}

	buf.Reset()
	buf.WriteString("!21\r\nSYNTAX invalid syntax\r\n")
	v, marker, err = reader.ReadValue()
	if err = isError(err, marker); err != nil {
		t.Errorf("failed to read: %v", err)
	}
	if v.Err != "SYNTAX invalid syntax" {
		t.Errorf("not expected, got %s", v.Err)
	}
}

func TestReader_Number(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	reader := NewReader(buf)
	buf.WriteString(":1234\r\n")
	v, marker, err := reader.ReadValue()
	if err = isError(err, marker); err != nil {
		t.Errorf("failed to read: %v", err)
	}
	if v.Integer != 1234 {
		t.Errorf("not expected, got %v", v.Integer)
	}
}
func TestReader_Double(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	reader := NewReader(buf)
	buf.WriteString(",1.23\r\n")
	v, marker, err := reader.ReadValue()
	if err = isError(err, marker); err != nil {
		t.Errorf("failed to read: %v", err)
	}
	if v.Double != 1.23 {
		t.Errorf("not expected, got %v", v.Double)
	}
	buf.Reset()
	buf.WriteString(",inf\r\n")
	v, marker, err = reader.ReadValue()
	if err = isError(err, marker); err != nil {
		t.Errorf("failed to read: %v", err)
	}
	if math.IsInf(v.Double, 1) && v.Double > math.MaxFloat64 {
		t.Errorf("not expected, got %v", v.Double)
	}
	buf.Reset()
	buf.WriteString(",-inf\r\n")
	v, marker, err = reader.ReadValue()
	if err = isError(err, marker); err != nil {
		t.Errorf("failed to read: %v", err)
	}
	if math.IsInf(v.Double, -1) && v.Double < math.MaxFloat64 {
		t.Errorf("not expected, got %v", v.Double)
	}
}
func TestReader_Null(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	reader := NewReader(buf)
	buf.WriteString("_\r\n")
	_, marker, err := reader.ReadValue()
	if err = isError(err, marker); err != nil {
		t.Errorf("failed to read: %v", err)
	}
}

func TestReader_Boolean(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	reader := NewReader(buf)
	buf.WriteString("#t\r\n")
	v, marker, err := reader.ReadValue()
	if err = isError(err, marker); err != nil {
		t.Errorf("failed to read: %v", err)
	}
	if !v.Boolean {
		t.Errorf("not expected, got %s", v.Err)
	}
	buf.Reset()
	buf.WriteString("#f\r\n")
	v, marker, err = reader.ReadValue()
	if err = isError(err, marker); err != nil {
		t.Errorf("failed to read: %v", err)
	}
	if v.Boolean {
		t.Errorf("not expected, got %v", v.Boolean)
	}
}
func TestReader_BigInt(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	reader := NewReader(buf)
	buf.WriteString("(3492890328409238509324850943850943825024385\r\n")
	v, marker, err := reader.ReadValue()
	if err = isError(err, marker); err != nil {
		t.Errorf("failed to read: %v", err)
	}
	if v.BigInt.String() != "3492890328409238509324850943850943825024385" {
		t.Errorf("not expected, got %v", v.BigInt.String())
	}
}

func TestReader_Array(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	reader := NewReader(buf)
	buf.WriteString("*3\r\n:1\r\n:2\r\n:3\r\n")
	v, marker, err := reader.ReadValue()
	if err = isError(err, marker); err != nil {
		t.Errorf("failed to read: %v", err)
	}
	if len(v.Elems) != 3 {
		t.Errorf("not expected, got %v", v.Elems)
	}
	for i, vv := range v.Elems {
		if vv.Type != TypeNumber || vv.Integer != int64(i+1) {
			t.Errorf("not expected, got %v", vv.Integer)
		}
	}
}
func TestReader_Map(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	reader := NewReader(buf)
	buf.WriteString("%2\r\n+first\r\n:1\r\n+second\r\n:2\r\n")
	v, marker, err := reader.ReadValue()
	if err = isError(err, marker); err != nil {
		t.Errorf("failed to read: %v", err)
	}
	if v.KV == nil || v.KV.Size() != 2 {
		t.Errorf("not expected, got %v", v.KV)
	}
	for _, key := range v.KV.Keys() {
		k := key.(*Value)
		if k.Str != "first" && k.Str != "second" {
			t.Errorf("not expected, got %v", k)
		}
		vv, ok := v.KV.Get(key)
		if !ok {
			t.Fatalf("not found")
		}
		if vv.(*Value).Integer != 1 && vv.(*Value).Integer != 2 {
			t.Errorf("not expected, got %v", vv)
		}
	}
}
func TestReader_Set(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	reader := NewReader(buf)
	buf.WriteString("~5\r\n+orange\r\n+apple\r\n#t\r\n:100\r\n:999\r\n")
	v, marker, err := reader.ReadValue()
	if err = isError(err, marker); err != nil {
		t.Errorf("failed to read: %v", err)
	}
	if v.Type != TypeSet || len(v.Elems) != 5 {
		t.Errorf("not expected, got %v", v.Elems)
	}

	if v.Elems[0].Str != "orange" {
		t.Errorf("not expected, got %v", v.Elems[0])
	}
	if v.Elems[1].Str != "apple" {
		t.Errorf("not expected, got %v", v.Elems[1])
	}
	if !v.Elems[2].Boolean {
		t.Errorf("not expected, got %v", v.Elems[2])
	}
	if v.Elems[3].Integer != 100 {
		t.Errorf("not expected, got %v", v.Elems[3])
	}
	if v.Elems[4].Integer != 999 {
		t.Errorf("not expected, got %v", v.Elems[4])
	}
}

func TestReader_Push(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	reader := NewReader(buf)
	buf.WriteString(">4\r\n+pubsub\r\n+message\r\n+somechannel\r\n+this is the message\r\nn")
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

func TestReader_Attribute(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	reader := NewReader(buf)
	buf.WriteString("|1\r\n+key-popularity\r\n%2\r\n$1\r\na\r\n,0.1923\r\n$1\r\nb\r\n,0.0012\r\n*2\r\n:2039123\r\n:9543892\r\n")
	v, marker, err := reader.ReadValue()
	if err = isError(err, marker); err != nil {
		t.Errorf("failed to read: %v", err)
	}

	// reply is an array
	if v.Type != TypeArray || len(v.Elems) != 2 {
		t.Errorf("not expected, got %c,%v", v.Type, v.Elems)
	}

	if v.Elems[0].Integer != 2039123 {
		t.Errorf("not expected, got %v", v.Elems[0])
	}
	if v.Elems[1].Integer != 9543892 {
		t.Errorf("not expected, got %v", v.Elems[1])
	}

}
