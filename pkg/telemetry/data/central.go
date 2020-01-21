package data

import (
	"bytes"

	"github.com/gogo/protobuf/jsonpb"
	licenseproto "github.com/stackrox/rox/generated/shared/license"
	"google.golang.org/grpc/codes"
)

// GRPCInvocationStats contains telemetry data about GRPC API calls
type GRPCInvocationStats struct {
	Code  codes.Code `json:"code"`
	Count int64      `json:"count"`
}

// HTTPInvocationStats contains telemetry data about HTTP API calls
type HTTPInvocationStats struct {
	StatusCode int   `json:"statusCode"` // HTTP status code, or -1 if there was a write error.
	Count      int64 `json:"count"`
}

// PanicStats contains telemetry data about API panics
type PanicStats struct {
	PanicDesc string `json:"panicDesc"` // Code location of panic, if the handler panicked.
	Count     int64  `json:"count"`
}

// APIStat contains telemetry data about different kinds of API calls
type APIStat struct {
	MethodName string `json:"methodName"`

	HTTP   []HTTPInvocationStats `json:"http,omitempty"`
	GRPC   []GRPCInvocationStats `json:"grpc,omitempty"`
	Panics []PanicStats          `json:"panics,omitempty"`
}

// APIInfo contains metrics about API calls and errors gathering those metrics
type APIInfo struct {
	APIStats []*APIStat `json:"apiStats,omitempty"`
}

// BucketStats contains telemetry data about a DB bucket
type BucketStats struct {
	Name        string `json:"name"`
	UsedBytes   int64  `json:"usedBytes"`
	Cardinality int    `json:"cardinality"`
}

// DatabaseStats contains telemetry data about a DB
type DatabaseStats struct {
	Type      string         `json:"type"`
	Path      string         `json:"path"`
	UsedBytes int64          `json:"usedBytes"`
	Buckets   []*BucketStats `json:"buckets,omitempty"`
	Errors    []string       `json:"errors,omitempty"`
}

// StorageInfo contains telemetry data about available disk, storage type, and the available databases
type StorageInfo struct {
	DiskCapacityBytes int64            `json:"diskCapacityBytes"`
	DiskUsedBytes     int64            `json:"diskUsedBytes"`
	StorageType       string           `json:"storageType,omitempty"`
	Databases         []*DatabaseStats `json:"dbs,omitempty"`
	Errors            []string         `json:"errors,omitempty"`
}

// LicenseJSON type encapsulates the License type and adds Marshal/Unmarshal methods
type LicenseJSON licenseproto.License

// MarshalJSON marshals license data to bytes, following jsonpb rules.
func (l *LicenseJSON) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	if err := (&jsonpb.Marshaler{}).Marshal(&buf, (*licenseproto.License)(l)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnmarshalJSON unmarshals license JSON bytes into a License object, following jsonpb rules.
func (l *LicenseJSON) UnmarshalJSON(data []byte) error {
	return jsonpb.Unmarshal(bytes.NewReader(data), (*licenseproto.License)(l))
}

// CentralInfo contains telemetry data specific to StackRox' Central deployment
type CentralInfo struct {
	*RoxComponentInfo

	License      *LicenseJSON      `json:"license,omitempty"`
	Storage      *StorageInfo      `json:"storage,omitempty"`
	APIStats     *APIInfo          `json:"apiStats,omitempty"`
	Orchestrator *OrchestratorInfo `json:"orchestrator,omitempty"`

	Errors []string `json:"errors,omitempty"`
}
