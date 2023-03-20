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
	// CreateTableComplianceOperatorRulesStmt holds the create statement for table `compliance_operator_rules`.
	CreateTableComplianceOperatorRulesStmt = &postgres.CreateStmts{
		GormModel: (*ComplianceOperatorRules)(nil),
		Children:  []*postgres.CreateStmts{},
	}

	// ComplianceOperatorRulesSchema is the go schema for table `compliance_operator_rules`.
	ComplianceOperatorRulesSchema = func() *walker.Schema {
		schema := GetSchemaForTable("compliance_operator_rules")
		if schema != nil {
			return schema
		}
		schema = walker.Walk(reflect.TypeOf((*storage.ComplianceOperatorRule)(nil)), "compliance_operator_rules")
		schema.ScopingResource = &resources.ComplianceOperator
		RegisterTable(schema, CreateTableComplianceOperatorRulesStmt)
		return schema
	}()
)

const (
	// ComplianceOperatorRulesTableName specifies the name of the table in postgres.
	ComplianceOperatorRulesTableName = "compliance_operator_rules"
)

// ComplianceOperatorRules holds the Gorm model for Postgres table `compliance_operator_rules`.
type ComplianceOperatorRules struct {
	ID         string `gorm:"column:id;type:varchar;primaryKey"`
	Serialized []byte `gorm:"column:serialized;type:bytea"`
}
