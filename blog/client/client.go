package main

import (
	"context"
	"fmt"
	"grpc-udemy/blog/blogpb"
	"io"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	conn, err := grpc.Dial(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to create dial: %v", err)
	}
	defer conn.Close()

	c := blogpb.NewBlogServiceClient(conn)

	blog := createBlog(c, "blog0")

	readBlog(c, blog.Id)

	blog.Title = fmt.Sprintf("updated - %s", blog.Title)
	updateBlog(c, blog)

	deleteBlog(c, blog.Id)

	createBlog(c, "blog1")
	createBlog(c, "blog2")
	createBlog(c, "blog3")
	listBlogs(c)
}

func createBlog(c blogpb.BlogServiceClient, title string) *blogpb.Blog {
	fmt.Println("Create blog")

	blog := blogpb.Blog{
		AuthorId: "John",
		Title:    title,
		Content:  "Content of John's blog.",
	}
	res, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: &blog})

	if err != nil {
		log.Fatalf("failed to create blog: %v", err)
		return nil
	}

	fmt.Println(res.Blog)
	return res.Blog
}

func readBlog(c blogpb.BlogServiceClient, id string) {
	fmt.Printf("Read blog: %v\n", id)
	res, err := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{Id: id})

	if err != nil {
		log.Fatalf("failed to read blog: %v", err)
		return
	}

	fmt.Println(res.Blog)
}

func updateBlog(c blogpb.BlogServiceClient, update *blogpb.Blog) {
	fmt.Printf("Update blog: %v\n", update.Id)

	res, err := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: update})
	if err != nil {
		log.Fatalf("failed to update blog: %v", err)
		return
	}

	fmt.Println(res.Blog)
}

func deleteBlog(c blogpb.BlogServiceClient, id string) {
	fmt.Printf("Delete blog: %v\n", id)

	res, err := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{Id: id})
	if err != nil {
		log.Fatalf("failed to delete blog: %v", err)
		return
	}

	fmt.Print(res.Id)
}

func listBlogs(c blogpb.BlogServiceClient) {
	fmt.Printf("list blogs\n")

	stream, err := c.ListBlog(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Fatalf("failed to open stream for logs: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		fmt.Println(res.Blog)
	}
}
