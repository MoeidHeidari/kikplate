package main

import (
	"github.com/joho/godotenv"
	"github.com/kickplate/api/bootstrap"
)

func main() {
	godotenv.Load()
	err := bootstrap.RootApp.Execute()
	if err != nil {
		return
	}
}
