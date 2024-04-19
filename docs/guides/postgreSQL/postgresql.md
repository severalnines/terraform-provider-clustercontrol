# Deploying PostgreSQL Database Clusters
- [PostgreSQL cluster deployment and configuration](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/postgres)
  - [PostgreSQL cluster deployment and configuration example](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/examples/postgres/main.tf)
  - [Enabling TimescaleDB](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/postgres#db_enable_timescale-enabling-timescaledb-extension)
  - [Enabling PgBackRest agent](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/postgres#db_enable_pgbackrest_agent-enabling-pgbackrest-agent)

### Scaling a database cluster
- [Scaling a cluster](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/postgres#addingremoving-nodes-to-an-existing-cluster---clustercontrol_db_cluster)

### Toggling Cluster Auto-Recovery
- [Toggling cluster auto-recovery](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples#toggling-cluster-auto-recovery-option)

### Backups
- [Backup methods supported for various database types](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/postgres#supported-backup-methods-for-postgres)
- [Scheduling backups](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/postgres#scheduling-backups-using-the---clustercontrol_db_cluster_backup_schedule-resource)
- [Taking adhoc backups](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/postgres#taking-adhoc-backups-using-the---clustercontrol_db_cluster_backup-resource)

### Maintenence window
- [Scheduling maintenance window](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples#setting-a-maintenance-window-using-the---clustercontrol_db_cluster_maintenance-resource)
