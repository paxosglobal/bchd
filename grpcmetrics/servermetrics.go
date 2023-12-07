package grpcmetrics

import (
	"sync"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

var (
	serverSyncOnce          sync.Once
	serverMetrics           *grpcprom.ServerMetrics
	unaryServerInterceptor  grpc.UnaryServerInterceptor
	streamServerInterceptor grpc.StreamServerInterceptor
)

func RegisterPrometheusServerMetrics() (*grpcprom.ServerMetrics, grpc.UnaryServerInterceptor, grpc.StreamServerInterceptor) {
	serverSyncOnce.Do(func() {
		serverMetrics = grpcprom.NewServerMetrics(
			grpcprom.WithServerHandlingTimeHistogram(),
		)
		prometheus.MustRegister(serverMetrics)
		unaryServerInterceptor = serverMetrics.UnaryServerInterceptor()
		streamServerInterceptor = serverMetrics.StreamServerInterceptor()
	})
	return serverMetrics, unaryServerInterceptor, streamServerInterceptor
}
