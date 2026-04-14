---
sidebar_position: 2
---

import { Tabs, TabPanel, TipBox } from '@site/src/components/Steps';

# CLI

The `ingest` command discovers metadata from configured data sources and catalogs them as assets in Marmot. It supports multiple data sources, can establish lineage relationships between assets, and can attach documentation to assets.

## Installation

<Tabs items={[
{ label: "Automatic", value: "auto", icon: "mdi:download" },
{ label: "Manual", value: "manual", icon: "mdi:folder-download" }
]}>
<TabPanel>

```bash
curl -fsSL get.marmotdata.io | sh
```

</TabPanel>
<TabPanel>

Download the latest binary for your platform from [GitHub Releases](https://github.com/marmotdata/marmot/releases), then:

```bash
chmod +x marmot && sudo mv marmot /usr/local/bin/
```

</TabPanel>
</Tabs>

See the [CLI Reference](/docs/cli) for configuring the host, API key and other global options.

## Configuration File

The ingest command requires a YAML configuration file. The top-level structure is:

```yaml
name: my_pipeline_name

connections:
  - name: <connection-name>
    config:
      # connection-specific fields (host, credentials, etc.)

runs:
  - source_type:
      connection: <connection-name>   # references a connection above
      # discovery-specific settings
```

`connections` stores credentials and host configuration. `runs` contains discovery settings for each plugin. The CLI merges them at runtime: connection fields are the base, and run fields take precedence on any key collision.

<TipBox variant="info" title="Pipeline Names">
Give your pipeline a unique name. This is used to track the state of the ingestion and identify stale assets for cleanup.
</TipBox>

You can find all [available source types and their configuration in the Plugins documentation.](/docs/Plugins)

## Connections

The `connections` block defines named credential sets that runs can reference. Each entry requires a `name` and a `config` map whose fields match the corresponding connection type.

```yaml
connections:
  - name: prod-kafka
    config:
      bootstrap_servers: "kafka-broker:9092"
      authentication:
        type: "sasl_ssl"
        username: "your-username"
        password: "your-password"
        mechanism: "PLAIN"

  - name: prod-postgres
    config:
      host: "db.example.com"
      port: 5432
      database: "analytics"
      user: "marmot_reader"
      password: "your-password"
      ssl_mode: "require"
```

A run references a connection by name using the `connection` key:

```yaml
runs:
  - kafka:
      connection: prod-kafka
      # discovery options only, no credentials here
      tags:
        - "kafka"
        - "production"
```

If a run omits `connection`, the plugin receives only the fields defined directly in its run block. This is useful for plugins that don't need external credentials (e.g. file-based plugins).

## Example: Ingesting Kafka Topics

```yaml
name: kafka-production

connections:
  - name: prod-kafka
    config:
      bootstrap_servers: "kafka-broker:9092"
      client_id: "marmot-kafka-plugin"
      client_timeout_seconds: 60
      authentication:
        type: "sasl_ssl"
        username: "your-username"
        password: "your-password"
        mechanism: "PLAIN"
      schema_registry:
        url: "http://schema-registry:8081"
        enabled: true
        config:
          basic.auth.user.info: "your-username:your-password"

runs:
  - kafka:
      connection: prod-kafka
      tags:
        - "kafka"
        - "production"
```

Run the ingestion:

```bash
marmot ingest -c config.yaml
```

## Example: Multiple Sources

You can define multiple connections and runs in a single pipeline file. Each run references its own connection:

```yaml
name: data-platform

connections:
  - name: postgres-prod
    config:
      host: "db.example.com"
      port: 5432
      database: "analytics"
      user: "marmot_reader"
      password: "your-password"
      ssl_mode: "require"

  - name: kafka-prod
    config:
      bootstrap_servers: "kafka-broker:9092"
      authentication:
        type: "sasl_ssl"
        username: "your-username"
        password: "your-password"
        mechanism: "PLAIN"

runs:
  - postgresql:
      connection: postgres-prod
      tags:
        - "postgres"

  - kafka:
      connection: kafka-prod
      tags:
        - "kafka"
```

## Destroying a Pipeline

To remove all assets, lineage, and documentation created by a pipeline, use the `--destroy` flag:

```bash
marmot ingest -c config.yaml --destroy
```

This will prompt for confirmation before deleting anything. The pipeline name in the config file identifies which resources to remove.
