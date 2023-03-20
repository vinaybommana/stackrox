// Code generated by pg-bindings generator. DO NOT EDIT.

package schema

import (
	"reflect"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/walker"
	"github.com/stackrox/rox/pkg/sac/resources"
)

var (
	// CreateTableSystemInfosStmt holds the create statement for table `system_infos`.
	CreateTableSystemInfosStmt = &postgres.CreateStmts{
		GormModel: (*SystemInfos)(nil),
		Children:  []*postgres.CreateStmts{},
	}

	// SystemInfosSchema is the go schema for table `system_infos`.
	SystemInfosSchema = func() *walker.Schema {
		schema := GetSchemaForTable("system_infos")
		if schema != nil {
			return schema
		}
		schema = walker.Walk(reflect.TypeOf((*storage.SystemInfo)(nil)), "system_infos")
		schema.ScopingResource = &resources.Administration
		RegisterTable(schema, CreateTableSystemInfosStmt)
		return schema
	}()
)

const (
	// SystemInfosTableName specifies the name of the table in postgres.
	SystemInfosTableName = "system_infos"
)

// SystemInfos holds the Gorm model for Postgres table `system_infos`.
type SystemInfos struct {
	BackupInfoRequestorName string `gorm:"column:backupinfo_requestor_name;type:varchar"`
	Serialized              []byte `gorm:"column:serialized;type:bytea"`
}
