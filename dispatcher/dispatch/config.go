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

package dispatch

import (
	"errors"
	"os"
)

const (
	projectIDVarName           = "GOOGLE_PROJECT_ID"
	topicVarName               = "DISPATCHER_PUBSUB_TOPIC"
	logLevelVarName            = "DISPATCHER_LOG_LEVEL"
	logJSONVarName             = "DISPATCHER_LOG_JSON"
	kubeConfigVarName          = "DISPATCHER_KUBE_CONFIG"
	folderVarName              = "DISPATCHER_FOLDER"
	queueNameVarName           = "COMPRESSOR_QUEUE_NAME"
	compressorImgVarName       = "COMPRESSOR_IMAGE"
	compressorNamespaceVarName = "COMPRESSOR_NAMESPACE"
	compressorSAVarName        = "COMPRESSOR_SA"
	compressorPrioClassVarName = "COMPRESSOR_PRIORITY_CLASS"
	compressorFolderVarName    = "COMPRESSOR_FOLDER"

	defaultFolderName = "inbox"
)

type Config struct {
	ProjectID           string
	Topic               string
	LogLevel            string
	LogJSON             string
	KubeConfig          string
	QueueName           string
	Folder              string
	CompressorImage     string
	CompressorNamespace string
	CompressorSA        string
	CompressorPrioClass string
	CompressorFolder    string
}

func (c *Config) Valid() error {
	if c.ProjectID == "" {
		return errors.New("project identifier is not set")
	}
	if c.Topic == "" {
		return errors.New("pub/sub topic is not set")
	}
	if c.QueueName == "" {
		return errors.New("queue name is not set")
	}
	if c.Folder == "" {
		return errors.New("folder is not set")
	}
	if c.CompressorImage == "" {
		return errors.New("compressor image is not set")
	}
	if c.CompressorNamespace == "" {
		return errors.New("compressor namespace is not set")
	}
	if c.CompressorSA == "" {
		return errors.New("compressor service account is not set")
	}
	return nil
}

func NewConfigFromEnv() *Config {
	config := &Config{}
	if val := os.Getenv(projectIDVarName); val != "" {
		config.ProjectID = val
	}
	if val := os.Getenv(topicVarName); val != "" {
		config.Topic = val
	}
	if val := os.Getenv(logLevelVarName); val != "" {
		config.LogLevel = val
	}
	if val := os.Getenv(logJSONVarName); val != "" {
		config.LogJSON = val
	}
	if val := os.Getenv(kubeConfigVarName); val != "" {
		config.KubeConfig = val
	}
	if val := os.Getenv(queueNameVarName); val != "" {
		config.QueueName = val
	}
	if val := os.Getenv(folderVarName); val != "" {
		config.Folder = val
	} else {
		config.Folder = defaultFolderName
	}
	if val := os.Getenv(compressorImgVarName); val != "" {
		config.CompressorImage = val
	}
	if val := os.Getenv(compressorNamespaceVarName); val != "" {
		config.CompressorNamespace = val
	}
	if val := os.Getenv(compressorSAVarName); val != "" {
		config.CompressorSA = val
	}
	if val := os.Getenv(compressorPrioClassVarName); val != "" {
		config.CompressorPrioClass = val
	}
	if val := os.Getenv(compressorFolderVarName); val != "" {
		config.CompressorFolder = val
	}
	return config
}
