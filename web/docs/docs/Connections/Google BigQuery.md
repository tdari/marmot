---
title: Google BigQuery
description: Google BigQuery data warehouse connection
status: stable
---

# Google BigQuery

<div class="flex flex-col gap-3 mb-6 pb-6 border-b border-gray-200">
<div class="flex items-center gap-3">
<span class="inline-flex items-center rounded-full px-4 py-2 text-sm font-medium bg-earthy-green-300 text-earthy-green-900">Stable</span>
<span class="inline-flex items-center rounded-full px-4 py-2 text-sm font-medium bg-earthy-blue-100 text-earthy-blue-900 border border-earthy-blue-300">database</span>
</div>
</div>



## Example Configuration

```yaml

project_id: my-gcp-project
dataset: my_dataset
credentials_path: /path/to/key.json
location: US

```

## Configuration

The following configuration options are available:

| Property | Type | Required | Sensitive | Default | Description |
|----------|------|----------|-----------|---------|-------------|
| Credentials Path | string | false | false |  | Path to service account JSON key file |
| Dataset | string | false | false |  | BigQuery dataset name |
| Location | string | false | false | US | BigQuery dataset location |
| Project ID | string | true | false |  | Google Cloud project ID |
| Service Account Key | string | false | true |  | JSON key file content for service account |
| Use Default Credentials | bool | false | false | false | Use application default credentials |