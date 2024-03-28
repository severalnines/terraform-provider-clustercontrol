package provider

const (
	RESOURCE_DB_CLUSTER                 = "clustercontrol_db_cluster"
	RESOURCE_DB_LOAD_BALANCER           = "clustercontrol_db_load_balancer"
	RESOURCE_DB_CLUSTER_MAINTENANCE     = "clustercontrol_db_cluster_maintenance"
	RESOURCE_DB_CLUSTER_BACKUP          = "clustercontrol_db_cluster_backup"
	RESOURCE_DB_CLUSTER_BACKUP_SCHEDULE = "clustercontrol_db_cluster_backup_schedule"
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
	VENDOR_PERCONA    = "percona"
	VENDOR_MARIADB    = "mariadb"
	VENDOR_MONGODB    = "10gen"
	VENDOR_ORACLE     = "oracle"
	VENDOR_ELASTIC    = "elasticsearch"
	VENDOR_REDIS      = "redis"
	VENDOR_MICROSOFT  = "microsoft"
	VENDOR_DEFAULT    = "default"
	VENDOR_POSTGRESQL = "postgresql"
	VENDOR_EDB        = "EDB"
	VENDOR_10GEN      = "10gen"
)

const (
	EXT_CLUSTER_TYPE_PG_REPLICAION      = "pg-replication"
	EXT_CLUSTER_TYPE_MYSQL_REPLICATIOIN = "mysql-replication"
	EXT_CLUSTER_TYPE_GALERA             = "galera"
	EXT_CLUSTER_TYPE_MONGODB            = "mongo"
	EXT_CLUSTER_TYPE_REDIS_SENTINEL     = "redis-sentinel"
	EXT_CLUSTER_TYPE_REDIS_CLUSTER      = "redis-cluster"
	EXT_CLUSTER_TYPE_ELASTICSEARH       = "elasticsearch"
	EXT_CLUSTER_TYPE_MSSQL_ASYN         = "mssql-async"
	//EXT_CLUSTER_TYPE_ = ""
)

const (
	EXT_VENDOR_PERCONA    = "percona"
	EXT_VENDOR_ORACLE     = "oracle"
	EXT_VENDOR_MARIADB    = "mariadb"
	EXT_VENDOR_MONGO      = "mongodb-community"
	EXT_VENDOR_MICROSOFT  = "microsoft"
	EXT_VENDOR_ELASTIC    = "elastic"
	EXT_VENDOR_REDIS      = "redis"
	EXT_VENDOR_POSTGRESQL = "postgresql"
	EXT_VENDOR_MONGO_ENT  = "mongodb-X"

	//EXT_VENDOR_ = ""
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
	DEFAULT_PROXYSQL_ADMIN_PORT         = "6032"
	DEFAULT_PROXYSQL_LISTEN_PORT        = "6033"
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
	CMON_JOB_DELETE_JOB = "deleteJobInstance"
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
	CMON_CLASS_NAME_CMON_HOST          = "CmonHost"
	CMON_CLASS_NAME_PROMETHEUS_HOST    = "CmonPrometheusHost"
	CMON_CLASS_NAME_MYSQL_HOST         = "CmonMySqlHost"
	CMON_CLASS_NAME_PROXYSQL_HOST      = "CmonProxySqlHost"
	CMON_CLASS_NAME_PROXYSQL_SRVR_HOST = "CmonProxySqlServer"
	CMON_CLASS_NAME_CMON_AGEN_HOST     = "CmonAgentHost"
	CMON_CLASS_NAME_GALERA_HOST        = "CmonGaleraHost"
	CMON_CLASS_NAME_PROSGRESQL_HOST    = "CmonPostgreSqlHost"
	CMON_CLASS_NAME_HAPROXY_HOST       = "CmonHaProxyHost"
	CMON_CLASS_NAME_PGBACKREST_HOST    = "CmonPgBackRestHost"
	CMON_CLASS_NAME_MONGO_HOST         = "CmonMongoHost"
	CMON_CLASS_NAME_PBM_AGENT_HOST     = "CmonPBMAgentHost"
)

