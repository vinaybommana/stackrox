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
)

var (
	// CreateTableTestParent2Stmt holds the create statement for table `test_parent2`.
	CreateTableTestParent2Stmt = &postgres.CreateStmts{
		GormModel: (*TestParent2)(nil),
		Children:  []*postgres.CreateStmts{},
	}

	// TestParent2Schema is the go schema for table `test_parent2`.
	TestParent2Schema = func() *walker.Schema {
		schema := walker.Walk(reflect.TypeOf((*storage.TestParent2)(nil)), "test_parent2")
		referencedSchemas := map[string]*walker.Schema{
			"storage.TestGrandparent": TestGrandparentsSchema,
		}

		schema.ResolveReferences(func(messageTypeName string) *walker.Schema {
			return referencedSchemas[fmt.Sprintf("storage.%s", messageTypeName)]
		})
		schema.ScopingResource = &resources.Namespace
		schema.SetOptionsMap(search.Walk(v1.SearchCategory(68), "testparent2", (*storage.TestParent2)(nil)))
		return schema
	}()
)

const (
	// TestParent2TableName specifies the name of the table in postgres.
	TestParent2TableName = "test_parent2"
)

// TestParent2 holds the Gorm model for Postgres table `test_parent2`.
type TestParent2 struct {
	ID                  string           `gorm:"column:id;type:uuid;primaryKey"`
	ParentID            string           `gorm:"column:parentid;type:varchar"`
	Val                 string           `gorm:"column:val;type:varchar"`
	Serialized          []byte           `gorm:"column:serialized;type:bytea"`
	TestGrandparentsRef TestGrandparents `gorm:"foreignKey:parentid;references:id;belongsTo;constraint:OnDelete:CASCADE"`
}
