package provider

const (
	RESOURCE_DB_CLUSTER       = "cc_db_cluster"
	RESOURCE_DB_LOAD_BALANCER = "cc_db_load_balancer"
)

const (
	API_USER       = "cc_api_user"
	API_USER_PW    = "cc_api_user_password"
	CONTROLLER_URL = "cc_api_url"
)

const (
	CLUSTER_TYPE_REPLICATION    = "replication"
	CLUSTER_TYPE_GALERA         = "galera"
	CLUSTER_TYPE_MOGNODB        = "mongodb"
	CLUSTER_TYPE_PG_SINGLE      = "postgresql_single"
	CLUSTER_TYPE_REDIS          = "redis"
	CLUSTER_TYPE_MSSQL_AO_ASYNC = "mssql_ao_async"
	CLUSTER_TYPE_MSSQL_SINGLE   = "mssql_single"
	CLUSTER_TYPE_ELASTIC        = "elastic"
)

const (
	VENDOR_PERCONA   = "percona"
	VENDOR_MARIADB   = "mariadb"
	VENDOR_MONGODB   = "10gen"
	VENDOR_ORACLE    = "oracle"
	VENDOR_ELASTIC   = "elasticsearch"
	VENDOR_REDIS     = "redis"
	VENDOR_MICROSOFT = "microsoft"
	VENDOR_DEFAULT   = "default"
	VENDOR_10GEN     = "10gen"
)

const (
	DEFAULT_MYSQL_PORT                  = "3306"
	DEFAULT_POSTGRES_PORT               = "5432"
	DEFAULT_MONGO_PORT                  = "27017"
	DEFAULT_MONGO_CONFIG_SRVR_PORT      = "27019"
	DEFAULT_MONGO_REDIS_PORT            = "6379"
	DEFAULT_MONGO_REDIS_SENTINEL_PORT   = "26379"
	DEFAULT_MONGO_ELASTIC_HTTP_PORT     = "9200"
	DEFAULT_MONGO_ELASTIC_TRANSFER_PORT = "9200"
	DEFAULT_MONGO_MSSQL_PORT            = "1433"
)

const (
	MYSQL_VERSION_8   = "8.0"
	MYSQL_VERSION_5_7 = "5.7"
)

const (
	MARIADB_VERSION_10_11 = "10.11"
	MARIADB_VERSION_10_10 = "10.10"
	MARIADB_VERSION_10_9  = "10.9"
	MARIADB_VERSION_10_8  = "10.8"
	MARIADB_VERSION_10_6  = "10.6"
	MARIADB_VERSION_10_5  = "10.5"
	MARIADB_VERSION_10_4  = "10.4"
	MARIADB_VERSION_10_3  = "10.3"
)

const (
	POSTGRESQL_VERSION_15 = "15"
	POSTGRESQL_VERSION_14 = "14"
	POSTGRESQL_VERSION_13 = "13"
	POSTGRESQL_VERSION_12 = "12"
	POSTGRESQL_VERSION_11 = "11"
)

const (
	MONGODB_VERSION_6_0 = "6.0"
	MONGODB_VERSION_5_0 = "5.0"
	MONGODB_VERSION_4_4 = "4.4"
	MONGODB_VERSION_4_2 = "4.2"
)

const (
	REDIS_VERSION_7 = "7"
	REDIS_VERSION_6 = "6"
)

const (
	MSSQL_VERSION_2019 = "2019"
	MSSQL_VERSION_2022 = "2022"
)

const (
	ELASTIC_VERSION_8_3_1  = "8.3.1"
	ELASTIC_VERSION_8_1_3  = "8.1.3"
	ELASTIC_VERSION_7_17_3 = "7.17.3"
)

// move to the CC client SDK library
const (
	CMON_OP_AUTHENTICATE_WITH_PW = "authenticateWithPassword"
)

const (
	CMON_MOCK_USER = "mock-user"
)

const (
	CMON_JOB_CREATE_JOB = "createJobInstance"
	CMON_JOB_GET_JOB    = "getJobInstance"
)

const (
	CMON_JOB_CLASS_NAME = "CmonJobInstance"
)

const (
	CMON_CLASS_NAME_REDIS_HOST         = "CmonRedisHost"
	CMON_CLASS_NAME_REDIS_SENTNEL_HOST = "CmonRedisSentinelHost"
	CMON_CLASS_NAME_MSSQL_HOST         = "CmonMsSqlHost"
	CMON_CLASS_NAME_ELASTIC_HOST       = "CmonElasticHost"
)

