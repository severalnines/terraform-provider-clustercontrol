provider "cc" {
  cc_api_user          = var.cc_api_user
  cc_api_user_password = var.cc_api_user_password
  cc_api_url           = var.cc_api_url
}


locals {
  is_db_create = (!var.db_cluster_import ? var.db_cluster_create : false)
  is_db_import = var.db_cluster_import
}

resource "cc_db_cluster" "this" {
  db_cluster_create        = true
  db_cluster_import        = false
  db_cluster_name          = "mydbcluster"
  db_cluster_type          = "replication"
  db_vendor                = "mariadb"
  db_version               = "10.11"
  db_admin_user_password   = "blah%blah"
  db_port                  = var.db_port
  db_data_directory        = var.db_data_directory
  disable_firewall         = var.disable_firewall
  db_install_software      = var.db_install_software
  ssh_user                 = var.ssh_user
  ssh_user_password        = var.ssh_user_password
  ssh_key_file             = var.ssh_key_file
  ssh_port                 = var.ssh_port
  db_tags                  = ["terra-deploy"]
  # db_tags                  = []
  # db_tags                  = var.db_tags

  db_host {
    hostname = "test-primary"
    # hostname_data = "foo"
    # hostname_internal = "foo"
    # port = "foo"
    # config_file = "foo"
    # data_dir = "foo"
  }
  # db_host {
  #   hostname          = "dev-replica"
  #   hostname_data     = "hnd-foo"
  #   hostname_internal = "hni-foo"
  #   port              = "p-foo"
  #   # config_file       = "cf-foo"
  #   data_dir          = "dd-foo"
  # }

  # db_topology {
  #   primary = "dev-primary"
  #   replica = "dev-replica"
  # }

  # timeouts = {
  #   create = lookup(var.timeouts, "create", null)
  #   import = lookup(var.timeouts, "import", null)
  #   delete = lookup(var.timeouts, "delete", null)
  #   update = lookup(var.timeouts, "update", null)
  # }

}
