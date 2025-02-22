// Code generated by pg-bindings generator. DO NOT EDIT.

package schema

import (
	"fmt"
	"reflect"
	"time"

	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/walker"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/search/postgres/mapping"
)

var (
	// CreateTableClusterHealthStatusesStmt holds the create statement for table `cluster_health_statuses`.
	CreateTableClusterHealthStatusesStmt = &postgres.CreateStmts{
		GormModel: (*ClusterHealthStatuses)(nil),
		Children:  []*postgres.CreateStmts{},
	}

	// ClusterHealthStatusesSchema is the go schema for table `cluster_health_statuses`.
	ClusterHealthStatusesSchema = func() *walker.Schema {
		schema := GetSchemaForTable("cluster_health_statuses")
		if schema != nil {
			return schema
		}
		schema = walker.Walk(reflect.TypeOf((*storage.ClusterHealthStatus)(nil)), "cluster_health_statuses")
		referencedSchemas := map[string]*walker.Schema{
			"storage.Cluster": ClustersSchema,
		}

		schema.ResolveReferences(func(messageTypeName string) *walker.Schema {
			return referencedSchemas[fmt.Sprintf("storage.%s", messageTypeName)]
		})
		schema.SetOptionsMap(search.Walk(v1.SearchCategory_CLUSTER_HEALTH, "clusterhealthstatus", (*storage.ClusterHealthStatus)(nil)))
		RegisterTable(schema, CreateTableClusterHealthStatusesStmt)
		mapping.RegisterCategoryToTable(v1.SearchCategory_CLUSTER_HEALTH, schema)
		return schema
	}()
)

const (
	// ClusterHealthStatusesTableName specifies the name of the table in postgres.
	ClusterHealthStatusesTableName = "cluster_health_statuses"
)

// ClusterHealthStatuses holds the Gorm model for Postgres table `cluster_health_statuses`.
type ClusterHealthStatuses struct {
	ID                           string                                        `gorm:"column:id;type:uuid;primaryKey"`
	SensorHealthStatus           storage.ClusterHealthStatus_HealthStatusLabel `gorm:"column:sensorhealthstatus;type:integer"`
	CollectorHealthStatus        storage.ClusterHealthStatus_HealthStatusLabel `gorm:"column:collectorhealthstatus;type:integer"`
	OverallHealthStatus          storage.ClusterHealthStatus_HealthStatusLabel `gorm:"column:overallhealthstatus;type:integer"`
	AdmissionControlHealthStatus storage.ClusterHealthStatus_HealthStatusLabel `gorm:"column:admissioncontrolhealthstatus;type:integer"`
	ScannerHealthStatus          storage.ClusterHealthStatus_HealthStatusLabel `gorm:"column:scannerhealthstatus;type:integer"`
	LastContact                  *time.Time                                    `gorm:"column:lastcontact;type:timestamp"`
	Serialized                   []byte                                        `gorm:"column:serialized;type:bytea"`
}
