// Code generated by pg-bindings generator. DO NOT EDIT.

package schema

import (
	"fmt"
	"reflect"

	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/walker"
	"github.com/stackrox/rox/pkg/sac/resources"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/search/postgres/mapping"
)

var (
	// CreateTableProcessIndicatorsStmt holds the create statement for table `process_indicators`.
	CreateTableProcessIndicatorsStmt = &postgres.CreateStmts{
		GormModel: (*ProcessIndicators)(nil),
		Children:  []*postgres.CreateStmts{},
	}

	// ProcessIndicatorsSchema is the go schema for table `process_indicators`.
	ProcessIndicatorsSchema = func() *walker.Schema {
		schema := GetSchemaForTable("process_indicators")
		if schema != nil {
			return schema
		}
		schema = walker.Walk(reflect.TypeOf((*storage.ProcessIndicator)(nil)), "process_indicators")
		schema.ScopingResource = &resources.DeploymentExtension
		referencedSchemas := map[string]*walker.Schema{
			"storage.Deployment": DeploymentsSchema,
		}

		schema.ResolveReferences(func(messageTypeName string) *walker.Schema {
			return referencedSchemas[fmt.Sprintf("storage.%s", messageTypeName)]
		})
		schema.SetOptionsMap(search.Walk(v1.SearchCategory_PROCESS_INDICATORS, "processindicator", (*storage.ProcessIndicator)(nil)))
		RegisterTable(schema, CreateTableProcessIndicatorsStmt)
		mapping.RegisterCategoryToTable(v1.SearchCategory_PROCESS_INDICATORS, schema)
		return schema
	}()
)

const (
	// ProcessIndicatorsTableName specifies the name of the table in postgres.
	ProcessIndicatorsTableName = "process_indicators"
)

// ProcessIndicators holds the Gorm model for Postgres table `process_indicators`.
type ProcessIndicators struct {
	ID                 string `gorm:"column:id;type:uuid;primaryKey"`
	DeploymentID       string `gorm:"column:deploymentid;type:uuid;index:processindicators_deploymentid,type:hash"`
	ContainerName      string `gorm:"column:containername;type:varchar"`
	PodID              string `gorm:"column:podid;type:varchar"`
	PodUID             string `gorm:"column:poduid;type:uuid;index:processindicators_poduid,type:hash"`
	SignalContainerID  string `gorm:"column:signal_containerid;type:varchar"`
	SignalName         string `gorm:"column:signal_name;type:varchar"`
	SignalArgs         string `gorm:"column:signal_args;type:varchar"`
	SignalExecFilePath string `gorm:"column:signal_execfilepath;type:varchar"`
	SignalUID          uint32 `gorm:"column:signal_uid;type:bigint"`
	ClusterID          string `gorm:"column:clusterid;type:uuid;index:processindicators_sac_filter,type:btree"`
	Namespace          string `gorm:"column:namespace;type:varchar;index:processindicators_sac_filter,type:btree"`
	Serialized         []byte `gorm:"column:serialized;type:bytea"`
}
