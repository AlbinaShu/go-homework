package main

import (
	"fmt"
	"os"
	"runtime"
)

func main() {
	// Имя пользователя из переменной окружения
	username := os.Getenv("USER")

	fmt.Println("Имя пользователя:", username)

	// Аргументы командной строки
	fmt.Println("Аргументы CLI:", os.Args[1:])

	// Текущая версия Go
	fmt.Println("Версия Go:", runtime.Version())
}