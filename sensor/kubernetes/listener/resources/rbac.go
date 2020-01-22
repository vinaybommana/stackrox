package resources

import (
	"github.com/stackrox/rox/generated/internalapi/central"
	v1 "k8s.io/api/rbac/v1"
)

// Handle RBAC-related events
type rbacDispatcher struct {
	store rbacUpdater
}

func newRBACDispatcher(store rbacUpdater) *rbacDispatcher {
	return &rbacDispatcher{
		store: store,
	}
}

func (r *rbacDispatcher) ProcessEvent(obj interface{}, action central.ResourceAction) []*central.SensorEvent {
	switch obj := obj.(type) {
	case *v1.Role:
		if action == central.ResourceAction_REMOVE_RESOURCE {
			return r.store.removeRole(obj)
		}
		return r.store.upsertRole(obj)
	case *v1.RoleBinding:
		if action == central.ResourceAction_REMOVE_RESOURCE {
			return r.store.removeBinding(obj)
		}
		return r.store.upsertBinding(obj)
	case *v1.ClusterRole:
		if action == central.ResourceAction_REMOVE_RESOURCE {
			return r.store.removeClusterRole(obj)
		}
		return r.store.upsertClusterRole(obj)
	case *v1.ClusterRoleBinding:
		if action == central.ResourceAction_REMOVE_RESOURCE {
			return r.store.removeClusterBinding(obj)
		}
		return r.store.upsertClusterBinding(obj)
	}
	return nil
}
