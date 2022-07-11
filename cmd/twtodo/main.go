package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/simachri/taskwarrior-ms-todo/internal/cli"
	"github.com/simachri/taskwarrior-ms-todo/internal/mstodo"
)

func main() {
	godotenv.Load()
	client, err := mstodo.Authenticate(os.Getenv("TENANT_ID"), os.Getenv("CLIENT_ID"))
	if err != nil {
		os.Exit(1)
	}
    cli.Sync(client)
}
