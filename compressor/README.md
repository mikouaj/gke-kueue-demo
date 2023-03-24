# Compressor

Compressor is a simple app that retrieves object from Google Cloud Storage,
compresses it and stores the compressed object in same Cloud Storage bucket.

The app was made for demonstration purposes only.
The goal is to demonstrate the computational heavy job.

## Configuration

The compressor is configured via environmental variables.

* `COMPRESSOR_BUCKET` - (required) - the Cloud Storage bucket to read the object from
* `COMPRESSOR_OBJECT` - (required) - the Cloud Storage object name to read and compress
* `COMPRESSOR_LOG_LEVEL` - (optional) - logging level (info, debug)
* `COMPRESSOR_LOG_JSON` - (optional) - set to `true` to make compressor format logs in JSON format

## Authentication

The application uses [Application Default Credentials](https://cloud.google.com/docs/authentication/application-default-credentials)
to authenticate in Google Cloud.

## Usage

### Running locally

1. Build the app by running `make`
2. Set environmental variables and run the app

   ```sh
   COMPRESSOR_BUCKET=my-bucket  COMPRESSOR_OBJECT=my-object ./compressor
   ```

### Running as a Kubernetes job

