package main

var Handlers = map[string]func([]Value) Value{
	"ping": ping,
}

func ping(args []Value) Value {
	return Value{typ: "string", str: "PONG"}
}
