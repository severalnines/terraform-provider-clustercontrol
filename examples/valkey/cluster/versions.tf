terraform {
  required_providers {
    clustercontrol = {
      source = "severalnines/clustercontrol"
      version = ">=0.2.21"
    }
  }
}
# terraform {
#   required_version = ">= 1.0"
#   required_providers {
#     clustercontrol = {
#       source = "severalnines.com/severalnines/clustercontrol"
#       version = ">= 0.2.21"
#     }
#   }
# }