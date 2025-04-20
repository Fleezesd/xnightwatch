package main

import (
	_ "go.uber.org/automaxprocs/maxprocs"

	"github.com/fleezesd/xnightwatch/cmd/gateway/app"
)

func main() {
	app.NewApp().Run()
}
