package main

import (
	"context"
	"fmt"
	"grpc-udemy/greet/greetpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to create dial: %v", err)
	}
	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)
	// doUnary(c)
	// doServerStreaming(c)
	// doClientStreaming(c)
	// doBidirectionalStreaming(c)
	doDeadlineExceed(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetRequest{Greeting: &greetpb.Greeting{
		FirstName: "Filip",
	}}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("failed to greet: %v", err)
	}
	fmt.Println(res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	req := &greetpb.GreeetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Filip",
		},
	}
	stream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("failed to greet: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Println(res.Result)
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling longGreet: %v", err)
	}

	requests := []*greetpb.LongGreetRequest{
		{Greeting: &greetpb.Greeting{FirstName: "Filip"}},
		{Greeting: &greetpb.Greeting{FirstName: "John"}},
		{Greeting: &greetpb.Greeting{FirstName: "Mark"}},
	}

	for _, req := range requests {
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	resp, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("failed to close stream: %v", err)
	}
	fmt.Println(resp.Result)
}

func doBidirectionalStreaming(c greetpb.GreetServiceClient) {
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("error while creating stream: %v", err)
	}

	requests := []*greetpb.GreetEveryoneRequest{
		{Greeting: &greetpb.Greeting{FirstName: "Filip"}},
		{Greeting: &greetpb.Greeting{FirstName: "John"}},
		{Greeting: &greetpb.Greeting{FirstName: "Mark"}},
	}

	waitc := make(chan struct{})

	//send
	go func() {
		for _, r := range requests {
			stream.Send(r)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	//recv
	go func() {
		for {
			resp, err := stream.Recv()

			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while recv: %v", err)
				break
			}
			fmt.Printf("recv: %v\n", resp.Result)
		}
		close(waitc)
	}()

	<-waitc
}

func doDeadlineExceed(c greetpb.GreetServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	res, err := c.GreetWithDeadline(ctx, &greetpb.GreetRequest{Greeting: &greetpb.Greeting{FirstName: "John"}})
	fmt.Println(res, err)
}
