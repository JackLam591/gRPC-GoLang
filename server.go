package main

import (
	"context"
	"fmt"
	"grpc-test/configs"
	"grpc-test/models"
	"grpc-test/pbservice/book"
	pb "grpc-test/pbservice/book"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type bookServiceServer struct {
	pb.UnimplementedBookServiceServer
}

var dbCollection *mongo.Collection = configs.GetCollection(configs.DB, "books")
var validate = validator.New()

func (s *bookServiceServer) SayHello(ctx context.Context, in *pb.Message) (*pb.Message, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.Message{Name: "Hello " + in.GetName()}, nil
}

func (s *bookServiceServer) CreateItem(ctx context.Context, in *pb.Book) (*pb.ID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	book := in
	defer cancel()

	newBook := models.Book{
		Id:       book.Id,
		Name:     book.Name,
		Category: book.Category,
	}

	result, err := dbCollection.InsertOne(ctx, newBook)
	if result != nil {
		fmt.Println("New Book created!")
	}
	if err != nil {
		return nil, status.Error(codes.Code(code.Code_ABORTED), "Insert new book fail")
	}

	return &pb.ID{Id: book.Id}, err
}

func (s *bookServiceServer) ReadItem(ctx context.Context, in *pb.ID) (*pb.Book, error) {
	if in.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "ID is empty, please try again")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var book pb.Book
	defer cancel()

	err := dbCollection.FindOne(ctx, bson.M{"id": in.Id}).Decode(&book)
	if err != nil {
		log.Printf("Error retrieving book with id: %s, error: %v", in.Id, err)
		return nil, status.Error(codes.NotFound, "Book not exist")
	}

	return &book, nil
}

// Get All User Interface
func (s *bookServiceServer) AllItem(ctx context.Context, in *emptypb.Empty) (*book.AllBook, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var books *pb.AllBook
	defer cancel()

	results, err := dbCollection.Find(ctx, bson.M{})
	if err != nil {
		panic(err.Error())
	}

	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleBook pb.Book
		var err = results.Decode(&singleBook)
		if err != nil {
			panic(err.Error())
		}

		// books.Books = new([]pb.Book)
		// *books.Books = append(*books.Books, Book{})
	}

	return books, nil
}

func main() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	book.RegisterBookServiceServer(grpcServer, &bookServiceServer{})
	reflection.Register(grpcServer)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %s", err)
		}
	}()

	fmt.Println("gRPC server started listening on port 9000........")

	c := make(chan os.Signal)

	signal.Notify(c, os.Interrupt)

	<-c
}
