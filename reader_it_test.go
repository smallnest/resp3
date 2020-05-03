package resp3

import (
	"net"
	"strconv"
	"testing"
	"time"
)

func TestReader_IT_Test(t *testing.T) {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:6379", 5*time.Second)
	if err != nil {
		t.Logf("can't found one of redis 6.0 server")
		return
	}
	defer conn.Close()

	w := NewWriter(conn)
	r := NewReader(conn)

	// send an unsupport protocol version
	w.WriteCommand("HELLO", "4")
	helloErrorResp, _, err := r.ReadValue()
	if err != nil {
		t.Fatalf("failed to send a unsupported protocol version")
	}
	if helloErrorResp.Err != "NOPROTO unsupported protocol version" {
		t.Fatalf("expect an error but got %+v", helloErrorResp)
	}

	// send protocol version 3, get a map result
	w.WriteCommand("HELLO", "3")
	helloResp, _, err := r.ReadValue()
	if err != nil {
		t.Fatalf("failed to send a HELLO 3")
	}
	if helloResp.KV.Size() == 0 {
		t.Fatalf("expect some info but got %+v", helloResp)
	}
	t.Logf("%c, %v", helloResp.Type, helloResp.SmartResult())

	// set and get
	w.WriteCommand("SET", "A", "123")
	resp, _, err := r.ReadValue()
	if err != nil {
		t.Fatalf("failed to SET")
	}
	if resp.SmartResult().(string) != "OK" {
		t.Fatalf("expect set success but got %+v", resp)
	}
	w.WriteCommand("GET", "A")
	resp, _, err = r.ReadValue()
	if err != nil {
		t.Fatalf("failed to SET")
	}
	if resp.SmartResult().(string) != "123" {
		t.Fatalf("expect get success but got %+v", resp)
	}

	// incr
	w.WriteCommand("INCR", "B")
	resp, _, err = r.ReadValue()
	if err != nil {
		t.Fatalf("failed to INCR")
	}
	t.Logf("INCR result: %c, %+v", resp.Type, resp.SmartResult())

	// mget
	w.WriteCommand("MGET", "A", "B")
	resp, _, err = r.ReadValue()
	if err != nil {
		t.Fatalf("failed to MGET: %v", err)
	}
	t.Logf("MGET result: %c, %+v", resp.Type, resp.SmartResult())

	// exists
	w.WriteCommand("EXISTS", "C")
	resp, _, err = r.ReadValue()
	if err != nil {
		t.Fatalf("failed to EXISTS: %v", err)
	}
	t.Logf("EXISTS result: %c, %+v", resp.Type, resp.SmartResult())

	// hset
	w.WriteCommand("HSET", "D", "f1", "123")
	resp, _, err = r.ReadValue()
	if err != nil {
		t.Fatalf("failed to HSET: %v", err)
	}
	t.Logf("HSET result: %c, %+v", resp.Type, resp.SmartResult())

	//hgetall
	w.WriteCommand("HGETALL", "D")
	resp, _, err = r.ReadValue()
	if err != nil {
		t.Fatalf("failed to HGETALL: %v", err)
	}
	t.Logf("HGETALL result: %c, %+v", resp.Type, resp.SmartResult())

	//PFADD
	w.WriteCommand("PFADD", "hll", "a", "b", "c", "d", "e", "f", "g")
	resp, _, err = r.ReadValue()
	if err != nil {
		t.Fatalf("failed to PFADD: %v", err)
	}
	t.Logf("PFADD result: %c, %+v", resp.Type, resp.SmartResult())

	//PFCOUNT
	w.WriteCommand("PFCOUNT", "hll")
	resp, _, err = r.ReadValue()
	if err != nil {
		t.Fatalf("failed to PFCOUNT: %v", err)
	}
	t.Logf("PFCOUNT result: %c, %+v", resp.Type, resp.SmartResult())

	//LPUSH
	w.WriteCommand("LPUSH", "mylist", "hello", "world")
	resp, _, err = r.ReadValue()
	if err != nil {
		t.Fatalf("failed to LPUSH: %v", err)
	}
	t.Logf("LPUSH result: %c, %+v", resp.Type, resp.SmartResult())

	//LRANGE
	w.WriteCommand("LRANGE", "mylist", "0", "-1")
	resp, _, err = r.ReadValue()
	if err != nil {
		t.Fatalf("failed to LRANGE: %v", err)
	}
	t.Logf("LRANGE result: %c, %+v", resp.Type, resp.SmartResult())

	//SADD
	w.WriteCommand("SADD", "myset", "hello", "world", "hello")
	resp, _, err = r.ReadValue()
	if err != nil {
		t.Fatalf("failed to SADD: %v", err)
	}
	t.Logf("SADD result: %c, %+v", resp.Type, resp.SmartResult())

	//SMEMBERS
	w.WriteCommand("SMEMBERS", "myset")
	resp, _, err = r.ReadValue()
	if err != nil {
		t.Fatalf("failed to SMEMBERS: %v", err)
	}
	t.Logf("SMEMBERS result: %c, %+v", resp.Type, resp.SmartResult())

	//ZADD
	w.WriteCommand("ZADD", "myzset", "1", "one")
	resp, _, err = r.ReadValue()
	if err != nil {
		t.Fatalf("failed to ZADD: %v", err)
	}
	t.Logf("ZADD result: %c, %+v", resp.Type, resp.SmartResult())

	//ZRANGE
	w.WriteCommand("ZRANGE", "myzset", "0", "-1", "WITHSCORES")
	resp, _, err = r.ReadValue()
	if err != nil {
		t.Fatalf("failed to ZRANGE: %v", err)
	}
	t.Logf("ZRANGE result: %c, %+v", resp.Type, resp.SmartResult())

	//GEOADD
	w.WriteCommand("GEOADD", "Sicily", "13.361389", "38.115556", "Palermo", "15.087269", "37.502669", "Catania")
	resp, _, err = r.ReadValue()
	if err != nil {
		t.Fatalf("failed to GEOADD: %v", err)
	}
	t.Logf("GEOADD result: %c, %+v", resp.Type, resp.SmartResult())

	//GEORADIUS
	w.WriteCommand("GEORADIUS", "Sicily", "15", "37", "100", "km")
	resp, _, err = r.ReadValue()
	if err != nil {
		t.Fatalf("failed to GEORADIUS: %v", err)
	}
	t.Logf("GEORADIUS result: %c, %+v", resp.Type, resp.SmartResult())

	//SUBSCRIBE
	w.WriteCommand("SUBSCRIBE", "news")
	resp, _, err = r.ReadValue()
	if err != nil {
		t.Fatalf("failed to SUBSCRIBE: %v", err)
	}
	t.Logf("SUBSCRIBE result: %c, %+v", resp.Type, resp.SmartResult())

	{
		conn, err := net.DialTimeout("tcp", "127.0.0.1:6379", 5*time.Second)
		if err != nil {
			t.Logf("can't found one of redis 6.0 server")
			return
		}

		w := NewWriter(conn)
		r := NewReader(conn)
		//PUBLISH
		w.WriteCommand("PUBLISH", "news", "resp3 lib is released")
		resp, _, err = r.ReadValue()
		if err != nil {
			t.Fatalf("failed to PUBLISH: %v", err)
		}
		t.Logf("PUBLISH result: %c, %+v", resp.Type, resp.SmartResult())
		conn.Close()
	}
	resp, _, err = r.ReadValue()
	if err != nil {
		t.Fatalf("failed to receive a message: %v", err)
	}
	t.Logf("received SUBSCRIBED result: %c, %+v", resp.Type, resp.SmartResult())
}

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

		//  uint64_t hash = crc64(0,(unsigned char*)sdskey,sdslen(sdskey))&(TRACKING_TABLE_SIZE-1);
		hash := crc64([]byte("a")) & (TRACKING_TABLE_SIZE - 1)
		t.Logf("calculated hash: %d", hash)

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
