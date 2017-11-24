package main

import (
	//	"./lib/ds"
	"encoding/gob"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
)

func main() {
	server_port_flag := flag.Int("port", 9875,
		"Port on which the server listens on")
	server_address_flag := flag.String("addr", "localhost",
		"Address of the server")
	payload_flag := flag.String("payload", "test",
		"Payload to add to blockchain")

	flag.Parse()

	conn, err := net.Dial("tcp", *server_address_flag+":"+strconv.Itoa(*server_port_flag))
	if err != nil {
		fmt.Println("Error connecting: ", err)
		os.Exit(1)
	}

	enc := gob.NewEncoder(conn)
	enc_err := enc.Encode(&payload_flag)

	if enc_err != nil {
		fmt.Println("Error sending to server: ", enc_err)
	}
}
