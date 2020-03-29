package main

import (
	"shorterUrl/src/app"
	"shorterUrl/src/env"
)

func main() {
	a := app.App{}
	a.Initialize(env.GetEnv())
	a.Run(":8080")
}
