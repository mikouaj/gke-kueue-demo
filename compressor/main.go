// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"cloud.google.com/go/storage"
	"github.com/mikouaj/kueue-demo/cloud-storage-compressor/compress"
	"github.com/mikouaj/kueue-demo/cloud-storage-compressor/log"
)

func main() {
	log.Infof("Cloud storage object compressor starting\n")
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("\nerror: %s\n", err)
			os.Exit(1)
		}
	}()

	config := compress.NewConfigFromEnv()
	if err := config.Valid(); err != nil {
		log.Fatalf("error: %s\n", err)
		os.Exit(1)
	}
	log.ConfigureLogger(config.LogLevel, config.LogJSON)

	ctx, cancel := context.WithCancel(context.Background())
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("error when creating cloud storage client: %s", err)
		os.Exit(1)
	}

	compressor := compress.NewCompressor(ctx, client, config.Bucket, config.Object, config.Folder, config.Delete)
	defer compressor.Close()
	done := make(chan bool, 1)
	go func() {
		if err := compressor.Compress(); err != nil {
			log.Errorf("failed to compress cloud storage object: %s", err)
		}
		done <- true
	}()

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		done <- true
	}()

	<-done
	cancel()
	log.Infof("Cloud storage object compressor finishing\n")
}
