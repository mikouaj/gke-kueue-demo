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

// Package compress contains login to compress objects from cloud storage bucket
package compress

import (
	"compress/gzip"
	"context"
	"io"
	"regexp"

	"cloud.google.com/go/storage"
	"github.com/mikouaj/kueue-demo/cloud-storage-compressor/log"
)

type compressor struct {
	ctx          context.Context
	client       *storage.Client
	bucket       string
	object       string
	folder       string
	deleteObject bool
}

func NewCompressor(ctx context.Context, client *storage.Client, bucket, object, folder string, deleteObject bool) *compressor {
	return &compressor{
		ctx:          ctx,
		client:       client,
		bucket:       bucket,
		object:       object,
		folder:       folder,
		deleteObject: deleteObject,
	}
}

func (c *compressor) Close() error {
	if c.client != nil {
		log.Info("Closing Cloud Storage client")
		return c.client.Close()
	}
	return nil
}

func (c *compressor) Compress() error {
	newObjName := getNewObjectName(c.folder, c.object)
	log.Infof("Compressing object %s from bucket %s to %s", c.object, c.bucket, newObjName)
	bucket := c.client.Bucket(c.bucket)
	object := bucket.Object(c.object)
	objReader, err := object.NewReader(c.ctx)
	if err != nil {
		return err
	}
	defer objReader.Close()
	objData, err := io.ReadAll(objReader)
	if err != nil {
		return err
	}
	newObj := bucket.Object(newObjName)
	newObjWriter := newObj.NewWriter(c.ctx)
	defer newObjWriter.Close()
	gzWriter, err := gzip.NewWriterLevel(newObjWriter, gzip.BestCompression)
	if err != nil {
		return err
	}
	defer gzWriter.Close()
	if _, err := gzWriter.Write(objData); err != nil {
		return err
	}
	if c.deleteObject {
		log.Infof("Deleting compressed object %s", c.object)
		if err := object.Delete(c.ctx); err != nil {
			return err
		}
	}
	return nil
}

func getNewObjectName(folder, objName string) string {
	r := regexp.MustCompile(`^[^/]+/(.+)$`)
	if !r.MatchString(objName) {
		return folder + "/" + objName + ".gz"
	}
	matches := r.FindStringSubmatch(objName)
	return folder + "/" + matches[1] + ".gz"
}
