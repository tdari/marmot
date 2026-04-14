---
title: Trino
description: Trino distributed SQL query engine connection
status: stable
---

# Trino

<div class="flex flex-col gap-3 mb-6 pb-6 border-b border-gray-200">
<div class="flex items-center gap-3">
<span class="inline-flex items-center rounded-full px-4 py-2 text-sm font-medium bg-earthy-green-300 text-earthy-green-900">Stable</span>
<span class="inline-flex items-center rounded-full px-4 py-2 text-sm font-medium bg-earthy-blue-100 text-earthy-blue-900 border border-earthy-blue-300">database</span>
</div>
</div>



## Example Configuration

```yaml

host: trino.company.com
port: 8080
user: marmot_reader
password: your-password
secure: false

```

## Configuration

The following configuration options are available:

| Property | Type | Required | Sensitive | Default | Description |
|----------|------|----------|-----------|---------|-------------|
| Access Token | string | false | true |  | JWT bearer token |
| Host | string | true | false |  | Trino coordinator hostname |
| Password | string | false | true |  | Password (requires HTTPS) |
| Port | int | false | false | 8080 | Trino coordinator port |
| SSL Cert Path | string | false | false |  | Path to TLS certificate file |
| Secure | bool | false | false | false | Use HTTPS |
| User | string | true | false |  | Username for authentication |