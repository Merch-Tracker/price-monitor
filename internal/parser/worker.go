package parser

import (
	log "github.com/sirupsen/logrus"
	"parsing-service/internal/network"
)

const zeroPrice = 0 //for debug purposes

func ProcessData(in <-chan network.Merch, out chan<- network.MerchResp) {
	for job := range in {
		log.Debugf("Processing job: %s", job)

		log.WithFields(log.Fields{
			"link":          job.Link,
			"cookie values": job.CookieValues,
		}).Debug("Parsing | Request data")

		page, err := getPage(job.Link, job.CookieValues)
		if err != nil {
			log.WithField("err", err).Error("Parsing | Can't get page")
			sendPrice(out, job.MerchUuid, zeroPrice)
			continue
		}
		log.WithField("Page fetched", page != nil).Debug("Parsing | Step 1")

		var data []string

		symbols := makeSymbolsList(job.ParseSubstring)
		for _, symbol := range symbols {
			data = append(data, findData(page, job.ParseTag, symbol)...)
		}

		if len(data) == 0 {
			log.Debug("Parsing | No data")
			sendPrice(out, job.MerchUuid, zeroPrice)
			continue
		}
		log.WithField("find data", len(data)).Debug("Parsing | Step 2")

		var price uint32 = 0
		price = uint32(findMinimal(data, job.Separator))
		log.WithField("find minimal price", price).Debug("Parsing | Step 3")

		if price > 0 {
			sendPrice(out, job.MerchUuid, price)
		} else {
			log.Debug("Price is not > 0, sending zero price")
			sendPrice(out, job.MerchUuid, zeroPrice)
		}
		log.WithField("price", price).Debug("Parsing | Step 4 END")
	}
	log.Error("Data processing aborted, return zero price")
}

func sendPrice(out chan<- network.MerchResp, merchUuid string, price uint32) {
	out <- network.MerchResp{MerchUuid: merchUuid, Price: price}
}
