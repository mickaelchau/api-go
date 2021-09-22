package main

import (
	"fmt"

	handler "gitlab.com/tcmlabs/api-webserver/pokeserver/controller"
)

func main() {
	fmt.Println("Application succesfully launched")
	handler.HandleRequests()
}
