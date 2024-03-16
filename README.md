<!-- markdownlint-disable first-line-h1 no-inline-html -->
<a href="https://terraform.io">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset=".github/terraform_logo_dark.svg">
    <source media="(prefers-color-scheme: light)" srcset=".github/terraform_logo_light.svg">
    <img src=".github/terraform_logo_light.svg" alt="Terraform logo" title="Terraform" align="right" height="50">
  </picture>
</a><a href="https://severalnines.com">


# Terraform ClusterControl Provider

This is the Terraform Provider for the Severalnines ClusterControl - Database Automation Tool.

- [ClusterControl Website](https://severalnines.com/clustercontrol/)
- [Documentation](https://docs.severalnines.com/docs/clustercontrol/)
- [Support](https://support.severalnines.com/hc/en-us/requests/new) -  (sign up before you create the request).

## Requirements

| Name | Version  |
|------|----------|
| <a name="requirement_terraform"></a> Terraform | >= 0.13.x   |
| <a name="requirement_cc"></a> ClusterControl | >= 1.9.7 |


## Providers

| Name | Version |
|------|---------|
| <a name="requirement_teraform_cc"></a> Terraform ClusterControl Provider | >= 0.1.0 |

## Resources

| Name |
|------|
| clustercontrol_db_cluster |

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

```curl -k 'https://10.0.0.5:9501/v2/clusters' -XPOST -d '{"operation": "getAllClusterInfo", "authenticate": {"password": "CHANGE-ME","username": "CHANGE-ME"}}'```

Where ``username`` and ``password`` are valid login credentials for ClusterControl.

### Deploying database clusters using terraform

Navigate to the examples sub-directory for more information.

