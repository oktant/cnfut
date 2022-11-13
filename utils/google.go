package utils

import (
	"cloud.google.com/go/storage"
	"context"
	"sync"
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
