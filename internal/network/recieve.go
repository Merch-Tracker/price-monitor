package network

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"
	"io"
	pb "parsing-service/pkg/pricewatcher"
)

type Merch struct {
	MerchUuid      string
	Link           string
	ParseTag       string
	ParseSubstring string
	CookieValues   string
	Separator      string
}

func GetMerch(ctx context.Context, client pb.PriceWatcherClient) []Merch {
	var merchList []Merch
	stream, err := client.GetMerch(ctx, &empty.Empty{})
	if err != nil {
		log.Error("Error calling GetMerch: %v", err)
		return nil
	}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Error("Error receiving response: %v", err)
			return nil
		}

		newEntry := Merch{
			MerchUuid:      response.MerchUuid,
			Link:           response.Link,
			ParseTag:       response.ParseTag,
			ParseSubstring: response.ParseSubs,
			CookieValues:   response.CookieValues,
			Separator:      response.Separator,
		}

		merchList = append(merchList, newEntry)
		log.WithField("entry added", newEntry).Debug("gRPC Receive")
	}
	return merchList
}
