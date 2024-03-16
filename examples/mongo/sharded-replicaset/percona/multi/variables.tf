variable "cc_api_user" {
  description = "API user"
  type        = string
  sensitive   = true
}

variable "cc_api_user_password" {
  description = "API user's password"
  type        = string
  sensitive   = true
}

variable "cc_api_url" {
  description = "ClusterControl controller coordinates"
  type        = string
}

variable "db_cluster_create" {
  description = "Whether to create this resource or not?"
  type        = bool
  default     = false
}

variable "db_cluster_import" {
  description = "Whether to create this resource or not?"
  type        = bool
  default     = false
}

variable "db_cluster_name" {
  description = "The name of the database cluster"
  type        = string
  default     = null
  # validation {
  #   condition     = length(var.db_cluster_name) > 0
  #   error_message = "The db_cluster_name value must not be an empty string."
  # }
}

variable "db_cluster_type" {
  description = "Type of cluster - MySQL Replication, Galera, Postgres, MongoDB, etc"
  type        = string
  default     = null
  # validation {
  #   condition     = length(var.db_cluster_type) > 0
  #   error_message = "The db_cluster_type value must not be an empty string."
  # }
}

variable "db_vendor" {
  description = "Database vendor - Oracle, Percona, MariaDB, Mongo/10Gen, Microsoft, etc"
  type        = string
  default     = null
  # validation {
  #   condition     = length(var.db_vendor) > 0
  #   error_message = "The db_vendor value must not be an empty string."
  # }
}

variable "db_version" {
  description = "The database version to use"
  type        = string
  default     = null
  # validation {
  #   condition     = length(var.db_version) > 0
  #   error_message = "The db_version value must not be an empty string."
  # }
}

variable "db_admin_username" {
  description = "Name for the admin/root user for the database"
  type        = string
  default     = "dbadminusr"
  # validation {
  #   condition     = length(var.db_admin_username) > 0
  #   error_message = "The db_admin_username value must not be an empty string."
  # }
}

variable "db_admin_user_password" {
  description = "Password for the admin/root user for the database. Note that this may show up in logs, and it will be stored in the state file"
  type        = string
  default     = null
  sensitive   = true
  # validation {
  #   condition     = length(var.db_admin_user_password) > 0
  #   error_message = "The db_admin_user_password value must not be an empty string."
  # }
}

variable "db_port" {
  description = "The port on which the DB accepts connections"
  type        = string
  default     = null
  # validation {
  #   condition     = length(var.db_port) > 0
  #   error_message = "The db_port value must not be an empty string."
  # }
}

variable "db_data_directory" {
  description = "TODO"
  type        = string
  default     = null
  # validation {
  #   condition     = length(var.db_data_directory) > 0
  #   error_message = "The db_data_directory value must not be an empty string."
  # }
}

# variable "db_config_template" {
#   description = "TODO"
#   type        = string
#   default     = null
#   # validation {
#   #   condition     = length(var.db_config_template) > 0
#   #   error_message = "The db_config_template value must not be an empty string."
#   # }
# }

variable "disable_firewall" {
  description = "TODO"
  type        = bool
  nullable    = false
  default     = true
}

variable "db_install_software" {
  description = "TODO"
  type        = bool
  nullable    = false
  default     = true
}

# variable "db_sync_replication" {
#   description = "TODO"
#   type        = bool
#   default     = false
# }

variable "db_semi_sync_replication" {
  description = "TODO"
  type        = bool
  default     = false
}

variable "ssh_user" {
  description = "TODO"
  type        = string
  default     = "ubuntu"
  # default     = "root"
  validation {
    condition     = length(var.ssh_user) > 0
    error_message = "The ssh_user value must not be an empty string."
  }
}

variable "ssh_user_password" {
  description = "Sudo user password"
  type        = string
  default     = null
  # validation {
  #   condition     = length(var.ssh_user_password) > 0
  #   error_message = "The ssh_user_password value must not be an empty string."
  # }
}

variable "ssh_key_file" {
  description = "TODO"
  type        = string
  default     = "/home/ubuntu/.ssh/id_rsa"
  # default     = "/root/.ssh/id_rsa"
  validation {
    condition     = length(var.ssh_key_file) > 0
    error_message = "The ssh_key_file value must not be an empty string."
  }
}

variable "ssh_port" {
  description = "TODO"
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
    # config_file       = string
  }))
  default = null
}

# variable "db_config_server" {
#   description = "The list of nodes/hosts that make up the cluster"
#   type = list(object({
#     hostname          = string
#     hostname_data     = string
#     hostname_internal = string
#     port              = string
#   }))
#   default = null
# }


variable "db_mongos_server" {
  description = "The list of nodes/hosts that make up the cluster"
  type = list(object({
    hostname          = string
    hostname_data     = string
    hostname_internal = string
    port              = string
  }))
  default = null
}

# TODO: need to fix this....
# variable "db_replica_set" {
#   description = "The list of nodes/hosts that make up the cluster"
#   type = list(object({
#     hostname          = string
#     hostname_data     = string
#     hostname_internal = string
#     port              = string
#   }))
#   default = null
# }

variable "db_topology" {
  description = "The list of nodes/hosts that make up the cluster"
  type = list(object({
    primary = string
    replica = string
  }))
  default = null
}

variable "db_tags" {
  description = "A mapping of tags to assign to all resources"
  type        = set(string)
  default     = []
  # type        = list(string)
}

variable "timeouts" {
  description = "Updated Terraform resource management timeouts. Applies to permit resource management times"
  type        = map(string)
  default     = {}
}

# Future stuff ...

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
