package main

import (
	"fmt"
	"net"
	"strings"
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

		if value.typ != "array" {
			fmt.Println("Expected array type, got:", value.typ)
			return
		}
		if len(value.array) == 0 {
			fmt.Println("Empty command array")
			return
		}
		cmd := strings.ToUpper(value.array[0].bulk)
		fmt.Println("Received value:", value)
		args := value.array[1:]
		fmt.Println("Received args:", args)

		writer := NewWriter(conn)

		handler, exists := Handlers[cmd]
		if !exists {
			fmt.Println("Invalid command:", cmd)
			writer.Write(Value{typ: "string", str: ""})
			continue
		}
		result := handler(args)
		writer.Write(result)
		// conn.Write([]byte("+OK\r\n"))
	}
}
