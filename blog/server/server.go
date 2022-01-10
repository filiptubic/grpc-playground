package main

import (
	"context"
	"fmt"
	"grpc-udemy/blog/blogpb"
	"log"
	"net"
	"os"
	"os/signal"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	collection *mongo.Collection
)

type server struct {
	blogpb.UnimplementedBlogServiceServer
}

func (s *server) loadBlogItem(id primitive.ObjectID) (*blogItem, error) {
	var blog blogItem

	if err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&blog); err != nil {
		log.Printf("failed read mongo object id: %v\n", err)
		return nil, err
	}

	return &blog, nil
}

func (s *server) ReadBlog(ctx context.Context, req *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {
	oId, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		log.Printf("failed to parse object id: %v\n", req.Id)
		return nil, status.Errorf(codes.Internal, "failed to parse object id %v", req.Id)
	}

	blog, err := s.loadBlogItem(oId)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &blogpb.ReadBlogResponse{Blog: blog.toPb()}, nil
}

func (s *server) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	blogReq := req.GetBlog()
	data := blogItem{
		AuthorId: blogReq.AuthorId,
		Title:    blogReq.Title,
		Content:  blogReq.Content,
	}

	res, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		log.Printf("failed to insert into mongo: %v\n", err)
		return nil, status.Errorf(codes.Internal, "internal error %v", err)
	}

	oId, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Printf("failed parse mongo object id: %v\n", res.InsertedID)
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	blogRes, err := s.loadBlogItem(oId)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &blogpb.CreateBlogResponse{
		Blog: blogRes.toPb(),
	}, nil
}

func (s *server) UpdateBlog(ctx context.Context, req *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error) {
	oId, err := primitive.ObjectIDFromHex(req.Blog.Id)
	if err != nil {
		log.Printf("failed to parse object id: %v\n", req.Blog.Id)
		return nil, status.Errorf(codes.Internal, "failed to parse object id %v", req.Blog.Id)
	}
	_, err = collection.ReplaceOne(context.Background(), bson.M{"_id": oId}, req.Blog)
	if err != nil {
		log.Printf("failed to update: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to update object id %v", req.Blog.Id)
	}

	blog, err := s.loadBlogItem(oId)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &blogpb.UpdateBlogResponse{Blog: blog.toPb()}, nil
}

func (s *server) DeleteBlog(ctx context.Context, req *blogpb.DeleteBlogRequest) (*blogpb.DeleteBlogResponse, error) {
	oId, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		log.Printf("failed to parse object id: %v\n", req.Id)
		return nil, status.Errorf(codes.Internal, "failed to parse object id %v", req.Id)
	}

	if _, err = collection.DeleteOne(context.Background(), &bson.M{"_id": oId}); err != nil {
		log.Printf("failed to deletew: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to delete object id %v", req.Id)
	}

	return &blogpb.DeleteBlogResponse{Id: oId.Hex()}, nil
}

func (s *server) ListBlog(in *emptypb.Empty, stream blogpb.BlogService_ListBlogServer) error {
	cursor, err := collection.Find(context.Background(), &bson.M{})
	defer func() {
		err := cursor.Close(context.Background())
		if err != nil {
			log.Printf("failed to close cursor: %v", err)
		}
	}()

	if err != nil {
		log.Fatalf("failed to list blogs: %v", err)
		return status.Error(codes.Internal, "failed list blogs")
	}

	queue := make([]*blogpb.Blog, 0)
	for cursor.Next(context.Background()) {
		var b blogItem
		if err := cursor.Decode(&b); err != nil {
			log.Printf("failed to decode blog: %v", err)
			return status.Error(codes.Internal, "failed to list blogs")
		}
		queue = append(queue, b.toPb())

		if len(queue) >= 2 {
			stream.Send(&blogpb.ListBlogResponse{
				Blog: queue,
			})
			queue = make([]*blogpb.Blog, 0)
		}
	}

	if len(queue) > 0 {
		stream.Send(&blogpb.ListBlogResponse{
			Blog: queue,
		})
	}

	return nil
}

type blogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorId string             `bson:"author_id"`
	Content  string             `bson:"content"`
	Title    string             `bson:"title"`
}

func (b *blogItem) toPb() *blogpb.Blog {
	return &blogpb.Blog{
		Id:       b.ID.Hex(),
		AuthorId: b.AuthorId,
		Title:    b.Title,
		Content:  b.Content,
	}
}

func connectToMongo() *mongo.Client {
	log.Println("Connecting to mongodb.")
	credentials := options.Credential{
		Username: "admin",
		Password: "admin",
	}
	client, err := mongo.NewClient(
		options.Client().ApplyURI("mongodb://localhost:27017"),
		options.Client().SetAuth(credentials),
	)
	if err != nil {
		log.Fatalf("Failed to create mongo client: %v", err)
	}
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatalf("Failed to connect to mongodb: %v", err)
	}
	collection = client.Database("mydb").Collection("blog")
	return client
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	mongoClient := connectToMongo()

	log.Println("Listening on port :50051.")
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}

	s := grpc.NewServer()
	blogpb.RegisterBlogServiceServer(s, &server{})

	reflection.Register(s)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to server: %v", err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	fmt.Println("Stoping server.")
	s.Stop()
	fmt.Println("Closing listener.")
	lis.Close()
	fmt.Println("Closing mongodb connection.")
	mongoClient.Disconnect(context.TODO())
	fmt.Println("Done.")
}
