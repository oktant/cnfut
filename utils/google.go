package utils

import (
	"bytes"
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	zlog "github.com/rs/zerolog/log"
	"io"
	"sync"
	"time"
)

var once sync.Once

var gClient *storage.Client

func GetInstance(ctx context.Context) (*storage.Client, error) {
	var err error
	if gClient == nil {
		once.Do(
			func() {
				gClient, err = storage.NewClient(ctx)
			})
	}

	return gClient, err
}

func UploadFileToGoogle(ctx context.Context, buf *bytes.Buffer, bucket, object string) error {
	client, err := GetInstance(ctx)
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	wc := client.Bucket(bucket).Object(object).NewWriter(ctx)
	wc.ChunkSize = 5

	if _, err = io.Copy(wc, buf); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	zlog.Info().Msg("Successfully copied objects")
	return nil

}
