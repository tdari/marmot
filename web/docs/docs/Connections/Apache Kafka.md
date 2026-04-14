---
title: Apache Kafka
description: Kafka streaming platform connection
status: stable
---

# Apache Kafka

<div class="flex flex-col gap-3 mb-6 pb-6 border-b border-gray-200">
<div class="flex items-center gap-3">
<span class="inline-flex items-center rounded-full px-4 py-2 text-sm font-medium bg-earthy-green-300 text-earthy-green-900">Stable</span>
<span class="inline-flex items-center rounded-full px-4 py-2 text-sm font-medium bg-earthy-blue-100 text-earthy-blue-900 border border-earthy-blue-300">messaging</span>
</div>
</div>



## Example Configuration

```yaml

bootstrap_servers: kafka-1:9092,kafka-2:9092
client_id: marmot-discovery
client_timeout_seconds: 30
authentication:
  type: sasl_ssl
  username: your-username
  password: your-password
  mechanism: PLAIN
tls:
  enabled: true

```

## Configuration

The following configuration options are available:

| Property | Type | Required | Sensitive | Default | Description |
|----------|------|----------|-----------|---------|-------------|
| Authentication | KafkaAuthConfig | false | false |  | Authentication configuration |
| Bootstrap Servers | string | true | false |  | Comma-separated list of bootstrap servers |
| Client ID | string | false | false |  | Client ID for the consumer |
| Client Timeout | int | false | false | 30 | Request timeout in seconds |
| Consumer Config | map[string]string | false | false |  | Additional consumer configuration |
| Schema Registry | KafkaSchemaRegistryConfig | false | false |  | Schema Registry configuration |
| TLS | KafkaTLSConfig | false | false |  | TLS configuration |