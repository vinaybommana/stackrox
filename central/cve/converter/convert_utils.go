package converter

import (
	"strings"
	"time"

	"github.com/facebookincubator/nvdtools/cvefeed/nvd/schema"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/cvss/cvssv2"
	"github.com/stackrox/rox/pkg/cvss/cvssv3"
	"github.com/stackrox/rox/pkg/protoconv"
	"github.com/stackrox/rox/pkg/scans"
)

const (
	timeFormat = "2006-01-02T15:04Z"
)

// CVEType is the type of a CVE fetched by fetcher
type CVEType int32

// K8s is type for k8s CVEs, Istio is type for istio CVEs
const (
	K8s = iota
	Istio
)

// NvdCveToProtoCVE converts a nvd.CVEEntry object to an proto CVE
func NvdCveToProtoCVE(nvdCVE *schema.NVDCVEFeedJSON10DefCVEItem, ct CVEType) (*storage.CVE, error) {
	protoCVE := &storage.CVE{
		Id: nvdCVE.CVE.CVEDataMeta.ID,
	}

	if ct == K8s {
		protoCVE.Type = storage.CVE_K8S_CVE
	} else if ct == Istio {
		protoCVE.Type = storage.CVE_ISTIO_CVE
	} else {
		return nil, errors.Errorf("unknown CVE type: %d", ct)
	}

	cvssv2, err := nvdCvssv2ToProtoCvssv2(nvdCVE.Impact.BaseMetricV2)
	if err != nil {
		return nil, err
	}
	protoCVE.CvssV2 = cvssv2

	cvssv3, err := nvdCvssv3ToProtoCvssv3(nvdCVE.Impact.BaseMetricV3)
	if err != nil {
		return nil, err
	}
	protoCVE.CvssV3 = cvssv3

	if nvdCVE.PublishedDate != "" {
		if ts, err := time.Parse(timeFormat, nvdCVE.PublishedDate); err == nil {
			protoCVE.PublishedOn = protoconv.ConvertTimeToTimestamp(ts)
		}
	}

	if nvdCVE.LastModifiedDate != "" {
		if ts, err := time.Parse(timeFormat, nvdCVE.LastModifiedDate); err == nil {
			protoCVE.LastModified = protoconv.ConvertTimeToTimestamp(ts)
		}
	}

	if len(nvdCVE.CVE.Description.DescriptionData) > 0 {
		protoCVE.Summary = nvdCVE.CVE.Description.DescriptionData[0].Value
	}

	protoCVE.Link = scans.GetVulnLink(protoCVE.Id)

	if nvdCVE.Impact.BaseMetricV3.CVSSV3.BaseScore != 0.0 {
		protoCVE.Cvss = float32(nvdCVE.Impact.BaseMetricV3.CVSSV3.BaseScore)
		protoCVE.ScoreVersion = storage.CVE_V3
	} else if nvdCVE.Impact.BaseMetricV2.CVSSV2.BaseScore != 0.0 {
		protoCVE.Cvss = float32(nvdCVE.Impact.BaseMetricV2.CVSSV2.BaseScore)
		protoCVE.ScoreVersion = storage.CVE_V2
	} else {
		protoCVE.ScoreVersion = storage.CVE_UNKNOWN
	}

	fixVersions := getFixedVersions(nvdCVE.Configurations)
	if len(fixVersions) > 0 {
		protoCVE.SetFixedBy = &storage.CVE_FixedBy{
			FixedBy: strings.Join(fixVersions, ","),
		}
	}

	return protoCVE, nil
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

// NvdCVEsToProtoCVEs converts  NVD CVEs to Proto CVEs
func NvdCVEsToProtoCVEs(cves []*schema.NVDCVEFeedJSON10DefCVEItem, ct CVEType) ([]*storage.CVE, error) {
	protoCVEs := make([]*storage.CVE, 0, len(cves))
	for _, cve := range cves {
		ev, err := NvdCveToProtoCVE(cve, ct)
		if err != nil {
			return nil, err
		}
		protoCVEs = append(protoCVEs, ev)
	}
	return protoCVEs, nil
}

// ProtoCVEsToEmbeddedCVEs coverts Proto CVEs to Embedded Vulns
func ProtoCVEsToEmbeddedCVEs(protoCVEs []*storage.CVE) ([]*storage.EmbeddedVulnerability, error) {
	embeddedVulns := make([]*storage.EmbeddedVulnerability, 0, len(protoCVEs))
	for _, protoCVE := range protoCVEs {
		em := ProtoCVEToEmbeddedCVE(protoCVE)
		em.SetFixedBy = &storage.EmbeddedVulnerability_FixedBy{
			FixedBy: protoCVE.GetFixedBy(),
		}
		embeddedVulns = append(embeddedVulns, em)
	}
	return embeddedVulns, nil
}

// ProtoCVEToEmbeddedCVE coverts a Proto CVEs to Embedded Vuln
// It converts all the fields except except Fixed By which gets set depending on the CVE
func ProtoCVEToEmbeddedCVE(protoCVE *storage.CVE) *storage.EmbeddedVulnerability {
	embeddedCVE := &storage.EmbeddedVulnerability{
		Cve:          protoCVE.GetId(),
		Cvss:         protoCVE.GetCvss(),
		Summary:      protoCVE.GetSummary(),
		Link:         protoCVE.GetLink(),
		CvssV2:       protoCVE.GetCvssV2(),
		CvssV3:       protoCVE.GetCvssV3(),
		PublishedOn:  protoCVE.GetPublishedOn(),
		LastModified: protoCVE.GetLastModified(),
		Suppressed:   protoCVE.GetSuppressed(),
	}
	if protoCVE.CvssV3 != nil {
		embeddedCVE.ScoreVersion = storage.EmbeddedVulnerability_V3
	} else {
		embeddedCVE.ScoreVersion = storage.EmbeddedVulnerability_V2
	}
	embeddedCVE.VulnerabilityType = protoToEmbeddedVulnType(protoCVE.Type)
	return embeddedCVE
}

func protoToEmbeddedVulnType(protoCVEType storage.CVE_CVEType) storage.EmbeddedVulnerability_VulnerabilityType {
	switch protoCVEType {
	case storage.CVE_IMAGE_CVE:
		return storage.EmbeddedVulnerability_IMAGE_VULNERABILITY
	case storage.CVE_K8S_CVE:
		return storage.EmbeddedVulnerability_K8S_VULNERABILITY
	case storage.CVE_ISTIO_CVE:
		return storage.EmbeddedVulnerability_ISTIO_VULNERABILITY
	default:
		return storage.EmbeddedVulnerability_UNKNOWN_VULNERABILITY
	}
}
