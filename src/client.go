package client

import (
	"flag"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "./public"
	"log"
)

var (
	cid = flag.String("cid", "", "The cid")
	serverAddr = flag.String("server_addr", "127.0.0.1:8999", "The server address in the format of host:port")
)

func main() {
	flag.Parse()
	conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewIMIClient(conn)
	response, err := client.Search(context.Background(), &pb.SearchRequest{*cid})
	if err != nil {
		log.Fatalf("request error %v", err)
	}
	log.Println(response)
}
