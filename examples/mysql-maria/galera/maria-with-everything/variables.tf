# --------------------------------------------
# Database (DB) Cluster specific variables ...
# --------------------------------------------
variable "cc_api_user" {
  description = "ClusterControl API user"
  type        = string
  sensitive   = true
}

variable "cc_api_user_password" {
  description = "API user's password"
  type        = string
  sensitive   = true
}

variable "cc_api_url" {
  description = "ClusterControl controller url e.g. (https://cc-host:9501/v2)"
  type        = string
}

# --------------------------------------------
# Database (DB) Cluster specific variables ...
# --------------------------------------------
variable "db_cluster_create" {
  description = "Whether to create this resource or not?"
  type        = bool
  default     = false
}

variable "db_cluster_import" {
  description = "Whether to import this resource or not?"
  type        = bool
  default     = false
}

variable "db_cluster_name" {
  description = "The name of the database cluster"
  type        = string
  default     = null
}

variable "db_cluster_type" {
  description = "Type of cluster - replication, galera, postgresql_single (single is a misnomer), etc"
  type        = string
  default     = null
}

variable "db_vendor" {
  description = "Database vendor - oracle, percona, mariadb, 10gen, microsoft, redis, elasticsearch, for postgresql it is `default` etc"
  type        = string
  default     = null
}

variable "db_version" {
  description = "The database version"
  type        = string
  default     = null
}

variable "db_admin_username" {
  description = "Name for the admin/root user for the database"
  type        = string
  default     = "dbadminusr"
}

variable "db_admin_user_password" {
  description = "Password for the admin/root user for the database. Note that this may show up in logs, and it will be stored in the state file"
  type        = string
  default     = null
  sensitive   = true
}

variable "db_port" {
  description = "The port on which the DB will accepts connections"
  type        = string
  default     = null
}

variable "db_data_directory" {
  description = "The data directory for the database data files. If not set explicily, the default for the respective DB vendor will be chosen"
  type        = string
  default     = null
}

variable "disable_firewall" {
  description = "Disable firewall on the host OS when installing DB packages."
  type        = bool
  nullable    = false
  default     = true
}

variable "disable_selinux" {
  description = "Disable SELinux on the host OS when installing DB packages."
  type        = bool
  nullable    = false
  default     = true
}

variable "db_install_software" {
  description = "Install DB packages from respective repos"
  type        = bool
  nullable    = false
  default     = true
}

variable "db_enable_uninstall" {
  description = "When removing DB cluster from ClusterControl, enable uinstalling DB packages."
  type        = bool
  nullable    = false
  default     = true
}

variable "db_semi_sync_replication" {
  description = "Semi-synchronous replication for MySQL and MariaDB non-galera clusters"
  type        = bool
  default     = false
}

variable "ssh_user" {
  description = "The SSH user ClusterControl will use to SSH to the DB host from the ClusterControl host"
  type        = string
  default     = "ubuntu"
  validation {
    condition     = length(var.ssh_user) > 0
    error_message = "The ssh_user value must not be an empty string."
  }
}

variable "ssh_user_password" {
  description = "Sudo user's password. If sudo user doesn't have a password, leave this field blank"
  type        = string
  default     = null
}

variable "ssh_key_file" {
  description = "The path to the private key file for the Sudo user on the ClusterControl host"
  type        = string
  default     = "/home/ubuntu/.ssh/id_rsa"
  validation {
    condition     = length(var.ssh_key_file) > 0
    error_message = "The ssh_key_file value must not be an empty string."
  }
}

variable "ssh_port" {
  description = "The ssh port."
  type        = string
  default     = "22"
  validation {
    condition     = length(var.ssh_port) > 0
    error_message = "The ssh_port value must not be an empty string."
  }
}

variable "db_host" {
  description = "The list of nodes/hosts that make up the cluster"
  type = list(object({
    hostname          = string
    hostname_data     = string
    hostname_internal = string
    port              = string
    data_dir          = string
    sync_replication  = bool
  }))
  default = null
}

variable "db_topology" {
  description = "Only applicable to MySQL/MariaDB non-galera clusters. A way to specify Master and Slave(s). See examples."
  type = list(object({
    primary = string
    replica = string
  }))
  default = null
}

variable "db_tags" {
  description = "Tags to associate with a DB cluster. The tags are only relevant in the ClusterControl domain."
  type        = set(string)
  default     = []
}

variable "db_deploy_agents" {
  description = "Automatically deploy prometheus and other relevant agents after setting up the intial DB cluster."
  type        = bool
  default     = false
}

variable "db_auto_recovery" {
  description = "Have cluster auto-recovery on (or off)"
  type        = bool
  default     = true
}

# --------------------------
# Load balancer variables ...
# --------------------------

variable "db_lb_create" {
  description = "Whether to create this resource or not?"
  type        = bool
  default     = false
}

variable "db_lb_import" {
  description = "Whether to import this resource or not?"
  type        = bool
  default     = false
}

variable "db_cluster_id" {
  description = "The ID of the DB cluster"
  type        = string
  default     = null
}

