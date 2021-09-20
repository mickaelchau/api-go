package main

import (
	handler "gitlab.com/tcmlabs/api-webserver/pokeserver/controller"
)

func main() {
	handler.HandleRequests()
}
