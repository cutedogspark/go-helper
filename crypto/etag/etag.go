package etag

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"io"
	"os"
)

type Etag struct {
	BlockBits int64
	BlockSize int64
}

func New() *Etag {
	return &Etag{
		BlockBits: 22, // Indicate that the block size is 4M
		BlockSize: 1 << 22,
	}
}

func (e *Etag) BlockCount(fsize int64) int {
	return int((fsize + (e.BlockSize - 1)) >> e.BlockBits)
}

func CalSha1(b []byte, r io.Reader) ([]byte, error) {
	h := sha1.New()
	_, err := io.Copy(h, r)
	if err != nil {
		return nil, err
	}
	return h.Sum(b), nil
}

func (e *Etag) GetEtag(filename string) (etag string, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return
	}

	size := fi.Size()
	blockCnt := e.BlockCount(size)
	sha1Buf := make([]byte, 0, 21)

	if blockCnt <= 1 {
		sha1Buf = append(sha1Buf, 0x16)
		sha1Buf, err = CalSha1(sha1Buf, f)
		if err != nil {
			return
		}
	} else {
		sha1Buf = append(sha1Buf, 0x96)
		sha1BlockBuf := make([]byte, 0, blockCnt*20)
		for i := 0; i < blockCnt; i++ {
			body := io.LimitReader(f, e.BlockSize)
			sha1BlockBuf, err = CalSha1(sha1BlockBuf, body)
			if err != nil {
				return
			}
		}
		sha1Buf, _ = CalSha1(sha1Buf, bytes.NewReader(sha1BlockBuf))
	}
	etag = base64.URLEncoding.EncodeToString(sha1Buf)
	return
}
