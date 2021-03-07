package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/jaimeteb/chatto/bot"
	"github.com/jaimeteb/chatto/internal/logger"
	"github.com/jaimeteb/chatto/version"
)

func main() {
	port := flag.Int("port", 4770, "Specify port to use.")
	path := flag.String("path", ".", "Path to YAML files.")
	debug := flag.Bool("debug", false, "Enable debug logging.")
	vers := flag.Bool("version", false, "Display version.")
	flag.Parse()

	if *vers {
		fmt.Println(version.BuildStr())
		return
	}

	if strings.EqualFold(os.Getenv("CHATTO_BOT_DEBUG"), "true") {
		*debug = true
	}

	logger.SetLogger(*debug)

	server := bot.NewServer(*path, *port)

	server.Run()
}
