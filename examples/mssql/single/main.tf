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
  db_cluster_create      = true
  db_cluster_import      = false
  db_cluster_name        = "mydbcluster"
  db_cluster_type        = "mssql-standalone"
  db_vendor              = "microsoft"
  db_version             = "2022"
  db_admin_username      = "SQLServerAdmin"
  db_admin_user_password = "Iule3cX$Bs5b#KI42~j5D7Yk"
  db_auto_recovery       = true
  db_mssqlserver_port    = var.db_mssqlserver_port
  db_data_directory      = var.db_data_directory
  disable_firewall       = var.disable_firewall
  disable_selinux        = var.disable_selinux
  db_enable_uninstall    = var.db_enable_uninstall
  db_install_software    = var.db_install_software
  db_deploy_agents       = var.db_deploy_agents
  db_enable_ssl          = var.db_enable_ssl
  ssh_user               = var.ssh_user
  ssh_user_password      = var.ssh_user_password
  ssh_key_file           = var.ssh_key_file
  ssh_port               = var.ssh_port
  db_tags                = ["terra-deploy"]

  db_host {
    hostname = "test-primary"
    # hostname_data = "foo"
    # hostname_internal = "foo"
  }
}

# resource "clustercontrol_db_cluster_backup" "full-1" {
#   depends_on                   = [clustercontrol_db_cluster.this]
#   db_cluster_id                = clustercontrol_db_cluster.this.id
#   db_backup_method             = "mssqlfull"
#   db_backup_dir                = "/var/lib/backups"
#   db_backup_subdir             = var.db_backup_subdir
#   db_backup_host               = "auto"
#   db_backup_storage_controller = var.db_backup_storage_controller
#   db_backup_compression        = var.db_backup_compression
#   db_backup_compression_level  = -1
#   db_backup_retention          = var.db_backup_retention
#   db_backup_system_db          = var.db_backup_system_db
# }

# resource "clustercontrol_db_cluster_backup_schedule" "daily-full" {
#   depends_on                   = [clustercontrol_db_cluster.this]
#   db_backup_sched_title        = "Daily snapshot"
#   db_backup_sched_time         = "TZ=UTC 0 0 * * *"
#   db_cluster_id                = clustercontrol_db_cluster.this.id
#   db_backup_method             = "mssqlfull"
#   db_backup_dir                = "/var/lib/backups"
#   db_backup_subdir             = var.db_backup_subdir
#   db_backup_host               = "auto"
#   db_backup_storage_controller = var.db_backup_storage_controller
#   db_backup_compression        = var.db_backup_compression
#   db_backup_compression_level  = -1
#   db_backup_retention          = var.db_backup_retention
#   db_backup_system_db          = var.db_backup_system_db
# }

# resource "clustercontrol_db_cluster_maintenance" "server-upgrade-03232024" {
#   depends_on          = [clustercontrol_db_cluster.this]
#   db_cluster_id       = clustercontrol_db_cluster.this.id
#   db_maint_start_time = "Mar-27-2024T22:00"
#   db_maint_stop_time  = "Mar-28-2024T23:30"
#   db_maint_reason     = "Hardware refresh March 27, 2024"
# }

