---
title: Apache Iceberg (REST Catalog)
description: Apache Iceberg REST catalog connection
status: stable
---

# Apache Iceberg (REST Catalog)

<div class="flex flex-col gap-3 mb-6 pb-6 border-b border-gray-200">
<div class="flex items-center gap-3">
<span class="inline-flex items-center rounded-full px-4 py-2 text-sm font-medium bg-earthy-green-300 text-earthy-green-900">Stable</span>
<span class="inline-flex items-center rounded-full px-4 py-2 text-sm font-medium bg-earthy-blue-100 text-earthy-blue-900 border border-earthy-blue-300">data-lake</span>
</div>
</div>



## Example Configuration

```yaml

uri: http://localhost:8181
warehouse: my-warehouse
token: your-api-secret

```

## Configuration

The following configuration options are available:

| Property | Type | Required | Sensitive | Default | Description |
|----------|------|----------|-----------|---------|-------------|
| Credential | string | false | true |  | Credential for OAuth2 client credentials authentication |
| Prefix | string | false | false |  | Optional prefix for the REST catalog |
| Properties | map[string]string | false | false |  | Additional catalog properties |
| Token | string | false | true |  | Bearer token for authentication |
| URI | string | true | false |  | REST catalog URI |
| Warehouse | string | true | false |  | Warehouse identifier |