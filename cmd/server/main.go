package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VibuRoshin25/vibrox-arena/internal/game"
	arenapb "github.com/VibuRoshin25/vibrox-arena/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type arenaServer struct {
	arenapb.UnimplementedArenaServer
}

func (arenaServer) PlayMove(_ context.Context, request *arenapb.PlayMoveRequest) (*arenapb.PlayMoveResponse, error) {
	decision, err := game.PlayMove(request.GetBoard(), request.GetPosition(), request.GetHumanMark())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &arenapb.PlayMoveResponse{
		Board:              decision.Board,
		BotPosition:        decision.BotPosition,
		Outcome:            decision.Outcome,
		Strategy:           decision.Strategy,
		Score:              decision.Score,
		NodesEvaluated:     decision.NodesEvaluated,
		SearchDepth:        decision.SearchDepth,
		DecisionTimeMicros: decision.DecisionTimeMicros,
	}, nil
}

func main() {
	grpcAddress := envOrDefault("GRPC_LISTEN_ADDR", ":8100")
	healthAddress := envOrDefault("HEALTH_LISTEN_ADDR", ":8054")

	listener, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		log.Fatal("listen for gRPC: ", err)
	}
	grpcServer := grpc.NewServer()
	arenapb.RegisterArenaServer(grpcServer, arenaServer{})

	healthMux := http.NewServeMux()
	healthMux.HandleFunc("/healthz", func(response http.ResponseWriter, _ *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(response).Encode(map[string]string{"service": "vibrox-arena", "status": "healthy"})
	})
	healthServer := &http.Server{Addr: healthAddress, Handler: healthMux, ReadHeaderTimeout: 3 * time.Second}

	go func() {
		log.Printf("vibrox-arena gRPC listening on %s", grpcAddress)
		if err := grpcServer.Serve(listener); err != nil {
			log.Printf("gRPC server stopped: %v", err)
		}
	}()
	go func() {
		log.Printf("vibrox-arena health listening on %s", healthAddress)
		if err := healthServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("health server stopped: %v", err)
		}
	}()

	shutdown, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-shutdown.Done()
	grpcServer.GracefulStop()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_ = healthServer.Shutdown(ctx)
}

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
