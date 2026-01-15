package bootstrap

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	giftserviceapi "gift_manager/internal/api/gift_service_api"
	"gift_manager/internal/pb/gift_manager_api"
	"gift_manager/internal/producer/kafkaproducer"

	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func AppRun(api *giftserviceapi.GiftServiceAPI, producer *kafkaproducer.KafkaProducer) {
	go func() {
		if err := runGRPCServer(api); err != nil {
			log.Fatalf("failed to run gRPC server: %v", err)
		}
	}()

	if err := runGatewayServer(); err != nil {
		log.Fatalf("failed to run gateway server: %v", err)
	}
}

func runGRPCServer(api *giftserviceapi.GiftServiceAPI) error {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	gift_manager_api.RegisterGiftManagerServiceServer(s, api)

	log.Println("gRPC server listening on :50051")
	return s.Serve(lis)
}

func runGatewayServer() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	swaggerPath := os.Getenv("swaggerPath")
	if _, err := os.Stat(swaggerPath); os.IsNotExist(err) {
		log.Printf("Warning: swagger file not found: %s", swaggerPath)
		swaggerPath = ""
	}

	r := chi.NewRouter()

	if swaggerPath != "" {
		r.Get("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, swaggerPath)
		})

		r.Get("/docs/*", httpSwagger.Handler(
			httpSwagger.URL("/swagger.json"),
		))

		r.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/docs/", http.StatusFound)
		})

		log.Println("Swagger UI available at /docs/")
	}

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := gift_manager_api.RegisterGiftManagerServiceHandlerFromEndpoint(
		ctx, mux, "localhost:50051", opts,
	)
	if err != nil {
		return fmt.Errorf("failed to register gateway: %v", err)
	}

	r.Mount("/", mux)

	log.Println("gRPC-Gateway server listening on :8080")
	return http.ListenAndServe(":8080", r)
}
