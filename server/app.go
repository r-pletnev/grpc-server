package server

import (
	"context"
	"fmt"
	"github.com/globalsign/mgo"
	pb "github.com/zeroed88/grpc-server/grpcserver"
	"github.com/zeroed88/grpc-server/server/misc"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	port    = ":50051"
	dbUrl   = "mongodb://root:example@mongo:27017"
	dbName  = "csv"
	dbTable = "products"
)

var fetchClient = misc.NewDownloader()

type server struct {
	pb.UnimplementedProductServiceServer
	Session *mgo.Session
}

func (s *server) InsertProducts(products []misc.Product) error {
	c := s.Session.DB(dbName).C(dbTable)
	return misc.CreateOrUpdate(c, products)
}

func (s *server) GetProducts(params misc.FilterParameters) (products []misc.Product, err error) {
	c := s.Session.DB(dbName).C(dbTable)
	return misc.GetList(c, params)
}

// Fetch receive csv url file, parse and save it to storage
func (s *server) Fetch(ctx context.Context, in *pb.FetchRequest) (*pb.FetchResponse, error) {
	buf, err := fetchClient.DownloadFile(in.GetUrl())
	if err != nil {
		return &pb.FetchResponse{StatusCode: 400, Message: err.Error()}, nil
	}
	products, err := misc.ReadCSV(buf)
	if err != nil {
		return &pb.FetchResponse{StatusCode: 400, Message: err.Error()}, nil
	}
	fmt.Printf("products: %v\n", products)
	if err := s.InsertProducts(products); err != nil {
		return &pb.FetchResponse{StatusCode: 500, Message: err.Error()}, nil
	}
	return &pb.FetchResponse{StatusCode: 200, Message: "success"}, nil
}

func (s *server) List(ctx context.Context, in *pb.ListRequest) (*pb.ListResponse, error) {
	params := misc.FilterParameters{
		PageNumber: in.Paging.GetPageNumber(),
		PerPage:    in.Paging.GetPerPage(),
		Field: misc.FieldType(in.Sorting.GetField()),
		Order: misc.OrderType(in.Sorting.GetOrder()),
	}
	products, err := s.GetProducts(params)
	if err != nil {
		return &pb.ListResponse{}, err
	}
	result := make([]*pb.Product, 0)
	for _, elm := range products {
		p := &pb.Product{
			Name:              elm.Name,
			Price:             uint32(elm.Price),
			PriceUpdatedCount: uint32(elm.PriceUpdatedCount),
			UpdatedAt:         elm.UpdatedAt.Unix(),
		}
		result = append(result, p)
	}
	return &pb.ListResponse{Products: result}, nil
}

func Run() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	dbSession := misc.NewDBSession(dbUrl)
	defer dbSession.Close()

	s := grpc.NewServer()
	pb.RegisterProductServiceServer(s, &server{Session: dbSession})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

