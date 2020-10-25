package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"
)

const SSH_PORT = 22
const BANNER_CHARSET = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	// listen on SSH port
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", SSH_PORT))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	ch := make(chan int)
	// reporting activity
	go Report(ch)

	// accepting connections
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		// making clients wait
		go func(c net.Conn) {
			defer c.Close()

			// signaling connection
			ch <- 1

			var err error = nil
			for err == nil {
				time.Sleep(10 * time.Second)
				banner := append(RandBytes(), byte(13), byte(10))
				_, err = c.Write(banner)
			}

			// signaling disconnection
			ch <- 0
		}(conn)
	}
}

// RandBytes return 10 to 250 random bytes
func RandBytes() []byte {
	n := rand.Intn(50) + 10
	buf := make([]byte, n)

	_, err := rand.Read(buf)
	if err != nil {
		return []byte("X")
	}

	for i, b := range buf {
		buf[i] = BANNER_CHARSET[int(b)%len(BANNER_CHARSET)]
	}

	return buf
}

// Report print a message on connection/disconnection
func Report(ch chan int) {
	var nb int64 = 0

	for {
		activity := <-ch

		if activity == 1 {
			// connection
			nb++
		} else {
			// disconnection
			nb--
		}

		log.Printf("%d clients trapped\n", nb)
	}
}
