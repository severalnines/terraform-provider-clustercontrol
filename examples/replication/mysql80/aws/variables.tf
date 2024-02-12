variable "API_USER" {
  description = "CMON API user's username"
  type        = string
  sensitive   = true
}

variable "CMON API_USER_PW" {
  description = "API user's password"
  type        = string
  sensitive   = true
}

variable "SSH_KEY_FILE" {
  description = "Path to SSH Key file (e.g /home/user/.ssh/id_rsa)"
  type        = string
  sensitive   = true
}

variable "SSH_PUBLIC_KEY" {
  description = "Content of id_rsa.pub (starting with 'ssh-rsa AAAA')"
  type        = string
  sensitive   = true
}

variable "AWS_CREDENTIALS_FILE" {
  description = "Path to AWS Credentialsfile (e.g /home/user/.aws/credentials)"
  type        = string
  sensitive   = true
}

variable "CONTROLLER_URL" {
  description = "ClusterControl controller coordinates"
  type        = string
  default = "https://127.0.0.1:9501/v2"
}
