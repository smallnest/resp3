# resp3

A redis RESP3 protocol library, written in Go.

`Value` represents a redis command or an redis response, which supports all RESP3 types:

## Types equivalent to RESP version 2

- **Array**: an ordered collection of N other types
- **Blob string**: binary safe strings
- **Simple string**: a space efficient non binary safe string
- **Simple error**: a space efficient non binary safe error code and message
- **Number**: an integer in the signed 64 bit range

## Types introduced by RESP3

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