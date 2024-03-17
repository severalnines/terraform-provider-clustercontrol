# Deploy a datastore on AWS
cd terraform-provider-clustercontrol/examples/replication/mysql80/aws
## Update aws.shared_credentials_files in infra.tf
```
  shared_credentials_files = ["~/.aws/credentials"]
```
## Update infra.tf
- vpc
- region
- ami
- subnet
- secgroupname
Please note that the VMs must be reachable from ClusterControl.

## Datastore configuration / Terraform schema
| Configuration Parameter     | Value                                              |
|-----------------------------|----------------------------------------------------|
| `db_cluster_name`           | `NAME_OF_DATABASE`                                 |
| `db_vendor`                 | `percona`                                          |
| `db_version`                | `8.0`                                              |
| `db_cluster_type`           | `replication`                                      |
| `ssh_key_file`              | `SSH_KEYFILE, e. /home/user/.ssh/id_rsa`           |
| `ssh_user`                  | default: `"ubuntu"`                                |
| `db_host.hostname`          | `join(",", aws_instance.project-iac.*.public_ip)`  |
| `db_host.hostname_internal` | `join(",", aws_instance.project-iac.*.private_ip)` |

Please note that installing software make take time and the operation may timeout.
