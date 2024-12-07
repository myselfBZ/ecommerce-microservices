package main

import "log"

func main() {
	s := NewAPI()
	log.Println("server running")
	log.Fatal("error server down: ", s.run(":8080"))
}
