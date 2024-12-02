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
  db_cluster_create                = true
  db_cluster_import                = false
  db_cluster_name                  = "mydbcluster"
  db_cluster_type                  = "valkey-sharded"
  db_vendor                        = "valkey"
  db_version                       = "8"
  db_admin_user_password           = "blah2blah"
  db_auto_recovery                 = true
  db_redis_port                    = var.db_redis_port
  db_redis_bus_port                = var.db_redis_bus_port
  db_redis_node_timeout_ms         = var.db_redis_node_timeout_ms
  db_redis_replica_validity_factor = var.db_redis_replica_validity_factor
  db_data_directory                = var.db_data_directory
  disable_firewall                 = var.disable_firewall
  disable_selinux                  = var.disable_selinux
  db_enable_uninstall              = var.db_enable_uninstall
  db_install_software              = var.db_install_software
  db_deploy_agents                 = var.db_deploy_agents
  db_enable_ssl                    = var.db_enable_ssl
  ssh_user                         = "rocky"
  ssh_user_password                = var.ssh_user_password
  ssh_key_file                     = var.ssh_key_file
  ssh_port                         = var.ssh_port
  db_tags                          = ["terra-deploy"]

  db_host {
    hostname = "valkey-4"
    # hostname_data = "foo"
    # hostname_internal = "foo"
    host_role = "primary"
  }

  db_host {
    hostname = "valkey-5"
    # hostname_data     = "hnd-foo"
    # hostname_internal = "hni-foo"
    host_role = "replica"
  }

  # db_host {
  #   hostname = "valkey-3"
  #   # hostname_data     = "hnd-foo"
  #   # hostname_internal = "hni-foo"
  #   host_role = "replica"
  # }

  #   db_host {
  #     hostname         = "test-primary-5"
  #     # hostname_data     = "hnd-foo"
  #     # hostname_internal = "hni-foo"
  #   }

}

# resource "clustercontrol_db_cluster_backup" "full-1" {
#   depends_on                   = [clustercontrol_db_cluster.this]
#   db_cluster_id                = clustercontrol_db_cluster.this.id
#   db_backup_method             = ""
#   db_backup_dir                = var.db_backup_dir
#   db_backup_subdir             = var.db_backup_subdir
#   db_backup_encrypt            = var.db_backup_encrypt
#   db_backup_host               = var.db_backup_host
#   db_backup_storage_controller = var.db_backup_storage_controller
#   db_backup_compression        = var.db_backup_compression
#   db_backup_compression_level  = var.db_backup_compression_level
#   db_backup_retention          = var.db_backup_retention
# }

# resource "clustercontrol_db_cluster_backup_schedule" "daily-full" {
#   depends_on                   = [clustercontrol_db_cluster.this]
#   db_backup_sched_title        = "Daily full backup"
#   db_backup_sched_time         = "TZ=UTC 0 0 * * *"
#   db_cluster_id                = clustercontrol_db_cluster.this.id
#   db_backup_method             = ""
#   db_backup_dir                = var.db_backup_dir
#   db_backup_subdir             = var.db_backup_subdir
#   db_backup_encrypt            = var.db_backup_encrypt
#   db_backup_host               = var.db_backup_host
#   db_backup_storage_controller = var.db_backup_storage_controller
#   db_backup_compression        = var.db_backup_compression
#   db_backup_compression_level  = var.db_backup_compression_level
#   db_backup_retention          = var.db_backup_retention
# }

# resource "clustercontrol_db_cluster_maintenance" "server-upgrade-12042024" {
#   depends_on          = [clustercontrol_db_cluster.this]
#   db_cluster_id       = clustercontrol_db_cluster.this.id
#   db_maint_start_time = "Dec-04-2024T22:00"
#   db_maint_stop_time  = "Dec-04-2024T22:30"
#   db_maint_reason     = "Hardware refresh March 27, 2024"
# }
