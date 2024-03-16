# ClusterControl Provider Examples

This directory contains a set of examples of deploying PostgreSQL database clusters 
using ClusterControl. 

## Resources

| Name |
|------|
| cc_db_instance |

## Attributes specific to PostgreSQL primary and hot-standby deployment

| Attribute                | Data Type   | Required | Description                                                                       |
|--------------------------|-------------|----------|-----------------------------------------------------------------------------------|
| db_cluster_type | string      | Yes      | Type of cluster. The valid is -``postgresql_single`` (note: single is a misnomer) |
| db_vendor                | string      | Yes      | Database vendor (default, edb). ``default`` indicates postgresql.org              |
| db_version               | string      | Yes      | DB version (11, 12, 13, 14, 15)                                                   |
| db_admin_user            | string      | Yes      | DB admin user (eg: pgadmin)                                                       |
| db_host                  | object      | Yes      | DB host specification. List of hosts for multi-host cluster                       |

