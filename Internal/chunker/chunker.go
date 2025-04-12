package chunker

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
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

func SizeChunk(fileSize int64) int64 {
	switch {
	case fileSize < 1024*1024*100:
		return 1024 * 1024
	case fileSize < 1024*1024*1024:
		return 1024 * 1024 * 4
	case fileSize < 1024*1024*1024*10:
		return 1024 * 1024 * 8
	default:
		return 1024 * 1024 * 16
	}
}
func ChunkFile(path string) (*FileMetaData, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	var chunks []ChunkMetaData
	buf := make([]byte, SizeChunk(stat.Size()))
	index := int64(0)
	offset := int64(0)

	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if n == 0 {
			break
		}

		h := sha256.Sum256(buf[:n])
		chunks = append(chunks, ChunkMetaData{
			Index:  index,
			Hash:   hex.EncodeToString(h[:]),
			Size:   int64(n),
			Offset: offset,
		})
		offset += int64(n)
		index++
	}
	return &FileMetaData{
		Name:       stat.Name(),
		Size:       stat.Size(),
		ChunkSize:  SizeChunk(stat.Size()),
		ChunkCount: int64(len(chunks)),
		Chunks:     chunks,
	}, nil

}
