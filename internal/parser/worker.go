package parser

import (
	"parsing-service/internal/network"
)

func ProcessData(in <-chan network.Merch, out chan<- network.MerchResp) {
	for job := range in {
		page, err := getPage(job.Link, job.CookieValues)
		if err != nil {
			return
		}

		results := findData(page, job.ParseTag, job.ParseSubstring)
		if results == nil {
			return
		}

		var price uint32 = 0

		if len(results) > 0 {
			price = uint32(findMinimal(results, job.Separator))
		}

		if price > 0 {
			out <- network.MerchResp{MerchUuid: job.MerchUuid, Price: price}
		}
	}
}
