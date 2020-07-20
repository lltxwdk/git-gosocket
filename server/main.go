package main

import (
	"GoSocket/git-gosocket/server/protocol"
	"fmt"
	"log"
	"net"
	"os"
)

func Log(v ...interface{}) {
	log.Println(v...)
}

func reader(readerChannel chan []byte) {
	for {
		select {
		case data := <-readerChannel:
			Log(string(data))
		}
	}
}

func handleProtocalConnection(conn net.Conn) {

	tmpBuffer := make([]byte, 0)
	readerChannel := make(chan []byte, 16)
	go reader(readerChannel)

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			Log(conn.RemoteAddr().String(), " connection error:%s", err.Error())
			return
		}
		tmpBuffer = protocol.Depack(append(tmpBuffer, buffer[:n]...), readerChannel)
	}

	defer conn.Close()
}

func main() {
	fmt.Println("begin")
	netListen, err := net.Listen("tcp", "127.0.0.1:6060")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error:$s", err.Error())
		os.Exit(1)
	}

	defer netListen.Close()
	Log("Waiting for clients connet")
	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}
		Log(conn.RemoteAddr().String(), "tcp connect success")

		go handleProtocalConnection(conn)
	}
}
