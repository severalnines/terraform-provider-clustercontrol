# terraform-provider-clustercontrol
Terraform Provider for Severalnines ClusterControl



## How to build
```
make
make install
```
The provider is installed in
```
~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
```
See the Makefile for details.


## Deploy a datastore on AWS
cd terraform-provider-clustercontrol/examples/replication/mysql80/aws
### Update aws.shared_credentials_files in infra.tf
```
  shared_credentials_files = ["~/.aws/credentials"]
```
### Update infra.tf
- vpc
- region
- ami
- subnet
- secgroupname
Please note that the VMs must be reachable from ClusterControl.
### Initialize terraform
terraform init
### Apply
terraform apply

## Datastore configuration / Terraform schema
| Configuration Parameter     | Value                                            |
|-----------------------------|--------------------------------------------------|
| `cluster_name`              | `NAME_OF_DATABASE`                                   |
| `database_type`             | `mysql`                                        |
| `database_vendor`           | `percona`                                     |
| `database_version`          | `8.0`                                  |
| `database_topology`         | `replication`                                  |
| `ssh_key_file`              | `SSH_KEYFILE, e. /home/user/.ssh/id_rsa`                      |
| `ssh_user`                  | `"ubuntu"`                                       |
| `install_software`          | `"truex"`                                        |
| `primary_database_host`     | `join(",", aws_instance.project-iac.*.public_ip)`|
| `hostname_internal`         | `join(",", aws_instance.project-iac.*.private_ip)`|

Please note that installing software make take time and the operation may timeout.

### Example datastore.tf
An example is located here:
```
terraform-provider-clustercontrol/examples/replication/mysql80/aws/datastore.tf
```
