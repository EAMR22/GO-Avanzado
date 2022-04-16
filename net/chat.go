package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
)

type Client chan<- string

var (
	incomingClients = make(chan Client)
	leavingClients  = make(chan Client)
	messages        = make(chan string)
)

var (
	host = flag.String("h", "localhost", "host")
	port = flag.Int("p", 3090, "port")
)

// Client1 -> Server -> HandleConnection(Client1)

func HandleConnection(conn net.Conn) {
	defer conn.Close()
	message := make(chan string)
	go MessageWrite(conn, message)
	// Client1:2560 Platzi.com, 38
	// platzi.com:38
	clientName := conn.RemoteAddr().String()

	message <- fmt.Sprintf("Welcome to the server, your name %s\n", clientName)
	messages <- fmt.Sprintf("New client is here, name %s\n", clientName)
	incomingClients <- message

	inputMessage := bufio.NewScanner(conn) // Empieza a leer todo lo que se esta escribiendo a travez de la terminal.
	for inputMessage.Scan() {              // Vamos a ser capaz de enviar todo lo que se a recibido a travez de messages.
		messages <- fmt.Sprintf("%s: %s\n", clientName, inputMessage.Text())
	}

	leavingClients <- message // Rompe el ciclo y cancelo el programa, por lo tanto el cliente a abandonado el chat.
	messages <- fmt.Sprintf("%s said goodbye", clientName)
}

func MessageWrite(conn net.Conn, messages <-chan string) {
	for message := range messages {
		fmt.Fprintln(conn, message)
	}
}
