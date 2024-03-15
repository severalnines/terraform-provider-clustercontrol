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

var (
	gDefultDbAdminUser = map[string]string{
		CLUSTER_TYPE_REPLICATION:    "root",
		CLUSTER_TYPE_GALERA:         "root",
		CLUSTER_TYPE_PG_SINGLE:      "pgadmin",
		CLUSTER_TYPE_MOGNODB:        "mongoadmin",
		CLUSTER_TYPE_REDIS:          "",
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
		CLUSTER_TYPE_REDIS:          "",
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

//var foo = map[string]map[string]map[string]string{}
