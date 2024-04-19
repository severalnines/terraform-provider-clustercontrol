# MySQL/MariaDB examples

This directory contains a set of examples for deploying MySQL or MariaDB database (Master/Slave or Galera multi-master) 
clusters using the terraform provider for ClusterControl. 

## Resources

| Name                                                                                                                                                                     |
|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [clustercontrol_db_cluster](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster.md)                                                 |
| [clustercontrol_db_cluster_backup](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_backup.md)                            |                                                                                                                                                                                    |
| [clustercontrol_db_cluster_backup_schedule](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_backup_schedule.md) |
| [clustercontrol_db_cluster_maintenance](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_maintenance.md)             |


## Choosing attribute values for MySQL and MariaDB (replication or galera)

### `db_cluster_type` - valid values for MySQL/MariaDB

| Cluster Type      | Description                                                               |
|-------------------|---------------------------------------------------------------------------|
| `mysql-replication` | Master/Slave replication cluster                                          |
| `galera`            | Multi-master cluster                                                      |

### `db_vendor` - valid values

| Vendors | Description                      |
|---------|----------------------------------|
| `percona` | Percona's MySQL distribution     |
| `oracle`  | Oracle's MySQL community edition |
| `mariadb` | MariaDB community edition        |

### `db_topology` - Specifying Master --> Slave replication topology
The `db_topology` field within the [clustercontrol_db_cluster](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster.md) should be used to specify the replication topology.

```text
resource "clustercontrol_db_cluster" "this" {
...

    db_host {
        hostname = "host-1"
    }

    db_host {
        hostname = "host-2"
    }

    db_host {
        hostname = "host-3"
    }

    db_topology {
      primary = "host-1"
      replica = "host-2"
    }

    db_topology {
      primary = "host-1"
      replica = "host-3"
    }

}
```

Above, `host-1` is the master and hosts `host-2` and `host-3` are slaves to `host-1`

## ProxySQL load balancer for MySQL/MariaDB
You can deploy a ProxySQL load balancer to your MySQL or MariaDB database cluster. 

```text
resource "clustercontrol_db_cluster" "this" {
...
     db_load_balancer {
       db_lb_type                  = "proxysql"
       db_lb_version               = var.db_lb_version
       db_lb_admin_username        = var.db_lb_admin_username
       db_lb_admin_user_password   = "blah%blah"
       db_lb_monitor_username      = var.db_lb_monitor_username
       db_lb_monitor_user_password = "blah%blah"
       db_lb_port                  = var.db_lb_port
       db_lb_admin_port            = var.db_lb_admin_port
       db_lb_use_clustering        = var.db_lb_use_clustering
       db_lb_use_rw_splitting      = var.db_lb_use_rw_splitting
       db_lb_install_software      = var.db_lb_install_software
       db_lb_enable_uninstall      = var.db_lb_enable_uninstall
       disable_firewall            = var.disable_firewall
       disable_selinux             = var.disable_selinux
       ssh_user                    = var.ssh_user
       ssh_user_password           = var.ssh_user_password
       ssh_key_file                = var.ssh_key_file
       ssh_port                    = var.ssh_port
         db_my_host {
           hostname = "lbhost-1"
         }
     }
}
```
The above will deploy a ProxySQL instance on host `lbhost-1` and will subsequently 
set up all the necessary configuration to front-end the backing database cluster (master/slave or galera).

### Adding/Removing nodes to an existing cluster - [clustercontrol_db_cluster](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster.md)

#### Adding a Replicaiton Slave to a Master/Slave cluster or an addition node to a Galera cluster

By adding an additional `db_host` block inside the `clsutercontrol_db_cluster` resource you can
either add a replication slave to a master/slave cluster or an additional node to a galera cluster.

Example:

```text
resource "clustercontrol_db_cluster" "this" {
    ...
    db_host {
        hostname = "host-3"
    }
    ...

}
```
The above block will add `host-3` as a replication slave to an existing master/slave cluster or as 
an additional node to a galera cluster.

#### Removing a node from a cluster

By removing a `db_host` block from inside the `clsutercontrol_db_cluster` resource you can
remove an existing node from a cluster.

Example: 

(**Current State**)

```text
resource "clustercontrol_db_cluster" "this" {
    ...
    db_host {
        hostname = "host-3"
    }
    ...
}
```

(**End State**)

```text
resource "clustercontrol_db_cluster" "this" {
    ...
    ...
}
```

In the above, the end state has removed the `db_host` block for host `host-3`. The result will be the 
removal of the corresponding `host-3` node from the cluster.

### Scheduling Backups using the - [clustercontrol_db_cluster_backup_schedule](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_backup_schedule.md) Resource
The backup schedule resource allows you to create a backup schedule for a cluster in ClusterControl through the terraform provider.

```hcl
resource "clustercontrol_db_cluster_backup_schedule" "daily-full" {
  depends_on                   = [clustercontrol_db_cluster.this]
  db_backup_sched_title        = "Daily full backup"
  db_backup_sched_time         = "TZ=UTC 0 0 * * *"
  db_cluster_id                = clustercontrol_db_cluster.this.id
  db_backup_method             = "xtrabackupfull"
  db_backup_dir                = var.db_backup_dir
  db_backup_subdir             = var.db_backup_subdir
  db_backup_encrypt            = var.db_backup_encrypt
  db_backup_host               = var.db_backup_host
  db_backup_storage_controller = var.db_backup_storage_controller
  db_backup_compression        = var.db_backup_compression
  db_backup_compression_level  = var.db_backup_compression_level
  db_backup_retention          = var.db_backup_retention
}
```

### Taking adhoc backups using the - [clustercontrol_db_cluster_backup](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_backup.md) resource
You can a maintenance window for a cluster using the `clustercontrol_db_cluster_backup` resource. Here's an example of it.

```hcl
resource "clustercontrol_db_cluster_backup" "full-1" {
  depends_on                   = [clustercontrol_db_cluster.this]
  db_cluster_id                = clustercontrol_db_cluster.this.id
  db_backup_method             = "xtrabackupfull"
  db_backup_dir                = var.db_backup_dir
  db_backup_subdir             = var.db_backup_subdir
  db_backup_encrypt            = var.db_backup_encrypt
  db_backup_host               = var.db_backup_host
  db_backup_storage_controller = var.db_backup_storage_controller
  db_backup_compression        = var.db_backup_compression
  db_backup_compression_level  = var.db_backup_compression_level
  db_backup_retention          = var.db_backup_retention
}
```

### Supported backup methods for MySQL and MariaDB

The following types are supported.

| Database type | Vendor         | Backup method                                                 |
|---------------|----------------|---------------------------------------------------------------|
| MySQL         | Oracle, Percona | `xtrabackupfull`, `xtrabackupincr`, `mysqldump`               |
| MariaDB       | MariaDB         | `mariabackupfull`, `mariabackupincr`, `mysqldump`             |

