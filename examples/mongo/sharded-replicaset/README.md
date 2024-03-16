# ClusterControl Provider Examples

This directory contains a set of examples of deploying MongoDb database (sharded replicaset) clusters 
using ClusterControl. 

## Resources

| Name |
|------|
| clustercontrol_db_cluster |

## Attributes specific to MongoDB sharded (replicaset) deployment

| Attribute                | Data Type   | Required | Description                                              |
|--------------------------|-------------|----------|----------------------------------------------------------|
| db_config_server         | object      | Yes      | Specify the mongodb config server for sharded deployment |
| db_mongos_server         | object      | Yes      | Specify the mongos server for sharded deployment         |

