# ClusterControl Provider Examples

This directory contains a set of examples of deploying Redis (Sentinel) database clusters 
using ClusterControl. 

## Resources

| Name |
|------|
| clustercontrol_db_cluster |

## Attributes specific to Redis sentinel

| Attribute                | Data Type   | Required             | Description                                   |
|--------------------------|-------------|----------------------|-----------------------------------------------|
| db_cluster_type | string      | Yes      | Type of cluster. The valid type is -``redis`` |
| db_vendor                | string      | Yes                  | Database vendor (redis)                       |
| db_version               | string      | Yes                  | DB version (6, 7)                             |
| db_host                  | object      | Yes                  | DB host specification                         |

