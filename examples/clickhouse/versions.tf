
# The below block is for the tf provider from TF registry...
terraform {
  required_providers {
    clustercontrol = {
      source = "severalnines/clustercontrol"
      version = ">=0.2.23"
    }
  }
}

# The below is for dev testing locally built tf provider...
# terraform {
#   required_version = ">= 1.0"
#
#   required_providers {
#     clustercontrol = {
#       source  = "severalnines.com/severalnines/clustercontrol"
#       version = ">= 0.2.23"
#     }
#   }
# }
