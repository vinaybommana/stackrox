package data

import (
	"bytes"

	"github.com/gogo/protobuf/jsonpb"
	licenseproto "github.com/stackrox/rox/generated/shared/license"
	"google.golang.org/grpc/codes"
)

// GRPCInvocationStats contains telemetry data about GRPC API calls
type GRPCInvocationStats struct {
	Code      codes.Code
	PanicDesc string `json:",omitempty"`

	Count uint64 `json:",omitempty"`
}

// HTTPInvocationStats contains telemetry data about HTTP API calls
type HTTPInvocationStats struct {
	Code      int    `json:",omitempty"` // HTTP status code, or -1 if there was a write error.
	PanicDesc string `json:",omitempty"` // Code location of panic, if the handler panicked.

	Count uint64
}

// APIStat contains telemetry data about different kinds of API calls
type APIStat struct {
	MethodName string
	IsGRPC     bool `json:"isGRPC,omitempty"`

	HTTP []HTTPInvocationStats `json:"http,omitempty"`
	GRPC []GRPCInvocationStats `json:"grpc,omitempty"`
}

// BucketStats contains telemetry data about a DB bucket
type BucketStats struct {
	Name        string
	UsedGB      int
	Cardinality int
}

// DatabaseStats contains telemetry data about a DB
type DatabaseStats struct {
	Type        string
	Path        string
	CapacityGB  int
	UsedGB      int
	StorageType string
	Buckets     []*BucketStats
}

// LicenseJSON type encapsulates the License type and adds Marshal/Unmarshal methods
type LicenseJSON licenseproto.License

// Marshal marshals license data to bytes
func (l *LicenseJSON) Marshal() ([]byte, error) {
	var buf bytes.Buffer
	if err := (&jsonpb.Marshaler{}).Marshal(&buf, (*licenseproto.License)(l)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Unmarshal unmarshals license bytes into a License object
func (l *LicenseJSON) Unmarshal(data []byte) error {
	return jsonpb.Unmarshal(bytes.NewReader(data), (*licenseproto.License)(l))
}

// CentralInfo contains telemetry data specific to StackRox' Central deployment
type CentralInfo struct {
	*RoxComponentInfo

	License      *LicenseJSON
	Database     []*DatabaseStats
	APIStats     []*APIStat
	Orchestrator *OrchestratorInfo

	Clusters []*ClusterInfo
}
