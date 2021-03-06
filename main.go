package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"
	"time"

	pb "git.local/go-app/models"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

var sampleapp *Sampleapp

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	sampleapp = NewSampleapp()
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{Name: "daemon", Usage: "run server", Action: daemon},
		{Name: "grpc-hello", Usage: "run grpc client", Action: grpcHello},
		{Name: "grpc-order-create", Action: grpcOrderCreate},
		{Name: "grpc-order-read", Action: grpcOrderRead},
	}
	app.RunAndExitOnError()
}

func daemon(ctx *cli.Context) {
	go func() { http.ListenAndServe("localhost:9002", nil) }()
	sampleapp.BatchAsync()
	//go sampleapp.ServeGrpc()
	sampleapp.ServeHTTP()
}

func grpcHello(clictx *cli.Context) {
	// Set up a connection to the server.
	target := "localhost" + sampleapp.Cf.GrpcPort
	fmt.Println(target)
	conn, err := grpc.Dial(target, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	fmt.Println("2")
	defer conn.Close()
	c := pb.NewSampleAPIClient(conn)

	// Contact the server and print out its response.
	name := ""
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	fmt.Println("1")
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}

func grpcOrderCreate(clictx *cli.Context) {
	// Set up a connection to the server.
	target := "localhost" + sampleapp.Cf.GrpcPort
	fmt.Println(target)
	conn, err := grpc.Dial(target, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewSampleAPIClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	productUrl := strings.TrimSpace(clictx.Args().Get(0))
	orderStatus := strings.TrimSpace(clictx.Args().Get(1))
	r, err := c.CreateOrder(ctx, &pb.Order{ProductUrl: productUrl, Status: orderStatus})
	if err != nil {
		log.Fatalf("Err: %v", err)
	}
	if r != nil {
		log.Printf("Order: %s", r)
	}
}

func grpcOrderRead(clictx *cli.Context) {
	// Set up a connection to the server.
	target := "localhost" + sampleapp.Cf.GrpcPort
	conn, err := grpc.Dial(target, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewSampleAPIClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	orderId := strings.TrimSpace(clictx.Args().Get(0))
	r, err := c.ReadOrder(ctx, &pb.Id{Id: orderId})
	if err != nil {
		log.Fatalf("Err: %v", err)
	}
	if r != nil {
		log.Printf("Order: %s", r)
	}
}
