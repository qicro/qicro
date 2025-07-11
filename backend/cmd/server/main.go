package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/qicro/qicro/backend/internal/app"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// 创建应用实例
	application, err := app.New()
	if err != nil {
		log.Fatal("Failed to initialize application:", err)
	}

	// 运行应用
	if err := application.Run(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}