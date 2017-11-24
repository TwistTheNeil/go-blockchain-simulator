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

func OpenListener(port int, new_connection chan net.Conn, what chan bool) {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println("OMG 1: Neil do something!")
		fmt.Println(err)
	}
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				fmt.Println("OMG 2: Neil do something!")
			}
			new_connection <- conn
		}
		what <- true
	}()

}

func GetPortFromPtr(port *int) int {
	x := fmt.Sprintf("%d", *port)
	p, _ := strconv.Atoi(x)
	return p
}

func main() {
	miner_port_flag := flag.Int("mport", 9876,
		"Port on which server listens on for miner")
	client_port_flag := flag.Int("cport", 9875,
		"Port on which server listens on client")

	flag.Parse()

	miners := make(map[net.Conn]int)
	clients := make(map[net.Conn]int)

	var miner_port int = GetPortFromPtr(miner_port_flag)
	var client_port int = GetPortFromPtr(client_port_flag)

	new_miner_connection := make(chan net.Conn)
	new_client_connection := make(chan net.Conn)
	remove_miner_connection := make(chan net.Conn)
	remove_client_connection := make(chan net.Conn)

	// TODO: Get rid of this channel eventually
	what := make(chan bool)

	OpenListener(miner_port, new_miner_connection, what)
	OpenListener(client_port, new_client_connection, what)

	for {
		select {
		case conn := <-new_miner_connection:
			miners[conn] = 0
			fmt.Printf("GOT CONN")
			fmt.Println(conn)

			go func(conn net.Conn) {
				dec := gob.NewDecoder(conn)
				for {
					var str string
					err := dec.Decode(&str)
					if err != nil {
						fmt.Println("DEC ERROR", err)
						break
					}
					fmt.Println("GOT HERE")
					fmt.Println(str)
					fmt.Println("FIN HERE")
				}
				remove_miner_connection <- conn
				fmt.Println("DEBUG: HIT THIS")
				fmt.Println(conn)
			}(conn)

		case conn := <-new_client_connection:
			clients[conn] = 0

		case conn := <-remove_miner_connection:
			delete(miners, conn)

		case conn := <-remove_client_connection:
			delete(clients, conn)
		}
	}

	// TODO: remove stuff beyond this point
	<-what
	<-what
	fmt.Println("Debug: Success")
	os.Exit(0)
}
