variable "API_USER" {
  description = "API user's username"
  type        = string
  sensitive   = true
}

variable "API_USER_PW" {
  description = "API user's password"
  type        = string
  sensitive   = true
}

variable "CONTROLLER_URL" {
  description = "ClusterControl controller coordinates"
  type        = string
  default = "https://127.0.0.1:9501/v2"
}
