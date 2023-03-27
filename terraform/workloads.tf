/**
 * Copyright 2022 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

resource "random_id" "bucket-prefix" {
  byte_length = 8
}

resource "google_storage_bucket" "bucket" {
  project                     = data.google_project.project.project_id
  name                        = "kueue-files-${random_id.bucket-prefix.hex}"
  location                    = "EU"
  uniform_bucket_level_access = true
  depends_on = [
    google_project_service.project
  ]
}

resource "google_pubsub_topic" "bucket-notifications" {
  project = data.google_project.project.project_id
  name    = "bucket-notifications"
  depends_on = [
    google_project_service.project
  ]
}

data "google_storage_project_service_account" "gcs_account" {
  project = data.google_project.project.project_id
}

resource "google_pubsub_topic_iam_member" "bucket-notifications-gcs-sa" {
  project = data.google_project.project.project_id
  topic   = google_pubsub_topic.bucket-notifications.id
  role    = "roles/pubsub.publisher"
  member  = "serviceAccount:${data.google_storage_project_service_account.gcs_account.email_address}"
}

resource "google_storage_notification" "notification" {
  bucket         = google_storage_bucket.bucket.name
  payload_format = "JSON_API_V1"
  topic          = google_pubsub_topic.bucket-notifications.id
  event_types    = ["OBJECT_FINALIZE"]
  depends_on     = [google_pubsub_topic_iam_member.bucket-notifications-gcs-sa]
}

resource "google_artifact_registry_repository" "kueue-handson" {
  project       = data.google_project.project.project_id
  location      = var.cluster_auto_region
  repository_id = "kueue-handson"
  description   = "Repository for Kueue hands on demo"
  format        = "docker"

  depends_on = [
    google_project_service.project
  ]
}

resource "google_artifact_registry_repository_iam_member" "cluster-auto-reader" {
  project    = data.google_project.project.project_id
  location   = google_artifact_registry_repository.kueue-handson.location
  repository = google_artifact_registry_repository.kueue-handson.name
  role       = "roles/artifactregistry.reader"
  member     = "serviceAccount:${google_service_account.cluster-auto.email}"
}

resource "google_service_account" "kueue-demo-job" {
  project      = data.google_project.project.project_id
  account_id   = "kueue-demo-job"
  display_name = "Kueue hands-on demo job"
}

resource "google_service_account_iam_binding" "admin-account-iam" {
  service_account_id = google_service_account.kueue-demo-job.name
  role               = "roles/iam.workloadIdentityUser"
  members = [
    "serviceAccount:${data.google_project.project.project_id}.svc.id.goog[team-a/dispatcher]",
    "serviceAccount:${data.google_project.project.project_id}.svc.id.goog[team-a/compressor]"
  ]
}

resource "google_storage_bucket_iam_member" "bucket-kueue-demo-job-writer" {
  bucket = google_storage_bucket.bucket.name
  role   = "roles/storage.admin"
  member = "serviceAccount:${google_service_account.kueue-demo-job.email}"
}

resource "google_project_iam_member" "notifications-topic-editor" {
  project = data.google_project.project.project_id
  role    = "roles/pubsub.edtor"
  member  = "serviceAccount:${google_service_account.kueue-demo-job.email}"
}
