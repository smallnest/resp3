// Package resp3 implements redis RESP3 protocol, which is used from redis 6.0.
//
// RESP (REdis Serialization Protocol) is the protocol used in the Redis database, however the protocol is designed to be used by other projects. With the version 3 of the protocol, currently a work in progress design, the protocol aims to become even more generally useful to other systems that want to implement a protocol which is simple, efficient, and with a very large landscape of client libraries implementations.
// That means you can use this library to access other RESP3 projects.
//
// This library contains three important components: Value, Reader and Writer.
//
// Value represents a redis command or a redis response. It is a common struct for all RESP3 types.
//
// Reader can parse redis responses from redis servers or commands from redis clients. You can use it to implement redis 6.0 clients,
// no need to pay attention to underlying parsing. Those new features of redis 6.0 can be implemented based on it.
//
// Writer is redis writer. You can use it to send commands to redis servers.
//
// RESP3 spec can be found at https://github.com/antirez/RESP3.
//
// A redis client based on it is just as the below:
//
//	conn, err := net.DialTimeout("tcp", "127.0.0.1:6379", 5*time.Second)
//	if err != nil {
//		t.Logf("can't found one of redis 6.0 server")
//		return
//	}
//	defer conn.Close()
//
//	w := NewWriter(conn)
//	r := NewReader(conn)
//
//	// send protocol version 3, get a map result
//	w.WriteCommand("HELLO", "3")
//	resp, _, _ := r.ReadValue()
//	log.Printf("%v", resp.SmartResult())
//
//	// set
//	w.WriteCommand("SET", "A", "123")
//	resp, _, _ := r.ReadValue()
//	log.Printf("%v", resp.SmartResult())
//
//	// get
//	w.WriteCommand("GET", "A")
//	resp, _, _ = r.ReadValue()
//	log.Printf("%v", resp.SmartResult())
package resp3
