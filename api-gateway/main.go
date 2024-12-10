package main

import "os"

func main() {
	s := NewAPI()
	s.run(os.Getenv("server_addr"))
}
