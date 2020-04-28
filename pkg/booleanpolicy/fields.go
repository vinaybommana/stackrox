package booleanpolicy

import (
	"github.com/stackrox/rox/pkg/search"
)

var (
	fieldsToQB = make(map[string]*metadataAndQB)
)

type option int

const (
	negationForbidden option = iota
	operatorsForbidden
)

type metadataAndQB struct {
	operatorsForbidden bool
	negationForbidden  bool
	qb                 queryBuilder
}

// This block enumerates field short names.
var (
	AddCaps                = newFieldWithFieldLabelQueryBuilder("Add Capabilities", search.AddCapabilities, negationForbidden)
	CVE                    = newFieldWithFieldLabelQueryBuilder("CVE", search.CVE)
	CVSS                   = newFieldWithFieldLabelQueryBuilder("CVSS", search.CVSS, negationForbidden)
	ContainerCPULimit      = newFieldWithFieldLabelQueryBuilder("Container CPU Limit", search.CPUCoresLimit, operatorsForbidden)
	ContainerCPURequest    = newFieldWithFieldLabelQueryBuilder("Container CPU Request", search.CPUCoresRequest, operatorsForbidden)
	ContainerMemLimit      = newFieldWithFieldLabelQueryBuilder("Container Memory Limit", search.MemoryLimit, operatorsForbidden)
	ContainerMemRequest    = newFieldWithFieldLabelQueryBuilder("Container Memory Request", search.MemoryRequest, operatorsForbidden)
	DisallowedAnnotation   = newFieldWithTODOQueryBuilder("Disallowed Annotation")
	DisallowedImageLabel   = newFieldWithTODOQueryBuilder("Disallowed Image Label")
	DockerfileLine         = newFieldWithTODOQueryBuilder("Dockerfile Line")
	DropCaps               = newFieldWithFieldLabelQueryBuilder("Drop Capabilities", search.DropCapabilities, negationForbidden)
	EnvironmentVariable    = newFieldWithTODOQueryBuilder("Environment Variable")
	FixedBy                = newFieldWithFieldLabelQueryBuilder("FixedBy", search.FixedBy)
	ImageAge               = newFieldWithTODOQueryBuilder("Image Age")
	ImageComponent         = newFieldWithFieldLabelQueryBuilder("Image Component", search.Component)
	ImageRegistry          = newFieldWithFieldLabelQueryBuilder("Image Registry", search.ImageRegistry)
	ImageRemote            = newFieldWithFieldLabelQueryBuilder("Image Remote", search.ImageRemote, negationForbidden)
	ImageScanAge           = newFieldWithTODOQueryBuilder("Image Scan Age")
	ImageTag               = newFieldWithFieldLabelQueryBuilder("Image Tag", search.ImageTag)
	MinimumRBACPermissions = newFieldWithTODOQueryBuilder("Minimum RBAC Permissions")
	Port                   = newFieldWithFieldLabelQueryBuilder("Port", search.Port)
	PortExposure           = newFieldWithTODOQueryBuilder("Port Exposure Method")
	Privileged             = newFieldWithFieldLabelQueryBuilder("Privileged", search.Privileged, negationForbidden, operatorsForbidden)
	ProcessAncestor        = newFieldWithFieldLabelQueryBuilder("Process Ancestor", search.ProcessAncestor)
	ProcessArguments       = newFieldWithFieldLabelQueryBuilder("Process Arguments", search.ProcessArguments)
	ProcessName            = newFieldWithFieldLabelQueryBuilder("Process Name", search.ProcessName)
	ProcessUID             = newFieldWithFieldLabelQueryBuilder("Process UID", search.ProcessUID)
	Protocol               = newFieldWithFieldLabelQueryBuilder("Protocol", search.PortProtocol)
	ReadOnlyRootFS         = newFieldWithFieldLabelQueryBuilder("Read-Only Root Filesystem", search.ReadOnlyRootFilesystem, negationForbidden, operatorsForbidden)
	RequiredAnnotation     = newFieldWithTODOQueryBuilder("Required Annotation")
	RequiredImageLabel     = newFieldWithTODOQueryBuilder("Required Image Label")
	RequiredLabel          = newFieldWithTODOQueryBuilder("Required Label")
	UnscannedImage         = newFieldWithTODOQueryBuilder("Unscanned Image")
	VolumeDestination      = newFieldWithFieldLabelQueryBuilder("Volume Destination", search.VolumeDestination)
	VolumeName             = newFieldWithFieldLabelQueryBuilder("Volume Name", search.VolumeName)
	VolumeSource           = newFieldWithFieldLabelQueryBuilder("Volume Source", search.VolumeSource)
	VolumeType             = newFieldWithFieldLabelQueryBuilder("Volume Type", search.VolumeType)
	WhitelistsEnabled      = newFieldWithTODOQueryBuilder("Unexpected Process Executed")
	WritableHostMount      = newFieldWithTODOQueryBuilder("Writable Host Mount")
	WritableVolume         = newFieldWithTODOQueryBuilder("Writable Volume")
)

func newFieldWithFieldLabelQueryBuilder(fieldName string, fieldLabel search.FieldLabel, options ...option) string {
	m := metadataAndQB{
		qb: &fieldLabelBasedQueryBuilder{fieldName: fieldName, fieldLabel: fieldLabel},
	}
	for _, o := range options {
		switch o {
		case negationForbidden:
			m.negationForbidden = true
		case operatorsForbidden:
			m.operatorsForbidden = true
		}
	}
	fieldsToQB[fieldName] = &m
	return fieldName
}

// TEMPORARY: this indicates fields for which we haven't indicated the relevant query builder yet.
func newFieldWithTODOQueryBuilder(fieldName string) string {
	return fieldName
}
