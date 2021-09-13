package services

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/dougtq/go-protobuf/pb"
)

type UserServiceServer interface {
	// AddUser(context.Context, *User) (*User, error)
	// AddUsers(UserService_AddUsersServer) error
	// AddUserVerbose(ctx context.Context, in *User, opts ...grpc.CallOption) (UserService_AddUserVerboseClient, error)
	// mustEmbedUnimplementedUserServiceServer()
}

type UserService struct {
	pb.UnimplementedUserServiceServer
}

func NewUserService() *UserService {
	return &UserService{}
}

func (*UserService) AddUser(ctx context.Context, req *pb.User) (*pb.User, error) {

	fmt.Println(req.Name)

	return &pb.User{
		Id:    "1",
		Name:  req.GetName(),
		Email: req.GetEmail(),
	}, nil
}

func (*UserService) AddUsers(stream pb.UserService_AddUsersServer) error {
	users := []*pb.User{}

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("Stream has ended")
			return stream.SendAndClose(&pb.Users{
				User: users,
			})
		}

		if err != nil {
			log.Fatalf("Unknown error ocurred on stream %v", err)
		}

		users = append(users, &pb.User{
			Id:    req.GetId(),
			Name:  req.GetName(),
			Email: req.GetEmail(),
		})

		fmt.Println("Adding", req.GetId())
		fmt.Println(users)
	}
}

func (*UserService) AddUserVerbose(req *pb.User, stream pb.UserService_AddUserVerboseServer) error {
	stream.Send(&pb.UserResultStream{
		Status: "init",
		User:   &pb.User{},
	})

	time.Sleep(time.Second * 2)

	stream.Send(&pb.UserResultStream{
		Status: "inserted",
		User: &pb.User{
			Id:    "1",
			Name:  req.GetName(),
			Email: req.GetEmail(),
		},
	})

	time.Sleep(time.Second * 2)

	return nil
}

func (*UserService) AddUserTwoWayStream(stream pb.UserService_AddUserTwoWayStreamServer) error {
	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error receiveng client streamed data: %v", err)
		}

		err = stream.Send(&pb.UserResultStream{
			Status: "inserted",
			User:   req,
		})

		if err != nil {
			log.Fatalf("Error sending data: %v", err)
		}
	}
}
