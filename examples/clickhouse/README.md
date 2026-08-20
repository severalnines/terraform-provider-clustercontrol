# ClickHouse Example

> **Status:** ClickHouse support does not exist yet in `terraform-provider-clustercontrol`
> or in the underlying `clustercontrol-client-sdk`. This example is scaffolding for
> the provider changes discussed for CLUS-8376 (a new `clickhouse.go` engine,
> `CLUSTER_TYPE_CLICKHOUSE`, etc). `db_cluster_type = "clickhouse"`, `db_vendor
> = "clickhouse"`, `db_clickhouse_native_port`, and `db_clickhouse_keeper_port`
> are **proposed**, following this repo's existing naming conventions - not yet
> confirmed against the real CMON `job_data` contract. The host `class_name`
> (`CmonClickHouseHost`, used for every node regardless of role) and the
> dedicated-Keeper `nodetype` value (`clickhouse_keeper`) ARE confirmed
> against real CMON job_data.
>
> **SSL is mandatory** for ClickHouse: `db_enable_ssl` is hardcoded to `true`.
>
> **Sharding is on ClusterControl's roadmap, not shipped yet.** The `shard`
> attribute below is accepted and validated by the provider but has no effect
> on the job sent to CMON today - see the `shard` row in the table below.

This directory contains a single example for deploying a ClickHouse cluster using
the terraform provider for ClusterControl. One `db_cluster` resource now covers
every topology - standalone, replicated, or replicated with dedicated Keeper
hosts - purely through how many `db_host` blocks you declare and what
`roles`/`shard` you give each one. There is no longer a separate `standalone/`
vs `replicated/` split.

## Resources

| Name                                                                                                                                               |
|----------------------------------------------------------------------------------------------------------------------------------------------------|
| [clustercontrol_db_cluster](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster.md)               |
| [clustercontrol_db_cluster_backup](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_backup.md) |
| [clustercontrol_db_cluster_backup_schedule](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_backup_schedule.md) |
| [clustercontrol_db_cluster_maintenance](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_maintenance.md) |

## Choosing attribute values for ClickHouse

### `db_cluster_type` / `db_vendor` - proposed values

| Attribute         | Value        | Description                          |
|--------------------|--------------|-----------------------------------------|
| `db_cluster_type`  | `clickhouse` | ClickHouse cluster, any topology below |
| `db_vendor`        | `clickhouse` | ClickHouse (open source)               |

### `db_host` attributes

| Attribute  | Applies to | Description |
|------------|------------|--------------|
| `hostname` | all        | Required. Hostname or IP of the node. |
| `roles`    | all        | One of `replica` (clickhouse-server; CMON manages embedded Keeper placement among replica hosts on its own - job_data can't distinguish a keeper-less replica from one that also runs embedded Keeper) or `keeper` (dedicated Keeper-only host, no clickhouse-server). Defaults to `replica` if omitted. Available now. |
| `shard`    | `replica` hosts only | Which shard this clickhouse-server host belongs to. **Not yet actionable** - sharded ClickHouse is coming to ClusterControl but isn't shipped yet, so this is validated (rejected on `keeper`-role hosts) but not currently sent to CMON. |

Topology examples using just these two attributes:

- **Standalone**: a single `db_host` block, no `roles`/`shard` set (defaults to `replica`).
- **Replicated** (what ClusterControl 2.5 ships today): 3+ `db_host` blocks, all `roles = "replica"`. CMON manages embedded Keeper placement among them automatically.
- **Replicated with dedicated Keeper hosts** (available now): some hosts `roles = "keeper"` (Keeper only, no ClickHouse data), others `roles = "replica"` (ClickHouse data).
- **Sharded** (roadmap, not shippable yet): multiple `replica` hosts carrying different `shard` values, as shown in `main.tf`'s 6-host example. Terraform will accept this shape today; CMON has nothing to act on it with until sharded ClickHouse ships.

### Ports

| Variable                     | Default | Purpose                                                 |
|-------------------------------|---------|-----------------------------------------------------------|
| `db_clickhouse_native_port`   | `9440`  | Secure (TLS) native TCP client protocol - the per-node `port` for every `replica` host |
| `db_clickhouse_keeper_port`   | `9281`  | Secure (TLS) Keeper client port - the per-node `port` for every dedicated `keeper` host |

### Backup methods

| Method                     | Type                | Description                          |
|-----------------------------|---------------------|----------------------------------------|
| `clickhouse-native`         | Full                | Full backup using ClickHouse-native tooling       |
| `clickhouse-native-incr`    | Incremental         | Incremental backup using ClickHouse-native tooling |

### Adding/removing nodes - [clustercontrol_db_cluster](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster.md)

Add or remove a `db_host` block to add/remove a node from an existing cluster, same
pattern as Redis/Elasticsearch. There is only one CMON host class for ClickHouse
(`replica` and `keeper` roles are distinguished purely by the per-node `nodetype`
and `port` values, not by class), so this follows the same single-class, one-node-
at-a-time restriction as other engines.

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
