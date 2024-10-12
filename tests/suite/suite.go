package suite

import (
	"context"
	"filemanager/internal/config"
	"fmt"
	gen "github.com/Goose47/go-grpc-filemanager.protos/gen/go/filemanager"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"strconv"
	"testing"
)

type Suite struct {
	*testing.T
	Cfg    *config.Config
	Client gen.FileManagerClient
}

const configPath = "../config/local_tests.yaml"

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadPath(configPath)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	client, err := createGRPCClient(cfg.GRPC.Host, cfg.GRPC.Port)
	if err != nil {
		t.Fatal(err)
	}

	return ctx, &Suite{
		T:      t,
		Cfg:    cfg,
		Client: client,
	}
}

func (s *Suite) NewClient() (gen.FileManagerClient, error) {
	return createGRPCClient(s.Cfg.GRPC.Host, s.Cfg.GRPC.Port)
}

func createGRPCClient(
	host string,
	port int,
) (gen.FileManagerClient, error) {
	gRPCAddress := net.JoinHostPort(host, strconv.Itoa(port))
	cc, err := grpc.NewClient(gRPCAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to grpc server: %v", err)
	}

	client := gen.NewFileManagerClient(cc)

	return client, nil
}
