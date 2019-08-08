package resp3

import (
	"math"
	"math/big"
	"testing"

	"github.com/emirpasic/gods/maps/linkedhashmap"
)

func TestBlobString(t *testing.T) {
	v := &Value{
		Type: TypeBlobString,
		Str:  "hello world",
	}

	s := v.ToRESP3String()
	expected := "$11\r\nhello world\r\n"
	if s != expected {
		t.Errorf("expected %s but got %s", expected, s)
	}

	v.Str = ""
	s = v.ToRESP3String()
	expected = "$0\r\n\r\n"
	if s != expected {
		t.Errorf("expected %s but got %s", expected, s)
	}
}

func TestSimpleString(t *testing.T) {
	v := &Value{
		Type: TypeSimpleString,
		Str:  "hello world",
	}

	s := v.ToRESP3String()
	expected := "+hello world\r\n"
	if s != expected {
		t.Errorf("expected %s but got %s", expected, s)
	}
}

func TestSimpleError(t *testing.T) {
	v := &Value{
		Type: TypeSimpleError,
		Err:  "ERR this is the error description",
	}

	s := v.ToRESP3String()
	expected := "-ERR this is the error description\r\n"
	if s != expected {
		t.Errorf("expected %s but got %s", expected, s)
	}
}

func TestBlobError(t *testing.T) {
	v := &Value{
		Type: TypeBlobError,
		Err:  "SYNTAX invalid syntax",
	}

	s := v.ToRESP3String()
	expected := "!21\r\nSYNTAX invalid syntax\r\n"
	if s != expected {
		t.Errorf("expected %s but got %s", expected, s)
	}
}

func TestNumber(t *testing.T) {
	v := &Value{
		Type:    TypeNumber,
		Integer: 1234,
	}

	s := v.ToRESP3String()
	expected := ":1234\r\n"
	if s != expected {
		t.Errorf("expected %s but got %s", expected, s)
	}
}

func TestNull(t *testing.T) {
	v := &Value{
		Type: TypeNull,
	}

	s := v.ToRESP3String()
	expected := "_\r\n"
	if s != expected {
		t.Errorf("expected %s but got %s", expected, s)
	}
}

func TestDouble(t *testing.T) {
	v := &Value{
		Type:   TypeDouble,
		Double: 1.23,
	}

	s := v.ToRESP3String()
	expected := ",1.23\r\n"
	if s != expected {
		t.Errorf("expected %s but got %s", expected, s)
	}

	v = &Value{
		Type:   TypeDouble,
		Double: math.NaN(),
	}

	s = v.ToRESP3String()
	expected = ",inf\r\n"
	if s != expected {
		t.Errorf("expected %s but got %s", expected, s)
	}

	v = &Value{
		Type:   TypeDouble,
		Double: -math.NaN(),
	}

	s = v.ToRESP3String()
	expected = ",-inf\r\n"
	if s != expected {
		t.Errorf("expected %s but got %s", expected, s)
	}
}

func TestBoolean(t *testing.T) {
	v := &Value{
		Type:    TypeBoolean,
		Boolean: true,
	}

	s := v.ToRESP3String()
	expected := "#t\r\n"
	if s != expected {
		t.Errorf("expected %s but got %s", expected, s)
	}

	v = &Value{
		Type:    TypeBoolean,
		Boolean: false,
	}

	s = v.ToRESP3String()
	expected = "#f\r\n"
	if s != expected {
		t.Errorf("expected %s but got %s", expected, s)
	}
}

func TestBolbError(t *testing.T) {
	v := &Value{
		Type: TypeBlobError,
		Err:  "SYNTAX invalid syntax",
	}

	s := v.ToRESP3String()
	expected := "!21\r\nSYNTAX invalid syntax\r\n"
	if s != expected {
		t.Errorf("expected %s but got %s", expected, s)
	}
}

func TestVerbatimString(t *testing.T) {
	v := &Value{
		Type:   TypeVerbatimString,
		Str:    "Some string",
		StrFmt: "txt",
	}

	s := v.ToRESP3String()
	expected := "=15\r\ntxt:Some string\r\n"
	if s != expected {
		t.Errorf("expected %s but got %s", expected, s)
	}
}

