package booleanpolicy

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/utils"
)

type individualFieldConverter func(fields *storage.PolicyFields) []*storage.PolicyGroup

var fieldsConverters = []individualFieldConverter{
	convertImageNamePolicy,
	convertImageAgeDays,
	convertDockerFileLineRule,
	convertCve,
	convertComponent,
	convertImageScanAge,
	convertNoScanExists,
	convertEnv,
	convertVolumePolicy,
	convertPortPolicy,
	convertRequiredLabel,
	convertRequiredAnnotation,
	convertDisallowedAnnotation,
	convertRequiredImageLabel,
	convertDisallowedImageLabel,
	convertPrivileged,
	convertProcessPolicy,
	convertHostMountPolicy,
	convertWhitelistEnabled,
	convertFixedBy,
	convertReadOnlyRootFs,
	convertCvss,
}

// CloneAndEnsureConverted returns a clone of the input that is upgraded if it is a legacy policy.
func CloneAndEnsureConverted(legacyPolicy *storage.Policy) *storage.Policy {
	p := legacyPolicy.Clone()
	if p.GetPolicyVersion() == legacyVersion {
		p.PolicyVersion = Version
		p.PolicySections = append(p.PolicySections, ConvertPolicyFieldsToSections(p.GetFields()))
		p.Fields.Reset()
	}
	return p
}

// ConvertPolicyFieldsToSections converts policy fields (version = "") to policy sections (version = "2.0").
func ConvertPolicyFieldsToSections(fields *storage.PolicyFields) *storage.PolicySection {
	var pgs []*storage.PolicyGroup
	for _, fieldConverter := range fieldsConverters {
		pgs = append(pgs, fieldConverter(fields)...)
	}

	if len(pgs) == 0 {
		return nil
	}

	return &storage.PolicySection{
		PolicyGroups: pgs,
	}
}

func convertImageScanAge(fields *storage.PolicyFields) []*storage.PolicyGroup {
	if fields.GetSetScanAgeDays() == nil {
		return nil
	}

	return []*storage.PolicyGroup{
		{
			FieldName: ImageScanAge,
			Values:    getPolicyValues(fields.GetScanAgeDays()),
		},
	}
}

func convertNoScanExists(fields *storage.PolicyFields) []*storage.PolicyGroup {
	if fields.GetSetNoScanExists() == nil {
		return nil
	}

	return []*storage.PolicyGroup{
		{
			FieldName: UnscannedImage,
			Values:    getPolicyValues(fields.GetNoScanExists()),
		},
	}
}

func convertEnv(fields *storage.PolicyFields) []*storage.PolicyGroup {
	p := fields.GetEnv()
	if p == nil {
		return nil
	}

	return []*storage.PolicyGroup{
		{
			FieldName: EnvironmentVariable,
			Values:    getPolicyValues(fmt.Sprintf("%s=%s=%s", p.GetEnvVarSource(), p.GetKey(), p.GetValue())),
		},
	}
}

func convertRequiredLabel(fields *storage.PolicyFields) []*storage.PolicyGroup {
	if p := convertKeyValuePolicy(fields.GetRequiredLabel(), RequiredLabel); p != nil {
		return []*storage.PolicyGroup{p}
	}

	return nil
}

func convertRequiredAnnotation(fields *storage.PolicyFields) []*storage.PolicyGroup {
	if p := convertKeyValuePolicy(fields.GetRequiredAnnotation(), RequiredAnnotation); p != nil {
		return []*storage.PolicyGroup{p}
	}

	return nil
}

func convertDisallowedAnnotation(fields *storage.PolicyFields) []*storage.PolicyGroup {
	if p := convertKeyValuePolicy(fields.GetDisallowedAnnotation(), DisallowedAnnotation); p != nil {
		return []*storage.PolicyGroup{p}
	}

	return nil
}

func convertRequiredImageLabel(fields *storage.PolicyFields) []*storage.PolicyGroup {
	if p := convertKeyValuePolicy(fields.GetRequiredImageLabel(), RequiredImageLabel); p != nil {
		return []*storage.PolicyGroup{p}
	}

	return nil
}

func convertDisallowedImageLabel(fields *storage.PolicyFields) []*storage.PolicyGroup {
	if p := convertKeyValuePolicy(fields.GetDisallowedImageLabel(), DisallowedImageLabel); p != nil {
		return []*storage.PolicyGroup{p}
	}

	return nil
}

