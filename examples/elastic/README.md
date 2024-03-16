# ClusterControl Provider Examples

This directory contains a set of examples of deploying Elasticsearch clusters 
using ClusterControl. 

## Resources

| Name |
|------|
| clustercontrol_db_cluster |

## Attributes specific to Elasticsearch

| Attribute                | Data Type   | Required             | Description                                                             |
|--------------------------|-------------|----------------------|-------------------------------------------------------------------------|
| db_cluster_type | string      | Yes      | Type of cluster. The valid is -``elastic`` |
| db_vendor                | string      | Yes                  | Database vendor (elasticsearch)                                         |
| db_version               | string      | Yes                  | DB version (7.17.3, 8.1.3, 8.3.1)                                       |
| db_admin_user            | string      | Yes      | DB admin user (eg: esadmin)                                             |
| db_snapshot_location            | string      | Yes      | Path to snapshot location                                               |
| db_snapshot_repository            | string      | Yes      | Name of snapshot repository                                             |
| db_host                  | object      | Yes                  | DB host specification                                                   |

