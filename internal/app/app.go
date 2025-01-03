package app

import (
	"context"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"parsing-service/config"
	"parsing-service/internal/network"
	"parsing-service/internal/parser"
	pb "parsing-service/pkg/pricewatcher"
	"runtime"
	"time"
)

type App struct {
	Address string
	NumCPUs int
}

func New(c *config.Config) *App {
	numCPUs := c.NumCPUs
	if numCPUs < 1 {
		numCPUs = runtime.NumCPU()
	}

	return &App{
		Address: c.Host + ":" + c.Port,
		NumCPUs: numCPUs,
	}
}

func (app *App) Run() {
	log.Info("Application start")
	log.WithFields(log.Fields{
		"Address":        app.Address,
		"Number of CPUs": app.NumCPUs,
	}).Debug("App settings")

	ctx := context.Background()
	var opts []grpc.DialOption
	insec := grpc.WithTransportCredentials(insecure.NewCredentials())
	opts = append(opts, insec)

	conn, err := grpc.NewClient(app.Address, opts...)
	if err != nil {
		log.Fatal(err)
	}

	client := pb.NewPriceWatcherClient(conn)

	receiver := make(chan network.Merch, app.NumCPUs*10)
	sender := make(chan network.MerchResp, app.NumCPUs*10)

	go func() {
		for {
			log.Info("Requesting data for parsing")
			receiveData := network.GetMerch(ctx, client)
			if receiveData != nil {
				for _, element := range receiveData {
					receiver <- element
				}
			}
			time.Sleep(time.Hour * 1)
		}
	}()

	go func() {
		ticker := time.NewTicker(time.Second * 2)
		defer ticker.Stop()

		for {
			var sendData []network.MerchResp

			for {
				select {
				case element := <-sender:
					sendData = append(sendData, element)
				case <-ticker.C:
					if len(sendData) > 0 {
						network.PostMerch(client, sendData)
						sendData = sendData[:0]
					}
					break
				}
			}
		}
	}()

	for i := 0; i < app.NumCPUs; i++ {
		go parser.ProcessData(receiver, sender)
	}

	select {}
}