func convertPrivileged(fields *storage.PolicyFields) []*storage.PolicyGroup {
	if fields.GetSetPrivileged() == nil {
		return nil
	}

	return []*storage.PolicyGroup{{
		FieldName: Privileged,
		Values:    getPolicyValues(fields.GetPrivileged()),
	},
	}
}

func convertWhitelistEnabled(fields *storage.PolicyFields) []*storage.PolicyGroup {
	if fields.GetSetWhitelist() == nil {
		return nil
	}

	return []*storage.PolicyGroup{{
		FieldName: WhitelistsEnabled,
		Values:    getPolicyValues(fields.GetWhitelistEnabled()),
	}}
}

func convertFixedBy(fields *storage.PolicyFields) []*storage.PolicyGroup {
	p := fields.GetFixedBy()
	if p == "" {
		return nil
	}

	return []*storage.PolicyGroup{{
		FieldName: FixedBy,
		Values:    getPolicyValues(p),
	}}
}

func convertReadOnlyRootFs(fields *storage.PolicyFields) []*storage.PolicyGroup {
	if fields.GetSetReadOnlyRootFs() == nil {
		return nil
	}

	return []*storage.PolicyGroup{{
		FieldName: ReadOnlyRootFS,
		Values:    getPolicyValues(fields.GetReadOnlyRootFs()),
	}}
}

func getPolicyValues(p ...interface{}) []*storage.PolicyValue {
	vs := make([]*storage.PolicyValue, 0, len(p))
	for _, v := range p {
		switch val := v.(type) {
		case string:
			vs = append(vs, &storage.PolicyValue{Value: val})
		case int64:
			vs = append(vs, &storage.PolicyValue{Value: strconv.FormatInt(val, 10)})
		case bool:
			vs = append(vs, &storage.PolicyValue{Value: strconv.FormatBool(val)})
		default:
			utils.Should(errors.Errorf("invalid policy type: %T", val))
		}
	}

	if len(vs) == 0 {
		return nil
	}

	return vs
}

func convertImageNamePolicy(fields *storage.PolicyFields) []*storage.PolicyGroup {
	p := fields.GetImageName()
	if p == nil {
		return nil
	}

	var res []*storage.PolicyGroup
	if p.GetRegistry() != "" {
		res = append(res, &storage.PolicyGroup{
			FieldName: ImageRegistry,
			Values:    getPolicyValues(p.GetRegistry()),
		})
	}

	if p.GetRemote() != "" {
		res = append(res, &storage.PolicyGroup{
			FieldName: ImageRemote,
			Values:    getPolicyValues(p.GetRemote()),
		})
	}

	if p.GetTag() != "" {
		res = append(res, &storage.PolicyGroup{
			FieldName: ImageTag,
			Values:    getPolicyValues(p.GetTag()),
		})
	}

	return res
}

func convertImageAgeDays(fields *storage.PolicyFields) []*storage.PolicyGroup {
	if fields.GetSetImageAgeDays() == nil {
		return nil
	}

	return []*storage.PolicyGroup{{
		FieldName: ImageAge,
		Values:    getPolicyValues(fields.GetImageAgeDays()),
	}}
}

func convertDockerFileLineRule(fields *storage.PolicyFields) []*storage.PolicyGroup {
	// TODO(ROX-4720): Update convertor for missing fields.
	return nil
}

func convertCvss(fields *storage.PolicyFields) []*storage.PolicyGroup {
	p := fields.GetCvss()
	if p == nil {
		return nil
	}

	return []*storage.PolicyGroup{convertNumericalPolicy(p, CVSS)}
}

func convertCve(fields *storage.PolicyFields) []*storage.PolicyGroup {
	p := fields.GetCve()
	if p == "" {
		return nil
	}

	return []*storage.PolicyGroup{{
		FieldName: CVE,
		Values:    getPolicyValues(p),
	}}
}

func convertNumericalPolicy(p *storage.NumericalPolicy, fieldName string) *storage.PolicyGroup {
	if p == nil {
		return nil
	}

	op := p.GetOp().String()
	switch p.GetOp() {
	case storage.Comparator_EQUALS:
		op = "="
	case storage.Comparator_GREATER_THAN:
		op = ">"
	case storage.Comparator_GREATER_THAN_OR_EQUALS:
		op = ">="
	case storage.Comparator_LESS_THAN:
		op = "<"
	case storage.Comparator_LESS_THAN_OR_EQUALS:
		op = "<="
	default:
		utils.Should(errors.Errorf("invalid op for numerical policy: %+v", p))
	}

	return &storage.PolicyGroup{
		FieldName: fieldName,
		Values: []*storage.PolicyValue{
			{
				Value: fmt.Sprintf("%s %f", op, p.GetValue()),
			},
		},
	}
}

