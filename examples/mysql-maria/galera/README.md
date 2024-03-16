# ClusterControl Provider Examples

This directory contains a set of examples of deploying MySQL (Percona) or MariaDB (Galera - multi-master) database clusters 
using ClusterControl. 

## Resources

| Name |
|------|
| cc_db_instance |

## Attributes specific to MySQL and MariaDB (replication or galera)

| Attribute                | Data Type   | Required             | Description                                                             |
|--------------------------|-------------|----------------------|-------------------------------------------------------------------------|
| db_cluster_type | string      | Yes      | Type of cluster - ``galera``                                            |
| db_vendor                | string      | Yes                  | Database vendor (mariadb, percona)                                |
| db_host                  | object      | Yes                  | DB host specification                                                   |

