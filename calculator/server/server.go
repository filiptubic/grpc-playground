package main

import (
	"context"
	"grpc-udemy/calculator/calculatorpb"
	"io"
	"log"
	"math"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type server struct {
	calculatorpb.UnimplementedCalculatorServer
}

func (s *server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	res := &calculatorpb.SumResponse{
		Sum: req.X + req.Y,
	}
	return res, nil
}

// PrimeNumberDecomposition algorithm pseudocode:
// k = 2
// N = 210
// while N > 1:
//     if N % k == 0:   // if k evenly divides into N
//         print k      // this is a factor
//         N = N / k    // divide N by k so that we have the rest of the number left.
//     else:
//         k = k + 1
func (s *server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.Calculator_PrimeNumberDecompositionServer) error {
	n := req.Number
	k := int32(2)

	for n > 1 {
		if n%k == 0 {
			stream.Send(&calculatorpb.PrimeNumberDecompositionResponse{
				Number: k,
			})
			n = n / k
		} else {
			k = k + 1
		}
	}
	return nil
}

func (s *server) ComputeAverage(stream calculatorpb.Calculator_ComputeAverageServer) error {
	sum := 0.0
	count := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&calculatorpb.AverageResponse{
				Number: sum / float64(count),
			})
		}
		sum += req.Number
		count++
	}
}

func (s *server) FindMaximum(stream calculatorpb.Calculator_FindMaximumServer) error {
	max := int32(math.MinInt32)

	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			log.Fatalf("failed to recv: %v", err)
			return err
		}
		if req.Number > max {
			max = req.Number
			time.Sleep(1000 * time.Millisecond)

			stream.Send(&calculatorpb.MaximumResponse{
				Number: max,
			})
		}
	}
}

func (s *server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	if req.Number < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "number %v is negative", req.Number)
	}

	return &calculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(req.Number)),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServer(s, &server{})

	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}
}
