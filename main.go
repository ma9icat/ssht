package main

import (
	"log"
	"os"
	"ssht/internal/app"
	"ssht/internal/config"
)

func main() {
	// 检查是否有参数，没有则显示帮助
	if len(os.Args) == 1 {
		config.ShowHelp()
		os.Exit(0)
	}

	if err := app.Run(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
