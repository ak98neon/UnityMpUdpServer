package main

import (
	"log"
	"math/rand"
	"net"
	"strings"
	"time"
)

const maxBufferSize = 1024

func main() {
	server(":4252")
}

func server(address string) {
	log.Println("Server is starting")
	s, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		log.Println(err)
		return
	}

	conn, err := net.ListenUDP("udp4", s)
	if err != nil {
		log.Println(err)
		return
	}

	defer conn.Close()

	doneChan := make(chan error, 1)
	buffer := make([]byte, maxBufferSize)
	rand.Seed(time.Now().Unix())

	go func() {
		for {
			n, addr, err := conn.ReadFromUDP(buffer)
			if err != nil {
				log.Println(err)
				doneChan <- err
				return
			}

			log.Print("-> ", string(buffer[0:n-1]))

			if strings.TrimSpace(string(buffer[0:n])) == "STOP" {
				log.Println("Exiting UDP server!")
				return
			}

			data := []byte("lol")
			log.Printf("data: %s\n", string(data))
			_, err = conn.WriteToUDP(data, addr)

			if err != nil {
				log.Println(err)
				doneChan <- err
				return
			}
		}
	}()

	err = <-doneChan
	return
}
