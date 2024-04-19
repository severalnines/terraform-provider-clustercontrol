# Deploying MongoDB Database Clusters
- [MongoDB cluster deployment and configuration ](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/mongo)
  - [MongoDB Sharded Replicaset cluster deployment and configuration](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/mongo/sharded-replicaset)
    - [MongoDB Sharded Replicaset cluster deployment and configuration example](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/examples/mongo/sharded-replicaset/main.tf)
    - [Enabling Percona-Backup-MongoDB agent](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/mongo#db_enable_pbm_agent-enabling-pbm-percona-backup-for-mongodb-agent)
  - [MongoDB Replicaset cluster deployment and configuration](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/mongo/replicaset)
    - [MongoDB Replicaset cluster deployment and configuration example](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/examples/mongo/replicaset/main.tf)

### Scaling a database cluster
- [Scaling a cluster](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/mongo#addingremoving-a-node-tofrom-a-replicaset---clustercontrol_db_cluster)

### Toggling Cluster Auto-Recovery
- [Toggling cluster auto-recovery](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples#toggling-cluster-auto-recovery-option)

### Backups
- [Backup methods supported for various database types](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/mongo#supported-backup-methods-are-supported-for-mongodb)
- [Scheduling backups](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/mongo#scheduling-backups-using-the---clustercontrol_db_cluster_backup_schedule-resource)
- [Taking adhoc backups](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/mongo#taking-adhoc-backups-using-the---clustercontrol_db_cluster_backup-resource)

### Maintenence window
- [Scheduling maintenance window](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples#setting-a-maintenance-window-using-the---clustercontrol_db_cluster_maintenance-resource)
