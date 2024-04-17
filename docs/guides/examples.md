---
page_title: "Examples"
---

# Examples Guide

A complete set of examples can be found in the [examples](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples) folder.

## Example topics

- [Getting Started](./getting-started.md)
### Deploying Database Clusters
  - [MySQL/MariaDB deployment and configuration](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/mysql-maria)
    - [MySQL/MariaDB replication cluster deployment and configuration example](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/mysql-maria/replication/main.tf)
    - [Percona XtraDB cluster or MariaDB Galera cluster deployment and configuration example](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/mysql-maria/galera/main.tf)
    - [ProxySQL deployment and configuration ](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/mysql-maria#proxysql-load-balancer-for-mysqlmariadb)
  - [PostgreSQL cluster deployment and configuration](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/postgres)
    - [PostgreSQL cluster deployment and configuration example](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/postgres/main.tf)
    - [Enabling TimescaleDB](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/postgres#db_enable_timescale-enabling-timescaledb-extension)
    - [Enabling PgBackRest agent](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/postgres#db_enable_pgbackrest_agent-enabling-pgbackrest-agent)
  - [MongoDB cluster deployment and configuration ](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/mongo)
    - [MongoDB Sharded Replicaset cluster deployment and configuration](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/mongo/sharded-replicaset)
      - [MongoDB Sharded Replicaset cluster deployment and configuration example](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/mongo/sharded-replicaset/main.tf)
      - [Enabling Percona-Backup-MongoDB agent](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/mongo/sharded-replicaset#db_enable_pbm_agent-enabling-pbm-percona-backup-for-mongodb-agent)
    - [MongoDB Replicaset cluster deployment and configuration](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/mongo/replicaset)
      - [MongoDB Replicaset cluster deployment and configuration example](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/mongo/replicaset/main.tf)
  - [Redis cluster deployment and configuration ](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/redis)
    - [Redis Sentinel cluster deployment and configuration example](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/redis/sentinel/main.tf) 
  - [Microsoft SQL server cluster deployment and configuration ](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/mssql)
    - [Microsoft SQL server standalone database deployment and configuration example](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/mssql/single/main.tf)
    - [Microsoft SQL server Primary/Standby cluster deployment and configuration example](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/mssql/multi/main.tf)
- [Elasticsearch cluster deployment and configuration ](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/elastic)
    - [Elasticsearch cluster deployment and configuration example](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/elastic/main.tf)

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
