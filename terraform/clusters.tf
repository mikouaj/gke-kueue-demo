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

resource "google_service_account" "cluster-auto" {
  project      = data.google_project.project.project_id
  account_id   = "cluster-auto"
  display_name = "Node SA for cluster auto"
}

resource "google_project_iam_member" "cluster-auto-node-sa" {
  project = data.google_project.project.project_id
  role    = "roles/container.nodeServiceAccount"
  member  = "serviceAccount:${google_service_account.cluster-auto.email}"
}

resource "google_container_cluster" "cluster-auto" {
  project    = data.google_project.project.project_id
  name       = "cluster-auto"
  location   = var.cluster_auto_region
  network    = google_compute_network.kueue-demo.id
  subnetwork = google_compute_subnetwork.kueue-demo-cluster-auto.id
  release_channel {
    channel = "RAPID"
  }
  ip_allocation_policy {
  }
  enable_autopilot = true
  cluster_autoscaling {
    auto_provisioning_defaults {
      service_account = google_service_account.cluster-auto.email
    }
  }
  depends_on = [
    google_project_service.project
  ]
}
