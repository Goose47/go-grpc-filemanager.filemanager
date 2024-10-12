package suite

import (
	"context"
	"filemanager/internal/config"
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

	gRPCAddress := net.JoinHostPort(cfg.GRPC.Host, strconv.Itoa(cfg.GRPC.Port))
	cc, err := grpc.NewClient(gRPCAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to connect to grpc server: %v", err)
	}

	client := gen.NewFileManagerClient(cc)

	return ctx, &Suite{
		T:      t,
		Cfg:    cfg,
		Client: client,
	}
}
