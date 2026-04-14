---
title: Azure Blob Storage
description: Azure Blob Storage connection
status: stable
---

# Azure Blob Storage

<div class="flex flex-col gap-3 mb-6 pb-6 border-b border-gray-200">
<div class="flex items-center gap-3">
<span class="inline-flex items-center rounded-full px-4 py-2 text-sm font-medium bg-earthy-green-300 text-earthy-green-900">Stable</span>
<span class="inline-flex items-center rounded-full px-4 py-2 text-sm font-medium bg-earthy-blue-100 text-earthy-blue-900 border border-earthy-blue-300">storage</span>
</div>
</div>



## Example Configuration

```yaml

account_name: mystorageaccount
account_key: your-api-secret
container_name: mycontainer

```

## Configuration

The following configuration options are available:

| Property | Type | Required | Sensitive | Default | Description |
|----------|------|----------|-----------|---------|-------------|
| Account Key | string | false | true |  | Storage account key |
| Account Name | string | false | false |  | Storage account name |
| Connection String | string | false | true |  | Azure Storage connection string |
| Container Name | string | false | false |  | Blob container name |
| Custom Endpoint | string | false | false |  | Custom blob endpoint URL |