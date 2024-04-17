# Elasticsearch Examples

This directory contains a set of examples for deploying Elasticsearch clusters using the terraform provider for ClusterControl.

## Resources

| Name                                                                                                                                                                     |
|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [clustercontrol_db_cluster](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster.md#clustercontrol_db_cluster-resource) |
| [clustercontrol_db_cluster_backup](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_backup.md#clustercontrol_db_cluster_backup-resource)|                                                                                                                                                                                    |
| [clustercontrol_db_cluster_backup_schedule](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_backup_schedule.md#clustercontrol_db_cluster_backup_schedule-resource) |
| [clustercontrol_db_cluster_maintenance](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_maintenance.md#clustercontrol_db_cluster_maintenance-resource)|


## Choosing attribute values for MySQL and MariaDB (replication or galera)

### `db_cluster_type` - valid values for MySQL/MariaDB

| Cluster Type   | Description                     |
|----------------|---------------------------------|
| `elasticsearch` | Elasticsearch community edition |

### `db_vendor` - valid values

| Vendors   | Description |
|-----------|-----------|
| `elastic` | Elastic |


### `db_host`
The `db_host` block inside the `clsutercontrol_db_cluster` resource specifies the hosts that make up the cluster. Each host
that makes up the DB cluster should have one of these blocks. The mandatory attribute for each `db_host` block are:

| Vendors   | Description |
|-----------|-----------|
| `hostname` | Host name of the host |
| `roles` | `master-data` indicating that the host will be both a master node and a data node (Elastisearch terms) |

Example:

```text
resource "clustercontrol_db_cluster" "this" {
    ...
    db_host {
        hostname = "host-1"
        roles    = "master-data"
    }
    db_host {
        hostname = "host-2"
        roles    = "master-data"
    }
    db_host {
        hostname = "host-3"
        roles    = "master-data"
    }

}
```

### Host/Node roles

| Role          | Description                                                      |
|---------------|------------------------------------------------------------------|
| `master-data` | The host will function as both a `Master` node and a `Data` node |
| `master`      | The host will function as both a `Master` node only              |
| `data`        | The host will function as both a `Data` node only                |
