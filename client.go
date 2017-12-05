package main

import (
	//	"./lib/ds"
	"encoding/gob"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

func OpenConnection(addr string, port int) (net.Conn, error) {
	var err error
	var conn net.Conn

	for i := 0; i < 5; i++ {
		conn, err = net.Dial("tcp", addr+":"+strconv.Itoa(port))
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	return conn, err
}

func main() {
	server_port_flag := flag.Int("port", 9875,
		"Port on which the server listens on")
	server_address_flag := flag.String("addr", "localhost",
		"Address of the server")
	payload_flag := flag.String("payload", "test",
		"Payload to add to blockchain")

	flag.Parse()

	server_connection, err := OpenConnection(*server_address_flag, *server_port_flag)
	if err != nil {
		fmt.Println("Error while connecting to server:", err)
		os.Exit(2)
	}
	defer server_connection.Close()

	enc := gob.NewEncoder(server_connection)
	enc_err := enc.Encode(&payload_flag)

	if enc_err != nil {
		fmt.Println("Error sending to server: ", enc_err)
		os.Exit(3)
	}
}
