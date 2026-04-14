---
title: AWS
description: AWS connection
status: stable
---

# AWS

<div class="flex flex-col gap-3 mb-6 pb-6 border-b border-gray-200">
<div class="flex items-center gap-3">
<span class="inline-flex items-center rounded-full px-4 py-2 text-sm font-medium bg-earthy-green-300 text-earthy-green-900">Stable</span>
<span class="inline-flex items-center rounded-full px-4 py-2 text-sm font-medium bg-earthy-blue-100 text-earthy-blue-900 border border-earthy-blue-300">cloud</span>
</div>
</div>



## Example Configuration

```yaml

region: us-east-1
id: your-access-key-id
secret: your-api-secret
use_default: false

```

## Configuration

The following configuration options are available:

| Property | Type | Required | Sensitive | Default | Description |
|----------|------|----------|-----------|---------|-------------|
| Access Key ID | string | false | true |  | AWS access key ID |
| Custom Endpoint | string | false | false |  | Custom S3-compatible endpoint URL |
| Profile | string | false | false |  | AWS profile name |
| Region | string | true | false |  | AWS region |
| Role ARN | string | false | false |  | AWS IAM role ARN to assume |
| Role External ID | string | false | false |  | External ID for role assumption |
| Secret Access Key | string | false | true |  | AWS secret access key |
| Session Token | string | false | true |  | AWS session token for temporary credentials |
| Use Default Credentials | bool | false | false | false | Use default AWS credentials from environment/config files |