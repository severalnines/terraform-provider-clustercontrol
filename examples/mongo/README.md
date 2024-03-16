# ClusterControl Provider Examples

This directory contains a set of examples of deploying MongoDb database (replicaset) clusters 
using ClusterControl. 

## Resources

| Name |
|------|
| clustercontrol_db_cluster |

## Attributes specific to  MongoDB sharded & replicaset deployment

| Attribute                | Data Type   | Required | Description                                                                                 |
|--------------------------|-------------|----------|---------------------------------------------------------------------------------------------|
| db_cluster_type | string      | Yes      | Type of cluster. The valid is -``mongodb``                                                  |
| db_vendor                | string      | Yes      | Database vendor (percona, 10gen).                                                           |
| db_version               | string      | Yes      | DB version (4.2, 4.4, 5.0, 6.0)                                                             |
| db_admin_user            | string      | Yes      | DB admin user (eg: mongoadmin)                                                              |
| db_replica_set           | object      | Yes      | Replicaset specification. List of replicasets and primary/standby hosts for each replicaset |

