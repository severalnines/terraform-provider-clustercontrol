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
  db_cluster_create           = true
  db_cluster_import           = false
  db_cluster_name             = "mydbcluster"
  db_cluster_type             = "mongo"
  db_vendor                   = "percona"
  db_version                  = "6.0"
  db_admin_username           = "mongoadmin"
  db_admin_user_password      = "blah%blah"
  db_auto_recovery            = true
  db_mongo_port               = var.db_mongo_port
  db_mongo_config_server_port = var.db_mongo_config_server_port
  db_data_directory           = var.db_data_directory
  disable_firewall            = var.disable_firewall
  disable_selinux             = var.disable_selinux
  db_enable_uninstall         = var.db_enable_uninstall
  db_install_software         = var.db_install_software
  db_deploy_agents            = var.db_deploy_agents
  db_enable_ssl               = var.db_enable_ssl
  db_mongo_auth_db            = var.db_mongo_auth_db
  db_enable_pbm_agent         = true
  db_pbm_backup_dir           = "/nfs/mongobackup"
  ssh_user                    = var.ssh_user
  ssh_user_password           = var.ssh_user_password
  ssh_key_file                = var.ssh_key_file
  ssh_port                    = var.ssh_port
  db_tags                     = ["terra-deploy"]

  db_config_server {
    rs = "rs_config"
    member {
      hostname = "test-primary"
    }
  }

  db_mongos_server {
    hostname = "test-primary"
  }

  db_replica_set {
    rs = "rs0"
    member {
      hostname = "test-primary-2"
    }
#    member {
#      hostname = "test-primary-3"
#    }
#     member {
#       hostname = "test-primary-4"
#     }
  }

}

# resource "clustercontrol_db_cluster_backup" "full-03-27-2024_1" {
#   depends_on       = [clustercontrol_db_cluster.this]
#   db_cluster_id    = clustercontrol_db_cluster.this.id
#   db_backup_method = "percona-backup-mongodb"
#   # db_backup_dir    = var.db_backup_dir
#   # db_backup_subdir             = var.db_backup_subdir
#   # db_backup_encrypt            = var.db_backup_encrypt
#   # db_backup_host               = var.db_backup_host
#   # db_backup_storage_host       = var.db_backup_storage_host
#   # db_enable_backup_failover    = var.db_enable_backup_failover
#   # db_backup_failover_host      = var.db_backup_failover_host
#   # db_backup_storage_controller = var.db_backup_storage_controller
#   # db_backup_compression        = var.db_backup_compression
#   # db_backup_compression_level  = var.db_backup_compression_level
#   db_backup_retention = var.db_backup_retention
# }

# resource "clustercontrol_db_cluster_backup_schedule" "daily-full-1" {
#   depends_on                   = [clustercontrol_db_cluster.this]
#   db_backup_sched_title        = "Daily full"
#   db_backup_sched_time         = "TZ=UTC 0 0 * * *"
#   db_cluster_id                = clustercontrol_db_cluster.this.id
#   db_backup_method             = "percona-backup-mongodb"
##   db_backup_dir                = var.db_backup_dir
##   db_backup_subdir             = var.db_backup_subdir
##   db_backup_encrypt            = var.db_backup_encrypt
##   db_backup_host               = var.db_backup_host
##   db_backup_storage_controller = var.db_backup_storage_controller
##   db_backup_compression        = var.db_backup_compression
##   db_backup_compression_level  = var.db_backup_compression_level
##   db_backup_retention          = var.db_backup_retention
# }

# resource "clustercontrol_db_cluster_maintenance" "server-upgrade-03232024" {
#   depends_on          = [clustercontrol_db_cluster.this]
#   db_cluster_id       = clustercontrol_db_cluster.this.id
#   db_maint_start_time = "Mar-27-2024T22:00"
#   db_maint_stop_time  = "Mar-28-2024T23:30"
#   db_maint_reason     = "Hardware refresh March 27, 2024"
# }
