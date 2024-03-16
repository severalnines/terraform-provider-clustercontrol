terraform {
  required_version = ">= 1.0"

  required_providers {
    cc = {
      source = "severalnines.com/severalnines/clustercontrol"
      version = ">= 0.0.1"
    }
    
  }
}