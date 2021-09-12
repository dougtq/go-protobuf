package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/dougtq/go-protobuf/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:5051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Couldn't not connect to the server: %v", err)
	}

	defer connection.Close()

	client := pb.NewUserServiceClient(connection)

	new_user := AddUser(client)
	fmt.Println(new_user)

	AddUserVerbose(client)

}

func AddUser(client pb.UserServiceClient) *pb.User {
	req := &pb.User{
		Name:  "Doug",
		Email: "d@d.com",
	}

	res, err := client.AddUser(context.Background(), req)

	if err != nil {
		log.Fatalf("Couldn't not send grpc request: %v", err)
	}

	return res
}

func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{
		Name:  "Doug",
		Email: "d@d.com",
	}

	responseStream, err := client.AddUserVerbose(context.Background(), req)

	if err != nil {
		log.Fatalf("Couldn't not send grpc request: %v", err)
	}

	for {
		stream, err := responseStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Failure receiveing data: %v", err)
		}
		fmt.Println("Status:", stream.Status, " - ", stream.GetUser())
	}
}
