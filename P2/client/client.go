package main

import (
	"context"
	"log"
	"os"
	"strconv"

	file_data "github.com/thegoodparticle/music-share-system/file-server/file-data"

	file_meta "github.com/thegoodparticle/music-share-system/file-server/file-meta"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	serverIP   = "127.0.0.1"
	serverPort = 8010
)

func main() {
	// extract command line arguements
	args := os.Args
	if len(args) < 2 {
		log.Fatalf("not enough arguements passed\nInput format go run client.go <filename>.<extension>")
	}

	// first arguement should be file name
	fileName := args[1]

	// create grpc connection to the server
	serverConnection := connectThroughGRPC(serverIP, serverPort)
	// defer function closes the open connections on exiting the main method
	defer serverConnection.Close()

	// using protobufs defined, create new client
	client := file_meta.NewFileServerClient(serverConnection)

	// create file meta request using defined Protocol format
	req := file_meta.FileMetaRequest{
		FileName: fileName,
	}

	// call the method to get the file metadata
	resp, err := client.GetFileMetaInfo(context.Background(), &req)
	if err != nil || resp == nil {
		// if method responds with error, exit with error message
		log.Fatalf("error while getting meta info. Error - %+v\nExiting client now...", err)
	}

	log.Printf("response received - %+v", resp)

	// if in case, peer IP is not found for the requested file, exit with failure message
	if resp.ClientIP == "" {
		log.Fatalf("requested file meta not found in the server. Exiting client now...")
	}

	// once the peer IP and Port are found, connect to that peer through grpc
	peerConnection := connectThroughGRPC(resp.ClientIP, int(resp.PortNumber))
	// defer function closes the open connections on exiting the main method
	defer peerConnection.Close()

	// using the protobufs defined for file data in peer, create a new peer as a client
	peerClient := file_data.NewFileDataClient(peerConnection)

	fileReq := file_data.FileDataRequest{
		FileName: fileName,
	}

	// call the method to get the actual file data
	fileData, err := peerClient.GetFileData(context.Background(), &fileReq)
	if err != nil || fileData == nil {
		// exit, in case of error from peer
		log.Fatalf("error while getting file from peer. Error - %v", err)
	}

	log.Printf("Requested file '%s' has below content \n\n '%s'\n\n", fileName, string(fileData.FileData))
}

func connectThroughGRPC(serverIP string, serverPort int) *grpc.ClientConn {
	var opts []grpc.DialOption

	// add transport layer security as options; currently disabled as it is not a production server
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	portNumber := strconv.Itoa(serverPort)

	// using grpc function, call the requested server/peer on the port defined
	conn, err := grpc.Dial(serverIP+":"+portNumber, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	// return the connection object on successfull connection
	return conn
}
