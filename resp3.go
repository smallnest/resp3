package resp3

import (
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/emirpasic/gods/maps/linkedhashmap"
)

const CRLF = "\r\n"

var CRLFByte = []byte(CRLF)
var StreamMarkerPrefix = []byte("$EOF:")

type floatInfo struct {
	mantbits uint
	expbits  uint
	bias     int
}

var float64info = floatInfo{52, 11, -1023}

// https://github.com/antirez/RESP3/blob/master/spec.md

// resp3 type char
const (
	// simple types

	TypeBlobString     = '$' // $<length>\r\n<bytes>\r\n
	TypeSimpleString   = '+' // +<string>\r\n
	TypeSimpleError    = '-' // -<string>\r\n
	TypeNumber         = ':' // :<number>\r\n
	TypeNull           = '_' // _\r\n
	TypeDouble         = ',' // ,<floating-point-number>\r\n
	TypeBoolean        = '#' // #t\r\n or #f\r\n
	TypeBlobError      = '!' // !<length>\r\n<bytes>\r\n
	TypeVerbatimString = '=' // =<length>\r\n<format(3 bytes):><bytes>\r\n
	TypeBigNumber      = '(' // (<big number>\n

	// Aggregate data types

	TypeArray     = '*' // *<elements number>\r\n... numelements other types ...
	TypeMap       = '%' // %<elements number>\r\n... numelements key/value pair of other types ...
	TypeSet       = '~' // ~<elements number>\r\n... numelements other types ...
	TypeAttribute = '|' // |~<elements number>\r\n... numelements map type ...
	TypePush      = '>' // ><elements number>\r\n<first item is String>\r\n... numelements-1 other types ...

	//special type
	TypeStream = "$EOF:" // $EOF:<40 bytes marker><CR><LF>... any number of bytes of data here not containing the marker ...<40 bytes marker>
)

// Value is a common struct for all RESP3 type.
// There is no exact field for NULL type because the type field is enough.
// It is not used for Stream type.
type Value struct {
	Type         byte
	Str          string
	StrFmt       string
	Err          string
	Integer      int64
	Boolean      bool
	Double       float64
	BigInt       *big.Int
	Elems        []*Value           // for array & set
	KV           *linkedhashmap.Map //TODO sorted map, for map & attr
	Attrs        *linkedhashmap.Map
	StreamMarker string
}

// SmartResult converts itself to a real object.
// Attributes are dropped.
// simple objects are converted their Go types.
// String -> go string
// Interger -> go int64
// Double -> go float64
// Boolean -> go bool
// Err -> go string
// BigInt -> big.Int
// Array -> go array
// Map --> github.com/emirpasic/gods/maps/linkedhashmap.Map
// Set -> go array
// Push -> go array
// NULL -> nil
func (r *Value) SmartResult() interface{} {
	switch r.Type {
	case TypeSimpleString:
		return r.Str
	case TypeBlobString:
		return r.Str
	case TypeVerbatimString:
		return r.Str
	case TypeSimpleError:
		return r.Err
	case TypeBlobError:
		return r.Err
	case TypeNumber:
		return r.Integer
	case TypeDouble:
		return r.Double
	case TypeBigNumber:
		return r.BigInt
	case TypeNull:
		return nil
	case TypeBoolean:
		return r.Boolean
	case TypeArray, TypeSet, TypePush:
		var rt []interface{}
		for _, elem := range r.Elems {
			rt = append(rt, elem.SmartResult())
		}
		return rt
	case TypeMap:
		var rt = linkedhashmap.New()
		if r.KV != nil {
			r.KV.Each(func(k, v interface{}) {
				rt.Put(k.(*Value).SmartResult(), v.(*Value).SmartResult())
			})
		}
		return rt
	}

	return nil
}

// ToRESP3String converts this value to redis RESP3 string.
func (r *Value) ToRESP3String() string {
	buf := new(strings.Builder)

	//check attributes
	if r.Attrs != nil && r.Attrs.Size() > 0 {
		buf.WriteByte(TypeAttribute)
		buf.WriteString(strconv.Itoa(r.Attrs.Size()))
		buf.Write(CRLFByte)

		r.Attrs.Each(func(key, val interface{}) {
			k := key.(*Value)
			v := val.(*Value)
			buf.WriteByte(k.Type)
			k.toRESP3String(buf)
			buf.WriteByte(v.Type)
			v.toRESP3String(buf)
		})
	}

	buf.WriteByte(r.Type)
	r.toRESP3String(buf)
	return buf.String()
}

func (r *Value) toRESP3String(buf *strings.Builder) {
	switch r.Type {
	case TypeSimpleString:
		buf.WriteString(r.Str)
	case TypeBlobString:
		buf.WriteString(strconv.Itoa(len(r.Str)))
		buf.Write(CRLFByte)
		buf.WriteString(r.Str)
	case TypeVerbatimString:
		buf.WriteString(strconv.Itoa(len(r.Str) + 4))
		buf.Write(CRLFByte)
		buf.WriteString(r.StrFmt)
		buf.WriteByte(':')
		buf.WriteString(r.Str)
	case TypeSimpleError:
		buf.WriteString(r.Err)
	case TypeBlobError:
		buf.WriteString(strconv.Itoa(len(r.Err)))
		buf.Write(CRLFByte)
		buf.WriteString(r.Err)
	case TypeNumber:
		buf.WriteString(strconv.FormatInt(r.Integer, 10))
	case TypeDouble:
		bits := math.Float64bits(r.Double)
		flt := &float64info
		neg := bits>>(flt.expbits+flt.mantbits) != 0
		exp := int(bits>>flt.mantbits) & (1<<flt.expbits - 1)

		switch exp {
		case 1<<flt.expbits - 1:
			if neg {
				buf.WriteString("-inf")
			} else {
				buf.WriteString("inf")
			}
		default:
			buf.WriteString(strconv.FormatFloat(r.Double, 'f', -1, 64))
		}

	case TypeBigNumber:
		buf.WriteString(r.BigInt.String())
	case TypeNull:

	case TypeBoolean:
		if r.Boolean {
			buf.WriteByte('t')
		} else {
			buf.WriteByte('f')
		}
	case TypeArray, TypeSet, TypePush:
		buf.WriteString(strconv.Itoa(len(r.Elems)))
		buf.Write(CRLFByte)

		for _, v := range r.Elems {
			buf.WriteByte(v.Type)
			v.toRESP3String(buf)
		}
		return
	case TypeMap:
		buf.WriteString(strconv.Itoa(r.KV.Size()))
		buf.Write(CRLFByte)
		if r.KV != nil {
			r.KV.Each(func(key, val interface{}) {
				k := key.(*Value)
				v := val.(*Value)
				buf.WriteByte(k.Type)
				k.toRESP3String(buf)
				buf.WriteByte(v.Type)
				v.toRESP3String(buf)
			})
		}
		return
	}

	buf.Write(CRLFByte)
}
