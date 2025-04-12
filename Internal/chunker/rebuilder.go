package chunker

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
)

func Rebuild(meta *FileMetaData, sourcePath string, outputPath string) error {
	if meta == nil {
		return errors.New("FileMeta is nil")
	}
	file, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer file.Close()
	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	for i := int64(0); i < meta.ChunkCount; i++ {
		offset := meta.Chunks[i].Offset
		buf := make([]byte, meta.Chunks[i].Size)
		_, err := file.ReadAt(buf, offset)
		if err != nil {
			return err
		}

		h := sha256.Sum256(buf[:])
		if hex.EncodeToString(h[:]) != meta.Chunks[i].Hash {
			return errors.New("hash mismatch")
		}

		_, err = out.Write(buf)

		if err != nil {
			return err
		}

		fmt.Printf("Wrote chunk %d\n", i)
	}
	return nil
}
