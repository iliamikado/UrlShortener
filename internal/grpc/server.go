package grpc

import (
	"context"
	"net"
	"strings"

	"github.com/iliamikado/UrlShortener/internal/config"
	"github.com/iliamikado/UrlShortener/internal/storage"
	pb "github.com/iliamikado/UrlShortener/proto"
	"google.golang.org/grpc"
)

var urlStorage storage.URLStorage
var defaultUser = "default"

func RunGRPCServer(port string, st storage.URLStorage) error {
	urlStorage = st
	listen, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	pb.RegisterUrlShortenerServer(s, &UrlShortenerServer{})
	err = s.Serve(listen)
	return err
}

type UrlShortenerServer struct {
	pb.UnimplementedUrlShortenerServer
}

func (s *UrlShortenerServer) AddURL(ctx context.Context, in *pb.LongURL) (*pb.ShortURL, error) {

	id, err := urlStorage.AddURL(in.Url, defaultUser)
	if err != nil {
		return nil, err
	}
	shortURL := config.ResultAddress + "/" + id
	var resp pb.ShortURL
	resp.Url = shortURL

	return &resp, nil
}

func (s *UrlShortenerServer) GetLongURL(ctx context.Context, in *pb.ShortURL) (*pb.LongURL, error) {
	id := strings.TrimPrefix(in.Url, config.ResultAddress+"/")
	longURL, err := urlStorage.GetURL(id)
	if err != nil {
		return nil, err
	}

	var resp pb.LongURL
	resp.Url = longURL
	return &resp, nil
}

func (s *UrlShortenerServer) PostManyURL(ctx context.Context, in *pb.PostManyURLRequest) (*pb.PostManyURLResponse, error) {
	longURLs := make([]string, len(in.ReqItems))
	for i, item := range in.ReqItems {
		longURLs[i] = item.OriginalUrl
	}
	ids := urlStorage.AddManyURLs(longURLs, defaultUser)
	var resp pb.PostManyURLResponse
	resp.ResItems = make([]*pb.PostManyURLResponse_ResponseBatchItem, len(in.ReqItems))
	for i, id := range ids {
		resp.ResItems[i] = &pb.PostManyURLResponse_ResponseBatchItem{
			CorrelationId: in.ReqItems[i].CorrelationId,
			ShortUrl:      config.ResultAddress + "/" + id,
		}
	}
	return &resp, nil
}
