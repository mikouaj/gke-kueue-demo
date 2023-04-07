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

	"cloud.google.com/go/pubsub"
	"github.com/mikouaj/kueue-demo/event-job-creator/dispatch"
	"github.com/mikouaj/kueue-demo/event-job-creator/log"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("\nerror: %s\n", err)
			os.Exit(1)
		}
	}()

	config := dispatch.NewConfigFromEnv()
	if err := config.Valid(); err != nil {
		log.Fatalf("error: %s\n", err)
		os.Exit(1)
	}
	log.ConfigureLogger(config.LogLevel, config.LogJSON)

	log.Infof("Event based job dispatcher starting")

	ctx, cancel := context.WithCancel(context.Background())
	pubsubClient, err := pubsub.NewClient(ctx, config.ProjectID)
	if err != nil {
		log.Fatalf("failed to create pub/sub client : %s", err)
		os.Exit(1)
	}
	k8sClient, err := dispatch.NewKubeClient(config.KubeConfig)
	if err != nil {
		log.Fatalf("failed to create kubernetes client : %s", err)
		os.Exit(1)
	}

	creator := dispatch.NewCreator(ctx, pubsubClient, k8sClient, config.Topic, config.Folder,
		config.QueueName, config.CompressorImage, config.CompressorNamespace, config.CompressorSA,
		config.CompressorPrioClass, config.CompressorFolder)
	defer creator.Close()
	done := make(chan bool, 1)
	go func() {
		if err := creator.Start(); err != nil {
			log.Errorf("failed to get pub/sub messages: %s", err)
			done <- true
		}
	}()

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		done <- true
	}()

	<-done
	cancel()
	log.Infof("Event based job dispatcher exiting")
}
