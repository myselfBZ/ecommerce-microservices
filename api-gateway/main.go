package main

import (
	"log"
	"os"
)

func main() {
	s := NewAPI()
	log.Println(s.run(os.Getenv("api_gateway")))
}
