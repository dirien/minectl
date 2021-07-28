![Google Cloud](https://img.shields.io/badge/google--cloud-4285F4?style=for-the-badge&logo=google-cloud&logoColor=white)
# Getting Started - GCE edition

## How to create a Service Account in GCP via gcloud cli

```bash
# Get current projectID
export PROJECTID=$(gcloud config get-value core/project 2>/dev/null)

# Create a service account
gcloud iam service-accounts create minctl \
--description "minectl-sa service account" \
--display-name "minctl"

# Get service account email
export SERVICEACCOUNT=$(gcloud iam service-accounts list | grep minctl | awk '{print $2}')

# Assign appropriate roles to minectl service account
gcloud projects add-iam-policy-binding $PROJECTID \
--member serviceAccount:$SERVICEACCOUNT \
--role roles/compute.admin
gcloud projects add-iam-policy-binding $PROJECTID \
--member serviceAccount:$SERVICEACCOUNT \
--role roles/iam.serviceAccountUser
gcloud projects add-iam-policy-binding $PROJECTID \
--member serviceAccount:$SERVICEACCOUNT \
--role roles/compute.osAdminLogin

# Create minectl service account key file
gcloud iam service-accounts keys create key.json \
--iam-account $SERVICEACCOUNT
```