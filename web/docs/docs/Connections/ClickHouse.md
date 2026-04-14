---
title: ClickHouse
description: ClickHouse columnar database connection
status: stable
---

# ClickHouse

<div class="flex flex-col gap-3 mb-6 pb-6 border-b border-gray-200">
<div class="flex items-center gap-3">
<span class="inline-flex items-center rounded-full px-4 py-2 text-sm font-medium bg-earthy-green-300 text-earthy-green-900">Stable</span>
<span class="inline-flex items-center rounded-full px-4 py-2 text-sm font-medium bg-earthy-blue-100 text-earthy-blue-900 border border-earthy-blue-300">database</span>
</div>
</div>



## Example Configuration

```yaml

host: localhost
port: 9000
user: default
password: your-password
database: default
secure: false

```

## Configuration

The following configuration options are available:

| Property | Type | Required | Sensitive | Default | Description |
|----------|------|----------|-----------|---------|-------------|
| Database | string | true | false |  | Database name |
| Host | string | true | false |  | ClickHouse server hostname or IP address |
| Password | string | false | true |  | ClickHouse password |
| Port | int | false | false | 9000 | ClickHouse server port |
| Secure Connection | bool | false | false | false | Use secure connection (TLS) |
| User | string | true | false |  | ClickHouse username |