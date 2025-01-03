package network

import (
	"context"
	log "github.com/sirupsen/logrus"
	pb "parsing-service/pkg/pricewatcher"
)

type MerchResp struct {
	MerchUuid string
	Price     uint32
}

func PostMerch(client pb.PriceWatcherClient, merchList []MerchResp) {
	stream, err := client.PostMerch(context.Background())
	if err != nil {
		log.Fatalf("Error calling PostMerch: %v", err)
	}

	merchResponses := make([]pb.MerchResponse, len(merchList))
	for i, merch := range merchList {
		merchResponses[i] = pb.MerchResponse{
			MerchUuid: merch.MerchUuid,
			Price:     merch.Price,
		}
	}

	for i := range merchResponses {
		response := &merchResponses[i]
		if err = stream.Send(response); err != nil {
			log.Fatalf("Error sending request: %v", err)
		}
	}

	if err = stream.CloseSend(); err != nil {
		log.Fatalf("Error closing stream: %v", err)
	}
}