const (
	CMON_JOB_CREATE_CLUSTER_COMMAND = "create_cluster"
	CMON_JOB_REMOVE_CLUSTER_COMMAND = "remove_cluster"
)

const (
	CMON_CLUSTERS_OPERATION_GET_CLUSTERS = "getclusterinfo"
)

// TODO: doesn't seem to be use anymore (03/07/2024)
//const (
//	MYSQL_DATA_DIR           = "/var/lib/mysql"
//	MySQL_DB_ADMIN_USER_NAME = "root"
//)

const (
	TF_FIELD_CLUSTER_CREATE              = "db_cluster_create"
	TF_FIELD_CLUSTER_IMPORT              = "db_cluster_import"
	TF_FIELD_CLUSTER_NAME                = "db_cluster_name"
	TF_FIELD_CLUSTER_TYPE                = "db_cluster_type"
	TF_FIELD_CLUSTER_VENDOR              = "db_vendor"
	TF_FIELD_CLUSTER_VERSION             = "db_version"
	TF_FIELD_CLUSTER_ADMIN_USER          = "db_admin_username"
	TF_FIELD_CLUSTER_ADMIN_PW            = "db_admin_user_password"
	TF_FIELD_CLUSTER_PORT                = "db_port"
	TF_FIELD_CLUSTER_DATA_DIR            = "db_data_directory"
	TF_FIELD_CLUSTER_CFG_TEMPLATE        = "db_config_template"
	TF_FIELD_CLUSTER_DISABLE_FW          = "disable_firewall"
	TF_FIELD_CLUSTER_INSTALL_SW          = "db_install_software"
	TF_FIELD_CLUSTER_SYNC_REP            = "sync_replication"
	TF_FIELD_CLUSTER_SEMISYNC_REP        = "db_semi_sync_replication"
	TF_FIELD_CLUSTER_SSH_USER            = "ssh_user"
	TF_FIELD_CLUSTER_SSH_PW              = "ssh_user_password"
	TF_FIELD_CLUSTER_SSH_KEY_FILE        = "ssh_key_file"
	TF_FIELD_CLUSTER_SSH_PORT            = "ssh_port"
	TF_FIELD_CLUSTER_SNAPSHOT_LOC        = "db_snapshot_location"
	TF_FIELD_CLUSTER_SNAPSHOT_REPO       = "db_snapshot_repository"
	TF_FIELD_CLUSTER_SNAPSHOT_HOST       = "db_snapshot_host"
	TF_FIELD_CLUSTER_HOST                = "db_host"
	TF_FIELD_CLUSTER_HOSTNAME            = "hostname"
	TF_FIELD_CLUSTER_HOSTNAME_DATA       = "hostname_data"
	TF_FIELD_CLUSTER_HOSTNAME_INT        = "hostname_internal"
	TF_FIELD_CLUSTER_HOST_PORT           = "port"
	TF_FIELD_CLUSTER_HOSTNAME_DD         = "data_dir"
	TF_FIELD_CLUSTER_HOST_PRIORITY       = "priority"
	TF_FIELD_CLUSTER_HOST_SLAVE_DELAY    = "slave_delay"
	TF_FIELD_CLUSTER_HOST_ARBITER_ONLY   = "arbiter_only"
	TF_FIELD_CLUSTER_HOST_HIDDEN         = "hidden"
	TF_FIELD_CLUSTER_HOST_PROTO          = "protocol"
	TF_FIELD_CLUSTER_HOST_ROLES          = "roles"
	TF_FIELD_CLUSTER_TOPOLOGY            = "db_topology"
	TF_FIELD_CLUSTER_PRIMARY             = "primary"
	TF_FIELD_CLUSTER_REPLICA             = "replica"
	TF_FIELD_CLUSTER_TAGS                = "db_tags"
	TF_FIELD_CLUSTER_REPLICA_SET         = "db_replica_set"
	TF_FIELD_CLUSTER_REPLICA_SET_RS      = "rs"
	TF_FIELD_CLUSTER_REPLICA_MEMBER      = "member"
	TF_FIELD_CLUSTER_MONGO_CONFIG_SERVER = "db_config_server"
	TF_FIELD_CLUSTER_MONGOS_SERVER       = "db_mongos_server"
	TF_FIELD_CLUSTER_TIMEOUTS            = "timeouts"
)

const (
	JOB_STATUS_DEFINED  = "DEFINED"
	JOB_STATUS_RUNNING  = "RUNNING"
	JOB_STATUS_FINISHED = "FINISHED"
)

//const (
//	MONGO_CONFIG_SERVER_DEFAULT_PORT = 27019
//)