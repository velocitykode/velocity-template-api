package main

import (
	"{{MODULE_NAME}}/internal/app"
	_ "{{MODULE_NAME}}/routes"
)

func main() {
	app.Run()
}
