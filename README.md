# resp3

[![License](https://img.shields.io/:license-apache%202-blue.svg)](https://opensource.org/licenses/Apache-2.0) [![GoDoc](https://godoc.org/github.com/smallnest/resp3?status.png)](http://godoc.org/github.com/smallnest/resp3)  [![travis](https://travis-ci.org/smallnest/resp3.svg?branch=master)](https://travis-ci.org/smallnest/resp3) [![Go Report Card](https://goreportcard.com/badge/github.com/smallnest/resp3)](https://goreportcard.com/report/github.com/smallnest/resp3) [![coveralls](https://coveralls.io/repos/smallnest/resp3/badge.svg?branch=master&service=github)](https://coveralls.io/github/smallnest/resp3?branch=master) 

A redis RESP3 protocol library, written in Go.

`Value` represents a redis command or an redis response, which supports all RESP3 types:

### Types equivalent to RESP version 2

- **Array**: an ordered collection of N other types
- **Blob string**: binary safe strings
- **Simple string**: a space efficient non binary safe string
- **Simple error**: a space efficient non binary safe error code and message
- **Number**: an integer in the signed 64 bit range

### Types introduced by RESP3

- **Null**: a single null value replacing RESP v2 *-1 and $-1 null values.
- **Double**: a floating point number
- **Boolean**: true or false
- **Blob error**: binary safe error code and message.
- **Verbatim string**: a binary safe string that should be displayed to humans without any escaping or filtering. For instance the output of LATENCY DOCTOR in Redis.
- **Map**: an ordered collection of key-value pairs. Keys and values can be any other RESP3 type.
- **Set**: an unordered collection of N other types.
- **Attribute**: Like the Map type, but the client should keep reading the reply ignoring the attribute type, and return it to the client as additional information.
- **Push**: Out of band data. The format is like the Array type, but the client should just check the first string element, stating the type of the out of band data, a call a callback if there is one registered for this specific type of push information. Push types are not related to replies, since they are information that the server may push at any time in the connection, so the client should keep reading if it is reading the reply of a command.
- **Hello**: Like the Map type, but is sent only when the connection between the client and the server is established, in order to welcome the client with different information like the name of the server, its version, and so forth.
- **Big number**: a large number non representable by the Number type

`Reader` is used to read a `Value` object from the connection. It can be used by both redis clients and redis servers.

## Examples

### Client Cache (Tracking)

you can use it to test the client cache feature in Redis 6.0.

```go
func TestReader_IT_Tracking(t *testing.T) {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:6379", 5*time.Second)
	if err != nil {
		t.Logf("can't found one of redis 6.0 server")
		return
	}
	defer conn.Close()

	w := NewWriter(conn)
	r := NewReader(conn)

	w.WriteCommand("HELLO", "3")
	helloResp, _, err := r.ReadValue()
	if err != nil {
		t.Fatalf("failed to send a HELLO 3")
	}
	if helloResp.KV.Size() == 0 {
		t.Fatalf("expect some info but got %+v", helloResp)
	}
	t.Logf("hello response: %c, %v", helloResp.Type, helloResp.SmartResult())

	w.WriteCommand("CLIENT", "TRACKING", "on")
	resp, _, err := r.ReadValue()
	if err != nil {
		t.Fatalf("failed to TRACKING: %v", err)
	}
	t.Logf("TRACKING result: %c, %+v", resp.Type, resp.SmartResult())

	w.WriteCommand("GET", "a")
	resp, _, err = r.ReadValue()
	if err != nil {
		t.Fatalf("failed to GET: %v", err)
	}
	t.Logf("GET result: %c, %+v", resp.Type, resp.SmartResult())

	go func() {
		conn, err := net.DialTimeout("tcp", "127.0.0.1:9999", 5*time.Second)
		if err != nil {
			t.Logf("can't found one of redis 6.0 server")
			return
		}
		defer conn.Close()
		w := NewWriter(conn)
		r := NewReader(conn)

		for i := 0; i < 10; i++ {
			//PUBLISH
			w.WriteCommand("set", "a", strconv.Itoa(i))
			resp, _, err = r.ReadValue()
			if err != nil {
				t.Fatalf("failed to set: %v", err)
			}
			t.Logf("set result: %c, %+v", resp.Type, resp.SmartResult())
			time.Sleep(200 * time.Millisecond)
		}

	}()

	for i := 0; i < 10; i++ {
		resp, _, err = r.ReadValue()
		if err != nil {
			t.Fatalf("failed to receive a message: %v", err)
		}
		if resp.Type == TypePush && len(resp.Elems) >= 2 && resp.Elems[0].SmartResult().(string) == "invalidate" {
			t.Logf("received TRACKING result: %c, %+v", resp.Type, resp.SmartResult())

			// refresh cache "a"
			w.WriteCommand("GET", "a")
			resp, _, err = r.ReadValue()
		}
	}
}
```