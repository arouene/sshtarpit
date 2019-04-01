package main

import (
	"bytes"
	"log"
	"math/rand"
	"net"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	// listen on SSH port
	l, err := net.Listen("tcp", ":2222")
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
				_, err = c.Write(RandBytes())
			}

			// signaling disconnection
			ch <- 0
		}(conn)
	}
}

// RandBytes return 10 to 250 random bytes
func RandBytes() []byte {
	var b bytes.Buffer

	l := rand.Intn(240) + 10
	for ; l > 0; l-- {
		b.WriteByte(byte(rand.Int()))
	}

	return b.Bytes()
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
