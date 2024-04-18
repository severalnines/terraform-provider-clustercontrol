# Deploying MySQL Database Clusters
- [MySQL/MariaDB deployment and configuration](../../examples/mysql-maria/README.md)
  - [MySQL/MariaDB replication cluster deployment and configuration example](../../examples/mysql-maria/replication/main.tf)
  - [Percona XtraDB cluster or MariaDB Galera cluster deployment and configuration example](../../examples/mysql-maria/galera/main.tf)
  - [ProxySQL deployment and configuration ](../../examples/mysql-maria/README.md#proxysql-load-balancer-for-mysqlmariadb)

### Scaling a database cluster
- [Scaling a cluster](../../examples/mysql-maria/README.md#addingremoving-nodes-to-an-existing-cluster---clustercontrol_db_cluster)

### Toggling Cluster Auto-Recovery
- [Toggling cluster auto-recovery](../../examples/README.md#toggling-cluster-auto-recovery-option)

### Backups
- [Backup methods supported for various database types](../../examples/README.md#supported-backup-methods-for-the-respective-database-types-and-vendors)
- [Scheduling backups](../../examples/README.md#scheduling-backups-using-the---clustercontrol_db_cluster_backup_schedule-resource)
- [Taking adhoc backups](../../examples/README.md#taking-adhoc-backups-using-the---clustercontrol_db_cluster_backup-resource)

### Maintenence window
- [Scheduling maintenance window](../../examples/README.md#setting-a-maintenance-window-using-the---clustercontrol_db_cluster_maintenance-resource)
