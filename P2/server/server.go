package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/thegoodparticle/music-share-system/db"
	file_meta "github.com/thegoodparticle/music-share-system/file-server/file-meta"
	"google.golang.org/grpc"
)

type napsterFileServer struct {
	file_meta.UnimplementedFileServerServer
	dataStore *db.DBStore
}

func newServer() *napsterFileServer {
	return &napsterFileServer{
		dataStore: db.New(),
	}
}

func (s *napsterFileServer) GetFileMetaInfo(ctx context.Context, request *file_meta.FileMetaRequest) (*file_meta.FileMetaResponse, error) {
	log.Printf("received metadata info request for file %s", request.FileName)
	// data service defined separately to show the isolation of functionality and data
	// data store method returns the file metadata for the requested file name
	// if file metadata not found, returns an empty object
	return s.dataStore.GetSpecificFileMetaData(request.FileName), nil
}

const (
	serverPort = 8010
)

func main() {
	// announce the transport layer protocol and the address
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", serverPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	// create new grpc server with default options
	// no. of concurrent requests - 4294967295
	// default max size of buffer - 1024 * 1024 * 4
	grpcServer := grpc.NewServer(opts...)

	// register the file meta server with the server implementation described above
	file_meta.RegisterFileServerServer(grpcServer, newServer())

	log.Printf("starting grpc server at port %d...", serverPort)

	// start listening on the port
	grpcServer.Serve(lis)
}
