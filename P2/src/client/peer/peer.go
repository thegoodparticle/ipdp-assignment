package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	file_data "github.com/thegoodparticle/music-share-system/file-server/file-data"
	"google.golang.org/grpc"
)

type napsterMusicPeer struct {
	file_data.UnimplementedFileDataServer
}

func newPeer() *napsterMusicPeer {
	return &napsterMusicPeer{}
}

func (m *napsterMusicPeer) GetFileData(ctx context.Context, req *file_data.FileDataRequest) (*file_data.FileDataResponse, error) {
	dir, _ := os.Getwd()
	fileName := req.FileName

	// read file provided in the request
	dat, err := os.ReadFile(dir + "/client/peer/fdata/" + fileName)
	if err != nil {
		log.Printf("failed to find the file in the peer. Error: %v", err)
		return nil, err
	}

	log.Printf("responded content for file - %s", fileName)

	return &file_data.FileDataResponse{
		FileData: dat,
	}, nil
}

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Fatalf("not enough arguements passed\nInput format go run peer.go <port_number>")
	}

	portNumber, err := strconv.Atoi(args[1])
	if err != nil {
		log.Fatalf("arguement passed is not an integer")
	}

	// announce the transport layer protocol and the address
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", portNumber))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	// create new grpc server with default options
	// no. of concurrent requests - 4294967295
	// default max size of buffer - 1024 * 1024 * 4
	grpcPeer := grpc.NewServer(opts...)

	// register the file data server with the peer implementation
	napsterPeer := newPeer()
	file_data.RegisterFileDataServer(grpcPeer, napsterPeer)

	log.Printf("starting grpc peer at port %d...", portNumber)

	// start listening on the port
	grpcPeer.Serve(lis)
}
