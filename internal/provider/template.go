package provider

import "strings"

// *************************************
// OBSOLETE
// *************************************

type DbClusterTypeBase struct {
	ClusterName           string
	ClusterType           string
	ClusterVendor         string
	ClusterVersion        string
	ClusterConfigTemplate string
}

func NewClusterTypeBase(clusterType string, vendor string, version string) *DbClusterTypeBase {
	ret := DbClusterTypeBase{}
	ret.ClusterType = clusterType
	ret.ClusterVendor = vendor
	ret.ClusterVersion = version

	switch clusterType {
	case CLUSTER_TYPE_REPLICATION:
		if strings.EqualFold(vendor, VENDOR_ORACLE) || strings.EqualFold(vendor, VENDOR_PERCONA) {
			switch version {
			case MYSQL_VERSION_8:
				ret.ClusterConfigTemplate = ""
				break
				//case MYSQL_VERSION_5_7:
				ret.ClusterConfigTemplate = ""
				break
			default:
				// TODO: log WARN
				break
			}
		} else if strings.EqualFold(vendor, VENDOR_MARIADB) {
			switch version {
			case MARIADB_VERSION_10_11:
			case MARIADB_VERSION_10_10:
			case MARIADB_VERSION_10_9:
			case MARIADB_VERSION_10_8:
			case MARIADB_VERSION_10_6:
			case MARIADB_VERSION_10_5:
			case MARIADB_VERSION_10_4:
				//case MARIADB_VERSION_10_3:
				ret.ClusterConfigTemplate = ""
				break
			default:
				// TODO: log WARN
				break
			}
		} else {
			// TODO: log WARN
		}
		break
	case CLUSTER_TYPE_GALERA:
		if strings.EqualFold(vendor, VENDOR_ORACLE) || strings.EqualFold(vendor, VENDOR_PERCONA) {
		} else if strings.EqualFold(vendor, VENDOR_MARIADB) {
		} else {
			// TODO: log WARN
		}
		break
	case CLUSTER_TYPE_PG_SINGLE:
		break
	case CLUSTER_TYPE_MOGNODB:
		break
	case CLUSTER_TYPE_REDIS:
		break
	case CLUSTER_TYPE_MSSQL_SINGLE:
		break
	case CLUSTER_TYPE_MSSQL_AO_ASYNC:
		break
	case CLUSTER_TYPE_ELASTIC:
		break
	default:
		// TODO: log WARN
		break
	}

	switch vendor {
	case VENDOR_MICROSOFT:
		break
	case VENDOR_ELASTIC:
		break
	case VENDOR_REDIS:
		break
	case VENDOR_MARIADB:
		break
	case VENDOR_MONGODB:
		break
	case VENDOR_ORACLE:
		break
	case VENDOR_PERCONA:
		break
	case VENDOR_DEFAULT:
		break
	default:
		// TODO: log WARN
		break
	}

	return &ret
}
