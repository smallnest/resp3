package resp3

import (
	"net"
	"testing"
	"time"
)

func TestReader_IT_Test(t *testing.T) {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:9999", 5*time.Second)
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
	t.Logf("%v", helloResp.SmartResult())

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
	t.Logf("INCR result: %+v", resp.SmartResult())

	// mget
	w.WriteCommand("MGET", "A", "B")
	resp, _, err = r.ReadValue()
	if err != nil {
		t.Fatalf("failed to MGET: %v", err)
	}
	t.Logf("MGET result: %+v", resp.SmartResult())

	// exists
	w.WriteCommand("EXISTS", "C")
	resp, _, err = r.ReadValue()
	if err != nil {
		t.Fatalf("failed to EXISTS: %v", err)
	}
	t.Logf("EXISTS result: %+v", resp.SmartResult())

	// hset
	w.WriteCommand("HSET", "D", "f1", "123")
	resp, _, err = r.ReadValue()
	if err != nil {
		t.Fatalf("failed to HSET: %v", err)
	}
	t.Logf("HSET result: %+v", resp.SmartResult())

	//hgetall
	w.WriteCommand("HGETALL", "D")
	resp, _, err = r.ReadValue()
	if err != nil {
		t.Fatalf("failed to HGETALL: %v", err)
	}
	t.Logf("HGETALL result: %+v", resp.SmartResult())

	//PFADD
	w.WriteCommand("PFADD", "hll", "a", "b", "c", "d", "e", "f", "g")
	resp, _, err = r.ReadValue()
	if err != nil {
		t.Fatalf("failed to PFADD: %v", err)
	}
	t.Logf("PFADD result: %+v", resp.SmartResult())

	//PFCOUNT
	w.WriteCommand("PFCOUNT", "hll")
	resp, _, err = r.ReadValue()
	if err != nil {
		t.Fatalf("failed to PFCOUNT: %v", err)
	}
	t.Logf("PFCOUNT result: %+v", resp.SmartResult())

	//LPUSH
	w.WriteCommand("LPUSH", "mylist", "hello", "world")
	resp, _, err = r.ReadValue()
	if err != nil {
		t.Fatalf("failed to LPUSH: %v", err)
	}
	t.Logf("LPUSH result: %+v", resp.SmartResult())

	//LRANGE
	w.WriteCommand("LRANGE", "mylist", "0", "-1")
	resp, _, err = r.ReadValue()
	if err != nil {
		t.Fatalf("failed to LRANGE: %v", err)
	}
	t.Logf("LRANGE result: %+v", resp.SmartResult())

	//SADD
	w.WriteCommand("SADD", "myset", "hello", "world", "hello")
	resp, _, err = r.ReadValue()
	if err != nil {
		t.Fatalf("failed to SADD: %v", err)
	}
	t.Logf("SADD result: %+v", resp.SmartResult())

	//SMEMBERS
	w.WriteCommand("SMEMBERS", "myset")
	resp, _, err = r.ReadValue()
	if err != nil {
		t.Fatalf("failed to SMEMBERS: %v", err)
	}
	t.Logf("SMEMBERS result: %+v", resp.SmartResult())

	//ZADD
	w.WriteCommand("ZADD", "myzset", "1", "one")
	resp, _, err = r.ReadValue()
	if err != nil {
		t.Fatalf("failed to ZADD: %v", err)
	}
	t.Logf("ZADD result: %+v", resp.SmartResult())

	//ZRANGE
	w.WriteCommand("ZRANGE", "myzset", "0", "-1", "WITHSCORES")
	resp, _, err = r.ReadValue()
	if err != nil {
		t.Fatalf("failed to ZRANGE: %v", err)
	}
	t.Logf("ZRANGE result: %+v", resp.SmartResult())

	//GEOADD
	w.WriteCommand("GEOADD", "Sicily", "13.361389", "38.115556", "Palermo", "15.087269", "37.502669", "Catania")
	resp, _, err = r.ReadValue()
	if err != nil {
		t.Fatalf("failed to GEOADD: %v", err)
	}
	t.Logf("GEOADD result: %+v", resp.SmartResult())

	//GEORADIUS
	w.WriteCommand("GEORADIUS", "Sicily", "15", "37", "100", "km")
	resp, _, err = r.ReadValue()
	if err != nil {
		t.Fatalf("failed to GEORADIUS: %v", err)
	}
	t.Logf("GEORADIUS result: %+v", resp.SmartResult())

	// //SUBSCRIBE
	// w.WriteCommand("SUBSCRIBE", "news")
	// resp, _, err = r.ReadValue()
	// if err != nil {
	// 	t.Fatalf("failed to SUBSCRIBE: %v", err)
	// }
	// t.Logf("SUBSCRIBE result: %+v", resp.SmartResult())

	// //PUBLISH
	// w.WriteCommand("PUBLISH", "news", "resp3 lib is released")
	// resp, _, err = r.ReadValue()
	// if err != nil {
	// 	t.Fatalf("failed to PUBLISH: %v", err)
	// }
	// t.Logf("PUBLISH result: %+v", resp.SmartResult())
	// resp, _, err = r.ReadValue()
	// if err != nil {
	// 	t.Fatalf("failed to receive a message: %v", err)
	// }
	// t.Logf("received SUBSCRIBED result: %+v", resp.SmartResult())
}
