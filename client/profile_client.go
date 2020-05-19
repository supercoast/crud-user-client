package client

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
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

func (profileClient *ProfileClient) CreateProfile(profile *pb.Profile) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	profileID, err := profileClient.service.CreateProfile(ctx, profile)
	if err != nil {
		return "", fmt.Errorf("Couldn't create profile: %v", err)
	}

	return profileID.GetId(), nil

}

func (profileClient *ProfileClient) UploadImage(imagePath string) (string, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("Not able to open image: ", err)
	}
	defer file.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := profileClient.service.CreateImage(ctx)
	if err != nil {
		return "", fmt.Errorf("Not able to upload image: ", err)
	}

	imageType := "." + strings.Split(imagePath, ".")[1]

	req := &pb.Image{
		ImageOneof: &pb.Image_ImageMetaData{
			ImageMetaData: &pb.ImageMetadata{
				Type: imageType,
			},
		},
	}

	err = stream.Send(req)
	if err != nil {
		return "", fmt.Errorf("Not able to send image metadata info to server: ", err, stream.RecvMsg(nil))
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}

		if err != nil {
			return "", fmt.Errorf("Not able to read chunk to buffer: ", err)
		}

		req := &pb.Image{
			ImageOneof: &pb.Image_ImageData{
				ImageData: &pb.ImageData{
					Data: buffer[:n],
				},
			},
		}

		err = stream.Send(req)
		if err != nil {
			return "", fmt.Errorf("Not able to send chunk to server: ", err, stream.RecvMsg(nil))
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		return "", fmt.Errorf("Not able to receive response: ", err)
	}

	return res.GetId(), nil

}
