package main

import (
	pb "bookshop/server/bookshop/pb"
	"context"
	"log"
	"net"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedInventoryServer
}

var SampleBooks = []*pb.Book{
	{
		Title:     "The Hitchhiker's Guide to the Galaxy",
		Author:    "Douglas Adams",
		PageCount: 42,
	},
	{
		Title:     "The Lord of the Rings",
		Author:    "J.R.R. Tolkien",
		PageCount: 1234,
	},
}

func getSampleBooks() []*pb.Book {

	return SampleBooks
}
func addSampleBook(book *pb.Book) {
	SampleBooks = append(SampleBooks, book)
}

func (s *server) GetBookList(ctx context.Context, in *pb.GetBookListRequest) (*pb.GetBookListResponse, error) {
	log.Printf("Received request: %v", in.ProtoReflect().Descriptor().FullName())

	books := getSampleBooks()

	filteredSampleBooks := []*pb.Book{}
	if in.Title != nil {
		for _, book := range SampleBooks {
			if strings.Contains(strings.ToLower(book.Title), strings.ToLower(*in.Title)) {
				filteredSampleBooks = append(filteredSampleBooks, book)
			}
		}
	} else {
		filteredSampleBooks = books
	}

	return &pb.GetBookListResponse{
		Books: filteredSampleBooks,
	}, nil
}

func (s *server) AddNewBook(context context.Context, newBook *pb.Book) (*pb.Response, error) {
	resp := &pb.Response{
		Status: "1",
		ErrMsg: nil,
	}

	if newBook.Title == "" || newBook.Author == "" || newBook.PageCount == 0 {
		resp.Status = "0"
		*resp.ErrMsg = "Fill the title,author,page count"
	} else {
		addSampleBook(newBook)
	}

	return resp, nil
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()

	reflection.Register(s)

	pb.RegisterInventoryServer(s, &server{})

	log.Println("Server running...")

	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
