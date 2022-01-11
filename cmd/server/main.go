package main

import (
	"fmt"

	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load(".env")
}

func main() {
	fmt.Println("OK")
}
