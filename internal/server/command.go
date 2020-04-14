package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/mjpitz/grpc-on-kubernetes/api/v1"

	"github.com/spf13/cobra"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var (
	Command = &cobra.Command{
		Use:   "server",
		Short: "Start a server process",
		RunE: func(cmd *cobra.Command, args []string) error {
			hostname, err := os.Hostname()
			if err != nil {
				return err
			}

			server := grpc.NewServer()

			v1.RegisterDemoServer(server, &demo{
				hostname: hostname,
			})

			// For a more complete example on how to use the health server, see the example.
			// https://github.com/grpc/grpc-go/tree/master/examples/features/health
			healthCheck := health.NewServer()
			healthpb.RegisterHealthServer(server, healthCheck)

			go func() {
				// inspect dependencies and toggle service status appropriately
				// you can toggle the global system health, e.g. ""
				// or you could toggle a specific services health, e.g. "svc1"
				healthCheck.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
			}()

			listener, err := net.Listen("tcp", "0.0.0.0:8080")
			if err != nil {
				return err
			}

			log.Println(fmt.Sprint("[info] starting server on 0.0.0.0:8080"))
			if err := server.Serve(listener); err != nil {
				return err
			}

			return nil
		},
	}
)

type demo struct {
	hostname string
}

func (d *demo) Echo(ctx context.Context, req *v1.EchoRequest) (*v1.EchoResponse, error) {
	return &v1.EchoResponse{
		Hostname: d.hostname,
	}, nil
}

var _ v1.DemoServer = &demo{}
