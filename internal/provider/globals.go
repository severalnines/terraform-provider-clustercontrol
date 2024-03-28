package provider

import (
	"context"
	"github.com/severalnines/clustercontrol-client-sdk/go/pkg/openapi"
)

var (
	gNewCtx    context.Context
	gCfg       *openapi.Configuration
	gApiClient *openapi.APIClient
)

var (
	gDefultDbPortMap = map[string]string{
		CLUSTER_TYPE_REPLICATION:    DEFAULT_MYSQL_PORT,
		CLUSTER_TYPE_GALERA:         DEFAULT_MYSQL_PORT,
		CLUSTER_TYPE_PG_SINGLE:      DEFAULT_POSTGRES_PORT,
		CLUSTER_TYPE_MOGNODB:        DEFAULT_MONGO_PORT,
		CLUSTER_TYPE_REDIS:          DEFAULT_MONGO_REDIS_PORT,
		CLUSTER_TYPE_ELASTIC:        DEFAULT_MONGO_ELASTIC_HTTP_PORT,
		CLUSTER_TYPE_MSSQL_SINGLE:   DEFAULT_MONGO_MSSQL_PORT,
		CLUSTER_TYPE_MSSQL_AO_ASYNC: DEFAULT_MONGO_MSSQL_PORT,
	}
)

var ( // TODO: have to think about this one....
	gExtClusterTypeToIntClusterTypeMap = map[string]string{
		CLUSTER_TYPE_REPLICATION: CLUSTER_TYPE_REPLICATION,
		CLUSTER_TYPE_GALERA:      CLUSTER_TYPE_GALERA,
	}
)

var (
	gDefultDbAdminUser = map[string]string{
		CLUSTER_TYPE_REPLICATION:    "root",
		CLUSTER_TYPE_GALERA:         "root",
		CLUSTER_TYPE_PG_SINGLE:      "pgadmin",
		CLUSTER_TYPE_MOGNODB:        "mongoadmin",
		CLUSTER_TYPE_REDIS:          "admin",
		CLUSTER_TYPE_ELASTIC:        "admin",
		CLUSTER_TYPE_MSSQL_SINGLE:   "SQLServerAdmin",
		CLUSTER_TYPE_MSSQL_AO_ASYNC: "SQLServerAdmin",
	}
)

var (
	gDefultDataDir = map[string]string{
		CLUSTER_TYPE_REPLICATION:    "/var/lib/mysql",
		CLUSTER_TYPE_GALERA:         "/var/lib/mysql",
		CLUSTER_TYPE_PG_SINGLE:      "",
		CLUSTER_TYPE_MOGNODB:        "/var/lib/mongodb",
		CLUSTER_TYPE_REDIS:          "/var/lib/redis",
		CLUSTER_TYPE_ELASTIC:        "",
		CLUSTER_TYPE_MSSQL_SINGLE:   "/var/opt/mssql/data",
		CLUSTER_TYPE_MSSQL_AO_ASYNC: "/var/opt/mssql/data",
	}
)

var (
	gDefaultHostConfigFile = map[string]string{
		CLUSTER_TYPE_MSSQL_SINGLE:   "/var/opt/mssql/mssql.conf",
		CLUSTER_TYPE_MSSQL_AO_ASYNC: "/var/opt/mssql/mssql.conf",
	}
)
var (
	gDbMongosConfigTemplate = map[string]string{
		VENDOR_PERCONA: "mongos.conf.org",
		VENDOR_10GEN:   "",
	}
)

var (
	gDbConfigTemplate = map[string]map[string]map[string]string{
		VENDOR_ORACLE: {
			CLUSTER_TYPE_REPLICATION: {
				MYSQL_VERSION_5_7: "my.cnf.repl57",
				MYSQL_VERSION_8:   "my.cnf.repl80",
			},
			//CLUSTER_TYPE_GALERA: {
			//	MYSQL_VERSION_8: "my.cnf.80-pxc",
			//},
		},
		VENDOR_PERCONA: {
			CLUSTER_TYPE_REPLICATION: {
				MYSQL_VERSION_5_7: "my.cnf.repl57",
				MYSQL_VERSION_8:   "my.cnf.repl80",
			},
			CLUSTER_TYPE_GALERA: {
				MYSQL_VERSION_5_7: "my57.cnf.galera",
				MYSQL_VERSION_8:   "my.cnf.80-pxc",
			},
			CLUSTER_TYPE_MOGNODB: {
				MONGODB_VERSION_6_0: "mongodb.conf.6.0.percona",
				MONGODB_VERSION_5_0: "mongodb.conf.5.0.percona",
				MONGODB_VERSION_4_4: "mongodb.conf.4.4.percona",
				MONGODB_VERSION_4_2: "mongodb.conf.4.2.percona",
			},
		},
		VENDOR_MARIADB: {
			CLUSTER_TYPE_REPLICATION: {
				MARIADB_VERSION_10_11: "my.cnf.mdb10x-replication",
				MARIADB_VERSION_10_10: "my.cnf.mdb10x-replication",
				MARIADB_VERSION_10_9:  "my.cnf.mdb10x-replication",
				MARIADB_VERSION_10_8:  "my.cnf.mdb10x-replication",
				MARIADB_VERSION_10_6:  "my.cnf.mdb10x-replication",
				MARIADB_VERSION_10_5:  "my.cnf.mdb10x-replication",
				MARIADB_VERSION_10_4:  "my.cnf.mdb10x-replication",
				MARIADB_VERSION_10_3:  "my.cnf.mdb10x-replication",
			},
			CLUSTER_TYPE_GALERA: {
				MARIADB_VERSION_10_11: "my.cnf.mdb106+-galera",
				MARIADB_VERSION_10_10: "my.cnf.mdb106+-galera",
				MARIADB_VERSION_10_9:  "my.cnf.mdb106+-galera",
				MARIADB_VERSION_10_8:  "my.cnf.mdb106+-galera",
				MARIADB_VERSION_10_6:  "my.cnf.mdb106+-galera",
				MARIADB_VERSION_10_5:  "my.cnf.mdb10x-galera",
				MARIADB_VERSION_10_4:  "my.cnf.mdb10x-galera",
				MARIADB_VERSION_10_3:  "my.cnf.mdb10x-galera",
			},
		},
		VENDOR_10GEN: {
			CLUSTER_TYPE_MOGNODB: {
				MONGODB_VERSION_6_0: "mongodb.conf.6.0.percona",
				MONGODB_VERSION_5_0: "mongodb.conf.5.0.percona",
				MONGODB_VERSION_4_4: "mongodb.conf.4.4.percona",
				//MONGODB_VERSION_4_2: "mongodb.conf.4.2.percona",
			},
		},
	}
)

