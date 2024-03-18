package main

import (
	"testing"
	"github.com/Toront0/poker/internal/api"
)

func TestMain(t *testing.T) {

	
	server := api.NewServer(":3000")

	server.Run()

}