func convertComponent(fields *storage.PolicyFields) []*storage.PolicyGroup {
	p := fields.GetComponent()
	if p == nil {
		return nil
	}

	return []*storage.PolicyGroup{{
		FieldName: "Image Component",
		Values: []*storage.PolicyValue{
			{
				Value: fmt.Sprintf("%s=%s", p.GetName(), p.GetVersion()),
			},
		},
	}}
}

func convertKeyValuePolicy(p *storage.KeyValuePolicy, fieldName string) *storage.PolicyGroup {
	if p == nil {
		return nil
	}

	return &storage.PolicyGroup{
		FieldName: fieldName,
		Values: []*storage.PolicyValue{
			{
				Value: fmt.Sprintf("%s=%s", p.GetKey(), p.GetValue()),
			},
		},
	}
}

func convertVolumePolicy(fields *storage.PolicyFields) []*storage.PolicyGroup {
	p := fields.GetVolumePolicy()
	if p == nil {
		return nil
	}

	var res []*storage.PolicyGroup
	if p.GetName() != "" {
		res = append(res, &storage.PolicyGroup{
			FieldName: "Volume Name",
			Values:    getPolicyValues(p.GetName()),
		})
	}

	if p.GetType() != "" {
		res = append(res, &storage.PolicyGroup{
			FieldName: "Volume Type",
			Values:    getPolicyValues(p.GetType()),
		})
	}

	if p.GetDestination() != "" {
		res = append(res, &storage.PolicyGroup{
			FieldName: "Volume Destination",
			Values:    getPolicyValues(p.GetDestination()),
		})
	}

	if p.GetSource() != "" {
		res = append(res, &storage.PolicyGroup{
			FieldName: "Volume Source",
			Values:    getPolicyValues(p.GetSource()),
		})
	}

	ro := p.GetSetReadOnly()
	if ro != nil {
		res = append(res, &storage.PolicyGroup{
			FieldName: "Writable Volume",
			Values:    getPolicyValues(!p.GetReadOnly()),
		})
	}

	return res
}

func convertPortPolicy(fields *storage.PolicyFields) []*storage.PolicyGroup {
	p := fields.GetPortPolicy()
	if p == nil {
		return nil
	}

	var res []*storage.PolicyGroup
	if p.GetPort() != 0 {
		res = append(res, &storage.PolicyGroup{
			FieldName: "Port",
			Values:    getPolicyValues(int64(p.GetPort())),
		})
	}

	if p.GetProtocol() != "" {
		res = append(res, &storage.PolicyGroup{
			FieldName: "Protocol",
			Values:    getPolicyValues(p.GetProtocol()),
		})
	}

	return res
}

func convertProcessPolicy(fields *storage.PolicyFields) []*storage.PolicyGroup {
	p := fields.GetProcessPolicy()
	if p == nil {
		return nil
	}

	var res []*storage.PolicyGroup
	if p.GetName() != "" {
		res = append(res, &storage.PolicyGroup{
			FieldName: "Process Name",
			Values:    getPolicyValues(p.GetName()),
		})
	}

	if p.GetAncestor() != "" {
		res = append(res, &storage.PolicyGroup{
			FieldName: "Process Ancestor",
			Values:    getPolicyValues(p.GetAncestor()),
		})
	}

	if p.GetArgs() != "" {
		res = append(res, &storage.PolicyGroup{
			FieldName: "Process Arguments",
			Values:    getPolicyValues(p.GetArgs()),
		})
	}

	if p.GetUid() != "" {
		res = append(res, &storage.PolicyGroup{
			FieldName: "Process UID",
			Values:    getPolicyValues(p.GetUid()),
		})
	}

	return res
}

func convertHostMountPolicy(fields *storage.PolicyFields) []*storage.PolicyGroup {
	p := fields.GetHostMountPolicy()
	if p.GetSetReadOnly() == nil {
		return nil
	}

	return []*storage.PolicyGroup{{
		FieldName: "Writable Host Mount",
		Values:    getPolicyValues(!p.GetReadOnly()),
	},
	}
}
