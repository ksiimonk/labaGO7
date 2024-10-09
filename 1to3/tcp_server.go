package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var wg sync.WaitGroup

func main() {
	port := ":8080"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("Сервер запущен на порту", port)

	// Канал для сигналов
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-done
		fmt.Println("\nЗавершение работы сервера...")
		listener.Close() // Закрываем слушатель
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			// Проверяем, был ли сервер закрыт
			if ne, ok := err.(*net.OpError); ok && ne.Op == "accept" {
				break // Выход из цикла, если сервер закрыт
			}
			fmt.Println("Ошибка при приеме соединения:", err)
			continue
		}
		wg.Add(1) // Увеличиваем счётчик горутин
		go handleConnection(conn)
	}

	wg.Wait() // Ожидаем завершения всех горутин
	fmt.Println("Все соединения завершены. Сервер остановлен.")
}

func handleConnection(conn net.Conn) {
	defer wg.Done()    // Уменьшаем счётчик горутин
	defer conn.Close() // Закрываем соединение после завершения обработки

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Ошибка при чтении данных:", err)
		return
	}

	message := string(buffer[:n])
	fmt.Println("Получено сообщение:", message)

	confirmation := "Сообщение получено"
	_, err = conn.Write([]byte(confirmation))
	if err != nil {
		fmt.Println("Ошибка при отправке данных:", err)
		return
	}
}