const (
	CMON_DB_HOST_ROLE_MASTER              = "master"
	CMON_DB_HOST_ROLE_SLAVE               = "slave"
	CMON_DB_HOST_ROLE_MULI                = "multi" // Master->IntermediateSlave(multi)->Slave
	CMON_DB_HOST_ROLE_MONGO_CFG_SERVER    = "configsvr"
	CMON_DB_HOST_ROLE_MONGO_SHARD_SERVER  = "shardsvr"
	CMON_DB_HOST_ROLE_MONGO_MONGOS_SERVER = "mongos"
	CMON_DB_HOST_ROLE_PRIMARY             = "PRIMARY"
	CMON_DB_HOST_ROLE_SECONDARY           = "SECONDARY"
)

const (
	CMON_JOB_CREATE_CLUSTER_COMMAND           = "create_cluster"
	CMON_JOB_REMOVE_CLUSTER_COMMAND           = "remove_cluster"
	CMON_JOB_CREATE_PROXYSQL_COMMAND          = "proxysql"
	CMON_JOB_CREATE_HAPROXY_COMMAND           = "haproxy"
	CMON_JOB_CREATE_BACKUP_COMMAND            = "backup"
	CMON_JOB_DELETE_BACKUP_COMMAND            = "delete_backup"
	CMON_JOB_ENABLE_CLUSTER_RECOVERY_COMMAND  = "enable_recovery"
	CMON_JOB_DISABLE_CLUSTER_RECOVERY_COMMAND = "disable_recovery"
	CMON_JOB_ADD_REPLICATION_SLAVE_COMMAND    = "add_replication_slave"
	CMON_JOB_ADD_NODE_COMMAND                 = "addnode"
	CMON_JOB_REMOVE_NODE_COMMAND              = "removenode"
	CMON_JOB_REGISTER_NODE_COMMAND            = "registernode"
	CMON_JOB_ADD_SHARD_COMMAND                = "add_shard"
	CMON_JOB_PROMOTE_REPLICAION_SLAVE_COMMAND = "promote_replication_slave"
	CMON_JOB_PBM_AGENT_COMMAND                = "pbmagent"
	CMON_JOB_PGBACKREST_COMMAND               = "pgbackrest"
)

const (
	CMON_CLUSTERS_OPERATION_GET_CLUSTERS = "getclusterinfo"
	CMON_CLUSTERS_OPERATION_SET_CONFIG   = "setConfig"
)

const (
	CMON_BACKUP_OPERATION_GET = "getBackups"
)

const BACKUP_RECORD_VERSION_2 = 2

const (
	BACKUP_ORDER_CREATED_DESC = "created DESC"
)

const (
	CMON_MAINTENANCE_OPERATION_ADD_MAINT    = "addMaintenance"
	CMON_MAINTENANCE_OPERATION_REMOVE_MAINT = "removeMaintenance"
)

const (
	LOAD_BLANCER_TYPE_PROXYSQL = "proxysql"
	LOAD_BLANCER_TYPE_HAPROXY  = "haproxy"
)

const (
	JOB_ACTION_SETUP_PROXYSQL = "setupProxySql"
	JOB_ACTION_SETUP_HAPROXY  = "setupHaProxy"
	JOB_ACTION_SETUP          = "setup"
)

const (
	BACKUP_METHOD_XTRABACKUP_FULL  = "xtrabackupfull"
	BACKUP_METHOD_XTRABACKUP_INCR  = "xtrabackupincr"
	BACKUP_METHOD_MARIABACKUP_FULL = "mariabackupfull"
	BACKUP_METHOD_MARIABACKUP_INCR = "mariabackupincr"
	BACKUP_METHOD_MYSQLDUMP        = "mysqldump"
	BACKUP_METHOD_PG_BASEBACKUP    = "pg_basebackup"
	BACKUP_METHOD_PG_BACKREST_FULL = "pgbackrestfull"
	BACKUP_METHOD_PG_BACKREST_INCR = "pgbackrestincr"
	BACKUP_METHOD_PG_BACKREST_DIRR = "pgbackrestdiff"
	BACKUP_METHOD_PGDUMPALL        = "pgdumpall"
	BACKUP_METHOD_MONGODUMP        = "mongodump"
	BACKUP_METHOD_PBM              = "percona-backup-mongodb"
	BACKUP_MSSQL_FULL              = "mssqlfull"
	BACKUP_MSSQL_DIFF              = "mssqldiff"
	BACKUP_MSSQL_TRANSACTION_LOG   = "mssqllog"
)

