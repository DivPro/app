package main

import (
	"flag"

	"github.com/DivPro/app/internal/app"
)

func main() {
	dev := flag.Bool("dev", false, "development mode")

	flag.Parse()
	app.Run(*dev)
}
