package main

import (
	"os"

	"github.com/gogo/protobuf/types"
	"github.com/stackrox/rox/generated/internalapi/compliance"
	"github.com/stackrox/rox/pkg/jsonutil"
	"github.com/stackrox/rox/pkg/utils"
)

const (
	scrapeFixturePath  = "/files/scrape.json"
	resultsFixturePath = "/files/results.json"
)

var (
	defaultCheckResults *compliance.ComplianceReturn
)

func init() {
	defaultCheckResults = loadComplianceReturn(resultsFixturePath)
}

func loadComplianceReturn(path string) *compliance.ComplianceReturn {
	complianceBytes, err := os.ReadFile(path)
	utils.CrashOnError(err)

	var complianceReturn compliance.ComplianceReturn
	utils.Must(jsonutil.JSONBytesToProto(complianceBytes, &complianceReturn))
	return &complianceReturn
}

func getCheckResults(scrapeID, nodeName string) *compliance.ComplianceReturn {
	cr := defaultCheckResults.Clone()
	cr.ScrapeId = scrapeID
	cr.NodeName = nodeName
	cr.Time = types.TimestampNow()
	return cr
}
