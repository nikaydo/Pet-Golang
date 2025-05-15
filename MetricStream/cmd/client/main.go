package main

import (
	"context"
	"crypto/tls"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"main/internal/models"

	"github.com/gosuri/uilive"
	"github.com/quic-go/quic-go"
)

const addr = "localhost:4242"

func main() {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
	}
	writer := uilive.New()
	writer.Start()
	defer writer.Stop()
	conn, err := quic.DialAddr(context.Background(), addr, tlsConf, nil)
	if err != nil {
		return
	}
	defer conn.CloseWithError(0, "")
	for {
		Hello(conn)
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			panic(err)
		}
		go f(stream, writer)
	}
}

func Hello(conn quic.Connection) {
	settings := models.Settings{UpdateTime: 500}
	s, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		panic(err)
	}
	jsonData, err := json.Marshal(settings)
	if err != nil {
		panic(err)
	}
	binary.Write(s, binary.BigEndian, uint32(len(jsonData)))
	s.Write(jsonData)
	s.Close()
}

func f(stream quic.Stream, w *uilive.Writer) {
	var length uint32
	if err := binary.Read(stream, binary.BigEndian, &length); err != nil {
		return
	}
	buf := make([]byte, length)
	if _, err := io.ReadFull(stream, buf); err != nil {
		return
	}
	var Data models.Data
	Data.Unmarshal(buf)
	form(w, Data)
}

func VramForm(n uint64) string {
	result := float64(n) / (1024 * 1024 * 1024)
	if result < 1 {
		result = float64(n) / (1024 * 1024)
		if result < 1 {
			result = float64(n) / 1024
			return fmt.Sprintf("%4.0f KB", result)
		}
		return fmt.Sprintf("%4.0f MB", result)
	}
	return fmt.Sprintf("%.2f GB", result)
}

func form(w *uilive.Writer, data models.Data) {
	fmt.Fprintf(w, "cpu : %.2f\n", data.Cpu.CpuUsage)
	fmt.Fprintf(w, "vram: %.2f%%  Total: %s  Free: %s\n", data.Vram.VramUsage, VramForm(data.Vram.VramTotal), VramForm(data.Vram.VramUsed))
	w.Flush()
}
