# PosgreSQL Examples

This directory contains a set of examples for deploying Postgres clusters using the terraform provider for ClusterControl.

## Resources

| Name                                                                                                                                                                     |
|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [clustercontrol_db_cluster](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster.md#clustercontrol_db_cluster-resource) |
| [clustercontrol_db_cluster_backup](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_backup.md#clustercontrol_db_cluster_backup-resource)|                                                                                                                                                                                    |
| [clustercontrol_db_cluster_backup_schedule](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_backup_schedule.md#clustercontrol_db_cluster_backup_schedule-resource) |
| [clustercontrol_db_cluster_maintenance](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_maintenance.md#clustercontrol_db_cluster_maintenance-resource)|


## Choosing attribute values for MySQL and MariaDB (replication or galera)

### `db_cluster_type` - valid values for MySQL/MariaDB

| Cluster Type   | Description                                                                                       |
|----------------|---------------------------------------------------------------------------------------------------|
| `pg-replication` | Primary with Hot Standby replication cluster. Single host clusters should also use the same value |

### `db_vendor` - valid values

| Vendors    | Description                  |
|------------|------------------------------|
| postgresql | PostgreSQL community edition |

### `db_enable_timescale` Enabling TimescaleDB extension
Use the `db_enable_timescale` attributed in the [clustercontrol_db_cluster](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster.md#clustercontrol_db_cluster-resource) to enable TimescaleDB.

```text
resource "clustercontrol_db_cluster" "this" {
...
  db_enable_timescale        = true
...
}
```


### `db_enable_pgbackrest_agent` Enabling PgBackRest agent
Use the `db_enable_pgbackrest_agent` attributed in the [clustercontrol_db_cluster](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster.md#clustercontrol_db_cluster-resource) to enable PgBackRest agent. Once
the agent is enabled, you can use the pgbackrest(full,incr,diff) backup methods either in
[adhoc backups](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/examples/README.md) or in 
[backup schedules](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/examples/README.md).

```text
resource "clustercontrol_db_cluster" "this" {
...
  db_enable_pgbackrest_agent = false
...
}
```