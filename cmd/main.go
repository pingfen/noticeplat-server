package main

import (
	cos "github.com/bingbaba/storage/qcloud-cos"
	"github.com/pingfen/notify/router"
	"github.com/pingfen/notify/storage"
	"log"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	storage.Init(cos.NewStorage(cos.NewConfigByEnv()))
}

// @title Tracelog API
// @version 1.0
// @description tracelog
// @termsOfService http://127.0.0.1:8080
// @license.name MIT
func main() {

	router.Run()
}
