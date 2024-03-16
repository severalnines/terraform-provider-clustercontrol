# ClusterControl Provider Examples

This directory contains a set of examples of using ClusterControl to deploy various DB types (MySQL/MariaDB, PostgreSQL, MongoDB, Redis, MSSQL, Elastic). This is the top level examples directory. The sub-directories contain specific database types and respective toplogies (primary/standby, multi-master).

## Resources

| Name |
|------|
| clustercontrol_db_cluster |

## Common attributes to all database types

| Attribute  | Data Type   | Required | Description                                                                  |
|------------|-------------|----------|------------------------------------------------------------------------------|
| db_cluster_create | boolean     | Optional | Create/Deploy a DB cluster (default: true)                                   |
| db_cluster_import | boolean      | Optional      | Import an existing cluster into ClusterControl (default: false)              |
| db_cluster_name | string      | Yes      | Name of cluster                                                              |
| db_cluster_type | string      | Yes      | Type of cluster (replication, galera, mongodb, etc)                          |
| db_vendor  | string      | Yes      | Database vendor (mariadb, oracle, percona, 10gen, etc)                       |
| db_version | string      | Yes      | DB version (MySQL 8.0, MariaDB 10.11, PosgreSQL 15, etc)                     |
| db_admin_user_password | string      | Yes      | Admin user's password                                                        |
| ssh_user  | string      | Yes      | The SSH user used to SSH into the DB host from ClusterControl host           |
| ssh_user_password | string      | Optional | The SSH user's password                                                      |
| ssh_key_file | string      | Yes      | Path of the private SSH key file on the ClusterControl host for the SSH user |
| db_tags | set(string) | Optional | Comma separated set of tags (e.g. "dev", "app-v1")                           |
| disable_firewall | boolean | Optional | Disable firewall on the DB host (default: true)                              |
| db_install_software | boolean | Optional | Install DB software using vendor repositories (default: true)                |

