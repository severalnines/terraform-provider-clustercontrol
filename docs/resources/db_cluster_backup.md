---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "clustercontrol_db_cluster_backup Resource - terraform-provider-clustercontrol"
subcategory: ""
description: |-
  
---

# clustercontrol_db_cluster_backup (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `db_backup_method` (String) mariabackup, xtrabackup, ...
- `db_cluster_id` (String) The database cluster ID for which this LB is being deployed to.

### Optional

- `db_backup_compression` (Boolean) Whether to compress backups or not
- `db_backup_compression_level` (Number) Compression level
- `db_backup_dir` (String) Base direcory where backups is to be stored
- `db_backup_encrypt` (Boolean) Whether to encrypt or not
- `db_backup_failover_host` (String) If backup failover is enabled, which host to use as backup host in the event of failure of first choice host.
- `db_backup_host` (String) Where there are multiple hosts, which host to choose to create backup from.
- `db_backup_retention` (Number) Backup retention period in days
- `db_backup_storage_controller` (Boolean) Whether to store the backup on CMON controller host or not
- `db_backup_storage_host` (String) Which host to store the backup on. Typically, used with mongodump backup method.
- `db_backup_subdir` (String) Sub-dir for backups - default: "BACKUP-%I"
- `db_backup_system_db` (Boolean) Whether to enable backup failover to another host in case the host crashes
- `db_enable_backup_failover` (Boolean) Whether to enable backup failover to another host in case the host crashes
- `db_snapshot_repository` (String) Elasticsearch snapshot repository

### Read-Only

- `db_resource_id` (String) The ID of the resource allocated by ClusterControl.
- `id` (String) The ID of this resource.
- `last_updated` (String) Last updated timestamp for the resource in question.