variable "db_lb_type" {
  description = "The load balancer type (e.g., proxysql, haproxy, etc)"
  type        = string
  default     = "proxysql"
}

variable "db_lb_version" {
  description = "The load balancer version to use"
  type        = string
  default     = "2"
}

variable "db_lb_admin_username" {
  description = "The load balancer admin user"
  type        = string
  default     = "proxysql-admin"
}

variable "db_lb_admin_user_password" {
  description = "The load balancer admin user's password"
  type        = string
  sensitive   = true
  default     = null
}

variable "db_lb_monitor_username" {
  description = "The load balancer monitor user (only applicable to proxysql)"
  type        = string
  default     = "proxysql-monitor"
}

variable "db_lb_monitor_user_password" {
  description = "The load balancer monitor user's password"
  type        = string
  sensitive   = true
  default     = null
}

variable "db_lb_port" {
  description = "The load balancer port that it will accept connections on behalf of the database it is front-ending."
  type        = string
  default     = "6033"
}

variable "db_lb_admin_port" {
  description = "The load balancer port that it will accept connections to manage its configuraiton"
  type        = string
  default     = "6032"
}

variable "db_lb_use_clustering" {
  description = "Whether to use ProxySQL clustering or not. Only applicable to ProxySQL at this time"
  type        = bool
  default     = true
}

variable "db_lb_use_rw_splitting" {
  description = "Whether to Read/Write splitting for queries or not?"
  type        = bool
  default     = true
}

variable "db_lb_install_software" {
  description = "Whether to setup repos and subsequently install load balancer software or not?"
  type        = bool
  default     = true
}

variable "db_lb_enable_uninstall" {
  description = "When removing load balancer from ClusterControl, enable uinstalling its packages."
  type        = bool
  nullable    = false
  default     = true
}


variable "db_my_host" {
  description = "Details regarding the load balancer host"
  type = object({
    hostname          = string
    port              = string
  })
  default = null
}

# --------------------------
# Maintenance variables ...
# --------------------------

variable "db_maint_start_time" {
  description = "Maintenance start time. See examples for format"
  type        = string
  default     = null
}

variable "db_maint_stop_time" {
  description = "Maintenance stop time"
  type        = string
  default     = null
}

variable "db_maint_reason" {
  description = "Reason for maintenance"
  type        = string
  default     = null
}

# --------------------------
# Backup variables ...
# --------------------------

variable "db_backup_method" {
  description = "Which backup to use - mariabackup, xtrabackup, mysqldump, pbm, etc"
  type        = string
  default     = null
}

variable "db_backup_dir" {
  description = "Where in the filesystem to store the backups"
  type        = string
  default     = "/home/ubuntu/backups"
}

variable "db_backup_subdir" {
  description = "Subdirectory for this backup"
  type        = string
  default     = "BACKUP-%I"
}

variable "db_backup_storage_controller" {
  description = "Whether to store backups on ClusterControl host."
  type        = bool
  default     = false
}

variable "db_backup_encrypt" {
  description = "Option to encrypt backups taken by ClusterControl"
  type        = bool
  default     = true
}

variable "db_backup_host" {
  description = "Which host to take backup on. Primary, Standby, Auto - meaning let ClusterControl decide which host to select"
  type        = string
  default     = "auto"
}

variable "db_backup_compression" {
  description = "Whether to compress backups"
  type        = bool
  default     = true
}

variable "db_backup_compression_level" {
  description = "Compression level"
  type        = number
  default     = 6
}

variable "db_backup_retention" {
  description = "DB backup retentions period (days)"
  type        = number
  default     = 7
}

# --------------------------
# Future stuff ...
# --------------------------

variable "timeouts" {
  description = "Updated Terraform resource management timeouts. Applies to permit resource management times"
  type        = map(string)
  default     = {}
}

variable "db_monitoring_interval" {
  description = "The interval, in seconds, between points when Enhanced Monitoring metrics are collected for the DB instance. To disable collecting Enhanced Monitoring metrics, specify 0. The default is 0. Valid Values: 0, 1, 5, 10, 15, 30, 60."
  type        = number
  default     = 0
}

variable "db_maintenance_window" {
  description = "The window to perform maintenance in. Syntax: 'ddd:hh24:mi-ddd:hh24:mi'. Eg: 'Mon:00:00-Mon:03:00'"
  type        = string
  default     = null
}

variable "db_backup_retention_period" {
  description = "The days to retain backups for"
  type        = number
  default     = null
}

variable "db_backup_window" {
  description = "The daily time range (in UTC) during which automated backups are created if they are enabled. Example: '09:46-10:16'. Must not overlap with maintenance_window"
  type        = string
  default     = null
}

variable "db_delete_automated_backups" {
  description = "Specifies whether to remove automated backups immediately after the DB instance is deleted"
  type        = bool
  default     = true
}

variable "db_restore_to_point_in_time" {
  description = "Restore to a point in time (MySQL is NOT supported)"
  type        = map(string)
  default     = null
}