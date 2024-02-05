package main

import (
	"context"
	"fmt"
	"homework-8/internal/app/comments"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"

	"homework-8/internal/app/posts"
	config2 "homework-8/internal/config"
	"homework-8/internal/pkg/logger"
	"homework-8/internal/pkg/repository/postgresql"
	"homework-8/internal/utils"
	"log"
	"net"
	"os"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	pb "homework-8/pkg/posts"
)

func main() {
	ctx := context.Background()

	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("error creating zap logger: %w", err)
	}
	logger.SetGlobal(zapLogger)

	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            false,
			BufferFlushInterval: 1 * time.Second,
		},
		ServiceName: "posts-service",
	}
	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		logger.Infof(ctx, "cannot create tracer: %v\n", err)
		os.Exit(1)
	}
	defer func(closer io.Closer) {
		err2 := closer.Close()
		if err2 != nil {
			logger.Infof(ctx, "cannot close tracer: %+v\n", err)
		}
	}(closer)

	opentracing.SetGlobalTracer(tracer)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database := utils.ConnectDB(ctx)
	defer database.GetConnectionsPool(ctx).Close()

	repo := postgresql.NewRepo(database)
	postsApp := posts.New(repo)
	commentsApp := comments.New(repo)

	server := grpc.NewServer()
	pb.RegisterPostServiceServer(server, postsApp)
	pb.RegisterCommentServiceServer(server, commentsApp)

	port := config2.GetConfigs().GetPort()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("error after call to Listen: %w", err)
	}

	logger.Infof(ctx, "service messages listening on %q", port)

	return server.Serve(lis)

}
