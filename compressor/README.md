# Compressor

Compressor is a simple app that retrieves object from Google Cloud Storage,
compresses it and stores the compressed object in same Cloud Storage bucket.

The app was made for demonstration purposes only.
The goal is to demonstrate the computational heavy job.

## Configuration

The compressor is configured via environmental variables.

* `COMPRESSOR_BUCKET` - (required) - the Cloud Storage bucket to read the object from
* `COMPRESSOR_OBJECT` - (required) - the Cloud Storage object name to read and compress
* `COMPRESSOR_FOLDER` - (required) - the Cloud Storage folder for compressed objects
* `COMPRESSOR_LOG_LEVEL` - (optional) - logging level (info, debug)
* `COMPRESSOR_LOG_JSON` - (optional) - set to `true` to make compressor format logs in JSON format

## Authentication

The application uses [Application Default Credentials](https://cloud.google.com/docs/authentication/application-default-credentials)
to authenticate in Google Cloud.

## Usage

### Container images

The compressor container images are hosted on [Github Container registry](https://github.com/mikouaj/gke-kueue-demo/pkgs/container/compressor).

```sh
docker pull ghcr.io/mikouaj/compressor:latest
```

### Building from the source code

1. Build the app by running `make`
2. Set environmental variables and run the app

   ```sh
   COMPRESSOR_BUCKET=my-bucket COMPRESSOR_OBJECT=my-object COMPRESSOR_FOLDER=compressed ./compressor
   ```

### Running as a Kubernetes job

1. Create Google Cloud service account with write permissions to the Cloud Storage bucket

   ```sh
   export PROJECT_ID=my-google-cloud-project-id
   export SERVICE_ACCOUNT_NAME=my-service-account
   gcloud iam service-accounts create $SERVICE_ACCOUNT_NAME --project=$PROJECT_ID
   ```

   ```sh
   export BUCKET_NAME=my-storage-bucket
   gcloud storage buckets add-iam-policy-binding gs://$BUCKET_NAME \
   --member=”serviceAccount:${SERVICE_ACCOUNT_NAME}@${PROJECT_ID}.iam.gserviceaccount.com” \
   --role="roles/storage.admin"
   ```

2. Create IAM policy binding for compressor k8s service account on Google Cloud service account

   ```sh
   gcloud iam service-accounts add-iam-policy-binding ${SERVICE_ACCOUNT_NAME}@${PROJECT_ID}.iam.gserviceaccount.com \
   --member=”serviceAccount:PROJECT_ID.svc.id.goog[kueue-demo/compressor]” \
   --role="roles/iam.workloadIdentityUser"
   ```

3. Create Kubernetes service account and annotate it with Google Cloud service account identifier

   ```sh
   apiVersion: v1
   kind: ServiceAccount
   metadata:
     name: compressor
     namespace: kueue-demo
   annotations:
     iam.gke.io/gcp-service-account: my-service-account@my-google-cloud-project-id.iam.gserviceaccount.com
   ```

4. Create Kubernetes job

   ```sh
   apiVersion: batch/v1
   kind: Job
   metadata:
     namespace: kueue-demo
     generateName: compressor-
   spec:
     template:
       spec:
         serviceAccountName: compressor
         containers:
         - name: compressor
           image: ghcr.io/mikouaj/compressor:latest
           resources:
             requests:
               cpu: "500m"
               memory: "512Mi"
             limits:
               cpu: "500m"
               memory: "512Mi"
   ```
