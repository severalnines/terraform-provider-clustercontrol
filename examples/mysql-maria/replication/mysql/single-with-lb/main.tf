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
  db_cluster_type          = "replication"
  db_vendor                = "percona"
  db_version               = "8.0"
  db_admin_user_password   = "blah%blah"
  db_port                  = var.db_port
  db_data_directory        = var.db_data_directory
  disable_firewall         = var.disable_firewall
  db_install_software      = var.db_install_software
  db_semi_sync_replication = var.db_semi_sync_replication
  ssh_user                 = var.ssh_user
  ssh_user_password        = var.ssh_user_password
  ssh_key_file             = var.ssh_key_file
  ssh_port                 = var.ssh_port
  db_tags                  = ["terra-deploy"]

  db_host {
    hostname = "test-primary"
    # hostname_data = "foo"
    # hostname_internal = "foo"
    # port = "foo"
    # config_file = "foo"
    # data_dir = "foo"
  }

}

resource "clustercontrol_db_load_balancer" "this" {
  depends_on = [clustercontrol_db_cluster.this]

  db_lb_create                = true
  db_lb_import                = false
  # db_cluster_id               = [clustercontrol_db_cluster.this.db_cluster_id]
  db_cluster_id               = clustercontrol_db_cluster.this.id
  db_lb_type                  = "proxysql"
  db_lb_version               = var.db_lb_version
  db_lb_admin_username        = var.db_lb_admin_username
  db_lb_admin_user_password   = "blah%blah"
  db_lb_monitor_username      = var.db_lb_monitor_username
  db_lb_monitor_user_password = "blah%blah"
  db_lb_port                  = var.db_lb_port
  db_lb_use_clustering        = var.db_lb_use_clustering
  db_lb_use_rw_splitting      = var.db_lb_use_rw_splitting
  db_lb_install_software      = var.db_lb_install_software
  disable_firewall            = var.disable_firewall
  ssh_user                    = var.ssh_user
  ssh_user_password           = var.ssh_user_password
  ssh_key_file                = var.ssh_key_file
  ssh_port                    = var.ssh_port

  db_host {
    hostname = "test-primary"
    port     = clustercontrol_db_cluster.this.db_port
  }

  db_my_host {
    hostname = "test-primary-2"
    # port     = "3306"
  }
}