package bootstrap

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"notification_service/internal/api/notificationserviceapi"
	"notification_service/internal/consumer/personeventconsumer"
	"notification_service/internal/pb/notification_api"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func AppRun(api *notificationserviceapi.NotificationServiceAPI, consumer *personeventconsumer.PersonEventConsumer) {
	go consumer.Consume(context.Background())

	go func() {
		if err := runGRPCServer(api); err != nil {
			panic(fmt.Sprintf("gRPC server error: %v", err))
		}
	}()

	if err := runGatewayServer(); err != nil {
		panic(fmt.Sprintf("gateway server error: %v", err))
	}
}

func runGRPCServer(api *notificationserviceapi.NotificationServiceAPI) error {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	notification_api.RegisterNotificationServiceServer(s, api)

	log.Println("gRPC server listening on :50052")
	return s.Serve(lis)
}

func runGatewayServer() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := notification_api.RegisterNotificationServiceHandlerFromEndpoint(
		ctx, mux, "localhost:50052", opts,
	)
	if err != nil {
		return err
	}

	log.Println("gRPC-Gateway server listening on :8081")
	return http.ListenAndServe(":8081", mux)
}
