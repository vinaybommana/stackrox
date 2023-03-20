// Code generated by pg-bindings generator. DO NOT EDIT.

package schema

import (
	"reflect"

	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/walker"
	"github.com/stackrox/rox/pkg/sac/resources"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/search/postgres/mapping"
)

var (
	// CreateTableRoleBindingsStmt holds the create statement for table `role_bindings`.
	CreateTableRoleBindingsStmt = &postgres.CreateStmts{
		GormModel: (*RoleBindings)(nil),
		Children: []*postgres.CreateStmts{
			&postgres.CreateStmts{
				GormModel: (*RoleBindingsSubjects)(nil),
				Children:  []*postgres.CreateStmts{},
			},
		},
	}

	// RoleBindingsSchema is the go schema for table `role_bindings`.
	RoleBindingsSchema = func() *walker.Schema {
		schema := GetSchemaForTable("role_bindings")
		if schema != nil {
			return schema
		}
		schema = walker.Walk(reflect.TypeOf((*storage.K8SRoleBinding)(nil)), "role_bindings")
		schema.ScopingResource = &resources.K8sRoleBinding
		schema.SetOptionsMap(search.Walk(v1.SearchCategory_ROLEBINDINGS, "k8srolebinding", (*storage.K8SRoleBinding)(nil)))
		RegisterTable(schema, CreateTableRoleBindingsStmt)
		mapping.RegisterCategoryToTable(v1.SearchCategory_ROLEBINDINGS, schema)
		return schema
	}()
)

const (
	// RoleBindingsTableName specifies the name of the table in postgres.
	RoleBindingsTableName = "role_bindings"
	// RoleBindingsSubjectsTableName specifies the name of the table in postgres.
	RoleBindingsSubjectsTableName = "role_bindings_subjects"
)

// RoleBindings holds the Gorm model for Postgres table `role_bindings`.
type RoleBindings struct {
	ID          string            `gorm:"column:id;type:uuid;primaryKey"`
	Name        string            `gorm:"column:name;type:varchar"`
	Namespace   string            `gorm:"column:namespace;type:varchar;index:rolebindings_sac_filter,type:btree"`
	ClusterID   string            `gorm:"column:clusterid;type:uuid;index:rolebindings_sac_filter,type:btree"`
	ClusterName string            `gorm:"column:clustername;type:varchar"`
	ClusterRole bool              `gorm:"column:clusterrole;type:bool"`
	Labels      map[string]string `gorm:"column:labels;type:jsonb"`
	Annotations map[string]string `gorm:"column:annotations;type:jsonb"`
	RoleID      string            `gorm:"column:roleid;type:uuid"`
	Serialized  []byte            `gorm:"column:serialized;type:bytea"`
}

// RoleBindingsSubjects holds the Gorm model for Postgres table `role_bindings_subjects`.
type RoleBindingsSubjects struct {
	RoleBindingsID  string              `gorm:"column:role_bindings_id;type:uuid;primaryKey"`
	Idx             int                 `gorm:"column:idx;type:integer;primaryKey;index:rolebindingssubjects_idx,type:btree"`
	Kind            storage.SubjectKind `gorm:"column:kind;type:integer"`
	Name            string              `gorm:"column:name;type:varchar"`
	RoleBindingsRef RoleBindings        `gorm:"foreignKey:role_bindings_id;references:id;belongsTo;constraint:OnDelete:CASCADE"`
}