func TestBigNumber(t *testing.T) {
	bigInt, ok := new(big.Int).SetString("3492890328409238509324850943850943825024385", 10)
	if !ok {
		t.Fatalf("failed to parse big.Int")
	}

	v := &Value{
		Type:   TypeBigNumber,
		BigInt: bigInt,
	}

	s := v.ToRESP3String()
	expected := "(3492890328409238509324850943850943825024385\r\n"
	if s != expected {
		t.Errorf("expected %s but got %s", expected, s)
	}
}

func TestArray(t *testing.T) {
	v := &Value{
		Type: TypeArray,
		Elems: []*Value{
			&Value{
				Type:    TypeNumber,
				Integer: 1,
			},
			&Value{
				Type:    TypeNumber,
				Integer: 2,
			},
			&Value{
				Type:    TypeNumber,
				Integer: 3,
			},
		},
	}

	s := v.ToRESP3String()
	expected := "*3\r\n:1\r\n:2\r\n:3\r\n"
	if s != expected {
		t.Errorf("expected %s but got %s", expected, s)
	}
}

func TestMap(t *testing.T) {

	m := linkedhashmap.New()
	m.Put(&Value{
		Type: TypeSimpleString,
		Str:  "first",
	}, &Value{
		Type:    TypeNumber,
		Integer: 1,
	})
	m.Put(&Value{
		Type: TypeSimpleString,
		Str:  "second",
	}, &Value{
		Type:    TypeNumber,
		Integer: 2,
	})

	v := &Value{
		Type: TypeMap,
		KV:   m,
	}

	s := v.ToRESP3String()
	expected := "%2\r\n+first\r\n:1\r\n+second\r\n:2\r\n"
	if s != expected {
		t.Errorf("expected %s but got %s", expected, s)
	}
}

func TestSet(t *testing.T) {
	v := &Value{
		Type: TypeSet,
		Elems: []*Value{
			&Value{
				Type: TypeSimpleString,
				Str:  "orange",
			},
			&Value{
				Type: TypeSimpleString,
				Str:  "apple",
			},
			&Value{
				Type:    TypeBoolean,
				Boolean: true,
			},
			&Value{
				Type:    TypeNumber,
				Integer: 100,
			},
			&Value{
				Type:    TypeNumber,
				Integer: 999,
			},
		},
	}

	s := v.ToRESP3String()
	expected := "~5\r\n+orange\r\n+apple\r\n#t\r\n:100\r\n:999\r\n"
	if s != expected {
		t.Errorf("expected %s but got %s", expected, s)
	}
}

func TestAttribute(t *testing.T) {
	keyPopularityMap := linkedhashmap.New()
	keyPopularityMap.Put(&Value{
		Type: TypeBlobString,
		Str:  "a",
	}, &Value{
		Type:   TypeDouble,
		Double: 0.1923,
	})
	keyPopularityMap.Put(&Value{
		Type: TypeBlobString,
		Str:  "b",
	}, &Value{
		Type:   TypeDouble,
		Double: 0.0012,
	})

	attrMap := linkedhashmap.New()
	attrMap.Put(&Value{
		Type: TypeSimpleString,
		Str:  "key-popularity",
	}, &Value{
		Type: TypeMap,
		KV:   keyPopularityMap,
	})

	v := &Value{
		Type: TypeArray,
		Elems: []*Value{
			&Value{
				Type:    TypeNumber,
				Integer: 2039123,
			},
			&Value{
				Type:    TypeNumber,
				Integer: 9543892,
			},
		},
		Attrs: attrMap,
	}

	s := v.ToRESP3String()
	expected := "|1\r\n+key-popularity\r\n%2\r\n$1\r\na\r\n,0.1923\r\n$1\r\nb\r\n,0.0012\r\n*2\r\n:2039123\r\n:9543892\r\n"
	if s != expected {
		t.Errorf("expected %s but got %s", expected, s)
	}
}

func TestPush(t *testing.T) {
	v := &Value{
		Type: TypePush,
		Elems: []*Value{
			&Value{
				Type: TypeSimpleString,
				Str:  "pubsub",
			},
			&Value{
				Type: TypeSimpleString,
				Str:  "message",
			},
			&Value{
				Type: TypeSimpleString,
				Str:  "somechannel",
			},
			&Value{
				Type: TypeSimpleString,
				Str:  "this is the message",
			},
		},
	}

	s := v.ToRESP3String()
	expected := ">4\r\n+pubsub\r\n+message\r\n+somechannel\r\n+this is the message\r\n"
	if s != expected {
		t.Errorf("expected %s but got %s", expected, s)
	}
}

func TestStream(t *testing.T) {

}
