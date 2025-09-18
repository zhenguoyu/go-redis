package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

// aof.go 数据持久化
type Aof struct {
	file *os.File
	rd   *bufio.Reader
	mu   sync.Mutex
}

func NewAof(filename string) (*Aof, error) {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}
	aof := &Aof{
		file: f,
		rd:   bufio.NewReader(f),
		mu:   sync.Mutex{},
	}
	go func() {
		for {
			aof.mu.Lock()
			aof.file.Sync()
			aof.mu.Unlock()
			// 每秒同步一次
			time.Sleep(time.Second)
		}
	}()
	return aof, nil
}

// Close 关闭AOF文件
func (a *Aof) Close() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.file.Close()
}

// Write 将数据写入AOF文件
func (a *Aof) Write(value Value) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	_, err := a.file.Write(value.Marshal())
	return err
}

func (a *Aof) Read(callback func(value Value)) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	resp := NewResp(a.file)
	for {
		value, err := resp.Read()
		if err == nil {
			callback(value)
		}
		if err == io.EOF {
			break
		}
	}
	return nil
}

func aofCallback(value Value) {
	cmd := strings.ToUpper(value.array[0].bulk)
	args := value.array[1:]

	handler, ok := Handlers[cmd]
	if !ok {
		fmt.Println("Invalid command in AOF:", cmd)
		return
	}
	handler(args)
}
