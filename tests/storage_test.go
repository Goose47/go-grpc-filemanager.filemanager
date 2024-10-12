package tests

import (
	"errors"
	"filemanager/tests/suite"
	gen "github.com/Goose47/go-grpc-filemanager.protos/gen/go/filemanager"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
	"time"
)

func TestStorage_UploadDownloadHappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	uploadStream, err := st.Client.Upload(ctx)
	require.NoError(t, err)

	img := gofakeit.ImagePng(100, 100)

	const NChunks = 5
	chunkSize := len(img) / NChunks

	for i := 0; i <= NChunks; i++ {
		l, r := i, i+chunkSize
		var nextChunk []byte

		if r > cap(img) {
			nextChunk = img[l:cap(img)]
		} else {
			nextChunk = img[l:r]
		}

		err := uploadStream.Send(&gen.UploadRequest{
			Chunk: nextChunk,
		})
		require.NoError(t, err)
	}

	uploadRes, err := uploadStream.CloseAndRecv()
	require.NoError(t, err)

	uploadTime := time.Now()

	filename := uploadRes.Filename
	assert.NotEmpty(t, filename)

	listRes, err := st.Client.List(ctx, &gen.ListRequest{})
	assert.NoError(t, err)

	var foundFile *gen.FileInfo
	for _, file := range listRes.Files {
		if file.Filename == filename {
			foundFile = file
		}
	}

	assert.NotEmpty(t, foundFile)

	const deltaSeconds = 1
	assert.InDelta(t, uploadTime.Unix(), foundFile.CreationDate, deltaSeconds)

	downloadStream, err := st.Client.File(ctx, &gen.FileRequest{
		Filename: filename,
	})
	assert.NoError(t, err)

	var receivedImg []byte

	for {
		chunk, err := downloadStream.Recv()

		if errors.Is(err, io.EOF) {
			break
		}
		assert.NoError(t, err)

		receivedImg = append(receivedImg, chunk.Chunk...)
	}

	assert.Equal(t, len(img), len(receivedImg))
	for i := range img {
		assert.Equal(t, img[i], receivedImg[i])
	}
}
