---
page_title: "Examples"
---

# Examples Guide

A complete set of examples can be found in the [examples](../../examples/README.md) folder.

## Example topics

- [Getting Started](./getting-started.md)
### Deploying Database Clusters
  - [MySQL/MariaDB deployment and configuration](../../examples/mysql-maria/README.md)
    - [MySQL/MariaDB replication cluster deployment and configuration example](../../examples/mysql-maria/replication/main.tf)
    - [Percona XtraDB cluster or MariaDB Galera cluster deployment and configuration example](../../examples/mysql-maria/galera/main.tf)
    - [ProxySQL deployment and configuration ](../../examples/mysql-maria#proxysql-load-balancer-for-mysqlmariadb)
  - [PostgreSQL cluster deployment and configuration](../../examples/postgres)
    - [PostgreSQL cluster deployment and configuration example](../../examples/postgres/main.tf)
    - [Enabling TimescaleDB](../..n/examples/postgres#db_enable_timescale-enabling-timescaledb-extension)
    - [Enabling PgBackRest agent](../../examples/postgres#db_enable_pgbackrest_agent-enabling-pgbackrest-agent)
  - [MongoDB cluster deployment and configuration ](../../examples/mongo)
    - [MongoDB Sharded Replicaset cluster deployment and configuration](../../examples/mongo/sharded-replicaset)
      - [MongoDB Sharded Replicaset cluster deployment and configuration example](../../examples/mongo/sharded-replicaset/main.tf)
      - [Enabling Percona-Backup-MongoDB agent](../../examples/mongo/sharded-replicaset#db_enable_pbm_agent-enabling-pbm-percona-backup-for-mongodb-agent)
    - [MongoDB Replicaset cluster deployment and configuration](../../examples/mongo/replicaset)
      - [MongoDB Replicaset cluster deployment and configuration example](../../examples/mongo/replicaset/main.tf)
  - [Redis cluster deployment and configuration ](../../examples/redis)
    - [Redis Sentinel cluster deployment and configuration example](../../examples/redis/sentinel/main.tf) 
  - [Microsoft SQL server cluster deployment and configuration ](../../examples/mssql)
    - [Microsoft SQL server standalone database deployment and configuration example](../../examples/mssql/single/main.tf)
    - [Microsoft SQL server Primary/Standby cluster deployment and configuration example](../../examples/mssql/multi/main.tf)
- [Elasticsearch cluster deployment and configuration ](../../examples/elastic)
    - [Elasticsearch cluster deployment and configuration example](../../examples/elastic/main.tf)

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
