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

output "kueue-demo-sa" {
  description = "Email identifier of the service account for kueue-demo job workloads"
  value       = google_service_account.kueue-demo-job.email
}

output "kueue-demo-bucket-name" {
  description = "Name of the bucket for kueue-demo files"
  value       = google_storage_bucket.bucket.name
}
