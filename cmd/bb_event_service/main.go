package main

import (
	"log"
	"os"

	"github.com/buildbarn/bb-event-service/pkg/proto/configuration/bb_event_service"
	blobstore "github.com/buildbarn/bb-storage/pkg/blobstore/configuration"
	bb_grpc "github.com/buildbarn/bb-storage/pkg/grpc"
	"github.com/buildbarn/bb-storage/pkg/global"
	"github.com/buildbarn/bb-storage/pkg/util"
	build "google.golang.org/genproto/googleapis/devtools/build/v1"

	"google.golang.org/grpc"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: bb_event_service bb_event_service.jsonnet")
	}
	var configuration bb_event_service.ApplicationConfiguration
	if err := util.UnmarshalConfigurationFromFile(os.Args[1], &configuration); err != nil {
		log.Fatalf("Failed to read configuration from %s: %s", os.Args[1], err)
	}

	lifecycleState, err := global.ApplyConfiguration(configuration.Global)
	if err != nil {
		log.Fatal("Failed to apply global configuration options: ", err)
	}

	// Storage access.
	grpcClientFactory := bb_grpc.DefaultClientFactory
	contentAddressableStorage, actionCache, err := blobstore.NewCASAndACBlobAccessFromConfiguration(
		configuration.Blobstore,
		grpcClientFactory,
		int(configuration.MaximumMessageSizeBytes))

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		log.Fatal(
			"gRPC server failure: ",
			bb_grpc.NewServersFromConfigurationAndServe(
				configuration.GrpcServers,
				func(s *grpc.Server) {
					build.RegisterPublishBuildEventServer(s, &buildEventServer{
						instanceName:              "bb-event-service",
						contentAddressableStorage: contentAddressableStorage,
						actionCache:               actionCache,

						streams: map[string]*streamState{},
					})
				}))
	}()

	lifecycleState.MarkReadyAndWait()
}
