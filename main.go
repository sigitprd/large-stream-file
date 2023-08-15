package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type FileServer struct {
}

func (fs *FileServer) start() {
	ln, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go fs.readLoop(conn)
	}
}

func (fs *FileServer) readLoop(conn net.Conn) {
	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		file := buf[:n]
		fmt.Println(file)
		fmt.Printf("received %d bytes over the network\n", n)
	}
}

func sendFile(size int) error {
	file := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, file)
	if err != nil {
		return err
	}

	conn, err := net.Dial("tcp", "127.0.0.1:3000") // Specify the port number
	if err != nil {
		return err
	}
	defer conn.Close() // Close the connection when done

	n, err := conn.Write(file)
	if err != nil {
		return err
	}

	fmt.Printf("written %d bytes over the network\n", n)
	return nil
}

func main() {
	fs := &FileServer{}
	go fs.start()

	time.Sleep(1 * time.Second) // Sleep briefly to allow the server to start

	go func() {
		time.Sleep(4 * time.Second)
		sendFile(4000)
	}()

	// Sleep to allow time for the file to be sent before the program exits
	time.Sleep(10 * time.Second)
}
