package client

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/supercoast/crud-user-client/pb"
	"google.golang.org/grpc"
)

type ProfileClient struct {
	service pb.ProfileServiceClient
}

func NewProfileClient(cc *grpc.ClientConn) *ProfileClient {
	service := pb.NewProfileServiceClient(cc)
	return &ProfileClient{service}
}

func (profileClient *ProfileClient) UploadImage(imagePath string) {
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatalf("Not able to open image: ", err)
	}
	defer file.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := profileClient.service.CreateProfile(ctx)
	if err != nil {
		log.Fatalf("Not able to upload image: ", err)
	}

	req := &pb.Profile{
		ProfileOneof: &pb.Profile_ProfileData{
			ProfileData: &pb.ProfileData{
				GivenName: "Valentin",
				LastName:  "Widmer",
				Birthday: &pb.Date{
					Day:   2,
					Month: 11,
					Year:  1994,
				},
				Email:     "valentin.widmer@protonmail.com",
				ImageType: filepath.Ext(imagePath),
			},
		},
	}

	err = stream.Send(req)
	if err != nil {
		log.Fatal("Not able to send profile info to server: ", err, stream.RecvMsg(nil))
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Not able to read chunk to buffer: ", err)
		}

		req := &pb.Profile{
			ProfileOneof: &pb.Profile_ImageData{
				ImageData: &pb.ImageData{
					Data: buffer[:n],
				},
			},
		}

		err = stream.Send(req)
		if err != nil {
			log.Fatal("Not able to send chunk to server: ", err, stream.RecvMsg(nil))
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Not able to receive response: ", err)
	}

	log.Printf("Image has been uploade with id: %s", res.GetId())

}
