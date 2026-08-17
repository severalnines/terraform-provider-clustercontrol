# ClickHouse Examples

> **Status:** ClickHouse support does not exist yet in `terraform-provider-clustercontrol`
> or in the underlying `clustercontrol-client-sdk`. These examples are scaffolding for
> the provider changes proposed earlier (a new `clickhouse.go` engine, `CLUSTER_TYPE_CLICKHOUSE`,
> etc). Attribute names below (`db_cluster_type = "clickhouse"`, `db_vendor = "clickhouse"`,
> `db_clickhouse_native_port`, `db_clickhouse_keeper_port`) are **proposed**, following this
> repo's existing naming conventions — they are not yet confirmed against the real CMON
> `job_data` contract. Treat this directory as a template to validate/adjust once that
> lands, not as ready-to-apply Terraform against a live ClusterControl instance.
>
> **SSL is mandatory** for ClickHouse: `db_enable_ssl` is hardcoded to `true` in both
> examples, and the port defaults are ClickHouse's secure (TLS) variants.

This directory contains examples for deploying ClickHouse clusters using the terraform
provider for ClusterControl, based on ClusterControl 2.5's ClickHouse support:
single-node (standalone) instances, or replicated clusters using ClickHouse's
**embedded Keeper** for replica coordination (no separate Keeper host tier, no
sharding yet — both are on the ClusterControl roadmap).

| Directory      | Topology                                                              |
|----------------|------------------------------------------------------------------------|
| `standalone/`  | A single ClickHouse node, no replication, no Keeper                   |
| `replicated/`  | 3+ ClickHouse nodes, `ReplicatedMergeTree`-style replication, embedded Keeper on each node |

## Resources

| Name                                                                                                                                               |
|----------------------------------------------------------------------------------------------------------------------------------------------------|
| [clustercontrol_db_cluster](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster.md)               |
| [clustercontrol_db_cluster_backup](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_backup.md) |
| [clustercontrol_db_cluster_backup_schedule](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_backup_schedule.md) |
| [clustercontrol_db_cluster_maintenance](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_maintenance.md) |

## Choosing attribute values for ClickHouse

### `db_cluster_type` - proposed values for ClickHouse

| Cluster Type  | Description                                                    |
|---------------|-----------------------------------------------------------------|
| `clickhouse`  | ClickHouse cluster — 1 node (standalone) or 3+ nodes (replicated, embedded Keeper) |

### `db_vendor` - proposed values

| Vendors      | Description        |
|--------------|---------------------|
| `clickhouse` | ClickHouse (open source) |

### `db_host`

Every node that makes up the cluster needs a `db_host` block. Unlike Elasticsearch
(which needs `roles`) or Redis Sentinel (which needs a separate sentinel companion
process), ClickHouse replicated nodes are symmetric peers — no role attribute is
required. Each node runs both `clickhouse-server` and the embedded `clickhouse-keeper`.

```hcl
resource "clustercontrol_db_cluster" "this" {
    ...
    db_host {
        hostname = "host-1"
    }
    db_host {
        hostname = "host-2"
    }
    db_host {
        hostname = "host-3"
    }
}
```

### Ports

| Variable                     | Default | Purpose                                                 |
|-------------------------------|---------|-----------------------------------------------------------|
| `db_clickhouse_native_port`   | `9440`  | Secure (TLS) native TCP client protocol (all topologies)  |
| `db_clickhouse_keeper_port`   | `9281`  | Secure (TLS) embedded Keeper client port (replicated topology only) |

### SSL is mandatory

These examples require SSL — `db_enable_ssl` is hardcoded to `true` in `main.tf`
(not exposed as a variable), and the port defaults above are ClickHouse's
**secure** variants rather than its plaintext defaults (`9000` and `9181`
respectively). Do not set `db_enable_ssl = false` or swap in the plaintext
port numbers unless you're intentionally deviating from this mandate.

### Adding/removing replica nodes - [clustercontrol_db_cluster](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster.md)

Same pattern as Redis/Elasticsearch: add or remove a `db_host` block inside the
`clustercontrol_db_cluster` resource to add/remove a replica node from an existing
replicated cluster. See `replicated/main.tf` for a 3-node starting point.

### Backup methods

| Method                     | Type                | Description                          |
|-----------------------------|---------------------|----------------------------------------|
| `clickhouse-native`         | Full                | Full backup using ClickHouse-native tooling       |
| `clickhouse-native-incr`    | Incremental         | Incremental backup using ClickHouse-native tooling |

### Scheduling Backups using the - [clustercontrol_db_cluster_backup_schedule](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_backup_schedule.md) Resource

```hcl
resource "clustercontrol_db_cluster_backup_schedule" "daily-full" {
  depends_on            = [clustercontrol_db_cluster.this]
  db_backup_sched_title = "Daily full backup"
  db_backup_sched_time  = "TZ=UTC 0 0 * * *"
  db_cluster_id         = clustercontrol_db_cluster.this.id
  db_backup_method      = "clickhouse-native"
  db_backup_dir         = var.db_backup_dir
  db_backup_subdir      = var.db_backup_subdir
  db_backup_encrypt     = var.db_backup_encrypt
  db_backup_host        = var.db_backup_host
  db_backup_compression = var.db_backup_compression
  db_backup_retention   = var.db_backup_retention
}

resource "clustercontrol_db_cluster_backup_schedule" "hourly-incremental" {
  depends_on            = [clustercontrol_db_cluster.this]
  db_backup_sched_title = "Hourly incremental backup"
  db_backup_sched_time  = "TZ=UTC 0 * * * *"
  db_cluster_id         = clustercontrol_db_cluster.this.id
  db_backup_method      = "clickhouse-native-incr"
  db_backup_dir         = var.db_backup_dir
  db_backup_subdir      = var.db_backup_subdir
  db_backup_encrypt     = var.db_backup_encrypt
  db_backup_host        = var.db_backup_host
  db_backup_compression = var.db_backup_compression
  db_backup_retention   = var.db_backup_retention
}
```

### Taking adhoc backups using the - [clustercontrol_db_cluster_backup](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_backup.md) resource

```hcl
resource "clustercontrol_db_cluster_backup" "full-1" {
  depends_on            = [clustercontrol_db_cluster.this]
  db_cluster_id         = clustercontrol_db_cluster.this.id
  db_backup_method      = "clickhouse-native"
  db_backup_dir         = var.db_backup_dir
  db_backup_subdir      = var.db_backup_subdir
  db_backup_encrypt     = var.db_backup_encrypt
  db_backup_host        = var.db_backup_host
  db_backup_compression = var.db_backup_compression
  db_backup_retention   = var.db_backup_retention
}

resource "clustercontrol_db_cluster_backup" "incr-1" {
  depends_on            = [clustercontrol_db_cluster.this]
  db_cluster_id         = clustercontrol_db_cluster.this.id
  db_backup_method      = "clickhouse-native-incr"
  db_backup_dir         = var.db_backup_dir
  db_backup_subdir      = var.db_backup_subdir
  db_backup_encrypt     = var.db_backup_encrypt
  db_backup_host        = var.db_backup_host
  db_backup_compression = var.db_backup_compression
  db_backup_retention   = var.db_backup_retention
}
```
