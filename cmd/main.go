package main

import (
	"flag"
	"git.friends.com/PetLand/UserService/v2/internal/config"
	"git.friends.com/PetLand/UserService/v2/internal/service"
)

func main() {
	flag.Parse()
	cfg := config.GetConfig()
	srv := service.NewService(cfg)

	srv.Logger.Fatal(srv.Start())
}
