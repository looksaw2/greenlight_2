package main

import "github.com/looksaw/greenlight_2/cmd/api"

func main() {
	app := api.ApiInit()
	app.RunHTTP()
}
