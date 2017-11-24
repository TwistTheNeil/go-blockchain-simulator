package main

import (
	//"./lib/ds"
	//	"crypto/sha256"
	"encoding/gob"
	"flag"
	"fmt"
	"net"
	"strconv"
)

func OpenConnection(addr string, port int) net.Conn {
	conn, err := net.Dial("tcp", addr+":"+strconv.Itoa(port))
	if err != nil {
		fmt.Printf("ERROR: ")
		fmt.Println(err)
	}
	return conn
}

func Appender() func(string) (int, string) {
	nonce := -1
	return func(x string) (int, string) {
		nonce += 1
		return nonce, x + strconv.Itoa(nonce)
	}
}

func GetPortFromPtr(port *int) int {
	x := fmt.Sprintf("%d", *port)
	p, _ := strconv.Atoi(x)
	return p
}

func main() {
	server_address_flag := flag.String("dest", "localhost",
		"Address of the server")
	server_port_flag := flag.Int("port", 9876,
		"Port used to connect to server")

	flag.Parse()

	var server_port int = GetPortFromPtr(server_port_flag)

	server_connection := OpenConnection(*server_address_flag, server_port)
	defer server_connection.Close()

	// Create a new Encoder for communicating with server
	enc := gob.NewEncoder(server_connection)

	str := "test string of a large length0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000\n"

	// Send encoded gob to server
	enc.Encode(&str)

	//	AppendNonce := Appender()
	//	h := sha256.New()
	//	target := "00000"
	//
	//	for {
	//		nonce, AppendedNonce := AppendNonce("hello")
	//		h.Write([]byte(AppendedNonce))
	//
	//		calculated_hash := h.Sum(nil)
	//
	//		fmt.Printf("%d %x\n", nonce, calculated_hash)
	//		calculated_hash_s := fmt.Sprintf("%x", calculated_hash)
	//
	//		if calculated_hash_s[:len(target)] == target {
	//			fmt.Println(calculated_hash[:len(target)/2])
	//			fmt.Println(target)
	//			break
	//		}
	//	}
}
