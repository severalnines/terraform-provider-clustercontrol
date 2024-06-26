---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "clustercontrol Provider"
subcategory: ""
description: "Supported databases: MySQL Replication, MySQL Galera, PostgreSQL, TimeScaleDB, Redis Sentinel, and MongoDB ReplicaSet and Shards"
  
---

# clustercontrol provider

ClusterControl is a database operations orchestration platform for 
MySQL replication, MariaDB server, MariaDB Galera cluster, Percona server, 
Percona XtraDB, PostgreSQL, TimescaleDB, MongoDB, MS SQL Server, Redis, and ElasticSearch, 
that can be deployed in on-premises, cloud, and hybrid environments.

It supports full-lifecycle database ops such as deployment, replication, high availability, 
monitoring (including at the query level), disaster recovery, security, and user management. 
all operations can be configured and managed from the gui, cli, or api.

ClusterControl is not tied to the underlying infrastructure, allowing users to deploy new databases and 
import current ones in multiple environments—their own and/or in their own cloud accounts. 
integrations with popular tools also help users drop it into their workflows relatively easily.

This enables users to place databases in the environment/s of their choice and configure them with more precision. 
it also means that they can predict and assert more control over costs.

ClusterControl is for teams supporting mid to enterprise-level database operations that underpin core products and internal DBaaS projects.
these uses generally require access, workload portability, and automation at scale.

This provider is used to create database resources provided by clustercontrol. 
Supported databases: MySQL replication, MySQL Galera, Postgresql, Timescaledb, Redis Sentinel, and MongoDB Replicaset and Shards.

Use the navigation to the left to read about the available resources.

## [Quick Start Guide](guides/Getting-Started.md)

## [Examples](guides/Examples.md)

### [Deploying MySQL database clusters](guides/MySQL.md)
### [Deploying MariaDB database clusters](guides/MariaDB.md)
### [Deploying Galera database clusters](guides/Galera.md)
### [Deploying ProxySQL load balancer](guides/ProxySQL.md)
### [Deploying PostgreSQL database clusters](guides/PostgreSQL.md)
### [Deploying MongoDB database clusters](guides/MongoDB.md)
### [Deploying Redis database clusters](guides/Redis.md)
### [Deploying Microsoft SQL Server database clusters](guides/MsSQL.md)
### [Deploying Elastisearch clusters](guides/Elasticsearch.md)

<!-- schema generated by tfplugindocs -->
# schema

## required

- `cc_api_url` (string) clustercontrol controller url e.g. (https://cc-host:9501/v2)
- `cc_api_user` (string) clustercontrol api user
- `cc_api_user_password` (string, sensitive) clustercontrol api user's password
