# ClusterControl Provider Examples

This directory contains a set of examples of deploying Microsoft SQL server database clusters 
using ClusterControl. 

## Resources

| Name |
|------|
| cc_db_instance |

## Attributes specific to Microsoft SQL Server

| Attribute                | Data Type   | Required             | Description                        |
|--------------------------|-------------|----------------------|------------------------------------|
| db_vendor                | string      | Yes                  | Database vendor (microsoft)        |
| db_version               | string      | Yes                  | DB version (2019, 2022)            |
| db_host                  | object      | Yes                  | DB host specification              |
| db_admin_user            | string      | Yes      | DB admin user (eg: SQLServerAdmin) |

