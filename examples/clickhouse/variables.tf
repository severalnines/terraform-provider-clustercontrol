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
  description = "Database vendor - oracle, percona, mariadb, 10gen, microsoft, redis, elasticsearch, clickhouse, etc"
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

variable "db_clickhouse_native_port" {
  description = "The port on which ClickHouse will accept native TCP protocol client connections"
  type        = string
  default     = "9440"
}

variable "db_clickhouse_keeper_port" {
  description = "The client port for the embedded ClickHouse Keeper instance used for replica coordination"
  type        = string
  default     = "9281"
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

variable "db_tags" {
  description = "Tags to associate with a DB cluster. The tags are only relevant in the ClusterControl domain."
  type        = set(string)
  default     = []
}

variable "db_deploy_agents" {
  description = "Automatically deploy prometheus and other relevant agents after setting up the intial DB cluster."
  type        = bool
  default     = true
}

variable "db_auto_recovery" {
  description = "Have cluster auto-recovery on (or off)"
  type        = bool
  default     = true
}

# --------------------------
# Backup variables ...
# --------------------------

variable "db_backup_method" {
  description = "Which backup to use for ClickHouse - \"clickhouse-native\" (full) or \"clickhouse-native-incr\" (incremental)"
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

# --------------------------------------------
# Backup schedule variables ...
# --------------------------------------------

variable "db_backup_sched_title" {
  description = "A title for the backup schedule (e.g., Daily full, Hourly incremental, etc)"
  type        = string
  default     = "Sample backup schedule title"
}

variable "db_backup_sched_time" {
  description = "The time to kick off a backup (e.g. 'TZ=UTC 0 0 * * *')"
  type        = string
  default     = null
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

variable "db_cluster_id" {
  description = "The ID of the DB cluster"
  type        = string
  default     = null
}

variable "timeouts" {
  description = "Updated Terraform resource management timeouts. Applies to permit resource management times"
  type        = map(string)
  default     = {}
}
