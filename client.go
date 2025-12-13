package main

import (
	"bufio"
	"fmt"
	"os"

	"chat/client"
)

func main() {
	// Valor por defecto
	name := "Anonymous"
	channel := "general"

	// os.Args[0] es el nombre del programa
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	if len(os.Args) > 2 {
		channel = os.Args[2]
	}

	fmt.Println("👤 Participant name:", name)
	fmt.Println("🛰  Channel:", channel)

	p := client.NewParticipant(name, channel)

	// Hilo que escucha mensajes del servidor
	go client.CheckMessages(p)

	// Leer texto desde consola
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		p.TextRead(scanner.Text())
	}
}
