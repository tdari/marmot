---
title: MongoDB
description: MongoDB NoSQL database connection
status: stable
---

# MongoDB

<div class="flex flex-col gap-3 mb-6 pb-6 border-b border-gray-200">
<div class="flex items-center gap-3">
<span class="inline-flex items-center rounded-full px-4 py-2 text-sm font-medium bg-earthy-green-300 text-earthy-green-900">Stable</span>
<span class="inline-flex items-center rounded-full px-4 py-2 text-sm font-medium bg-earthy-blue-100 text-earthy-blue-900 border border-earthy-blue-300">database</span>
</div>
</div>



## Example Configuration

```yaml

host: localhost
port: 27017
user: admin
password: your-password
auth_source: admin
tls: false

```

## Configuration

The following configuration options are available:

| Property | Type | Required | Sensitive | Default | Description |
|----------|------|----------|-----------|---------|-------------|
| Auth Source | string | false | false | admin | Authentication database |
| Connection URI | string | false | false |  | MongoDB connection URI (mongodb:// or mongodb+srv://) |
| Enable TLS | bool | false | false | false | Enable TLS/SSL connection |
| Host | string | false | false |  | MongoDB server hostname or IP |
| Password | string | false | true |  | MongoDB password |
| Port | int | false | false | 27017 | MongoDB server port |
| TLS Insecure | bool | false | false | false | Skip TLS certificate verification |
| User | string | false | false |  | MongoDB username |