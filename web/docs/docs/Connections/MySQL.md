---
title: MySQL
description: MySQL database connection
status: stable
---

# MySQL

<div class="flex flex-col gap-3 mb-6 pb-6 border-b border-gray-200">
<div class="flex items-center gap-3">
<span class="inline-flex items-center rounded-full px-4 py-2 text-sm font-medium bg-earthy-green-300 text-earthy-green-900">Stable</span>
<span class="inline-flex items-center rounded-full px-4 py-2 text-sm font-medium bg-earthy-blue-100 text-earthy-blue-900 border border-earthy-blue-300">database</span>
</div>
</div>



## Example Configuration

```yaml

host: localhost
port: 3306
database: mydb
user: admin
password: your-password
tls: "false"

```

## Configuration

The following configuration options are available:

| Property | Type | Required | Sensitive | Default | Description |
|----------|------|----------|-----------|---------|-------------|
| Database | string | true | false |  | Database name to connect to |
| Host | string | true | false |  | MySQL server hostname or IP address |
| Password | string | true | true |  | Database password for authentication |
| Port | int | false | false | 3306 | MySQL server port |
| TLS Mode | string | false | false | false | TLS/SSL mode |
| User | string | true | false |  | Database username for authentication |