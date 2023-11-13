package db

import (
	"encoding/json"
	"log"

	file_meta "github.com/thegoodparticle/music-share-system/file-server/file-meta"
	"golang.org/x/exp/slices"
)

type FileMetaData struct {
	ClientIPAddress  string   `json:"client_ip"`
	ClientPortNumber int32    `json:"client_port"`
	FilesStored      []string `json:"file_names"`
}

// update the peer IPs in this list with the node IPs
const data = `
[
    {
        "client_ip": "127.0.0.1",
        "client_port": 9010,
        "file_names": [
            "payphone.mp3",
            "memories.mp3"
        ]
    },
    {
        "client_ip": "127.0.0.1",
        "client_port": 9010,
        "file_names": [
            "bones.mp3",
            "whatever_it_takes.mp3"
        ]
    },
    {
        "client_ip": "127.0.0.1",
        "client_port": 9010,
        "file_names": [
            "kala_chashma.mp3",
            "ilahi.mp3"
        ]
    }
]
`

type DBStore struct {
	allFileMetaInfo []FileMetaData
}

func New() *DBStore {
	var allFileInfo []FileMetaData

	if err := json.Unmarshal([]byte(data), &allFileInfo); err != nil {
		log.Print("error marshalling files info")
		return nil
	}

	return &DBStore{allFileMetaInfo: allFileInfo}
}

func (db *DBStore) GetSpecificFileMetaData(fileName string) *file_meta.FileMetaResponse {
	var response file_meta.FileMetaResponse

	allFileInfo := db.allFileMetaInfo
	for _, client := range allFileInfo {
		if slices.Contains(client.FilesStored, fileName) {
			response.ClientIP = client.ClientIPAddress
			response.PortNumber = client.ClientPortNumber
		}
	}

	return &response
}
