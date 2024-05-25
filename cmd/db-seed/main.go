package main

import api "github.com/randuck-dev/rd-api/pkg"

func main() {
	api.InitDb()
	api.SeedDb()
}
