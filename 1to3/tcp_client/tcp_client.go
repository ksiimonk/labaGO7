package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// Указываем адрес сервера и порт
	serverAddress := "localhost:8080"
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		fmt.Println("Ошибка при подключении к серверу:", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Считываем сообщение от пользователя
	fmt.Print("Введите сообщение для отправки: ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		message := scanner.Text()

		// Отправляем сообщение серверу
		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Ошибка при отправке сообщения:", err)
			return
		}
	}

	// Читаем ответ от сервера
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Ошибка при чтении ответа:", err)
		return
	}

	fmt.Println("Ответ от сервера:", string(buffer[:n]))
}
