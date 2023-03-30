package main

import (
	"context"
	"net/http"
	"sync"

	"github.com/vibeitco/apikey-service/model"
	"github.com/vibeitco/apikey-service/service"

	"github.com/vibeitco/go-utils/common"
	"github.com/vibeitco/go-utils/config"
	"github.com/vibeitco/go-utils/locker"
	"github.com/vibeitco/go-utils/log"
	"github.com/vibeitco/go-utils/server"
	"github.com/vibeitco/go-utils/storage/mongo"
)

func main() {
	ctx := context.Background()
	// read config
	conf := &service.Config{}
	err := config.Populate(conf)
	if err != nil {
		log.Fatal(ctx, err, nil, "failed populating config")
	}
	log.Info(ctx, log.Data{"conf": conf}, "config")

	// dao
	storageCfg := conf.MongoDB.ToStorageConfig()
	mongoCfg := mongo.Config{
		//Indexes: mongo.ToIndexes(service.GetIndexFields()),
	}
	dao, err := mongo.NewHandler(storageCfg, mongoCfg)
	if err != nil {
		log.Fatal(ctx, err, nil, "failed creating dao")
	}
	// locker
	locker := locker.NewMemLocker()

	// handler
	handler, err := service.NewHandler(conf.Core, dao, locker)
	if err != nil {
		log.Fatal(ctx, err, nil, "failed creating handler")
	}
	// create server and listener
	srv, lis, err := server.NewGRPC(&conf.Core)
	if err != nil {
		log.Fatal(ctx, err, nil, "failed creating REST server")
	}
	model.RegisterApiKeyServiceServer(srv, handler)

	// gateway
	mux, host, err := server.NewGRPCGatewayV2(ctx, &conf.Core, model.RegisterApiKeyServiceHandlerFromEndpoint)
	if err != nil {
		log.Fatal(ctx, err, nil, "failed creating grpc-gateway server")
	}

	var wg sync.WaitGroup
	// serve GRPC
	wg.Add(2)
	go server.Run(ctx, "grpc",
		func() {
			err = srv.Serve(lis)
			if err != nil {
				log.Fatal(ctx, err, nil, common.EventServerFatal)
			}
		},
		func() {
			defer wg.Done()
			srv.GracefulStop()
		})
	// serve Gateway
	go server.Run(ctx, "grpc-gateway",
		func() {
			err = http.ListenAndServe(host, mux)
			if err != nil {
				log.Fatal(ctx, err, nil, common.EventServerFatal)
			}
		},
		func() {
			defer wg.Done()
			srv.GracefulStop()
		})
	wg.Wait()
}
