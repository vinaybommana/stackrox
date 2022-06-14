// Code generated by pg-bindings generator. DO NOT EDIT.

package schema

import (
	"fmt"
	"reflect"

	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/walker"
	"github.com/stackrox/rox/pkg/search"
)

var (
	// CreateTableNodeComponentEdgesStmt holds the create statement for table `node_component_edges`.
	CreateTableNodeComponentEdgesStmt = &postgres.CreateStmts{
		Table: `
               create table if not exists node_component_edges (
                   Id varchar,
                   NodeId varchar,
                   NodeComponentId varchar,
                   serialized bytea,
                   PRIMARY KEY(Id),
                   CONSTRAINT fk_parent_table_0 FOREIGN KEY (NodeId) REFERENCES nodes(Id) ON DELETE CASCADE
               )
               `,
		GormModel: (*NodeComponentEdges)(nil),
		Indexes:   []string{},
		Children:  []*postgres.CreateStmts{},
	}

	// NodeComponentEdgesSchema is the go schema for table `node_component_edges`.
	NodeComponentEdgesSchema = func() *walker.Schema {
		schema := GetSchemaForTable("node_component_edges")
		if schema != nil {
			return schema
		}
		schema = walker.Walk(reflect.TypeOf((*storage.NodeComponentEdge)(nil)), "node_component_edges")
		referencedSchemas := map[string]*walker.Schema{
			"storage.Node":          NodesSchema,
			"storage.NodeComponent": NodeComponentsSchema,
		}

		schema.ResolveReferences(func(messageTypeName string) *walker.Schema {
			return referencedSchemas[fmt.Sprintf("storage.%s", messageTypeName)]
		})
		schema.SetOptionsMap(search.Walk(v1.SearchCategory_NODE_COMPONENT_EDGE, "nodecomponentedge", (*storage.NodeComponentEdge)(nil)))
		RegisterTable(schema, CreateTableNodeComponentEdgesStmt)
		return schema
	}()
)

const (
	NodeComponentEdgesTableName = "node_component_edges"
)

// NodeComponentEdges holds the Gorm model for Postgres table `node_component_edges`.
type NodeComponentEdges struct {
	Id              string `gorm:"column:id;type:varchar;primaryKey"`
	NodeId          string `gorm:"column:nodeid;type:varchar"`
	NodeComponentId string `gorm:"column:nodecomponentid;type:varchar"`
	Serialized      []byte `gorm:"column:serialized;type:bytea"`
	NodesRef        Nodes  `gorm:"foreignKey:nodeid;references:id;belongsTo;constraint:OnDelete:CASCADE"`
}
