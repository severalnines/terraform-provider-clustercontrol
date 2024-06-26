provider "clustercontrol" {
  cc_api_user          = var.cc_api_user
  cc_api_user_password = var.cc_api_user_password
  cc_api_url           = var.cc_api_url
}

locals {
  is_db_create = (!var.db_cluster_import ? var.db_cluster_create : false)
  is_db_import = var.db_cluster_import
}


resource "clustercontrol_db_cluster" "mydatabase" {
  db_cluster_create        = true
  db_cluster_import        = false
  db_cluster_name          = "mydatabase"
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
  #  ssh_key_file             = var.ssh_key_file
  ssh_key_file             = var.SSH_KEY_FILE
  ssh_port                 = var.ssh_port
  db_tags                  = ["terra-deploy"]

  db_host {
    hostname          = join(",", aws_instance.project-iac.*.public_ip)
    hostname_internal = join(",", aws_instance.project-iac.*.private_ip)
  }

}
