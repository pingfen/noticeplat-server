package main

import (
	"github.com/pingfen/noticeplat-server/server/app"
	_ "github.com/pingfen/noticeplat-server/server/app/v1"
)

func main() {
	err := app.Get().Run()
	if err != nil {
		panic(err)
	}
}
