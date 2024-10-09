package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Middleware для логирования
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)

		log.Printf("Метод: %s, URL: %s, Время выполнения: %s\n", r.Method, r.URL, duration)
	})
}

// Обработчик для маршрута GET /hello
func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Привет!"))
}

// Обработчик для маршрута POST /data
func dataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Ошибка при разборе JSON", http.StatusBadRequest)
		return
	}

	// Выводим полученные данные в консоль
	fmt.Printf("Полученные данные: %+v\n", data)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Данные получены"))
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/hello", loggingMiddleware(http.HandlerFunc(helloHandler)))
	mux.Handle("/data", loggingMiddleware(http.HandlerFunc(dataHandler)))

	port := ":8080"
	fmt.Println("Сервер запущен на порту", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal("Ошибка при запуске сервера:", err)
	}
}

// curl http://localhost:8080/hello
// cd "C:\dev\projects\labPP7\4-5"
// cd curl -X POST http://localhost:8080/data -H "Content-Type: application/json" -d @data.json
