package main

import (
	"context"
	"fmt"
	"grpc-udemy/calculator/calculatorpb"
	"io"
	"log"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func main() {
	conn, err := grpc.Dial(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to create dial: %v", err)
	}
	defer conn.Close()

	c := calculatorpb.NewCalculatorClient(conn)
	// doUnary(c)
	// doPrimeNumberDecomposition(c)
	// doComputeAverage(c)
	// doFindMaximum(c)
	doErrorUnary(c)
}

func doUnary(c calculatorpb.CalculatorClient) {
	req := &calculatorpb.SumRequest{X: 3, Y: 5}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("failed to greet: %v", err)
	}
	fmt.Println(res.Sum)
}

func doPrimeNumberDecomposition(c calculatorpb.CalculatorClient) {
	req := &calculatorpb.PrimeNumberDecompositionRequest{Number: 120}
	stream, err := c.PrimeNumberDecomposition(context.Background(), req)

	if err != nil {
		log.Fatalf("failed to do prime number decomposition stream: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		fmt.Println(res.Number)
	}
}

func doComputeAverage(c calculatorpb.CalculatorClient) {
	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("failed to open stream: %v", err)
	}
	for i := 1; i <= 4; i++ {
		stream.Send(&calculatorpb.AverageRequest{Number: float64(i)})
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("failed to close stream: %v", err)
	}
	fmt.Println(res.Number)
}

func doFindMaximum(c calculatorpb.CalculatorClient) {
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("failed to open stream: %v", err)
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	// send
	go func() {
		defer wg.Done()
		for _, num := range []int32{1, 5, 3, 6, 2, 20} {
			stream.Send(&calculatorpb.MaximumRequest{Number: num})
		}
		stream.CloseSend()
		log.Println("send done")
	}()

	// recv
	go func() {
		defer wg.Done()
		for {
			res, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					return
				}
				log.Fatalf("failed to recv: %v", err)
				return
			}
			fmt.Printf("max: %d\n", res.Number)
		}
	}()

	wg.Wait()
}

func doErrorUnary(c calculatorpb.CalculatorClient) {
	res, err := c.SquareRoot(
		context.Background(),
		&calculatorpb.SquareRootRequest{Number: -1},
	)
	if err != nil {
		if respErr, ok := status.FromError(err); ok {
			fmt.Println(respErr.Message())
			fmt.Println(respErr.Code())
		} else {
			log.Fatalf("unknown error %v", err)
		}
	} else {
		fmt.Println(res.NumberRoot)
	}
}
