package xnightwatch

import (
	_ "go.uber.org/automaxprocs/maxprocs"

	"github.com/fleezesd/xnightwatch/cmd/xnightwatch/app"
)

func main() {
	app.NewApp().Run()
}
