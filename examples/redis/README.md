# Redis Examples

This directory contains an example for deploying Redis (sentinel) clusters using the terraform provider for ClusterControl.

## Resources

| Name                                                                                                                                                                     |
|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [clustercontrol_db_cluster](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster.md#clustercontrol_db_cluster-resource) |
| [clustercontrol_db_cluster_backup](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_backup.md#clustercontrol_db_cluster_backup-resource)|                                                                                                                                                                                    |
| [clustercontrol_db_cluster_backup_schedule](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_backup_schedule.md#clustercontrol_db_cluster_backup_schedule-resource) |
| [clustercontrol_db_cluster_maintenance](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_maintenance.md#clustercontrol_db_cluster_maintenance-resource)|


## Choosing attribute values for MySQL and MariaDB (replication or galera)

### `db_cluster_type` - valid values for MySQL/MariaDB

| Cluster Type     | Description                        |
|------------------|------------------------------------|
| `redis-sentinel` | 1 or 3 node redis sentinel cluster |

### `db_vendor` - valid values

| Vendors | Description             |
|---------|-------------------------|
| redis   | Redis community edition |


### Adding/Removing nodes to an existing cluster - [clustercontrol_db_cluster](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster.md#clustercontrol_db_cluster-resource)

#### Adding a Replicaiton Slave to a cluster

By adding an additional `db_host` block inside the `clsutercontrol_db_cluster` resource you can
add a replication slave to an existing cluster.

Example:

```text
resource "clustercontrol_db_cluster" "this" {
    ...
    db_host {
        hostname = "host-3"
    }
    ...

}
```
The above block will add `host-3` as a replication slave to an existing cluster.

#### Removing a node from a cluster

By removing a `db_host` block from inside the `clsutercontrol_db_cluster` resource you can
remove an existing node from a cluster.

Example:

(**Current State**)

```text
resource "clustercontrol_db_cluster" "this" {
    ...
    db_host {
        hostname = "host-3"
    }
    ...
}
```

(**End State**)

```text
resource "clustercontrol_db_cluster" "this" {
    ...
    ...
}
```

In the above, the end state has removed the `db_host` block for host `host-3`. The result will be the
removal of the corresponding `host-3` node from the cluster.