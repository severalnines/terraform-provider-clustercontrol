# Deploying MongoDB Database Clusters
- [MongoDB cluster deployment and configuration ](../../examples/mongo)
  - [MongoDB Sharded Replicaset cluster deployment and configuration](../../examples/mongo/sharded-replicaset)
    - [MongoDB Sharded Replicaset cluster deployment and configuration example](../../examples/mongo/sharded-replicaset/main.tf)
    - [Enabling Percona-Backup-MongoDB agent](../../examples/mongo/sharded-replicaset#db_enable_pbm_agent-enabling-pbm-percona-backup-for-mongodb-agent)
  - [MongoDB Replicaset cluster deployment and configuration](../../examples/mongo/replicaset)
    - [MongoDB Replicaset cluster deployment and configuration example](../../examples/mongo/replicaset/main.tf)

### Scaling a database cluster
- [Scaling a cluster](../../examples/mysql-maria#addingremoving-nodes-to-an-existing-cluster---clustercontrol_db_cluster)

### Toggling Cluster Auto-Recovery
- [Toggling cluster auto-recovery](../../examples#toggling-cluster-auto-recovery-option)

### Backups
- [Backup methods supported for various database types](../../examples#supported-backup-methods-for-the-respective-database-types-and-vendors)
- [Scheduling backups](../../examples#scheduling-backups-using-the---clustercontrol_db_cluster_backup_schedule-resource)
- [Taking adhoc backups](../../examples#taking-adhoc-backups-using-the---clustercontrol_db_cluster_backup-resource)

### Maintenence window
- [Scheduling maintenance window](../../examples#setting-a-maintenance-window-using-the---clustercontrol_db_cluster_maintenance-resource)
