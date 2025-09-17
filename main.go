package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("Listening on port:6379...")
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer l.Close()
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting: ", err.Error())
		return
	}
	defer conn.Close()
	fmt.Println("Connection established")

	// 处理连接
	for {
		resp := NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			fmt.Println("Error reading from client:", err.Error())
			return
		}
		fmt.Println("Input:", value)
		_ = value
		// 回复客户端
		writer := NewWriter(conn)
		writer.Write(Value{typ: "string", str: "OK"})
		// conn.Write([]byte("+OK\r\n"))
	}
}
