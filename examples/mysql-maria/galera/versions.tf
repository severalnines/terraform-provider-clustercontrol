terraform {
  required_version = ">= 1.0"

  required_providers {
    clustercontrol = {
      source = "severalnines.com/severalnines/clustercontrol"
      version = ">= 0.1.0"
    }
    
  }
}