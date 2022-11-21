package main

import (
	"flag"
	"git.friends.com/PetLand/UserService/v2/internal/config"
	"git.friends.com/PetLand/UserService/v2/internal/service"
)

func main() {
	isLocalFlag := flag.Bool("use_local_config", true, "use for starting locally in debug mode")
	flag.Parse()
	cfg := config.GetConfig(*isLocalFlag)
	srv := service.NewService(cfg)

	srv.Logger.Fatal(srv.Start())
}
