package xnightwatch

import (
	"github.com/fleezesd/xnightwatch/cmd/xnightwatch/app"
	_ "go.uber.org/automaxprocs/maxprocs"
)

func main() {
	app.NewApp().Run()
}
