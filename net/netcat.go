package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

var (
	port = flag.Int("p", 3090, "port")
	host = flag.String("h", "localhost", "host")
)

// host:port
// Escribir -> host:port
// Leer -> host:port
// -> [Hola] -> host:port -> [Hola]

func main() {
	flag.Parse()
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		io.Copy(os.Stdout, conn) // El os.stdout nos va a permitir escribir en la consola los resultados que estamos recibiendo.
		done <- struct{}{}       // y conn va a ser el lector.
	}()
	CopyContent(conn, os.Stdin)
	conn.Close()
	<-done
}

func CopyContent(dst io.Writer, src io.Reader) {
	_, err := io.Copy(dst, src) // io.Copy espera recibir una inerfaz que pueda escribir y otra interfaz que pueda leer.
	if err != nil {
		log.Fatal(err)
	}
}
