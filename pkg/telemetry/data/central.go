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
	UsedBytes   int64
	Cardinality int
}

// DatabaseStats contains telemetry data about a DB
type DatabaseStats struct {
	Type      string
	Path      string
	UsedBytes int64
	Buckets   []*BucketStats
	Errors    []string
}

// StorageInfo contains telemetry data about available disk, storage type, and the available databases
type StorageInfo struct {
	DiskCapacityBytes int64
	DiskUsedBytes     int64
	StorageType       string
	Database          []*DatabaseStats
	Errors            []string
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
	Storage      *StorageInfo
	APIStats     []*APIStat
	Orchestrator *OrchestratorInfo

	Clusters []*ClusterInfo

	Errors []string
}
