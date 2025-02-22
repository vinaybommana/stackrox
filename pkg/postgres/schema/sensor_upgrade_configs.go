// Code generated by pg-bindings generator. DO NOT EDIT.

package schema

import (
	"reflect"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/walker"
)

var (
	// CreateTableSensorUpgradeConfigsStmt holds the create statement for table `sensor_upgrade_configs`.
	CreateTableSensorUpgradeConfigsStmt = &postgres.CreateStmts{
		GormModel: (*SensorUpgradeConfigs)(nil),
		Children:  []*postgres.CreateStmts{},
	}

	// SensorUpgradeConfigsSchema is the go schema for table `sensor_upgrade_configs`.
	SensorUpgradeConfigsSchema = func() *walker.Schema {
		schema := GetSchemaForTable("sensor_upgrade_configs")
		if schema != nil {
			return schema
		}
		schema = walker.Walk(reflect.TypeOf((*storage.SensorUpgradeConfig)(nil)), "sensor_upgrade_configs")
		RegisterTable(schema, CreateTableSensorUpgradeConfigsStmt)
		return schema
	}()
)

const (
	// SensorUpgradeConfigsTableName specifies the name of the table in postgres.
	SensorUpgradeConfigsTableName = "sensor_upgrade_configs"
)

// SensorUpgradeConfigs holds the Gorm model for Postgres table `sensor_upgrade_configs`.
type SensorUpgradeConfigs struct {
	Serialized []byte `gorm:"column:serialized;type:bytea"`
}
