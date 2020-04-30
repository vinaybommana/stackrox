package booleanpolicy

import (
	"testing"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stretchr/testify/suite"
)

func TestBPLUtil(t *testing.T) {
	suite.Run(t, new(BPLUtilTestSuite))
}

type BPLUtilTestSuite struct {
	suite.Suite
}

func (s *BPLUtilTestSuite) SetupTest() {
}

func (s *BPLUtilTestSuite) TearDownTest() {
}

func (s *BPLUtilTestSuite) TestIsWhitelistEnabled() {
	whitelistEnabled := &storage.Policy{
		PolicySections: []*storage.PolicySection{
			{
				PolicyGroups: []*storage.PolicyGroup{
					{
						FieldName: WhitelistsEnabled,
					},
					{
						FieldName: CVSS,
					},
				},
			},
		},
	}
	noWhitelistEnabled := &storage.Policy{
		PolicySections: []*storage.PolicySection{
			{
				PolicyGroups: []*storage.PolicyGroup{
					{
						FieldName: CVSS,
					},
				},
			},
		},
	}

	isWhitelistEnabled := IsWhitelistEnabled(whitelistEnabled)
	s.True(isWhitelistEnabled)

	isWhitelistEnabled = IsWhitelistEnabled(noWhitelistEnabled)
	s.False(isWhitelistEnabled)
}
