---
title: Apache Airflow
description: Apache Airflow workflow orchestration connection
status: stable
---

# Apache Airflow

<div class="flex flex-col gap-3 mb-6 pb-6 border-b border-gray-200">
<div class="flex items-center gap-3">
<span class="inline-flex items-center rounded-full px-4 py-2 text-sm font-medium bg-earthy-green-300 text-earthy-green-900">Stable</span>
<span class="inline-flex items-center rounded-full px-4 py-2 text-sm font-medium bg-earthy-blue-100 text-earthy-blue-900 border border-earthy-blue-300">orchestration</span>
</div>
</div>



## Example Configuration

```yaml

host: http://localhost:8080
username: admin
password: your-password

```

## Configuration

The following configuration options are available:

| Property | Type | Required | Sensitive | Default | Description |
|----------|------|----------|-----------|---------|-------------|
| API Token | string | false | true |  | API token for authentication (alternative to basic auth) |
| Host | string | true | false |  | Airflow webserver URL |
| Password | string | false | true |  | Password for basic authentication |
| Username | string | false | false |  | Username for basic authentication |