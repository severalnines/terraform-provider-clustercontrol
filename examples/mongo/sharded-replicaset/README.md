# MongoDB Shard Example

This directory contains an example for deploying MongoDB shard cluster using the terraform provider for ClusterControl.

## Resources

| Name                                                                                                                                               |
|----------------------------------------------------------------------------------------------------------------------------------------------------|
| [clustercontrol_db_cluster](../../docs/resources/db_cluster.md#clustercontrol_db_cluster-resource)                                                 |
| [clustercontrol_db_cluster_backup](../../docs/resources/db_cluster_backup.md#clustercontrol_db_cluster_backup-resource)                            |                                                                                                                                                                                    |
| [clustercontrol_db_cluster_backup_schedule](../../docs/resources/db_cluster_backup_schedule.md#clustercontrol_db_cluster_backup_schedule-resource) |
| [clustercontrol_db_cluster_maintenance](../../docs/resources/db_cluster_maintenance.md#clustercontrol_db_cluster_maintenance-resource)             |


### Specifying MongoDB Shards

#### Specifing MongoDB Config servers and Mongos server for shard clusters along with replicaset members

```text
resource "clustercontrol_db_cluster" "this" {
...

    db_config_server {
        rs = "rs_config"
        member {
          hostname = "config-server"
        }
    }

    db_mongos_server {
        hostname = "config-server"
    }

    db_replica_set {
        rs = "rs0"
        member {
          hostname = "shard0-host-1"
        }
        member {
          hostname = "shard0-host-2"
        }
    }

}
```

The `db_config_server` and `db_mongos_server` fields within the [clustercontrol_db_cluster](../../docs/resources/db_cluster.md#clustercontrol_db_cluster-resource) should be used to specify
the Mongo config server and mongos server.

Above, the `db_replica_set` specifies a shard with two hosts in the replicaset.
