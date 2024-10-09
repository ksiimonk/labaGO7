package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Разрешить все источники (не рекомендуется в продакшене)
		},
	}
	clients = make(map[*websocket.Conn]bool) // Подключенные клиенты
	mu      sync.Mutex                       // Для синхронизации доступа к clients
)

func main() {
	http.HandleFunc("/ws", handleConnections)

	port := ":8080"
	fmt.Println("Сервер запущен на порту", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
	}
}

// Обработчик соединений
func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Ошибка при обновлении соединения:", err)
		return
	}
	defer conn.Close()

	mu.Lock()
	clients[conn] = true
	mu.Unlock()

	for {
		_, msg, err := conn.ReadMessage() // Чтение сообщения
		if err != nil {
			fmt.Println("Ошибка при чтении сообщения:", err)
			mu.Lock()
			delete(clients, conn)
			mu.Unlock()
			break
		}

		fmt.Printf("Получено сообщение: %s\n", msg)
		broadcastMessage(string(msg)) // Отправка сообщения всем клиентам
	}
}

// Функция для рассылки сообщений всем клиентам
func broadcastMessage(message string) {
	mu.Lock()
	defer mu.Unlock()

	for client := range clients {
		if err := client.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			fmt.Println("Ошибка при отправке сообщения:", err)
			client.Close()
			delete(clients, client)
		}
	}
}
