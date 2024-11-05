package main

import (
	"context"
	"fmt"

	pb "github.com/iliamikado/UrlShortener/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// устанавливаем соединение с сервером
	conn, _ := grpc.Dial(":8088", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	// получаем переменную интерфейсного типа UsersClient,
	// через которую будем отправлять сообщения
	c := pb.NewUrlShortenerClient(conn)

	resp, _ := c.AddURL(context.Background(), &pb.LongURL{
		Url: "https://ya.ru",
	})
	fmt.Println(resp.Url)
	resp2, _ := c.GetLongURL(context.Background(), &pb.ShortURL{
		Url: resp.Url,
	})
	fmt.Println(resp2.Url)

	reqItems := make([]*pb.PostManyURLRequest_RequestBatchItem, 2)
	reqItems[0] = &pb.PostManyURLRequest_RequestBatchItem{
		CorrelationId: "1",
		OriginalUrl:   "https://google.com",
	}
	reqItems[1] = &pb.PostManyURLRequest_RequestBatchItem{
		CorrelationId: "2",
		OriginalUrl:   "https://google2.com",
	}
	resp3, _ := c.PostManyURL(context.Background(), &pb.PostManyURLRequest{
		ReqItems: reqItems,
	})
	for _, item := range resp3.ResItems {
		ans, _ := c.GetLongURL(context.Background(), &pb.ShortURL{
			Url: item.ShortUrl,
		})
		fmt.Println(ans.Url)
	}

}
