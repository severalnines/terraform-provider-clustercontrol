# ClusterControl Provider Examples

This directory contains a set of examples of deploying MySQL or MariaDB replicaton database clusters 
using ClusterControl. 

## Resources

| Name |
|------|
| cc_db_instance |

## Attributes specific to MySQL and MariaDB (replication)

| Attribute                | Data Type   | Required             | Description                                      |
|--------------------------|-------------|----------------------|--------------------------------------------------|
| db_cluster_type | string      | Yes      | Type of cluster -``replication``                 |
| db_semi_sync_replication | boolean     | Optional             | True implies semi-synchronous (default: false (asynchronous)) |
| db_topology             | object      | Yes (for multi-host) | For a multi-host replication cluster, specifies master and slaves |

