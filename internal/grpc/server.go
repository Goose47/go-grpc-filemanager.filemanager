package filemanagergrpc

import (
	"context"
	"errors"
	"filemanager/internal/models"
	"filemanager/internal/storage"
	gen "github.com/Goose47/go-grpc-filemanager.protos/gen/go/filemanager"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
)

type serverAPI struct {
	gen.UnimplementedFileManagerServer
	storage Storage
}

type Storage interface {
	FileWriter(ctx context.Context) (closer io.WriteCloser, filename string, err error)
	ListFiles(ctx context.Context) ([]models.File, error)
	FileReader(ctx context.Context, filename string) (io.ReadCloser, error)
}

func Register(gRPCServer *grpc.Server, storage Storage) {
	gen.RegisterFileManagerServer(gRPCServer, &serverAPI{storage: storage})
}

func (s *serverAPI) Upload(
	stream gen.FileManager_UploadServer,
) error {
	writer, filename, err := s.storage.FileWriter(stream.Context())
	if err != nil {
		return status.Error(codes.Internal, "failed to create file")
	}
	defer writer.Close()

	for {
		req, err := stream.Recv()

		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return status.Error(codes.Internal, "error while reading")
		}

		_, err = writer.Write(req.GetChunk())
		if err != nil {
			return status.Error(codes.Internal, "error while writing")
		}
	}

	err = stream.SendAndClose(&gen.UploadResponse{
		Filename: filename,
	})
	if err != nil {
		return status.Error(codes.Internal, "failed to send response")
	}

	return nil
}

func (s *serverAPI) List(
	ctx context.Context,
	in *gen.ListRequest,
) (*gen.ListResponse, error) {
	files, err := s.storage.ListFiles(ctx)

	if err != nil {
		return nil, status.Error(codes.Internal, "failed to retrieve files")
	}

	filesInfo := make([]*gen.FileInfo, len(files))
	for i, file := range files {
		filesInfo[i] = &gen.FileInfo{
			Filename:     file.Name,
			CreationDate: file.CreationDate,
			UpdateDate:   file.UpdateDate,
		}
	}

	return &gen.ListResponse{
		Files: filesInfo,
	}, nil
}

type streamWriter struct {
	stream gen.FileManager_FileServer
}

func (w streamWriter) Write(p []byte) (int, error) {
	err := w.stream.Send(&gen.FileResponse{
		Chunk: p,
	})
	if err != nil {
		return 0, errors.New("error while sending file chunk")
	}
	return len(p), nil
}

func (s *serverAPI) File(
	in *gen.FileRequest,
	stream gen.FileManager_FileServer,
) error {
	if in.Filename == "" {
		return status.Error(codes.InvalidArgument, "filename is required")
	}

	reader, err := s.storage.FileReader(stream.Context(), in.Filename)
	if err != nil {
		if errors.Is(err, storage.ErrFileNotFound) {
			return status.Error(codes.InvalidArgument, err.Error())
		}

		return status.Error(codes.Internal, "failed to retrieve file")
	}
	defer reader.Close()

	writer := streamWriter{
		stream: stream,
	}

	_, err = io.Copy(writer, reader)
	if err != nil {
		return status.Error(codes.Internal, "error while writing file")
	}

	return nil
}
