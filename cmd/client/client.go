package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/phelliperodrigues/full-cycle-grpc/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to gRPC Server %v", err)
	}

	defer connection.Close()
	client := pb.NewUserServiceClient(connection)
	// AddUser(client)
	// AddUserVerbose(client)
	// AddUsers(client)
	AddUserStreamBoth(client)

}

func AddUser(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Phellipe",
		Email: "email@email.com",
	}

	res, err := client.AddUser(context.Background(), req)

	if err != nil {
		log.Fatalf("Could not make gRPC Request %v", err)
	}

	fmt.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Phellipe",
		Email: "email@email.com",
	}

	responseStream, err := client.AddUserVerbose(context.Background(), req)

	if err != nil {
		log.Fatalf("Could not make gRPC Request %v", err)
	}
	for {
		stream, err := responseStream.Recv()

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not receive the msg: %v", err)
		}

		fmt.Println("Status: ", stream.Status, "User: ", stream.GetUser())
	}
}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		&pb.User{
			Id:    "p1",
			Name:  "Phellipe 1",
			Email: "1@email.com",
		},
		&pb.User{
			Id:    "p2",
			Name:  "Phellipe 2",
			Email: "2@email.com",
		},
		&pb.User{
			Id:    "p3",
			Name:  "Phellipe 3",
			Email: "3@email.com",
		},
		&pb.User{
			Id:    "p4",
			Name:  "Phellipe 4",
			Email: "4@email.com",
		},
		&pb.User{
			Id:    "p5",
			Name:  "Phellipe 5",
			Email: "5@email.com",
		},
	}
	stream, err := client.AddUsers(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}

	fmt.Println(res)
}

func AddUserStreamBoth(client pb.UserServiceClient) {
	stream, err := client.AddUserStreamBoth(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	reqs := []*pb.User{
		&pb.User{
			Id:    "p1",
			Name:  "Phellipe 1",
			Email: "1@email.com",
		},
		&pb.User{
			Id:    "p2",
			Name:  "Phellipe 2",
			Email: "2@email.com",
		},
		&pb.User{
			Id:    "p3",
			Name:  "Phellipe 3",
			Email: "3@email.com",
		},
		&pb.User{
			Id:    "p4",
			Name:  "Phellipe 4",
			Email: "4@email.com",
		},
		&pb.User{
			Id:    "p5",
			Name:  "Phellipe 5",
			Email: "5@email.com",
		},
	}

	wait := make(chan int)
	go func() {
		for _, req := range reqs {
			fmt.Println("Send user: ", req.Name)
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}

		stream.CloseSend()

	}()

	go func() {
		for {
			res, err := stream.Recv()

			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Could not receive the msg: %v", err)
				break
			}
			fmt.Printf("Recebendo user: %v com status: %v\n", res.GetUser().Name, res.Status)

		}
		close(wait)
	}()
	<-wait
}
