package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/supercoast/crud-user-client/client"
	"github.com/supercoast/crud-user-client/pb"
	"google.golang.org/grpc"
)

const (
	address = "127.0.0.1"
	port    = "8080"
)

func main() {
	email := flag.String("email", "", "Email of user")
	lastName := flag.String("lastname", "", "Lastname of user")
	givenName := flag.String("givenname", "", "Given name of user")
	imagePath := flag.String("imagepath", "", "Path to image for the upload")
	flag.Parse()

	conn, err := grpc.Dial(strings.Join([]string{address, port}, ":"), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Not able to connect: %w", err)
	}
	defer conn.Close()

	c := client.NewProfileClient(conn)
	profile := &pb.Profile{
		Email:     *email,
		LastName:  *lastName,
		GivenName: *givenName,
	}
	profileID, err := c.CreateProfile(profile)
	if err != nil {
		log.Fatalf("Profile creation failed: %v", err)
	}

	log.Printf("Profile ID: %s\n", profileID)

	if *imagePath != "" {
		if !strings.HasPrefix(*imagePath, "/") {
			fmt.Println("Please provide absolut path for image: ", *imagePath)
		}
		imageID, err := c.UploadImage(*imagePath)
		if err != nil {
			log.Fatalf("Image uploade failed: %v", err)
		}
		log.Printf("Image ID: %s\n", imageID)
	}

}
