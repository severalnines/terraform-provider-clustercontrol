provider "clustercontrol" {
  cc_api_user          = var.cc_api_user
  cc_api_user_password = var.cc_api_user_password
  cc_api_url           = var.cc_api_url
}


locals {
  is_db_create = (!var.db_cluster_import ? var.db_cluster_create : false)
  is_db_import = var.db_cluster_import
}

resource "clustercontrol_db_cluster" "this" {
  db_cluster_create        = true
  db_cluster_import        = false
  db_cluster_name          = "mydbcluster"
  db_cluster_type          = "mysql-replication"
  db_vendor                = "percona"
  db_version               = "8.4"
  db_admin_user_password   = "blah%blah"
  db_auto_recovery         = true
  db_mysql_port            = var.db_mysql_port
  db_data_directory        = var.db_data_directory
  disable_firewall         = var.disable_firewall
  disable_selinux          = var.disable_selinux
  db_enable_uninstall      = var.db_enable_uninstall
  db_install_software      = var.db_install_software
  db_semi_sync_replication = var.db_semi_sync_replication
  db_deploy_agents         = var.db_deploy_agents
  db_enable_ssl            = var.db_enable_ssl
  ssh_user                 = var.ssh_user
  ssh_user_password        = var.ssh_user_password
  ssh_key_file             = var.ssh_key_file
  ssh_port                 = var.ssh_port
  db_tags                  = ["terra-deploy"]

  db_host {
    hostname = "test-primary-1"
    # hostname_data = "foo"
    # hostname_internal = "foo"
    # config_file = "foo"
    # datadir = "foo"
  }

  # db_host {
  #   hostname = "test-primary-2"
  #   # hostname_data = "foo"
  #   # hostname_internal = "foo"
  #   # config_file = "foo"
  #   # datadir = "foo"
  # }

  #  db_topology {
  #    primary = "test-primary"
  #    replica = "test-primary-2"
  #  }

  #   db_load_balancer {
  #     db_lb_type                  = "proxysql"
  #     db_lb_version               = var.db_lb_version
  #     db_lb_admin_username        = var.db_lb_admin_username
  #     db_lb_admin_user_password   = "blah%blah"
  #     db_lb_monitor_username      = var.db_lb_monitor_username
  #     db_lb_monitor_user_password = "blah%blah"
  #     db_lb_port                  = var.db_lb_port
  #     db_lb_admin_port            = var.db_lb_admin_port
  #     db_lb_use_clustering        = var.db_lb_use_clustering
  #     db_lb_use_rw_splitting      = var.db_lb_use_rw_splitting
  #     db_lb_install_software      = var.db_lb_install_software
  #     db_lb_enable_uninstall      = var.db_lb_enable_uninstall
  #     disable_firewall            = var.disable_firewall
  #     disable_selinux             = var.disable_selinux
  #     ssh_user                    = var.ssh_user
  #     ssh_user_password           = var.ssh_user_password
  #     ssh_key_file                = var.ssh_key_file
  #     ssh_port                    = var.ssh_port
  #       db_my_host {
  #         hostname = "test-primary-3"
  #       }
  #   }

}

# resource "clustercontrol_db_cluster_backup_schedule" "full-1" {
#   depends_on                   = [clustercontrol_db_cluster.this]
#   db_backup_sched_title        = "Daily full"
#   db_backup_sched_time         = "TZ=UTC 0 0 * * *"
#   db_cluster_id                = clustercontrol_db_cluster.this.id
#   db_backup_method             = "xtrabackupfull"
#   db_backup_dir                = var.db_backup_dir
#   db_backup_subdir             = var.db_backup_subdir
#   db_backup_encrypt            = var.db_backup_encrypt
#   db_backup_host               = var.db_backup_host
#   db_backup_storage_controller = var.db_backup_storage_controller
#   db_backup_compression        = var.db_backup_compression
#   db_backup_compression_level  = var.db_backup_compression_level
#   db_backup_retention          = var.db_backup_retention
# }

# resource "clustercontrol_db_cluster_backup" "full-1" {
#   depends_on                   = [clustercontrol_db_cluster.this]
#   db_cluster_id                = clustercontrol_db_cluster.this.id
#   db_backup_method             = "xtrabackupfull"
#   db_backup_dir                = var.db_backup_dir
#   db_backup_subdir             = var.db_backup_subdir
#   db_backup_encrypt            = var.db_backup_encrypt
#   db_backup_host               = var.db_backup_host
#   db_backup_storage_controller = var.db_backup_storage_controller
#   db_backup_compression        = var.db_backup_compression
#   db_backup_compression_level  = var.db_backup_compression_level
#   db_backup_retention          = var.db_backup_retention
# }

# resource "clustercontrol_db_cluster_backup" "full-2" {
#   depends_on = [clustercontrol_db_cluster.this]
#   db_cluster_id                = clustercontrol_db_cluster.this.id
#   db_backup_method             = "mysqldump"
#   db_backup_dir                = var.db_backup_dir
#   db_backup_subdir             = var.db_backup_subdir
#   db_backup_encrypt            = var.db_backup_encrypt
#   db_backup_host               = var.db_backup_host
#   db_backup_storage_controller = var.db_backup_storage_controller
#   db_backup_compression        = var.db_backup_compression
#   db_backup_compression_level  = var.db_backup_compression_level
#   db_backup_retention          = var.db_backup_retention
# }

# resource "clustercontrol_db_cluster_backup" "full-3" {
#   depends_on = [clustercontrol_db_cluster.this]
#   db_cluster_id                = clustercontrol_db_cluster.this.id
#   db_backup_method             = "xtrabackupfull"
#   db_backup_dir                = var.db_backup_dir
#   db_backup_subdir             = var.db_backup_subdir
#   db_backup_encrypt            = var.db_backup_encrypt
#   db_backup_host               = var.db_backup_host
#   db_backup_storage_controller = var.db_backup_storage_controller
#   db_backup_compression        = var.db_backup_compression
#   db_backup_compression_level  = var.db_backup_compression_level
#   db_backup_retention          = var.db_backup_retention
# }

# resource "clustercontrol_db_cluster_maintenance" "server-upgrade-03232024" {
#   depends_on = [clustercontrol_db_cluster.this]
#   db_cluster_id       = clustercontrol_db_cluster.this.id
#   db_maint_start_time = "Mar-31-2024T00:00"
#   db_maint_stop_time  = "Mar-31-2024T23:30"
#   db_maint_reason     = "Hardware refresh March 31, 2024"
# }
