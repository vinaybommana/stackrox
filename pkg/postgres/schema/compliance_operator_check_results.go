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
	// CreateTableComplianceOperatorCheckResultsStmt holds the create statement for table `compliance_operator_check_results`.
	CreateTableComplianceOperatorCheckResultsStmt = &postgres.CreateStmts{
		GormModel: (*ComplianceOperatorCheckResults)(nil),
		Children:  []*postgres.CreateStmts{},
	}

	// ComplianceOperatorCheckResultsSchema is the go schema for table `compliance_operator_check_results`.
	ComplianceOperatorCheckResultsSchema = func() *walker.Schema {
		schema := GetSchemaForTable("compliance_operator_check_results")
		if schema != nil {
			return schema
		}
		schema = walker.Walk(reflect.TypeOf((*storage.ComplianceOperatorCheckResult)(nil)), "compliance_operator_check_results")
		schema.ScopingResource = &resources.ComplianceOperator
		RegisterTable(schema, CreateTableComplianceOperatorCheckResultsStmt)
		return schema
	}()
)

const (
	// ComplianceOperatorCheckResultsTableName specifies the name of the table in postgres.
	ComplianceOperatorCheckResultsTableName = "compliance_operator_check_results"
)

// ComplianceOperatorCheckResults holds the Gorm model for Postgres table `compliance_operator_check_results`.
type ComplianceOperatorCheckResults struct {
	ID         string `gorm:"column:id;type:varchar;primaryKey"`
	Serialized []byte `gorm:"column:serialized;type:bytea"`
}
