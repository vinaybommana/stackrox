// Code generated by pg-bindings generator. DO NOT EDIT.

package schema

import (
	"fmt"
	"reflect"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/walker"
)

var (
	// CreateTableProcessListeningOnPortsStmt holds the create statement for table `process_listening_on_ports`.
	CreateTableProcessListeningOnPortsStmt = &postgres.CreateStmts{
		GormModel: (*ProcessListeningOnPorts)(nil),
		Children:  []*postgres.CreateStmts{},
	}

	// ProcessListeningOnPortsSchema is the go schema for table `process_listening_on_ports`.
	ProcessListeningOnPortsSchema = func() *walker.Schema {
		schema := GetSchemaForTable("process_listening_on_ports")
		if schema != nil {
			return schema
		}
		schema = walker.Walk(reflect.TypeOf((*storage.ProcessListeningOnPortStorage)(nil)), "process_listening_on_ports")
		referencedSchemas := map[string]*walker.Schema{
			"storage.ProcessIndicator": ProcessIndicatorsSchema,
		}

		schema.ResolveReferences(func(messageTypeName string) *walker.Schema {
			return referencedSchemas[fmt.Sprintf("storage.%s", messageTypeName)]
		})
		RegisterTable(schema, CreateTableProcessListeningOnPortsStmt)
		return schema
	}()
)

const (
	ProcessListeningOnPortsTableName = "process_listening_on_ports"
)

// ProcessListeningOnPorts holds the Gorm model for Postgres table `process_listening_on_ports`.
type ProcessListeningOnPorts struct {
	Id                 string `gorm:"column:id;type:uuid;primaryKey"`
	ProcessIndicatorId string `gorm:"column:processindicatorid;type:varchar;index:processlisteningonports_processindicatorid,type:btree"`
	Serialized         []byte `gorm:"column:serialized;type:bytea"`
}