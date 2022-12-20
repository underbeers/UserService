package main

import (
	"flag"
	"git.friends.com/PetLand/UserService/v2/internal/config"
	"git.friends.com/PetLand/UserService/v2/internal/service"
)

func main() {
	debugMode := flag.Bool("use_db_config", true, "use for starting locally in debug mode")
	flag.Parse()
	cfg := config.GetConfig(*debugMode)
	srv := service.NewService(cfg)

	srv.Logger.Fatal(srv.Start())
}
