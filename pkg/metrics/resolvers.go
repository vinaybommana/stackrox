package metrics

// Resolver represents a graphql resolver that we want to time.
//go:generate stringer -type=Resolver
type Resolver int

// The following is the list of graphql resolvers that we want to time.
const (
	Cluster Resolver = iota
	Compliance
	Deployments
	Groups
	Images
	K8sRoles
	Namespaces
	Nodes
	Notifiers
	Policies
	Roles
	Secrets
	ServiceAccounts
	Subjects
	Tokens
	Violations
)
