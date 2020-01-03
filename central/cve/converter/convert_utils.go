package converter

import (
	"fmt"
	"strings"
	"time"

	"github.com/facebookincubator/nvdtools/cvefeed/nvd/schema"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/cvss/cvssv2"
	"github.com/stackrox/rox/pkg/cvss/cvssv3"
	"github.com/stackrox/rox/pkg/protoconv"
	"github.com/stackrox/rox/pkg/scans"
)

const (
	timeFormat = "2006-01-02T15:04Z"
)

// CveType is the type of a CVE fetched by fetcher
type CveType int32

// K8s is type for k8s CVEs, Istio is type for istio CVEs
const (
	K8s = iota
	Istio
)

// NvdCveToEmbeddedVulnerability converts a nvd.CVEEntry object to an EmbeddedVulnerability which is used elsewhere.
func NvdCveToEmbeddedVulnerability(cve *schema.NVDCVEFeedJSON10DefCVEItem, ct CveType) (*storage.EmbeddedVulnerability, error) {
	ev := &storage.EmbeddedVulnerability{
		Cve: cve.CVE.CVEDataMeta.ID,
	}

	if ct == K8s {
		ev.VulnerabilityType = storage.EmbeddedVulnerability_K8S_VULNERABILITY
	} else if ct == Istio {
		ev.VulnerabilityType = storage.EmbeddedVulnerability_ISTIO_VULNERABILITY
	} else {
		return nil, fmt.Errorf("unknown CVE type: %d", ct)
	}

	cvssv2, err := nvdCvssv2ToProtoCvssv2(cve.Impact.BaseMetricV2)
	if err != nil {
		return nil, err
	}
	ev.CvssV2 = cvssv2

	cvssv3, err := nvdCvssv3ToProtoCvssv3(cve.Impact.BaseMetricV3)
	if err != nil {
		return nil, err
	}
	ev.CvssV3 = cvssv3

	if cve.PublishedDate != "" {
		if ts, err := time.Parse(timeFormat, cve.PublishedDate); err == nil {
			ev.PublishedOn = protoconv.ConvertTimeToTimestamp(ts)
		}
	}

	if cve.LastModifiedDate != "" {
		if ts, err := time.Parse(timeFormat, cve.LastModifiedDate); err == nil {
			ev.LastModified = protoconv.ConvertTimeToTimestamp(ts)
		}
	}

	if len(cve.CVE.Description.DescriptionData) > 0 {
		ev.Summary = cve.CVE.Description.DescriptionData[0].Value
	}

	ev.Link = scans.GetVulnLink(ev.Cve)

	if cve.Impact.BaseMetricV3.CVSSV3.BaseScore != 0.0 {
		ev.Cvss = float32(cve.Impact.BaseMetricV3.CVSSV3.BaseScore)
		ev.ScoreVersion = storage.EmbeddedVulnerability_V3
	} else {
		ev.Cvss = float32(cve.Impact.BaseMetricV2.CVSSV2.BaseScore)
		ev.ScoreVersion = storage.EmbeddedVulnerability_V2
	}

	fixVersions := getFixedVersions(cve.Configurations)
	if len(fixVersions) > 0 {
		ev.SetFixedBy = &storage.EmbeddedVulnerability_FixedBy{
			FixedBy: strings.Join(fixVersions, ","),
		}
	}

	return ev, nil
}

func getFixedVersions(configurations *schema.NVDCVEFeedJSON10DefConfigurations) []string {
	var versions []string
	if configurations == nil {
		return versions
	}
	for _, node := range configurations.Nodes {
		for _, cpeMatch := range node.CPEMatch {
			if cpeMatch.VersionEndExcluding != "" {
				versions = append(versions, cpeMatch.VersionEndExcluding)
			}
		}
	}
	return versions
}

func nvdCvssv2ToProtoCvssv2(baseMetricV2 *schema.NVDCVEFeedJSON10DefImpactBaseMetricV2) (*storage.CVSSV2, error) {
	cvssV2, err := cvssv2.ParseCVSSV2(baseMetricV2.CVSSV2.VectorString)
	if err != nil {
		return nil, err
	}

	if baseMetricV2.Severity != "" {
		k := strings.ToUpper(baseMetricV2.Severity[:1])
		sv, err := cvssv2.GetSeverityMapProtoVal(k)
		if err != nil {
			return nil, err
		}
		cvssV2.Severity = sv
	}

	cvssV2.Score = float32(baseMetricV2.CVSSV2.BaseScore)
	cvssV2.ExploitabilityScore = float32(baseMetricV2.ExploitabilityScore)
	cvssV2.ImpactScore = float32(baseMetricV2.ImpactScore)

	return cvssV2, nil
}

func nvdCvssv3ToProtoCvssv3(baseMetricV3 *schema.NVDCVEFeedJSON10DefImpactBaseMetricV3) (*storage.CVSSV3, error) {
	cvssV3, err := cvssv3.ParseCVSSV3(baseMetricV3.CVSSV3.VectorString)
	if err != nil {
		return nil, err
	}
	if baseMetricV3.CVSSV3.BaseSeverity != "" {
		k := strings.ToUpper(baseMetricV3.CVSSV3.BaseSeverity[:1])
		sv, err := cvssv3.GetSeverityMapProtoVal(k)
		if err != nil {
			return nil, err
		}
		cvssV3.Severity = sv
	}

	cvssV3.Score = float32(baseMetricV3.CVSSV3.BaseScore)
	cvssV3.ExploitabilityScore = float32(baseMetricV3.ExploitabilityScore)
	cvssV3.ImpactScore = float32(baseMetricV3.ImpactScore)

	return cvssV3, nil
}

// NvdCVEsToEmbeddedVulnerabilities converts  NVD cves to EmbeddedVulnerabilities
func NvdCVEsToEmbeddedVulnerabilities(cves []*schema.NVDCVEFeedJSON10DefCVEItem, ct CveType) ([]*storage.EmbeddedVulnerability, error) {
	evs := make([]*storage.EmbeddedVulnerability, 0, len(cves))
	for _, cve := range cves {
		ev, err := NvdCveToEmbeddedVulnerability(cve, ct)
		if err != nil {
			return nil, err
		}
		evs = append(evs, ev)
	}
	return evs, nil
}
