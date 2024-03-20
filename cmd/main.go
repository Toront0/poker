package main

import (
	"github.com/Toront0/poker/internal/api"

)

func main() {
	

	server := api.NewServer(":3000")

	server.Run()

}