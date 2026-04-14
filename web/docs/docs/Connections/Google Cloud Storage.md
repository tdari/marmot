---
title: Google Cloud Storage
description: GCS object storage connection
status: stable
---

# Google Cloud Storage

<div class="flex flex-col gap-3 mb-6 pb-6 border-b border-gray-200">
<div class="flex items-center gap-3">
<span class="inline-flex items-center rounded-full px-4 py-2 text-sm font-medium bg-earthy-green-300 text-earthy-green-900">Stable</span>
<span class="inline-flex items-center rounded-full px-4 py-2 text-sm font-medium bg-earthy-blue-100 text-earthy-blue-900 border border-earthy-blue-300">storage</span>
</div>
</div>



## Example Configuration

```yaml

project_id: my-gcp-project
bucket_name: my-bucket
credentials_file: /path/to/key.json

```

## Configuration

The following configuration options are available:

| Property | Type | Required | Sensitive | Default | Description |
|----------|------|----------|-----------|---------|-------------|
| Bucket Name | string | false | false |  | GCS bucket name |
| Credentials File | string | false | false |  | Path to service account JSON key file |
| Credentials JSON | string | false | true |  | Service account JSON key content |
| Custom Endpoint | string | false | false |  | Custom GCS endpoint URL |
| Disable Authentication | bool | false | false | false | Disable authentication (for public buckets) |
| Project ID | string | true | false |  | Google Cloud project ID |