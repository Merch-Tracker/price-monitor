package parser

import (
	log "github.com/sirupsen/logrus"
	"parsing-service/internal/network"
)

func ProcessData(in <-chan network.Merch, out chan<- network.MerchResp) {
	for job := range in {
		log.Debugf("Processing job: %s", job)

		log.WithFields(log.Fields{
			"link":          job.Link,
			"cookie values": job.CookieValues,
		}).Debug("Parsing | Request data")

		page, err := getPage(job.Link, job.CookieValues)
		if err != nil {
			log.WithField("err", err).Debug("Parsing | Can't get page")
			break
		}
		log.WithField("Page fetched", page != nil).Debug("Parsing | Step 1")

		var data []string

		symbols := makeSymbolsList(job.ParseSubstring)
		for _, symbol := range symbols {
			data = append(data, findData(page, job.ParseTag, symbol)...)
		}
		log.WithField("find data", len(data)).Debug("Parsing | Step 2")

		var price uint32 = 0

		if len(data) > 0 {
			price = uint32(findMinimal(data, job.Separator))
		}
		log.WithField("find minimal price", price).Debug("Parsing | Step 3")

		if price > 0 {
			out <- network.MerchResp{MerchUuid: job.MerchUuid, Price: price}
		}
		log.WithField("price", price).Debug("Parsing | Step 4 END")
	}
	log.Error("Data processing aborted")
}
