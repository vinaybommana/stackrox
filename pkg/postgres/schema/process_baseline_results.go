// Code generated by pg-bindings generator. DO NOT EDIT.

package schema

import (
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
	// CreateTableProcessBaselineResultsStmt holds the create statement for table `process_baseline_results`.
	CreateTableProcessBaselineResultsStmt = &postgres.CreateStmts{
		GormModel: (*ProcessBaselineResults)(nil),
		Children:  []*postgres.CreateStmts{},
	}

	// ProcessBaselineResultsSchema is the go schema for table `process_baseline_results`.
	ProcessBaselineResultsSchema = func() *walker.Schema {
		schema := GetSchemaForTable("process_baseline_results")
		if schema != nil {
			return schema
		}
		schema = walker.Walk(reflect.TypeOf((*storage.ProcessBaselineResults)(nil)), "process_baseline_results")
		schema.ScopingResource = &resources.DeploymentExtension
		schema.SetOptionsMap(search.Walk(v1.SearchCategory_PROCESS_BASELINE_RESULTS, "processbaselineresults", (*storage.ProcessBaselineResults)(nil)))
		RegisterTable(schema, CreateTableProcessBaselineResultsStmt)
		mapping.RegisterCategoryToTable(v1.SearchCategory_PROCESS_BASELINE_RESULTS, schema)
		return schema
	}()
)

const (
	// ProcessBaselineResultsTableName specifies the name of the table in postgres.
	ProcessBaselineResultsTableName = "process_baseline_results"
)

// ProcessBaselineResults holds the Gorm model for Postgres table `process_baseline_results`.
type ProcessBaselineResults struct {
	DeploymentID string `gorm:"column:deploymentid;type:uuid;primaryKey"`
	ClusterID    string `gorm:"column:clusterid;type:uuid;index:processbaselineresults_sac_filter,type:btree"`
	Namespace    string `gorm:"column:namespace;type:varchar;index:processbaselineresults_sac_filter,type:btree"`
	Serialized   []byte `gorm:"column:serialized;type:bytea"`
}
