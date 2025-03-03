package main

import "log"

func main() {
	TryLoadConfig()

	log.Println("Starting server at port:", LoadedConfig.Port)

	StartHttpServer()
}
