package mappings

import (
	"github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/pkg/search"
)

// OptionsMap is exposed for e2e test
var OptionsMap = map[search.FieldLabel]*v1.SearchField{
	search.ProcessID:    search.NewField(v1.SearchCategory_PROCESS_INDICATORS, "process_indicator.id", v1.SearchDataType_SEARCH_STRING, search.OptionHidden|search.OptionStore),
	search.DeploymentID: search.NewField(v1.SearchCategory_PROCESS_INDICATORS, "process_indicator.deployment_id", v1.SearchDataType_SEARCH_STRING, search.OptionHidden|search.OptionStore),
	search.ContainerID:  search.NewField(v1.SearchCategory_PROCESS_INDICATORS, "process_indicator.signal.container_id", v1.SearchDataType_SEARCH_STRING, search.OptionHidden),

	search.ProcessArguments: search.NewStringField(v1.SearchCategory_PROCESS_INDICATORS, "process_indicator.signal.args"),
	search.ProcessExecPath:  search.NewStringField(v1.SearchCategory_PROCESS_INDICATORS, "process_indicator.signal.exec_file_path"),
	search.ProcessName:      search.NewStringField(v1.SearchCategory_PROCESS_INDICATORS, "process_indicator.signal.name"),
	search.ProcessAncestor:  search.NewStringField(v1.SearchCategory_PROCESS_INDICATORS, "process_indicator.signal.lineage"),
}
