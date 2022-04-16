package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
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

func Broadcast() {
	clients := make(map[Client]bool)
	for {
		select {
		case message := <-messages: // Se le notifica a todos los clientes que llego el mensaje.
			for client := range clients {
				client <- message
			}
		case newClient := <-incomingClients: // Cuando un cliente nuevo se conecta, que lo agregamos a map.
			clients[newClient] = true
		case leavingClient := <-leavingClients: // Cuando un cliente se desconecta
			delete(clients, leavingClient)
			close(leavingClient)
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		log.Fatal(err)
	}
	go Broadcast()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go HandleConnection(conn)
	}
}
