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
  db_cluster_create          = true
  db_cluster_import          = false
  db_cluster_name            = "mydbcluster-pg"
  db_cluster_type            = "postgresql_single"
  db_vendor                  = "default"
  db_version                 = "15"
  db_admin_username          = "pgadmin"
  db_admin_user_password     = "blah%blah"
  db_auto_recovery           = true
  db_port                    = var.db_port
  db_data_directory          = var.db_data_directory
  disable_firewall           = var.disable_firewall
  disable_selinux            = var.disable_selinux
  db_enable_uninstall        = var.db_enable_uninstall
  db_install_software        = var.db_install_software
  db_deploy_agents           = var.db_deploy_agents
  db_enable_ssl              = var.db_enable_ssl
  db_enable_timescale        = false
  db_enable_pgbackrest_agent = true
  ssh_user                   = var.ssh_user
  ssh_user_password          = var.ssh_user_password
  ssh_key_file               = var.ssh_key_file
  ssh_port                   = var.ssh_port
  db_tags                    = ["terra-deploy"]

  db_host {
    hostname = "test-primary-4"
    # hostname_data = "foo"
    # hostname_internal = "foo"
    # port = "foo"
    # config_file = "foo"
    # data_dir = "foo"
  }

  # db_host {
  #   hostname = "test-primary-2"
  #   # hostname_data = "foo"
  #   # hostname_internal = "foo"
  #   # port = "foo"
  #   # config_file = "foo"
  #   # data_dir = "foo"
  # }

}

# resource "clustercontrol_db_cluster_backup" "full-03-27-2024_1" {
#   depends_on                   = [clustercontrol_db_cluster.this]
#   db_cluster_id                = clustercontrol_db_cluster.this.id
#   db_backup_method             = "pg_basebackup"
#   db_backup_dir                = var.db_backup_dir
#   db_backup_subdir             = var.db_backup_subdir
#   db_backup_encrypt            = var.db_backup_encrypt
#   db_backup_host               = var.db_backup_host
#   db_backup_storage_controller = var.db_backup_storage_controller
#   db_backup_compression        = var.db_backup_compression
#   db_backup_compression_level  = var.db_backup_compression_level
#   db_backup_retention          = var.db_backup_retention
# }

resource "clustercontrol_db_cluster_backup" "full-03-28-2024_1" {
  depends_on                   = [clustercontrol_db_cluster.this]
  db_cluster_id                = clustercontrol_db_cluster.this.id
  db_backup_method             = "pgbackrestfull"
  db_backup_dir                = var.db_backup_dir
  db_backup_subdir             = var.db_backup_subdir
  db_backup_encrypt            = var.db_backup_encrypt
  db_backup_host               = var.db_backup_host
  db_backup_storage_controller = var.db_backup_storage_controller
  db_backup_compression        = var.db_backup_compression
  db_backup_compression_level  = var.db_backup_compression_level
  db_backup_retention          = var.db_backup_retention
}

resource "clustercontrol_db_cluster_backup" "full-03-28-2024_1_1" {
  depends_on                   = [clustercontrol_db_cluster.this]
  db_cluster_id                = clustercontrol_db_cluster.this.id
  db_backup_method             = "pgbackrestincr"
  db_backup_dir                = var.db_backup_dir
  db_backup_subdir             = var.db_backup_subdir
  db_backup_encrypt            = var.db_backup_encrypt
  db_backup_host               = var.db_backup_host
  db_backup_storage_controller = var.db_backup_storage_controller
  db_backup_compression        = var.db_backup_compression
  db_backup_compression_level  = var.db_backup_compression_level
  db_backup_retention          = var.db_backup_retention
}

# resource "clustercontrol_db_cluster_maintenance" "server-upgrade-03232024" {
#   depends_on          = [clustercontrol_db_cluster.this]
#   db_cluster_id       = clustercontrol_db_cluster.this.id
#   db_maint_start_time = "Mar-27-2024T22:00"
#   db_maint_stop_time  = "Mar-28-2024T23:30"
#   db_maint_reason     = "Hardware refresh March 27, 2024"
# }

# resource "clustercontrol_db_cluster_backup" "full-dump-03-27-2024_1" {
#   depends_on                   = [clustercontrol_db_cluster.this]
#   db_cluster_id                = clustercontrol_db_cluster.this.id
#   db_backup_method             = "pgdumpall"
#   db_backup_dir                = var.db_backup_dir
#   db_backup_subdir             = var.db_backup_subdir
#   db_backup_encrypt            = var.db_backup_encrypt
#   db_backup_host               = var.db_backup_host
#   db_backup_storage_controller = var.db_backup_storage_controller
#   db_backup_compression        = var.db_backup_compression
#   db_backup_compression_level  = var.db_backup_compression_level
#   db_backup_retention          = var.db_backup_retention
# }
