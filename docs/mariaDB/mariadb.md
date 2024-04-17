# Deploying MariaDB Database Clusters
- [MySQL/MariaDB deployment and configuration](../../examples/mysql-maria)
  - [MySQL/MariaDB replication cluster deployment and configuration example](../../examples/mysql-maria/replication/main.tf)
  - [Percona XtraDB cluster or MariaDB Galera cluster deployment and configuration example](../../examples/mysql-maria/galera/main.tf)
  - [ProxySQL deployment and configuration ](../../examples/mysql-maria#proxysql-load-balancer-for-mysqlmariadb)

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
