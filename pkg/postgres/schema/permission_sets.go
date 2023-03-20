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
	// CreateTablePermissionSetsStmt holds the create statement for table `permission_sets`.
	CreateTablePermissionSetsStmt = &postgres.CreateStmts{
		GormModel: (*PermissionSets)(nil),
		Children:  []*postgres.CreateStmts{},
	}

	// PermissionSetsSchema is the go schema for table `permission_sets`.
	PermissionSetsSchema = func() *walker.Schema {
		schema := GetSchemaForTable("permission_sets")
		if schema != nil {
			return schema
		}
		schema = walker.Walk(reflect.TypeOf((*storage.PermissionSet)(nil)), "permission_sets")
		schema.ScopingResource = &resources.Role
		RegisterTable(schema, CreateTablePermissionSetsStmt)
		return schema
	}()
)

const (
	// PermissionSetsTableName specifies the name of the table in postgres.
	PermissionSetsTableName = "permission_sets"
)

// PermissionSets holds the Gorm model for Postgres table `permission_sets`.
type PermissionSets struct {
	ID         string `gorm:"column:id;type:uuid;primaryKey"`
	Name       string `gorm:"column:name;type:varchar;unique"`
	Serialized []byte `gorm:"column:serialized;type:bytea"`
}
