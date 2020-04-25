package main

import (
	"log"
	"strings"

	"github.com/supercoast/crud-user-client/client"
	"google.golang.org/grpc"
)

const (
	address     = "127.0.0.1"
	port        = "8080"
	imageFolder = "tmp"
)

func main() {
	conn, err := grpc.Dial(strings.Join([]string{address, port}, ":"), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Not able to connect: %w", err)
	}
	defer conn.Close()

	c := client.NewProfileClient(conn)
	c.UploadImage(imageFolder + "/sample-image.jpeg")

}
