
provider "cc" {
  #controller_url = var.CONTROLLER_URL
  controller_url = "https://127.0.0.1:9501/v2"
  user = var.API_USER
  password = var.API_USER_PW
}


resource "cc_mysql_maria_cluster" "mydatabase" {
  cluster_name = "mydatabase"
  database_type = "mysql"
  database_vendor = "percona"
  database_version = "8.0"
  database_topology = "replication"
  ssh_key_file = var.SSH_KEY_FILE
  ssh_user = "ubuntu"
  install_software = "false"
  primary_database_host = join(",",aws_instance.project-iac.*.public_ip)
  hostname_internal = join(",",aws_instance.project-iac.*.private_ip)  			
}
