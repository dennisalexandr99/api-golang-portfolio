package main

import (
	"example.com/try-echo/db"
	"example.com/try-echo/routes"
)

func main() {
	db.Init()
	e := routes.Init()

	e.Logger.Fatal(e.Start(":6969"))
}
