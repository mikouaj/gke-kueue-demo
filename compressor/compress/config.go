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

package compress

import (
	"errors"
	"os"
	"strings"
)

const (
	logLevelVarName = "COMPRESSOR_LOG_LEVEL"
	logJSONVarName  = "COMPRESSOR_LOG_JSON"
	bucketVarName   = "COMPRESSOR_BUCKET"
	objectVarName   = "COMPRESSOR_OBJECT"
	folderVarName   = "COMPRESSOR_FOLDER"
	deleteVarName   = "COMPRESSOR_DELETE_OBJECT"

	defaultFolder = "processed"
)

type Config struct {
	LogLevel string
	LogJSON  string
	Bucket   string
	Object   string
	Folder   string
	Delete   bool
}

func (c *Config) Valid() error {
	if c.Bucket == "" {
		return errors.New("cloud storage bucket is not set")
	}
	if c.Object == "" {
		return errors.New("cloud storage object is not set")
	}
	if c.Folder == "" {
		return errors.New("folder is not set")
	}
	return nil
}

func NewConfigFromEnv() *Config {
	config := &Config{}
	if val := os.Getenv(logLevelVarName); val != "" {
		config.LogLevel = val
	}
	if val := os.Getenv(logJSONVarName); val != "" {
		config.LogJSON = val
	}
	if val := os.Getenv(bucketVarName); val != "" {
		config.Bucket = val
	}
	if val := os.Getenv(objectVarName); val != "" {
		config.Object = val
	}
	if val := os.Getenv(folderVarName); val != "" {
		config.Folder = val
	} else {
		config.Folder = defaultFolder
	}
	if val := os.Getenv(deleteVarName); strings.ToLower(val) == "false" {
		config.Delete = false
	} else {
		config.Delete = true
	}
	return config
}
