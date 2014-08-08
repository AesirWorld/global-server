package main

import "log"
import "net"
import "os"
import "github.com/AesirWorld/global-server/account"

type Client struct {
	conn       net.Conn
	account_id int
	login_id1  uint32
	login_id2  uint32
}

type AuthDB struct {
	login_id1 uint32
	login_id2 uint32
}

var auth_db = make(map[int]*AuthDB)

func main() {
	log.Println("Starting global-server...")

	ln, err := net.Listen("tcp", ":6900")

	if err != nil {
		log.Println("Failed to listen...")
		log.Println(err)
		os.Exit(1)
	}

	log.Println("Global-server listening on port", 6900)

	// Connect to account database
	account.Connect()

	// Accept connections
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		// Send conenction to some goroutine/handler
		go initConnection(conn)
	}
}

func initConnection(c net.Conn) {
	// TODO:
	// Check for ipbans, etc.

	defer func() {
		c.Close()
		log.Printf("Connection from %v closed.\n", c.RemoteAddr())
	}()

	// Receive packets
	for {
		packet := make([]byte, 1024)

		length, _ := c.Read(packet)

		// Connection-closed
		if length == 0 {
			log.Println("Len 0, conenction closed")
			break
		}

		// First 2 bytes represent the packet_id
		packet_id := int(packet[0])<<0 | int(packet[1])<<8

		log.Printf("Received packed_id: %d - %#04x\n", packet_id, packet_id)

		switch packet_id {
		// Heartbeet packet
		case PKT_HEARTBEAT:
			log.Println("Heartbeet request")
			c.Write([]byte("pong"))
			break
		// Login packets
		case
			PKT_ENTER,  // Request client login
			PKT_ENTER2: // Request client login with encrypt
			clientEnter(c, packet)
			break
		// Pass connection to charServer handler if sucessful
		case PKT_CHR_ENTER:
			charServerEnter(c, packet)
			return
		default:
			log.Printf("Abnormal end of connection (ip: %s): Unknown packet 0x%x\n", c.RemoteAddr(), packet_id)
			return
		}
	}
}
