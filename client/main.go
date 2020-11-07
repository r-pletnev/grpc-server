package main

import (
	"context"
	pb "github.com/zeroed88/grpc-server/grpcserver"
	"google.golang.org/grpc"
	"log"
	"os"
	"time"
)

const (
	address = "localhost:50051"
	defaultUrl = "http://ya.ru"
)


func main(){
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewProductServiceClient(conn)

	url := defaultUrl
	if len(os.Args) > 1 {
		url = os.Args[1]
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Fetch(ctx, &pb.FetchRequest{Url: url})
	if err != nil {
		log.Fatalf("could not fetch: %v", err)
	}
	if r.StatusCode != 200 {
		log.Printf("error: %s", r.Message)
	}
}