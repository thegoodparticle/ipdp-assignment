## Prerequisites

### Install Go on all nodes

Follow steps mentioned here to install go on all the nodes
https://go.dev/doc/install

### Download source code

- Source code is uploaded at https://github.com/thegoodparticle/ipdp-assignment/tree/main/P2

- Use git command line to clone the repository on to the local nodes

```
git clone https://github.com/thegoodparticle/ipdp-assignment.git
```

- move to P2/ folder

```
cd P2/
```

## Commands to execute

1. Create protocol interfaces for File Metadata Server from protocol file

``` 
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative file-server/protobufs/file-meta.proto
```

2. Create protocol interfaces for File Data Peer from protocol file

``` 
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative file-server/protobufs/file-data.proto
```

3. Note down the node IPs and update them in the file present in db/store.go

```
vi db/store.go
```

4. Start the server

```
go run server/server.go
```

5. Start the peers in each node

```
go run client/peer/peer.go 9010

go run client/peer/peer.go 9011

go run client/peer/peer.go 9012
```

6. Start the client and send the request from client for a particular file

```
go run client/client.go memories.mp3
```