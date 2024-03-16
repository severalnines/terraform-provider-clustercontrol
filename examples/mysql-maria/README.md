# ClusterControl Provider Examples

This directory contains a set of examples of deploying MySQL or MariaDB database clusters 
using ClusterControl. 

## Resources

| Name |
|------|
| cc_db_instance |

## Attributes specific to MySQL and MariaDB (replication or galera)

| Attribute                | Data Type   | Required             | Description                                                               |
|--------------------------|-------------|----------------------|---------------------------------------------------------------------------|
| db_cluster_type | string      | Yes      | Type of cluster. The valid types are -``replication`` or ``galera``       |
| db_vendor                | string      | Yes                  | Database vendor (mariadb, oracle, percona)                                |
| db_version               | string      | Yes                  | DB version (MySQL 8.0, MariaDB 10.3,10.4,10.5,10.6,10.8,10.9,10.10,10.11) |
| db_admin_user            | string      | Optional             | DB admin user (default: root)                                             |
| db_host                  | object      | Yes                  | DB host specification                                                     |

