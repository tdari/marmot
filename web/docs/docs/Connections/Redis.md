---
title: Redis
description: Redis in-memory data store connection
status: stable
---

# Redis

<div class="flex flex-col gap-3 mb-6 pb-6 border-b border-gray-200">
<div class="flex items-center gap-3">
<span class="inline-flex items-center rounded-full px-4 py-2 text-sm font-medium bg-earthy-green-300 text-earthy-green-900">Stable</span>
<span class="inline-flex items-center rounded-full px-4 py-2 text-sm font-medium bg-earthy-blue-100 text-earthy-blue-900 border border-earthy-blue-300">database</span>
</div>
</div>



## Example Configuration

```yaml

host: localhost
port: 6379
password: your-password
db: 0
tls: false

```

## Configuration

The following configuration options are available:

| Property | Type | Required | Sensitive | Default | Description |
|----------|------|----------|-----------|---------|-------------|
| Database | int | false | false | 0 | Default database number |
| Enable TLS | bool | false | false | false | Enable TLS connection |
| Host | string | true | false |  | Redis server hostname or IP address |
| Password | string | false | true |  | Password for authentication |
| Port | int | false | false | 6379 | Redis server port |
| TLS Insecure | bool | false | false | false | Skip TLS certificate verification |
| Username | string | false | false |  | Username for ACL authentication |