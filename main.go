package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/supercoast/crud-user-client/client"
	"google.golang.org/grpc"
)

const (
	address = "127.0.0.1"
	port    = "8080"
)

func main() {
	imagePath := flag.String("imagepath", "", "Path to image for the upload")
	flag.Parse()

	if !strings.HasPrefix(*imagePath, "/") {
		fmt.Println("Please provide absolut path for image: ", *imagePath)
	}

	conn, err := grpc.Dial(strings.Join([]string{address, port}, ":"), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Not able to connect: %w", err)
	}
	defer conn.Close()

	c := client.NewProfileClient(conn)
	imageId, err := c.UploadImage(*imagePath)
	if err != nil {
		log.Fatalf("Image uploade failed: %v", err)
	}

	log.Printf("Image ID: %s\n", imageId)

}
