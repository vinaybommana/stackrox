package booleanpolicy

import (
	"regexp"

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
	valueRegex         *regexp.Regexp
}

// This block enumerates field short names.
var (
	AddCaps                = newFieldWithFieldLabelQueryBuilder("Add Capabilities", search.AddCapabilities, capabilitiesValueRegex, negationForbidden)
	CVE                    = newFieldWithFieldLabelQueryBuilder("CVE", search.CVE, stringValueRegex)
	CVSS                   = newFieldWithFieldLabelQueryBuilder("CVSS", search.CVSS, comparatorDecimalValueRegex, operatorsForbidden)
	ContainerCPULimit      = newFieldWithFieldLabelQueryBuilder("Container CPU Limit", search.CPUCoresLimit, integerValueRegex, operatorsForbidden)
	ContainerCPURequest    = newFieldWithFieldLabelQueryBuilder("Container CPU Request", search.CPUCoresRequest, integerValueRegex, operatorsForbidden)
	ContainerMemLimit      = newFieldWithFieldLabelQueryBuilder("Container Memory Limit", search.MemoryLimit, integerValueRegex, operatorsForbidden)
	ContainerMemRequest    = newFieldWithFieldLabelQueryBuilder("Container Memory Request", search.MemoryRequest, integerValueRegex, operatorsForbidden)
	DisallowedAnnotation   = newFieldWithFieldLabelQueryBuilder("Disallowed Annotation", "", keyValueValueRegex, negationForbidden)
	DisallowedImageLabel   = newFieldWithFieldLabelQueryBuilder("Disallowed Image Label", "", keyValueValueRegex, negationForbidden)
	DockerfileLine         = newFieldWithFieldLabelQueryBuilder("Dockerfile Line", "", dockerfileLineValueRegex, negationForbidden)
	DropCaps               = newFieldWithFieldLabelQueryBuilder("Drop Capabilities", search.DropCapabilities, capabilitiesValueRegex, negationForbidden)
	EnvironmentVariable    = newFieldWithFieldLabelQueryBuilder("Environment Variable", "", environmentVariableWithSourceRegex, negationForbidden)
	FixedBy                = newFieldWithFieldLabelQueryBuilder("Fixed By", search.FixedBy, stringValueRegex)
	ImageAge               = newFieldWithFieldLabelQueryBuilder("Image Age", "", integerValueRegex, negationForbidden, operatorsForbidden)
	ImageComponent         = newFieldWithFieldLabelQueryBuilder("Image Component", search.Component, keyValueValueRegex)
	ImageRegistry          = newFieldWithFieldLabelQueryBuilder("Image Registry", search.ImageRegistry, stringValueRegex)
	ImageRemote            = newFieldWithFieldLabelQueryBuilder("Image Remote", search.ImageRemote, stringValueRegex, negationForbidden)
	ImageScanAge           = newFieldWithFieldLabelQueryBuilder("Image Scan Age", "", integerValueRegex, negationForbidden, operatorsForbidden)
	ImageTag               = newFieldWithFieldLabelQueryBuilder("Image Tag", search.ImageTag, stringValueRegex)
	MinimumRBACPermissions = newFieldWithFieldLabelQueryBuilder("Minimum RBAC Permissions", "", rbacPermissionValueRegex, operatorsForbidden)
	Port                   = newFieldWithFieldLabelQueryBuilder("Port", search.Port, integerValueRegex)
	PortExposure           = newFieldWithFieldLabelQueryBuilder("Port Exposure Method", "", portExposureValueRegex)
	Privileged             = newFieldWithFieldLabelQueryBuilder("Privileged", search.Privileged, booleanValueRegex, negationForbidden, operatorsForbidden)
	ProcessAncestor        = newFieldWithFieldLabelQueryBuilder("Process Ancestor", search.ProcessAncestor, stringValueRegex)
	ProcessArguments       = newFieldWithFieldLabelQueryBuilder("Process Arguments", search.ProcessArguments, stringValueRegex)
	ProcessName            = newFieldWithFieldLabelQueryBuilder("Process Name", search.ProcessName, stringValueRegex)
	ProcessUID             = newFieldWithFieldLabelQueryBuilder("Process UID", search.ProcessUID, stringValueRegex)
	Protocol               = newFieldWithFieldLabelQueryBuilder("Protocol", search.PortProtocol, stringValueRegex)
	ReadOnlyRootFS         = newFieldWithFieldLabelQueryBuilder("Read-Only Root Filesystem", search.ReadOnlyRootFilesystem, booleanValueRegex, negationForbidden, operatorsForbidden)
	RequiredAnnotation     = newFieldWithFieldLabelQueryBuilder("Required Annotation", "", keyValueValueRegex, negationForbidden)
	RequiredImageLabel     = newFieldWithFieldLabelQueryBuilder("Required Image Label", "", keyValueValueRegex, negationForbidden)
	RequiredLabel          = newFieldWithFieldLabelQueryBuilder("Required Label", "", keyValueValueRegex, negationForbidden)
	UnscannedImage         = newFieldWithFieldLabelQueryBuilder("Unscanned Image", "", booleanValueRegex)
	VolumeDestination      = newFieldWithFieldLabelQueryBuilder("Volume Destination", search.VolumeDestination, stringValueRegex)
	VolumeName             = newFieldWithFieldLabelQueryBuilder("Volume Name", search.VolumeName, stringValueRegex)
	VolumeSource           = newFieldWithFieldLabelQueryBuilder("Volume Source", search.VolumeSource, stringValueRegex)
	VolumeType             = newFieldWithFieldLabelQueryBuilder("Volume Type", search.VolumeType, stringValueRegex)
	WhitelistsEnabled      = newFieldWithFieldLabelQueryBuilder("Unexpected Process Executed", "", booleanValueRegex, negationForbidden, operatorsForbidden)
	WritableHostMount      = newFieldWithFieldLabelQueryBuilder("Writable Host Mount", "", booleanValueRegex, negationForbidden, operatorsForbidden)
	WritableVolume         = newFieldWithFieldLabelQueryBuilder("Writable Volume", "", booleanValueRegex, negationForbidden, operatorsForbidden)
)

func newFieldWithFieldLabelQueryBuilder(fieldName string, fieldLabel search.FieldLabel, valueRegex *regexp.Regexp, options ...option) string {
	m := metadataAndQB{
		valueRegex: valueRegex,
	}
	// TEMPORARY: this indicates fields for which we haven't indicated the relevant query builder yet.
	if fieldLabel != "" {
		m.qb = &fieldLabelBasedQueryBuilder{fieldName: fieldName, fieldLabel: fieldLabel}
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
