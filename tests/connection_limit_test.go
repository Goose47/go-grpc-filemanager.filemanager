package tests

import (
	"filemanager/tests/suite"
	gen "github.com/Goose47/go-grpc-filemanager.protos/gen/go/filemanager"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConnectionLimit_MaxConcurrentClientStreamCalls(t *testing.T) {
	ctx, st := suite.New(t)

	done := make(chan bool, st.Cfg.GRPC.MaxStreamConnections)

	for range st.Cfg.GRPC.MaxStreamConnections {
		go func() {
			c, err := st.NewClient()
			require.NoError(t, err)

			stream, err := c.File(ctx, &gen.FileRequest{})
			require.NoError(t, err)

			<-done

			err = stream.CloseSend()
			require.NoError(t, err)
		}()
	}

	c, err := st.NewClient()
	require.NoError(t, err)
	require.NotEmpty(t, c)

	_, err = c.File(ctx, &gen.FileRequest{})
	require.Error(t, err)

	for range st.Cfg.GRPC.MaxStreamConnections {
		done <- true
	}
}
