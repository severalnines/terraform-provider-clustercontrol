<!-- markdownlint-disable first-line-h1 no-inline-html -->
<a href="https://terraform.io">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="logos/hashicorp/terraform_logo_dark.svg">
    <source media="(prefers-color-scheme: light)" srcset="logos/hashicorp/terraform_logo_light.svg">
    <img src="logos/hashicorp/terraform_logo_light.svg" alt="Terraform logo" title="Terraform" align="right" height="50">
  </picture>
</a>

<a href="https://severalnines.com">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="logos/severalnines/severalnies.png">
    <source media="(prefers-color-scheme: light)" srcset="logos/severalnines/severalnies.png">
    <img src="logos/severalnines/severalnies.png" alt="Terraform logo" title="Terraform" align="left" height="50">
  </picture>
</a>

# Terraform ClusterControl Provider

This is the Terraform Provider for the Severalnines ClusterControl - Database Automation Tool.

- [ClusterControl Website](https://severalnines.com/clustercontrol/)
- [Documentation](https://docs.severalnines.com/docs/clustercontrol/)
- [Support](https://support.severalnines.com/hc/en-us/requests/new) -  (sign up before you create the request).

## Requirements

| Name | Version   |
|------|-----------|
| <a name="requirement_terraform"></a> Terraform | >= 0.13.x |
| <a name="requirement_cc"></a> ClusterControl | >= 1.9.8  |


## Providers

| Name | Version  |
|------|----------|
| <a name="requirement_teraform_cc"></a> Terraform ClusterControl Provider | >= 0.2.0 |

## Resources

| Name                                                                                                                                                                     |
|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [clustercontrol_db_cluster](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster.md#clustercontrol_db_cluster-resource) |
| [clustercontrol_db_cluster_backup](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_backup.md#clustercontrol_db_cluster_backup-resource)|                                                                                                                                                                                    |
| [clustercontrol_db_cluster_backup_schedule](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_backup_schedule.md#clustercontrol_db_cluster_backup_schedule-resource) |
| [clustercontrol_db_cluster_maintenance](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_maintenance.md#clustercontrol_db_cluster_maintenance-resource)|


## Quick Start
### Installing and configuring ClusterControl for API access
[NOTE:] If you already have ClusterControl running in your env, skip to #4
1. [Getting started with ClusterControl](https://docs.severalnines.com/docs/clustercontrol/getting-started/)
2. [ClusterControl installation requirements](https://docs.severalnines.com/docs/clustercontrol/getting-started/)
3. [Install ClusterControl](https://docs.severalnines.com/docs/clustercontrol/installation/automatic-installation/)
4. Configure ClusterControl - Enable ClusterControl for API access by (on the ClusterControl host)

``sudo vi /etc/default/cmon`` and set the ``RPC_BIND_ADDRESSES`` as shown below

``RPC_BIND_ADDRESSES="10.0.0.15,127.0.0.1"``

Where ``10.0.0.5`` is the private IP of the ClusterControl host. Restart ClusterControl service.

``sudo systemctl restart cmon``

5. Run a quick test to make sure you can access ClusterControl via its REST API (curl or postman)

```shell
curl -c cookie.jar.txt -k 'https://10.0.0.5:9501/v2/auth' -XPOST -d '{"operation":"authenticateWithPassword","user_name":"CHANGE_ME","password":"CHANGE_ME"}'
```

Where ``username`` and ``password`` are valid login credentials for ClusterControl.

Make sure you get a response like the following. If you don't get a similar `JSON` response, it means something has gone wrong.

```json
{
    "controller_id": "76936896-794a-4f5f-95ac-5c44cefb6830",
    "is_superuser": true,
    "request_processed": "2024-05-24T15:00:35.400Z",
    "request_status": "Ok",
    "debug_messages": 
    [
        "RPC V2 authenticated user is 'username'."
    ],
    "user": 
    {
        "class_name": "CmonUser",
        "owner_user_name": "username",
        "groups": 
        [
            {
                "class_name": "CmonGroup",
                "cdt_path": "/groups",
                "group_name": "admins"
            }
        ],
        "timezone": 
        {
            "class_name": "CmonTimeZone",
            "name": "UTC",
            "abbreviation": "UTC",
            "offset": 0,
            "use_dst": false
        }
    }
}
```

A cookie will be returned by ClusterControl back-end and saved in the `cookie.jar.txt` per the curl command.
Execute the following additional curl command and check the response to make sure you are properly authenticated with ClusterControl.

```shell
curl -k -b cookie.jar.txt -k 'https://10.0.0.5:9501/v2/clusters' -XPOST -d '{"operation": "getAllClusterInfo"}'
```
Response should look like the following `JSON`:
```json
{
    "controller_id": "76936896-794a-4f5f-95ac-5c44cefb6830",
    "request_processed": "2024-05-24T15:07:25.837Z",
    "request_status": "Ok",
    "request_user_id": 5,
    "total": 0,
    "debug_messages": [
        "RPC V2 authenticated user is 'username'."
    ]
}
```

### Deploying database clusters using terraform for ClusterControl

**Navigate** to the [examples](./examples/README.md) folder 
for concrete examples on deploying database clusters of various types (MySQL/MariaDB replication or galera with ProxySQL, 
PostgreSql replication, MongoDB replicaset and/or sharded, Redis sentinel, Microsoft SQL server, and Elasticsearch)

**Navigate** to the [docs](./docs/index.md) folder for generated documentation on the terraform provider plugin for ClusterControl

### Setup ``terraform.tfvars`` file with the following secrets.


```editor
cc_api_url="https://<cc-host-or-ip>:9501/v2"
cc_api_user="CHANGE-ME"
cc_api_user_password="CHANGE-ME"
```

#### Running terraform to deploy database clusters

```shell
terraform init
terraform validate
terraform plan -var-file="terraform.tfvars"
terraform apply -var-file="terraform.tfvars"
```

#### Destroying a deployed database cluster.

After navigating to the appropriate directory which was used to deploy a cluster:

```shell
terraform destroy -var-file="terraform.tfvars"
```
