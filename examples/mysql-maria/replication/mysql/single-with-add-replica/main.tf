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
  db_cluster_name          = "mydbcluster-1"
  db_cluster_type          = "replication"
  db_vendor                = "percona"
  db_version               = "8.0"
  db_admin_user_password   = "blah%blah"
  db_port                  = var.db_port
  db_data_directory        = var.db_data_directory
  db_auto_recovery         = true
  disable_firewall         = var.disable_firewall
  db_install_software      = var.db_install_software
  db_semi_sync_replication = var.db_semi_sync_replication
  ssh_user                 = var.ssh_user
  ssh_user_password        = var.ssh_user_password
  ssh_key_file             = var.ssh_key_file
  ssh_port                 = var.ssh_port
  # db_tags                  = ["terra-deploy"]
  db_tags = ["terra-deploy", "terra-foo"]

  db_host {
    hostname = "test-primary"
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

# resource "clustercontrol_db_cluster_backup" "full-03-23-2024_1" {
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

# resource "clustercontrol_db_cluster_backup" "full-03-23-2024_2" {
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
#   db_maint_start_time = "Mar-23-2024T22:00"
#   db_maint_stop_time  = "Mar-23-2024T23:30"
#   db_maint_reason     = "Hardware refresh March 23, 2024"
# }
