package chunker

import (
	"encoding/json"
	"os"
)

type ChunkMetaData struct {
	Index  int64
	Size   int64
	Hash   string
	Offset int64
}
type FileMetaData struct {
	Name       string
	Size       int64
	ChunkSize  int64
	ChunkCount int64
	Chunks     []ChunkMetaData
}

func SaveMetaData(fileMeta *FileMetaData, outPath string) error {
	file, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(fileMeta)
	return err
}

func LoadMetaData(path string) (*FileMetaData, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var fileMeta FileMetaData
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&fileMeta)
	return &fileMeta, err
}
