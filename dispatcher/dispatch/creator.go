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

// Package dispatch contains pub/sub message based Kueue job creator
package dispatch

import (
	"context"
	"encoding/json"
	"strings"

	"cloud.google.com/go/pubsub"
	"github.com/mikouaj/kueue-demo/event-job-creator/log"
	"google.golang.org/api/iterator"
	storage "google.golang.org/api/storage/v1"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/applyconfigurations/core/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	subscriptionID = "kueue-dispatcher"
	saName         = "compressor"
)

type Creator struct {
	ctx                 context.Context
	pubsubClient        *pubsub.Client
	k8sClient           *kubernetes.Clientset
	topic               string
	folder              string
	queueName           string
	compressorImage     string
	compressorNamespace string
	compressorSA        string
}

func NewCreator(ctx context.Context, pubsubClient *pubsub.Client, k8sClient *kubernetes.Clientset,
	topic, folder, queueName, compressorImage, compressorNamespace, compressorSA string) *Creator {
	return &Creator{
		ctx:                 ctx,
		pubsubClient:        pubsubClient,
		k8sClient:           k8sClient,
		topic:               topic,
		queueName:           queueName,
		folder:              folder,
		compressorImage:     compressorImage,
		compressorNamespace: compressorNamespace,
		compressorSA:        compressorSA,
	}
}

func (c *Creator) Close() error {
	log.Infof("closing pub/sub client")
	return c.pubsubClient.Close()
}

func (c *Creator) Start() error {
	sub, err := c.getSubscription()
	if err != nil {
		return err
	}
	log.Info("Applying kubernetes service account for compressor job")
	if err := c.createServiceAccountForJob(); err != nil {
		return err
	}
	log.Infof("Receiving pub/sub messages")
	err = sub.Receive(c.ctx, func(ctx context.Context, m *pubsub.Message) {
		log.Debugf("got message with id=%s", m.ID)
		var obj storage.Object
		if err := json.Unmarshal(m.Data, &obj); err != nil {
			log.Errorf("failed to decode pub/sub message with id=%s :%s", m.ID, err)
			m.Nack()
		}
		if strings.HasPrefix(obj.Name, c.folder) && !strings.HasSuffix(obj.Name, "/") {
			c.createJobForObject(&obj)
		} else {
			log.Debugf("skipping object with name %s (is a folder or has missing %s prefix)", obj.Name, c.folder)
		}
		m.Ack()
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *Creator) getSubscription() (*pubsub.Subscription, error) {
	it := c.pubsubClient.Topic(c.topic).Subscriptions(c.ctx)
	for {
		sub, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		if sub.ID() == subscriptionID {
			log.Infof("found subscription with id=%s", subscriptionID)
			return sub, nil
		}
	}
	log.Infof("subscription with id=%s not found, creating new one", subscriptionID)
	return c.pubsubClient.CreateSubscription(c.ctx, subscriptionID, pubsub.SubscriptionConfig{
		Topic: c.pubsubClient.Topic(c.topic),
	})
}

func (c *Creator) createJobForObject(obj *storage.Object) error {
	var ttlSecondsAfterFinished int32 = 60
	var parallelism int32 = 1
	var completions int32 = 1
	suspend := true
	resources := v1.ResourceList{
		"cpu":               resource.MustParse("500m"),
		"memory":            resource.MustParse("512Mi"),
		"ephemeral-storage": resource.MustParse("1Gi"),
	}
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:    c.compressorNamespace,
			GenerateName: "compressor-job-",
			Annotations: map[string]string{
				"kueue.x-k8s.io/queue-name": c.queueName,
			},
		},
		Spec: batchv1.JobSpec{
			TTLSecondsAfterFinished: &ttlSecondsAfterFinished,
			Parallelism:             &parallelism,
			Completions:             &completions,
			Suspend:                 &suspend,
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					ServiceAccountName: saName,
					Containers: []v1.Container{
						{
							Name:  "compressor",
							Image: c.compressorImage,
							Resources: v1.ResourceRequirements{
								Requests: resources,
								Limits:   resources,
							},
							Env: []v1.EnvVar{
								{
									Name:  "COMPRESSOR_LOG_JSON",
									Value: "true",
								},
								{
									Name:  "COMPRESSOR_BUCKET",
									Value: obj.Bucket,
								},
								{
									Name:  "COMPRESSOR_OBJECT",
									Value: obj.Name,
								},
							},
						},
					},
					RestartPolicy: v1.RestartPolicyNever,
				},
			},
		},
	}
	log.Infof("Creating Kubernetes job to compress object %s in bucket %s in namespace %s, image %s", obj.Name, obj.Bucket, c.compressorNamespace, c.compressorImage)
	_, err := c.k8sClient.BatchV1().Jobs(c.compressorNamespace).Create(c.ctx, job, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (c *Creator) createServiceAccountForJob() error {
	sa := corev1.ServiceAccount(saName, c.compressorNamespace).
		WithAnnotations(map[string]string{
			"iam.gke.io/gcp-service-account": c.compressorSA,
		})
	_, err := c.k8sClient.CoreV1().ServiceAccounts(c.compressorNamespace).Apply(c.ctx, sa, metav1.ApplyOptions{FieldManager: "job-creator", Force: false})
	if err != nil {
		return err
	}
	return nil
}
