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
	SaveFile(ctx context.Context, fileData []byte) (bytes int, filename string, err error)
	ListFiles(ctx context.Context) ([]models.File, error)
	FileReader(ctx context.Context, filename string) (io.ReadCloser, error)
}

func Register(gRPCServer *grpc.Server, storage Storage) {
	gen.RegisterFileManagerServer(gRPCServer, &serverAPI{storage: storage})
}

func (s *serverAPI) Upload(
	stream gen.FileManager_UploadServer,
) error {
	var fileData []byte

	for {
		req, err := stream.Recv()

		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return status.Error(codes.Internal, "error while reading")
		}

		fileData = append(fileData, req.GetChunk()...)
	}

	// todo simultaneously read and write
	_, filename, err := s.storage.SaveFile(stream.Context(), fileData)
	if err != nil {
		return status.Error(codes.Internal, "failed to save file")
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
func (s *serverAPI) File(
	in *gen.FileRequest,
	stream gen.FileManager_FileServer,
) error {
	if in.Filename == "" {
		return status.Error(codes.InvalidArgument, "filename is required")
	}

	reader, err := s.storage.FileReader(stream.Context(), in.Filename)
	defer reader.Close()

	if err != nil {
		if errors.Is(err, storage.ErrFileNotFound) {
			return status.Error(codes.InvalidArgument, err.Error())
		}

		return status.Error(codes.Internal, "failed to retrieve file")
	}

	const chunkSize = 1024 // chunk size in bytes
	buf := make([]byte, chunkSize)

	for {
		n, err := reader.Read(buf)
		if err != nil && !errors.Is(err, io.EOF) {
			return status.Error(codes.Internal, "error while writing file")
		}
		if n == 0 {
			break
		}

		err = stream.Send(&gen.FileResponse{
			Chunk: buf[:n],
		})
		if err != nil {
			return status.Error(codes.Internal, "error while writing file")
		}
	}

	return nil
}
