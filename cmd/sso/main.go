package main

import (
	"fmt"

	"github.com/GazpachoGit/sso/internal/config"
)

func main() {
	//init config
	cfg := config.MustLoad()
	fmt.Printf("Loaded config: %+v\n", cfg)
	//init log(slog)

	//init app

	//run gRPC server

}
