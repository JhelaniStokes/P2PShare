package chunker

import (
	"encoding/json"
	"fmt"
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

type ProvidedFiles struct {
	Files []FileMetaData
}

type DataNotFound struct {
	Reason string
}

func (e DataNotFound) Error() string {
	return fmt.Sprintf("%s Not Found", e.Reason)
}

func SaveMetaData(fileMeta *FileMetaData) error {
	var file *os.File
	var err error
	_, err = os.Stat("ProvidedFiles")
	if os.IsNotExist(err) {
		file, err = os.Create("ProvidedFiles")
		if err != nil {
			return err
		}
	} else {
		file, err = os.Open("ProvidedFiles")
		if err != nil {
			return err
		}
	}
	defer file.Close()

	var providedFiles ProvidedFiles
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&providedFiles)
	if err != nil {
		return err
	}

	providedFiles.Files = append(providedFiles.Files, *fileMeta)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(providedFiles)
	return err
}

func LoadMetaData(name string) (*FileMetaData, error) {
	var file *os.File
	var err error
	var providedFiles ProvidedFiles
	_, err = os.Stat("ProvidedFiles")
	if os.IsNotExist(err) {
		return nil, DataNotFound{Reason: name}
	} else {
		file, err = os.Open("ProvidedFiles")
		if err != nil {
			return nil, err
		}
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&providedFiles)
	if err != nil {
		return nil, err
	}

	if providedFiles.Files == nil {
		return nil, DataNotFound{Reason: name}
	}

	for _, fileMeta := range providedFiles.Files {
		if fileMeta.Name == name {
			return &fileMeta, nil
		}
	}

	return nil, DataNotFound{Reason: name}
}
