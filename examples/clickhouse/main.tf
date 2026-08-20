provider "clustercontrol" {
  cc_api_user          = var.cc_api_user
  cc_api_user_password = var.cc_api_user_password
  cc_api_url            = var.cc_api_url
}

locals {
  is_db_create = (!var.db_cluster_import ? var.db_cluster_create : false)
  is_db_import = var.db_cluster_import
}

resource "clustercontrol_db_cluster" "this" {
  db_cluster_create         = true
  db_cluster_import         = false
  db_cluster_name           = "mydbcluster"
  db_cluster_type           = "clickhouse"
  db_vendor                 = "clickhouse"
  db_version                = "24.8"
  db_admin_username         = "chadmin"
  db_admin_user_password    = "blah%blah"
  db_auto_recovery          = true
  db_clickhouse_native_port = var.db_clickhouse_native_port
  db_clickhouse_keeper_port = var.db_clickhouse_keeper_port
  db_data_directory         = var.db_data_directory
  disable_firewall          = var.disable_firewall
  disable_selinux           = var.disable_selinux
  db_enable_uninstall       = var.db_enable_uninstall
  db_install_software       = var.db_install_software
  db_deploy_agents          = var.db_deploy_agents
  db_enable_ssl             = true # SSL is mandatory for ClickHouse - not user-configurable
  ssh_user                  = var.ssh_user
  ssh_user_password         = var.ssh_user_password
  ssh_key_file              = var.ssh_key_file
  ssh_port                  = var.ssh_port
  db_tags                   = ["terra-deploy"]

  # db_host {
  #   hostname = "clickhouse-1"
  #   roles    = "keeper"
  #   # shard    = "01" # keeper can't have shard!!!
  # }

  db_host {
    hostname = "clickhouse-1"
    roles    = "replica"
    shard    = "01"
  }
  db_host {
    hostname = "clickhouse-2"
    roles    = "replica"
    shard    = "01"
  }
  db_host {
    hostname = "clickhouse-3"
    roles    = "replica"
    shard    = "01"
  }

  # db_host {
  #   hostname = "test-primary-3"
  #   roles    = "replica"
  #   shard    = "01"
  # }
  #
  # db_host {
  #   hostname = "test-primary-4"
  #   roles    = "keeper"
  # }
  #
  # db_host {
  #   hostname = "test-primary-5"
  #   roles    = "keeper"
  # }
  #
  # db_host {
  #   hostname = "test-primary-6"
  #   roles    = "replica"
  # }

  # timeouts = {
  #   create = lookup(var.timeouts, "create", null)
  #   import = lookup(var.timeouts, "import", null)
  #   delete = lookup(var.timeouts, "delete", null)
  #   update = lookup(var.timeouts, "update", null)
  # }
}

# resource "clustercontrol_db_cluster_backup_schedule" "daily-full" {
#   depends_on            = [clustercontrol_db_cluster.this]
#   db_backup_sched_title = "Daily full backup"
#   db_backup_sched_time  = "TZ=UTC 0 0 * * *"
#   db_cluster_id         = clustercontrol_db_cluster.this.id
#   db_backup_method      = "clickhouse-native"
#   db_backup_retention   = var.db_backup_retention
# }

# Tested (Aug 20, 2026)
# resource "clustercontrol_db_cluster_backup" "full-1" {
#   depends_on          = [clustercontrol_db_cluster.this]
#   db_cluster_id       = clustercontrol_db_cluster.this.id
#   db_backup_method    = "clickhouse-native"
#   db_backup_retention = var.db_backup_retention
# }

# Tested (Aug 20, 2026)
# resource "clustercontrol_db_cluster_backup" "incr-1" {
#   depends_on          = [clustercontrol_db_cluster.this]
#   db_cluster_id       = clustercontrol_db_cluster.this.id
#   db_backup_method    = "clickhouse-native-incr"
#   db_backup_retention = var.db_backup_retention
# }

# resource "clustercontrol_db_cluster_maintenance" "server-upgrade-03232024" {
#   depends_on          = [clustercontrol_db_cluster.this]
#   db_cluster_id       = clustercontrol_db_cluster.this.id
#   db_maint_start_time = "Mar-27-2024T22:00"
#   db_maint_stop_time  = "Mar-28-2024T23:30"
#   db_maint_reason     = "Hardware refresh March 27, 2024"
# }
