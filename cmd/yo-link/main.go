package main

import (
	"fmt"

	"github.com/yowie645/Yo-Link/internal/config"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)

	//TODO: init logger : slog import "log/slog"
	//TODO: init storage : sqlite
	//TODO: init router : chi, "chi render"
	//TODO: run server
}
