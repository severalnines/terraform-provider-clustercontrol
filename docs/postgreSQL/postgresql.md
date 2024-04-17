# Deploying PostgreSQL Database Clusters
- [PostgreSql cluster deployment and configuration](../../examples/postgres)
  - [PostgreSql cluster deployment and configuration example](../../examples/postgres/main.tf)
  - [Enabling TimescaleDB](../../examples/postgres#db_enable_timescale-enabling-timescaledb-extension)
  - [Enabling PgBackRest agent](../../examples/postgres#db_enable_pgbackrest_agent-enabling-pgbackrest-agent)

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
