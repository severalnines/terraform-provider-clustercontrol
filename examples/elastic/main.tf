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
  db_cluster_create              = true
  db_cluster_import              = false
  db_cluster_name                = "mydbcluster"
  db_cluster_type                = "elasticsearch"
  db_vendor                      = "elastic"
  db_version                     = "8.3.1"
  db_admin_username              = "esadmin"
  db_admin_user_password         = "blah%blah"
  db_auto_recovery               = true
  db_elasticsearch_http_port     = var.db_elasticsearch_http_port
  db_elasticsearch_transfer_port = var.db_elasticsearch_transfer_port
  db_data_directory              = var.db_data_directory
  db_snapshot_location           = var.db_snapshot_location
  db_snapshot_repository         = var.db_snapshot_repository
  db_snapshot_repository_type    = var.db_snapshot_repository_type
  db_snapshot_storage_host       = "test-primary-1"
  disable_firewall               = var.disable_firewall
  disable_selinux                = var.disable_selinux
  db_enable_uninstall            = var.db_enable_uninstall
  db_install_software            = var.db_install_software
  db_deploy_agents               = var.db_deploy_agents
  db_enable_ssl                  = var.db_enable_ssl
  ssh_user                       = var.ssh_user
  ssh_user_password              = var.ssh_user_password
  ssh_key_file                   = var.ssh_key_file
  ssh_port                       = var.ssh_port
  db_tags                        = ["terra-deploy"]

  db_host {
    hostname = "test-primary-1"
    roles    = "master-data"
    # hostname_data = "foo"
    # hostname_internal = "foo"
    # port = "foo"
  }

  # db_host {
  #   hostname = "test-primary-2"
  #   roles    = "master-data"
  #   # hostname_data     = "hnd-foo"
  #   # hostname_internal = "hni-foo"
  #   # port              = "p-foo"
  # }

  # db_host {
  #   hostname = "test-primary-3"
  #   roles    = "master-data"
  #   # hostname_data     = "hnd-foo"
  #   # hostname_internal = "hni-foo"
  #   # port              = "p-foo"
  # }

  # db_host {
  #   hostname         = "test-primary-3"
  #   # hostname_data     = "hnd-foo"
  #   # hostname_internal = "hni-foo"
  #   # port              = "p-foo"
  # }

  # timeouts = {
  #   create = lookup(var.timeouts, "create", null)
  #   import = lookup(var.timeouts, "import", null)
  #   delete = lookup(var.timeouts, "delete", null)
  #   update = lookup(var.timeouts, "update", null)
  # }

}

# resource "clustercontrol_db_cluster_backup_schedule" "daily-snap" {
#   depends_on             = [clustercontrol_db_cluster.this]
#   db_backup_sched_title  = "Daily snapshot"
#   db_backup_sched_time   = "TZ=UTC 0 0 * * *"
#   db_cluster_id          = clustercontrol_db_cluster.this.id
#   db_backup_method       = ""
#   db_backup_retention    = var.db_backup_retention
#   db_snapshot_repository = var.db_snapshot_repository
# }

# resource "clustercontrol_db_cluster_backup" "snap-1" {
#   depends_on             = [clustercontrol_db_cluster.this]
#   db_cluster_id          = clustercontrol_db_cluster.this.id
#   db_backup_method       = ""
#   db_snapshot_repository = var.db_snapshot_repository
#   db_backup_retention    = var.db_backup_retention
# }

# resource "clustercontrol_db_cluster_maintenance" "server-upgrade-03232024" {
#   depends_on          = [clustercontrol_db_cluster.this]
#   db_cluster_id       = clustercontrol_db_cluster.this.id
#   db_maint_start_time = "Mar-27-2024T22:00"
#   db_maint_stop_time  = "Mar-28-2024T23:30"
#   db_maint_reason     = "Hardware refresh March 27, 2024"
# }