var (
	//gDbAvailableBackupMethods = map[string]map[string]map[string]string{
	gDbAvailableBackupMethods = map[string]map[string][]string{
		CLUSTER_TYPE_REPLICATION: {
			VENDOR_ORACLE: {
				BACKUP_METHOD_XTRABACKUP_FULL,
				BACKUP_METHOD_XTRABACKUP_INCR,
				BACKUP_METHOD_MYSQLDUMP,
			},
			VENDOR_PERCONA: {
				BACKUP_METHOD_XTRABACKUP_FULL,
				BACKUP_METHOD_XTRABACKUP_INCR,
				BACKUP_METHOD_MYSQLDUMP,
			},
			VENDOR_MARIADB: {
				BACKUP_METHOD_MARIABACKUP_FULL,
				BACKUP_METHOD_MARIABACKUP_INCR,
				BACKUP_METHOD_MYSQLDUMP,
			},
		},
		CLUSTER_TYPE_GALERA: {
			VENDOR_PERCONA: {
				BACKUP_METHOD_XTRABACKUP_FULL,
				BACKUP_METHOD_XTRABACKUP_INCR,
				BACKUP_METHOD_MYSQLDUMP,
			},
			VENDOR_MARIADB: {
				BACKUP_METHOD_MARIABACKUP_FULL,
				BACKUP_METHOD_MARIABACKUP_INCR,
				BACKUP_METHOD_MYSQLDUMP,
			},
		},
		CLUSTER_TYPE_MOGNODB: {
			VENDOR_PERCONA: {
				BACKUP_METHOD_PBM,
				BACKUP_METHOD_MONGODUMP,
			},
			VENDOR_10GEN: {
				BACKUP_METHOD_PBM,
				BACKUP_METHOD_MONGODUMP,
			},
		},
		CLUSTER_TYPE_MSSQL_AO_ASYNC: { // VENDOR_MICROSOFT
			VENDOR_MICROSOFT: {
				BACKUP_MSSQL_FULL,
				BACKUP_MSSQL_DIFF,
				BACKUP_MSSQL_TRANSACTION_LOG,
			},
		},
		CLUSTER_TYPE_MSSQL_SINGLE: { // VENDOR_MICROSOFT
			VENDOR_MICROSOFT: {
				BACKUP_MSSQL_FULL,
				BACKUP_MSSQL_DIFF,
				BACKUP_MSSQL_TRANSACTION_LOG,
			},
		},
		CLUSTER_TYPE_REDIS: { // VENDOR_REDIS
			VENDOR_REDIS: {
				"",
			},
		},
		CLUSTER_TYPE_PG_SINGLE: { // VENDOR_REDIS
			VENDOR_DEFAULT: {
				BACKUP_METHOD_PG_BASEBACKUP,
				BACKUP_METHOD_PGDUMPALL,
				BACKUP_METHOD_PG_BACKREST_FULL,
				BACKUP_METHOD_PG_BACKREST_INCR,
				BACKUP_METHOD_PG_BACKREST_DIRR,
			},
			VENDOR_POSTGRESQL: {
				BACKUP_METHOD_PG_BASEBACKUP,
				BACKUP_METHOD_PGDUMPALL,
				BACKUP_METHOD_PG_BACKREST_FULL,
				BACKUP_METHOD_PG_BACKREST_INCR,
				BACKUP_METHOD_PG_BACKREST_DIRR,
			},
		},
	}
)

var (
	CMON_CLUSTERS_OPERATION_SET_NAME        = "name"
	CMON_CLUSTERS_OPERATION_SET_CLUSTER_TAG = "tags"
)

//var foo = map[string]map[string]map[string]string{}