const (
	TF_FIELD_RESOURCE_ID                     = "db_resource_id"
	TF_FIELD_LAST_UPDATED                    = "last_updated"
	TF_FIELD_CLUSTER_CREATE                  = "db_cluster_create"
	TF_FIELD_CLUSTER_IMPORT                  = "db_cluster_import"
	TF_FIELD_CLUSTER_ID                      = "db_cluster_id"
	TF_FIELD_CLUSTER_NAME                    = "db_cluster_name"
	TF_FIELD_CLUSTER_TYPE                    = "db_cluster_type"
	TF_FIELD_CLUSTER_VENDOR                  = "db_vendor"
	TF_FIELD_CLUSTER_VERSION                 = "db_version"
	TF_FIELD_CLUSTER_ADMIN_USER              = "db_admin_username"
	TF_FIELD_CLUSTER_ADMIN_PW                = "db_admin_user_password"
	TF_FIELD_CLUSTER_PORT                    = "db_port"
	TF_FIELD_CLUSTER_DATA_DIR                = "db_data_directory"
	TF_FIELD_CLUSTER_CFG_TEMPLATE            = "db_config_template"
	TF_FIELD_CLUSTER_DISABLE_FW              = "disable_firewall"
	TF_FIELD_CLUSTER_DISABLE_SELINUX         = "disable_selinux"
	TF_FIELD_CLUSTER_INSTALL_SW              = "db_install_software"
	TF_FIELD_CLUSTER_ENABLE_UNINSTALL        = "db_enable_uninstall"
	TF_FIELD_CLUSTER_SYNC_REP                = "sync_replication"
	TF_FIELD_CLUSTER_SEMISYNC_REP            = "db_semi_sync_replication"
	TF_FIELD_CLUSTER_PG_TIMESALE_EXT         = "db_enable_timescale"
	TF_FIELD_CLUSTER_SSH_USER                = "ssh_user"
	TF_FIELD_CLUSTER_SSH_PW                  = "ssh_user_password"
	TF_FIELD_CLUSTER_SSH_KEY_FILE            = "ssh_key_file"
	TF_FIELD_CLUSTER_SSH_PORT                = "ssh_port"
	TF_FIELD_CLUSTER_SNAPSHOT_LOC            = "db_snapshot_location"
	TF_FIELD_CLUSTER_SNAPSHOT_REPO           = "db_snapshot_repository"
	TF_FIELD_CLUSTER_SNAPSHOT_HOST           = "db_snapshot_host"
	TF_FIELD_CLUSTER_HOST                    = "db_host"
	TF_FIELD_CLUSTER_HOSTNAME                = "hostname"
	TF_FIELD_CLUSTER_HOSTNAME_DATA           = "hostname_data"
	TF_FIELD_CLUSTER_HOSTNAME_INT            = "hostname_internal"
	TF_FIELD_CLUSTER_HOST_PORT               = "port"
	TF_FIELD_CLUSTER_HOST_DD                 = "datadir"
	TF_FIELD_CLUSTER_HOST_PRIORITY           = "priority"
	TF_FIELD_CLUSTER_HOST_SLAVE_DELAY        = "slave_delay"
	TF_FIELD_CLUSTER_HOST_ARBITER_ONLY       = "arbiter_only"
	TF_FIELD_CLUSTER_HOST_HIDDEN             = "hidden"
	TF_FIELD_CLUSTER_HOST_PROTO              = "protocol"
	TF_FIELD_CLUSTER_HOST_ROLES              = "roles"
	TF_FIELD_CLUSTER_TOPOLOGY                = "db_topology"
	TF_FIELD_CLUSTER_PRIMARY                 = "primary"
	TF_FIELD_CLUSTER_REPLICA                 = "replica"
	TF_FIELD_CLUSTER_TAGS                    = "db_tags"
	TF_FIELD_CLUSTER_REPLICA_SET             = "db_replica_set"
	TF_FIELD_CLUSTER_REPLICA_SET_RS          = "rs"
	TF_FIELD_CLUSTER_REPLICA_MEMBER          = "member"
	TF_FIELD_CLUSTER_MONGO_CONFIG_SERVER     = "db_config_server"
	TF_FIELD_CLUSTER_MONGOS_SERVER           = "db_mongos_server"
	TF_FIELD_CLUSTER_TIMEOUTS                = "timeouts"
	TF_FIELD_CLUSTER_DEPLOY_AGENTS           = "db_deploy_agents"
	TF_FIELD_CLUSTER_AUTO_RECOVERY           = "db_auto_recovery"
	TF_FIELD_CLUSTER_SSL                     = "db_enable_ssl"
	TF_FIELD_CLUSTER_MONGO_AUTH_DB           = "db_mongo_auth_db"
	TF_FIELD_CLUSTER_SENTINEL_PORT           = "db_sentinel_port"
	TF_FIELD_CLUSTER_ENABLE_PGM_AGENT        = "db_enable_pbm_agent"
	TF_FIELD_CLUSTER_PBM_BACKUP_DIR          = "db_pbm_backup_dir"
	TF_FIELD_CLUSTER_ENABLE_PGBACKREST_AGENT = "db_enable_pgbackrest_agent"

	// Load balancer fields
	TF_FIELD_LB_CREATE           = "db_lb_create"
	TF_FIELD_LB_IMPORT           = "db_lb_import"
	TF_FIELD_LB_TYPE             = "db_lb_type"
	TF_FIELD_LB_VERSION          = "db_lb_version"
	TF_FIELD_LB_ADMIN_USER       = "db_lb_admin_username"
	TF_FIELD_LB_ADMIN_USER_PW    = "db_lb_admin_user_password"
	TF_FIELD_LB_MONITOR_USER     = "db_lb_monitor_username"
	TF_FIELD_LB_MONITOR_USER_PW  = "db_lb_monitor_user_password"
	TF_FIELD_LB_PORT             = "db_lb_port"
	TF_FIELD_LB_ADMIN_PORT       = "db_lb_admin_port"
	TF_FIELD_LB_USE_CLUSTERING   = "db_lb_use_clustering"
	TF_FIELD_LB_USE_RW_SPLITTING = "db_lb_use_rw_splitting"
	TF_FIELD_LB_INSTALL_SW       = "db_lb_install_software"
	TF_FIELD_LB_MY_HOST          = "db_my_host"
	TF_FIELD_LB_ENABLE_UNINSTALL = "db_lb_enable_uninstall"

	// Maintenance fields
	TF_FIELD_MAINT_START_TIME = "db_maint_start_time"
	TF_FIELD_MAINT_STOP_TIME  = "db_maint_stop_time"
	TF_FIELD_MAINT_REASON     = "db_maint_reason"

	// Backup fields
	TF_FIELD_BACKUP_METHOD            = "db_backup_method"
	TF_FIELD_BACKUP_DIR               = "db_backup_dir"
	TF_FIELD_BACKUP_SUBDIR            = "db_backup_subdir"
	TF_FIELD_BACKUP_ENCRYPT           = "db_backup_encrypt"
	TF_FIELD_BACKUP_HOST              = "db_backup_host"
	TF_FIELD_BACKUP_COMPRESSION       = "db_backup_compression"
	TF_FIELD_BACKUP_COMPRESSION_LEVEL = "db_backup_compression_level"
	TF_FIELD_BACKUP_RETENTION         = "db_backup_retention"
	TF_FIELD_BACKUP_ON_CONTROLLER     = "db_backup_storage_controller"
	TF_FIELD_BACKUP_FAILOVER_HOST     = "db_backup_failover_host"   // MONGODUMP only
	TF_FIELD_BACKUP_FAILOVER          = "db_enable_backup_failover" // MONGODUMP only
	TF_FIELD_BACKUP_STORAGE_HOST      = "db_backup_storage_host"
	TF_FIELD_BACKUP_SYSTEM_DB         = "db_backup_system_db"
	//TF_FIELD_BACKUP_                  = "db_backup_"

	// Backup schedule fields
	TF_FIELD_BACKUP_SCHED_TITLE = "db_backup_sched_title"
	TF_FIELD_BACKUP_SCHED_TIME  = "db_backup_sched_time"
	//TF_FIELD_BACKUP_SCHED_                  = "db_backup_"

)

const (
	MONGO_DEFAULT_AUTH_DB = "admin"
)

const (
	STINRG_AUTO = "auto"
)

const (
	TIME_FORMAT = "Jan-02-2006T15:04"
)

const (
	JOB_STATUS_DEFINED   = "DEFINED"
	JOB_STATUS_RUNNING   = "RUNNING"
	JOB_STATUS_FINISHED  = "FINISHED"
	JOB_STATUS_SCHEDULED = "SCHEDULED"
)

//const (
//	MONGO_CONFIG_SERVER_DEFAULT_PORT = 27019
//)
