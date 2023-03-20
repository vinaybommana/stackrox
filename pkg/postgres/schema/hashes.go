// Code generated by pg-bindings generator. DO NOT EDIT.

package schema

import (
	"reflect"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/walker"
	"github.com/stackrox/rox/pkg/sac/resources"
)

var (
	// CreateTableHashesStmt holds the create statement for table `hashes`.
	CreateTableHashesStmt = &postgres.CreateStmts{
		GormModel: (*Hashes)(nil),
		Children:  []*postgres.CreateStmts{},
	}

	// HashesSchema is the go schema for table `hashes`.
	HashesSchema = func() *walker.Schema {
		schema := GetSchemaForTable("hashes")
		if schema != nil {
			return schema
		}
		schema = walker.Walk(reflect.TypeOf((*storage.Hash)(nil)), "hashes")
		schema.ScopingResource = &resources.Hash
		RegisterTable(schema, CreateTableHashesStmt, features.StoreEventHashes.Enabled)
		return schema
	}()
)

const (
	// HashesTableName specifies the name of the table in postgres.
	HashesTableName = "hashes"
)

// Hashes holds the Gorm model for Postgres table `hashes`.
type Hashes struct {
	ClusterID  string `gorm:"column:clusterid;type:varchar;primaryKey"`
	Serialized []byte `gorm:"column:serialized;type:bytea"`
}
