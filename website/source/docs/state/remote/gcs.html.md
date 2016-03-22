---
layout: "remotestate"
page_title: "Remote State Backend: gcs"
sidebar_current: "docs-state-remote-gcs"
description: |-
  Terraform can store the state remotely, making it easier to version and work with in a team.
---

# gcs

Stores the state as a given key in a given bucket on [Google Cloud Storage](https://cloud.google.com/storage/).

-> **Note:** Passing credentials directly via config options will
make them included in cleartext inside the persisted state.
Use of environment variables or config file is recommended.

## Example Usage

```
terraform remote config \
	-backend=gcs \
	-backend-config="bucket=terraform-state-prod" \
	-backend-config="path=network/terraform.tfstate" \
	-backend-config="project=goopro"
```

## Example Referencing

```
resource "terraform_remote_state" "foo" {
	backend = "gcs"
	config {
		bucket = "terraform-state-prod"
		path = "network/terraform.tfstate"
		project = "goopro"
	}
}
```


## Configuration variables

The following configuration options / environment variables are supported:

 * `bucket` - (Required) The name of the Google Cloud Storage bucket.
 * `path` - (Required) The path within the bucket to store the statefile
 * `project` - (Required) A valid API project identifier.
