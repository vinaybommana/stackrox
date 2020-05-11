package booleanpolicy

import (
	"regexp"

	"github.com/stackrox/rox/pkg/booleanpolicy/augmentedobjs"
	"github.com/stackrox/rox/pkg/booleanpolicy/querybuilders"
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
	qb                 querybuilders.QueryBuilder
	valueRegex         *regexp.Regexp
}

// This block enumerates field short names.
var (
	AddCaps                = newField("Add Capabilities", querybuilders.ForFieldLabelExact(search.AddCapabilities), capabilitiesValueRegex, negationForbidden)
	CVE                    = newField("CVE", querybuilders.ForFieldLabelRegex(search.CVE), stringValueRegex)
	CVSS                   = newField("CVSS", querybuilders.ForFieldLabel(search.CVSS), comparatorDecimalValueRegex, operatorsForbidden)
	ContainerCPULimit      = newField("Container CPU Limit", querybuilders.ForFieldLabel(search.CPUCoresLimit), comparatorDecimalValueRegex, operatorsForbidden)
	ContainerCPURequest    = newField("Container CPU Request", querybuilders.ForFieldLabel(search.CPUCoresRequest), comparatorDecimalValueRegex, operatorsForbidden)
	ContainerMemLimit      = newField("Container Memory Limit", querybuilders.ForFieldLabel(search.MemoryLimit), comparatorDecimalValueRegex, operatorsForbidden)
	ContainerMemRequest    = newField("Container Memory Request", querybuilders.ForFieldLabel(search.MemoryRequest), comparatorDecimalValueRegex, operatorsForbidden)
	DisallowedAnnotation   = newField("Disallowed Annotation", nil, keyValueValueRegex, negationForbidden)
	DisallowedImageLabel   = newField("Disallowed Image Label", nil, keyValueValueRegex, negationForbidden)
	DockerfileLine         = newField("Dockerfile Line", querybuilders.ForCompound(augmentedobjs.DockerfileLineCustomTag), dockerfileLineValueRegex, negationForbidden)
	DropCaps               = newField("Drop Capabilities", nil, capabilitiesValueRegex, negationForbidden)
	EnvironmentVariable    = newField("Environment Variable", nil, environmentVariableWithSourceRegex, negationForbidden)
	FixedBy                = newField("Fixed By", querybuilders.ForFieldLabelRegex(search.FixedBy), stringValueRegex)
	ImageAge               = newField("Image Age", nil, integerValueRegex, negationForbidden, operatorsForbidden)
	ImageComponent         = newField("Image Component", querybuilders.ForCompound(augmentedobjs.ComponentAndVersionCustomTag), keyValueValueRegex, negationForbidden)
	ImageRegistry          = newField("Image Registry", querybuilders.ForFieldLabelRegex(search.ImageRegistry), stringValueRegex)
	ImageRemote            = newField("Image Remote", querybuilders.ForFieldLabelRegex(search.ImageRemote), stringValueRegex, negationForbidden)
	ImageScanAge           = newField("Image Scan Age", nil, integerValueRegex, negationForbidden, operatorsForbidden)
	ImageTag               = newField("Image Tag", querybuilders.ForFieldLabelRegex(search.ImageTag), stringValueRegex)
	MinimumRBACPermissions = newField("Minimum RBAC Permissions", nil, rbacPermissionValueRegex, operatorsForbidden)
	Port                   = newField("Port", querybuilders.ForFieldLabel(search.Port), integerValueRegex)
	PortExposure           = newField("Port Exposure Method", nil, portExposureValueRegex)
	Privileged             = newField("Privileged", querybuilders.ForFieldLabel(search.Privileged), booleanValueRegex, negationForbidden, operatorsForbidden)
	ProcessAncestor        = newField("Process Ancestor", querybuilders.ForFieldLabelRegex(search.ProcessAncestor), stringValueRegex)
	ProcessArguments       = newField("Process Arguments", querybuilders.ForFieldLabelRegex(search.ProcessArguments), stringValueRegex)
	ProcessName            = newField("Process Name", querybuilders.ForFieldLabelRegex(search.ProcessName), stringValueRegex)
	ProcessUID             = newField("Process UID", querybuilders.ForFieldLabel(search.ProcessUID), stringValueRegex)
	Protocol               = newField("Protocol", querybuilders.ForFieldLabelUpper(search.PortProtocol), stringValueRegex)
	ReadOnlyRootFS         = newField("Read-Only Root Filesystem", querybuilders.ForFieldLabel(search.ReadOnlyRootFilesystem), booleanValueRegex, negationForbidden, operatorsForbidden)
	RequiredAnnotation     = newField("Required Annotation", nil, keyValueValueRegex, negationForbidden)
	RequiredImageLabel     = newField("Required Image Label", nil, keyValueValueRegex, negationForbidden)
	RequiredLabel          = newField("Required Label", nil, keyValueValueRegex, negationForbidden)
	UnscannedImage         = newField("Unscanned Image", nil, booleanValueRegex)
	VolumeDestination      = newField("Volume Destination", querybuilders.ForFieldLabelRegex(search.VolumeDestination), stringValueRegex)
	VolumeName             = newField("Volume Name", querybuilders.ForFieldLabelRegex(search.VolumeName), stringValueRegex)
	VolumeSource           = newField("Volume Source", querybuilders.ForFieldLabelRegex(search.VolumeSource), stringValueRegex)
	VolumeType             = newField("Volume Type", querybuilders.ForFieldLabelRegex(search.VolumeType), stringValueRegex)
	WhitelistsEnabled      = newField("Unexpected Process Executed", nil, booleanValueRegex, negationForbidden, operatorsForbidden)
	WritableHostMount      = newField("Writable Host Mount", nil, booleanValueRegex, negationForbidden, operatorsForbidden)
	WritableVolume         = newField("Writable Volume", nil, booleanValueRegex, negationForbidden, operatorsForbidden)
)

func newField(fieldName string, qb querybuilders.QueryBuilder, valueRegex *regexp.Regexp, options ...option) string {
	m := metadataAndQB{
		valueRegex: valueRegex,
		qb:         qb,
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
