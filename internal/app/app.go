package app

import (
	"context"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"parsing-service/config"
	"parsing-service/internal/network"
	"parsing-service/internal/parser"
	pb "parsing-service/pkg/pricewatcher"
	"runtime"
	"time"
)

var LastCheck uint64

type server struct {
	pb.UnimplementedPriceWatcherServer
	CheckPeriod uint32
	NumCPUs     uint32
	StartTime   uint64
}

type App struct {
	ClientAddress string
	ServerAddress string
	NumCPUs       int
	CheckPeriod   time.Duration
	StartTime     time.Time
}

func New(c *config.Config) *App {
	numCPUs := c.NumCPUs
	if numCPUs < 1 {
		numCPUs = runtime.NumCPU()
	}

	return &App{
		ClientAddress: c.Host + ":" + c.ClientPort,
		ServerAddress: c.Host + ":" + c.ServerPort,
		NumCPUs:       numCPUs,
		CheckPeriod:   time.Duration(c.CheckPeriod),
		StartTime:     time.Now(),
	}
}

func (app *App) Run() {
	log.Info("Application start")
	log.WithFields(log.Fields{
		"ClientAddress":  app.ClientAddress,
		"Number of CPUs": app.NumCPUs,
	}).Debug("App settings")

	ctx := context.Background()
	var opts []grpc.DialOption
	insec := grpc.WithTransportCredentials(insecure.NewCredentials())
	opts = append(opts, insec)

	conn, err := grpc.NewClient(app.ClientAddress, opts...)
	if err != nil {
		log.Fatal(err)
	}

	client := pb.NewPriceWatcherClient(conn)

	s := grpc.NewServer()
	sv := &server{
		CheckPeriod: uint32(app.CheckPeriod),
		NumCPUs:     uint32(app.NumCPUs),
		StartTime:   uint64(app.StartTime.Unix()),
	}
	pb.RegisterPriceWatcherServer(s, sv)

	receiver := make(chan network.Merch, app.NumCPUs*10)
	sender := make(chan network.MerchResp, app.NumCPUs*10)

	//request data for parsing
	go func() {
		for {
			log.Info("Requesting data for parsing")
			receiveData := network.GetMerch(ctx, client)
			log.WithField("length", len(receiveData)).Debug("End receiving")
			if receiveData != nil {
				for _, element := range receiveData {
					receiver <- element
				}
			}
			LastCheck = uint64(time.Now().Unix())
			time.Sleep(time.Hour * app.CheckPeriod)
		}
	}()

	// send parsed data
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
					l := len(sendData)
					if l > 0 {
						log.WithField("length", l).Debug("Sending parsed data")
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

	//gRPC server for status response
	go func() {
		listener, err := net.Listen("tcp", app.ServerAddress)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		log.Infof("gRPC server listening at %v", app.ServerAddress)
		if err := s.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	select {}
}

func (s *server) ParserInfo(ctx context.Context, req *pb.StatusRequest) (*pb.StatusResponse, error) {
	resp := &pb.StatusResponse{
		CheckPeriod: s.CheckPeriod,
		NumCpus:     s.NumCPUs,
		StartTime:   s.StartTime,
		LastCheck:   LastCheck,
	}
	return resp, nil
}
