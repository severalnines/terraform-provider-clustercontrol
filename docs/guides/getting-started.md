# Getting Started
## Installing and configuring ClusterControl for API access
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

```curl -k 'https://10.0.0.5:9501/v2/clusters' -XPOST -d '{"operation": "getAllClusterInfo", "authenticate": {"password": "CHANGE-ME","username": "CHANGE-ME"}}'```

Where ``username`` and ``password`` are valid login credentials for ClusterControl.

## Deploying database clusters using terraform for ClusterControl

**Navigate** to the [examples](guides/examples.md) folder
for concrete examples on deploying database clusters of various types (MySQL/MariaDB replication or galera with ProxySQL,
PostgreSql replication, MongoDB replicaset and/or sharded, Redis sentinel, Microsoft SQL server, and Elasticsearch)

### Setup ``terraform.tfvars`` file with the following secrets.


```editor
cc_api_url="https://<cc-host-or-ip>:9501/v2"
cc_api_user="CHANGE-ME"
cc_api_user_password="CHANGE-ME"
```

### Running terraform to deploy database clusters

```shell
terraform init
terraform validate
terraform plan -var-file="terraform.tfvars"
terraform apply -var-file="terraform.tfvars"
```

## Destroying a deployed database cluster.

After navigating to the appropriate directory which was used to deploy a cluster:

```shell
terraform destroy -var-file="terraform.tfvars"
```

