# Deploying Galera (multi-master) Database Clusters
- [MySQL/MariaDB deployment and configuration](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/mysql-maria)
  - [Percona XtraDB cluster or MariaDB Galera cluster deployment and configuration example](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/mysql-maria/galera/main.tf)
  - [ProxySQL deployment and configuration ](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/mysql-maria#proxysql-load-balancer-for-mysqlmariadb)

### Scaling a database cluster
- [Scaling a cluster](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/mysql-maria#addingremoving-nodes-to-an-existing-cluster---clustercontrol_db_cluster)

### Toggling Cluster Auto-Recovery
- [Toggling cluster auto-recovery](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples#toggling-cluster-auto-recovery-option)

### Backups
- [Backup methods supported for various database types](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples#supported-backup-methods-for-the-respective-database-types-and-vendors)
- [Scheduling backups](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples#scheduling-backups-using-the---clustercontrol_db_cluster_backup_schedule-resource)
- [Taking adhoc backups](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples#taking-adhoc-backups-using-the---clustercontrol_db_cluster_backup-resource)
### Maintenence window
- [Scheduling maintenance window](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples#setting-a-maintenance-window-using-the---clustercontrol_db_cluster_maintenance-resource)
