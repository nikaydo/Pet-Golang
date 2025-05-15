package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	TLS "main/internal/TLS"
	"main/internal/metric"
	"main/internal/models"
	"sync"
	"time"

	"github.com/quic-go/quic-go"
)

const addr = "localhost:4242"

var (
	connect = make(map[quic.Connection]models.Settings)
)

func main() {
	listener, err := quic.ListenAddr(addr, TLS.GenerateTLSConfig(), nil)
	if err != nil {
		return
	}
	defer listener.Close()
	mutex := sync.Mutex{}
	for {
		conn, err := listener.Accept(context.Background())
		if err != nil {
			return
		}
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			fmt.Println("Failed to accept stream:", err)
			return
		}
		var length uint32
		if err := binary.Read(stream, binary.BigEndian, &length); err != nil {
			fmt.Println("Failed to read settings length:", err)
			return
		}
		buf := make([]byte, length)
		if _, err := io.ReadFull(stream, buf); err != nil {
			fmt.Println("Failed to read settings:", err)
			return
		}
		var settings models.Settings
		if err := json.Unmarshal(buf, &settings); err != nil {
			fmt.Println("Failed to unmarshal settings:", err)
			return
		}
		mutex.Lock()
		connect[conn] = settings
		mutex.Unlock()
		go f(conn)
	}
}

func f(conn quic.Connection) {
	upd := connect[conn].UpdateTime
	for {
		stream, err := conn.OpenStreamSync(context.Background())
		if err != nil {
			return
		}
		defer stream.Close()
		n, _ := metric.Metrics()
		b, _ := n.Marshal()
		binary.Write(stream, binary.BigEndian, uint32(len(b)))
		stream.Write([]byte(b))
		stream.Close()
		time.Sleep(time.Duration(upd) * time.Millisecond)
	}
}
