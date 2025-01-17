package main

import (
	"fmt"
	"log"
	"match/internal/app"
	"match/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	application, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("Ошибка инициализации приложения: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
	fmt.Println("Сервер остановлен.")

}
