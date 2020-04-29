package booleanpolicy

import (
	"testing"

	"github.com/stackrox/rox/generated/storage"
	"gotest.tools/assert"
)

type testcase struct {
	desc                  string
	policyFields          *storage.PolicyFields
	expectedPolicySection *storage.PolicySection
}

func TestConvertPolicyFieldsToSections(t *testing.T) {
	tcs := []*testcase{
		{
			desc: "cvss",
			policyFields: &storage.PolicyFields{
				Cvss: &storage.NumericalPolicy{
					Op:    storage.Comparator_GREATER_THAN_OR_EQUALS,
					Value: 7.0,
				},
			},
			expectedPolicySection: &storage.PolicySection{
				PolicyGroups: []*storage.PolicyGroup{
					{
						FieldName: "cvss",
						Values: []*storage.PolicyValue{
							{
								Value: ">= 7.000000",
							},
						},
					},
				},
			},
		},

		{
			desc: "fixed by",
			policyFields: &storage.PolicyFields{
				FixedBy: "pkg=4",
			},
			expectedPolicySection: &storage.PolicySection{
				PolicyGroups: []*storage.PolicyGroup{
					{
						FieldName: "Fixed by",
						Values: []*storage.PolicyValue{
							{
								Value: "pkg=4",
							},
						},
					},
				},
			},
		},

		{
			desc: "process policy",
			policyFields: &storage.PolicyFields{
				ProcessPolicy: &storage.ProcessPolicy{
					Name:     "process",
					Args:     "--arg 1",
					Ancestor: "parent",
					Uid:      "123",
				},
			},
			expectedPolicySection: &storage.PolicySection{
				PolicyGroups: []*storage.PolicyGroup{
					{
						FieldName: "Process Name",
						Values:    []*storage.PolicyValue{{Value: "process"}},
					},

					{
						FieldName: "Process Ancestor",
						Values:    []*storage.PolicyValue{{Value: "parent"}},
					},

					{
						FieldName: "Process Args",
						Values:    []*storage.PolicyValue{{Value: "--arg 1"}},
					},

					{
						FieldName: "Process Uid",
						Values:    []*storage.PolicyValue{{Value: "123"}},
					},
				},
			},
		},

		{
			desc: "disallowed image label",
			policyFields: &storage.PolicyFields{
				DisallowedImageLabel: &storage.KeyValuePolicy{
					Key:   "k",
					Value: "v",
				},
			},
			expectedPolicySection: &storage.PolicySection{
				PolicyGroups: []*storage.PolicyGroup{
					{
						FieldName: "Disallowed Image Label",
						Values: []*storage.PolicyValue{
							{
								Value: "k=v",
							},
						},
					},
				},
			},
		},

		{
			desc: "required image label",
			policyFields: &storage.PolicyFields{
				RequiredImageLabel: &storage.KeyValuePolicy{
					Key:   "k",
					Value: "v",
				},
			},
			expectedPolicySection: &storage.PolicySection{
				PolicyGroups: []*storage.PolicyGroup{
					{
						FieldName: "Required Image Label",
						Values: []*storage.PolicyValue{
							{
								Value: "k=v",
							},
						},
					},
				},
			},
		},

		{
			desc: "disallowed annotation",
			policyFields: &storage.PolicyFields{
				DisallowedAnnotation: &storage.KeyValuePolicy{
					Key:   "k",
					Value: "v",
				},
			},
			expectedPolicySection: &storage.PolicySection{
				PolicyGroups: []*storage.PolicyGroup{
					{
						FieldName: "Disallowed Annotation",
						Values: []*storage.PolicyValue{
							{
								Value: "k=v",
							},
						},
					},
				},
			},
		},

		{
			desc: "required annotation",
			policyFields: &storage.PolicyFields{
				RequiredAnnotation: &storage.KeyValuePolicy{
					Key:   "k",
					Value: "v",
				},
			},
			expectedPolicySection: &storage.PolicySection{
				PolicyGroups: []*storage.PolicyGroup{
					{
						FieldName: "Required Annotation",
						Values: []*storage.PolicyValue{
							{
								Value: "k=v",
							},
						},
					},
				},
			},
		},

		{
			desc: "required label",
			policyFields: &storage.PolicyFields{
				RequiredLabel: &storage.KeyValuePolicy{
					Key:   "k",
					Value: "v",
				},
			},
			expectedPolicySection: &storage.PolicySection{
				PolicyGroups: []*storage.PolicyGroup{
					{
						FieldName: "Required Label",
						Values: []*storage.PolicyValue{
							{
								Value: "k=v",
							},
						},
					},
				},
			},
		},

		{
			desc: "env",
			policyFields: &storage.PolicyFields{
				Env: &storage.KeyValuePolicy{
					Key:          "k",
					Value:        "v",
					EnvVarSource: storage.ContainerConfig_EnvironmentConfig_RESOURCE_FIELD,
				},
			},
			expectedPolicySection: &storage.PolicySection{
				PolicyGroups: []*storage.PolicyGroup{
					{
						FieldName: "Environment Variable",
						Values: []*storage.PolicyValue{
							{
								Value: "RESOURCE_FIELD=k=v",
							},
						},
					},
				},
			},
		},

		{
			desc: "port policy",
			policyFields: &storage.PolicyFields{
				PortPolicy: &storage.PortPolicy{
					Port:     1234,
					Protocol: "protocol",
				},
			},
			expectedPolicySection: &storage.PolicySection{
				PolicyGroups: []*storage.PolicyGroup{
					{
						FieldName: "Port",
						Values: []*storage.PolicyValue{
							{
								Value: "1234",
							},
						},
					},

					{
						FieldName: "Protocol",
						Values: []*storage.PolicyValue{
							{
								Value: "protocol",
							},
						},
					},
				},
			},
		},

		{
			desc: "volume policy",
			policyFields: &storage.PolicyFields{
				VolumePolicy: &storage.VolumePolicy{
					Name:        "v",
					Source:      "s",
					Destination: "d",
					Type:        "fs",
				},
			},
			expectedPolicySection: &storage.PolicySection{
				PolicyGroups: []*storage.PolicyGroup{
					{
						FieldName: "Volume Name",
						Values: []*storage.PolicyValue{
							{
								Value: "v",
							},
						},
					},

					{
						FieldName: "Volume Type",
						Values: []*storage.PolicyValue{
							{
								Value: "fs",
							},
						},
					},

					{
						FieldName: "Volume Destination",
						Values: []*storage.PolicyValue{
							{
								Value: "d",
							},
						},
					},

					{
						FieldName: "Volume Source",
						Values: []*storage.PolicyValue{
							{
								Value: "s",
							},
						},
					},
				},
			},
		},

		{
			desc: "image name policy",
			policyFields: &storage.PolicyFields{
				ImageName: &storage.ImageNamePolicy{
					Registry: "r",
					Remote:   "r",
					Tag:      "t",
				},
			},
			expectedPolicySection: &storage.PolicySection{
				PolicyGroups: []*storage.PolicyGroup{
					{
						FieldName: "Image Registry",
						Values: []*storage.PolicyValue{
							{
								Value: "r",
							},
						},
					},

					{
						FieldName: "Image Remote",
						Values: []*storage.PolicyValue{
							{
								Value: "r",
							},
						},
					},

					{
						FieldName: "Image Tag",
						Values: []*storage.PolicyValue{
							{
								Value: "t",
							},
						},
					},
				},
			},
		},

		{
			desc: "cve",
			policyFields: &storage.PolicyFields{
				Cve: "cve",
			},
			expectedPolicySection: &storage.PolicySection{
				PolicyGroups: []*storage.PolicyGroup{
					{
						FieldName: "cve",
						Values: []*storage.PolicyValue{
							{
								Value: "cve",
							},
						},
					},
				},
			},
		},

		{
			desc: "component",
			policyFields: &storage.PolicyFields{
				Component: &storage.Component{
					Name:    "n",
					Version: "v",
				},
			},
			expectedPolicySection: &storage.PolicySection{
				PolicyGroups: []*storage.PolicyGroup{
					{
						FieldName: "Image Component",
						Values: []*storage.PolicyValue{
							{
								Value: "n=v",
							},
						},
					},
				},
			},
		},

		{
			desc: "image age days",
			policyFields: &storage.PolicyFields{
				SetImageAgeDays: &storage.PolicyFields_ImageAgeDays{ImageAgeDays: 30},
			},
			expectedPolicySection: &storage.PolicySection{
				PolicyGroups: []*storage.PolicyGroup{
					{
						FieldName: "Image Age",
						Values: []*storage.PolicyValue{
							{
								Value: "30",
							},
						},
					},
				},
			},
		},

		{
			desc: "scan age days",
			policyFields: &storage.PolicyFields{
				SetScanAgeDays: &storage.PolicyFields_ScanAgeDays{ScanAgeDays: 30},
			},
			expectedPolicySection: &storage.PolicySection{
				PolicyGroups: []*storage.PolicyGroup{
					{
						FieldName: "Image Scan Age",
						Values: []*storage.PolicyValue{
							{
								Value: "30",
							},
						},
					},
				},
			},
		},

		{
			desc: "unscanned image",
			policyFields: &storage.PolicyFields{
				SetNoScanExists: &storage.PolicyFields_NoScanExists{NoScanExists: true},
			},
			expectedPolicySection: &storage.PolicySection{
				PolicyGroups: []*storage.PolicyGroup{
					{
						FieldName: "Unscanned Image",
						Values: []*storage.PolicyValue{
							{
								Value: "true",
							},
						},
					},
				},
			},
		},

		{
			desc: "unscanned image",
			policyFields: &storage.PolicyFields{
				SetNoScanExists: &storage.PolicyFields_NoScanExists{NoScanExists: true},
			},
			expectedPolicySection: &storage.PolicySection{
				PolicyGroups: []*storage.PolicyGroup{
					{
						FieldName: "Unscanned Image",
						Values: []*storage.PolicyValue{
							{
								Value: "true",
							},
						},
					},
				},
			},
		},

		{
			desc: "privileged",
			policyFields: &storage.PolicyFields{
				SetPrivileged: &storage.PolicyFields_Privileged{Privileged: true},
			},
			expectedPolicySection: &storage.PolicySection{
				PolicyGroups: []*storage.PolicyGroup{
					{
						FieldName: "Privileged",
						Values: []*storage.PolicyValue{
							{
								Value: "true",
							},
						},
					},
				},
			},
		},

		{
			desc: "read only root fs",
			policyFields: &storage.PolicyFields{
				SetReadOnlyRootFs: &storage.PolicyFields_ReadOnlyRootFs{ReadOnlyRootFs: true},
			},
			expectedPolicySection: &storage.PolicySection{
				PolicyGroups: []*storage.PolicyGroup{
					{
						FieldName: "Read-Only Root Filesystem",
						Values: []*storage.PolicyValue{
							{
								Value: "true",
							},
						},
					},
				},
			},
		},

		{
			desc: "whitelist enabled",
			policyFields: &storage.PolicyFields{
				SetWhitelist: &storage.PolicyFields_WhitelistEnabled{WhitelistEnabled: true},
			},
			expectedPolicySection: &storage.PolicySection{
				PolicyGroups: []*storage.PolicyGroup{
					{
						FieldName: "Whitelist enabled",
						Values: []*storage.PolicyValue{
							{
								Value: "true",
							},
						},
					},
				},
			},
		},

		{
			desc: "writable host mount",
			policyFields: &storage.PolicyFields{
				HostMountPolicy: &storage.HostMountPolicy{SetReadOnly: &storage.HostMountPolicy_ReadOnly{ReadOnly: true}},
			},
			expectedPolicySection: &storage.PolicySection{
				PolicyGroups: []*storage.PolicyGroup{
					{
						FieldName: "Writable Host Mount",
						Values: []*storage.PolicyValue{
							{
								Value: "false",
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			got := ConvertPolicyFieldsToSections(tc.policyFields)
			assert.DeepEqual(t, tc.expectedPolicySection, got)
		})
	}
}
