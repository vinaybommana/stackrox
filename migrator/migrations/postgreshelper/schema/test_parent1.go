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
	// CreateTableTestParent1Stmt holds the create statement for table `test_parent1`.
	CreateTableTestParent1Stmt = &postgres.CreateStmts{
		GormModel: (*TestParent1)(nil),
		Children: []*postgres.CreateStmts{
			&postgres.CreateStmts{
				GormModel: (*TestParent1Childrens)(nil),
				Children:  []*postgres.CreateStmts{},
			},
		},
	}

	// TestParent1Schema is the go schema for table `test_parent1`.
	TestParent1Schema = func() *walker.Schema {
		schema := walker.Walk(reflect.TypeOf((*storage.TestParent1)(nil)), "test_parent1")
		referencedSchemas := map[string]*walker.Schema{
			"storage.TestGrandparent": TestGrandparentsSchema,
			"storage.TestChild1":      TestChild1Schema,
		}

		schema.ResolveReferences(func(messageTypeName string) *walker.Schema {
			return referencedSchemas[fmt.Sprintf("storage.%s", messageTypeName)]
		})
		schema.ScopingResource = &resources.Namespace
		schema.SetOptionsMap(search.Walk(v1.SearchCategory(62), "testparent1", (*storage.TestParent1)(nil)))
		return schema
	}()
)

const (
	// TestParent1TableName specifies the name of the table in postgres.
	TestParent1TableName = "test_parent1"
	// TestParent1ChildrensTableName specifies the name of the table in postgres.
	TestParent1ChildrensTableName = "test_parent1_childrens"
)

// TestParent1 holds the Gorm model for Postgres table `test_parent1`.
type TestParent1 struct {
	ID                  string           `gorm:"column:id;type:varchar;primaryKey"`
	ParentID            string           `gorm:"column:parentid;type:varchar"`
	Val                 string           `gorm:"column:val;type:varchar"`
	Serialized          []byte           `gorm:"column:serialized;type:bytea"`
	TestGrandparentsRef TestGrandparents `gorm:"foreignKey:parentid;references:id;belongsTo;constraint:OnDelete:CASCADE"`
}

// TestParent1Childrens holds the Gorm model for Postgres table `test_parent1_childrens`.
type TestParent1Childrens struct {
	TestParent1ID  string      `gorm:"column:test_parent1_id;type:varchar;primaryKey"`
	Idx            int         `gorm:"column:idx;type:integer;primaryKey;index:testparent1childrens_idx,type:btree"`
	ChildID        string      `gorm:"column:childid;type:varchar"`
	TestParent1Ref TestParent1 `gorm:"foreignKey:test_parent1_id;references:id;belongsTo;constraint:OnDelete:CASCADE"`
}
