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
	// CreateTableReportConfigurationsStmt holds the create statement for table `report_configurations`.
	CreateTableReportConfigurationsStmt = &postgres.CreateStmts{
		GormModel: (*ReportConfigurations)(nil),
		Children:  []*postgres.CreateStmts{},
	}

	// ReportConfigurationsSchema is the go schema for table `report_configurations`.
	ReportConfigurationsSchema = func() *walker.Schema {
		schema := GetSchemaForTable("report_configurations")
		if schema != nil {
			return schema
		}
		schema = walker.Walk(reflect.TypeOf((*storage.ReportConfiguration)(nil)), "report_configurations")
		schema.ScopingResource = &resources.VulnerabilityReports
		schema.SetOptionsMap(search.Walk(v1.SearchCategory_REPORT_CONFIGURATIONS, "reportconfiguration", (*storage.ReportConfiguration)(nil)))
		RegisterTable(schema, CreateTableReportConfigurationsStmt)
		mapping.RegisterCategoryToTable(v1.SearchCategory_REPORT_CONFIGURATIONS, schema)
		return schema
	}()
)

const (
	// ReportConfigurationsTableName specifies the name of the table in postgres.
	ReportConfigurationsTableName = "report_configurations"
)

// ReportConfigurations holds the Gorm model for Postgres table `report_configurations`.
type ReportConfigurations struct {
	ID                        string                                 `gorm:"column:id;type:varchar;primaryKey"`
	Name                      string                                 `gorm:"column:name;type:varchar"`
	Type                      storage.ReportConfiguration_ReportType `gorm:"column:type;type:integer"`
	ScopeID                   string                                 `gorm:"column:scopeid;type:varchar"`
	ResourceScopeCollectionID string                                 `gorm:"column:resourcescope_collectionid;type:varchar"`
	Serialized                []byte                                 `gorm:"column:serialized;type:bytea"`
}
