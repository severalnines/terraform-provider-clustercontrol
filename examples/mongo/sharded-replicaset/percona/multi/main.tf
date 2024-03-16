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
  db_cluster_create      = true
  db_cluster_import      = false
  db_cluster_name        = "mydbcluster"
  db_cluster_type        = "mongodb"
  db_vendor              = "percona"
  db_version             = "6.0"
  db_admin_username      = "mongoadmin"
  db_admin_user_password = "blah%blah"
  db_port                = var.db_port
  db_data_directory      = var.db_data_directory
  disable_firewall       = var.disable_firewall
  db_install_software    = var.db_install_software
  ssh_user               = var.ssh_user
  ssh_user_password      = var.ssh_user_password
  ssh_key_file           = var.ssh_key_file
  ssh_port               = var.ssh_port
  db_tags                = ["terra-deploy"]

  db_config_server {
    rs = "replica_set_config"
    member {
      hostname = "test-primary"
      port     = "27019"
    }
  }

  db_mongos_server {
    hostname = "test-primary"
    port     = "27107"
  }

  db_replica_set {
    rs = "rs0"
    member {
      hostname = "test-primary-2"
      port     = "27107"
    }
    member {
      hostname = "test-primary-3"
      port     = "27107"
    }
  }

  # timeouts = {
  #   create = lookup(var.timeouts, "create", null)
  #   import = lookup(var.timeouts, "import", null)
  #   delete = lookup(var.timeouts, "delete", null)
  #   update = lookup(var.timeouts, "update", null)
  # }

}
