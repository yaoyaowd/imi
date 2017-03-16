package main

import (
	"log"
	"net"
	"google.golang.org/grpc/reflection"
	"golang.org/x/net/context"
	pb "./public"
	"google.golang.org/grpc"
	data "./data"
	core "./core"
	"sort"
	"flag"
)

var (
	indexFile = flag.String("index", "", "the index cid file")
	dictFile = flag.String("dict", "", "the dict cid file")
)

type ImiServiceImpl struct {
	index *core.Index
	dictionary *data.Dataset
}

type Docs []*pb.SearchDoc

func (c Docs) Len() int {
	return len(c)
}

func (c Docs) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c Docs) Less(i, j int) bool {
	return c[i].Score < c[j].Score
}

func (s *ImiServiceImpl) Search(ctx context.Context, in *pb.SearchRequest) (*pb.SearchResult, error) {
	log.Printf("search %s", in.Cid)
	out := make(chan int)
	defer close(out)

	id, ok := s.dictionary.Index[in.Cid]
	if !ok {
		return &pb.SearchResult{}, nil
	}
	p := s.dictionary.Points[id]

	go func() {
		s.index.Query(p, 2000, out)
	}()

	docs := []*pb.SearchDoc{}
	for {
		key := <- out
		if key == -1 {
			break
		}
		docs = append(docs, &pb.SearchDoc{s.index.D.Points[key].Id, p.L2(s.index.D.Points[key]), "", })
	}

	sort.Sort(Docs(docs))
	return &pb.SearchResult{docs, int32(len(docs)), }, nil
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", ":8999")
	if err != nil {
		log.Fatal("failed to listen: %v", err)
	}

	indexData, err := data.NewDatasetWithHeader(*indexFile, "cid", "col_")
	if err != nil {
		log.Fatalln("failed loading index")
	}
	index := core.NewIndex(indexData, 100, 100)
	log.Printf("Finish building index %d", index.D.Rows)

	dictionary, err := data.NewDatasetWithHeader(*dictFile, "cid", "col_")
	serviceImpl := &ImiServiceImpl{index, dictionary}

	s := grpc.NewServer()
	pb.RegisterIMIServer(s, serviceImpl)
	reflection.Register(s)

	log.Printf("start serving, port%s\n", ":8999")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
