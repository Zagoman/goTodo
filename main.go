package main

func main() {
	server := NewApiServer(":3001")
	server.Run()
}
