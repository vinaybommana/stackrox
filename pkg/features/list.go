package features

var (
	// ConfigMgmtUI enables the config management UI.
	// NB: When removing this feature flag, remove references in ui/src/utils/featureFlags.js
	ConfigMgmtUI = registerFeature("Enable Config Mgmt UI", "ROX_CONFIG_MGMT_UI", true)

	// VulnMgmtUI enables the vulnerability management UI.
	// NB: When removing this feature flag, remove references in ui/src/utils/featureFlags.js
	VulnMgmtUI = registerFeature("Enable Vulnerability Management UI", "ROX_VULN_MGMT_UI", false)

	// Dackbox enables the id graph layer on top of badger.
	Dackbox = registerFeature("Use DackBox layer for the embedded Badger DB", "ROX_DACKBOX", false)

	// Telemetry enables the telemetry features
	Telemetry = registerFeature("Enable support for telemetry", "ROX_TELEMETRY", false)

	// IQTAnalystNotesUI enables the IQT Analyst Notes UI.
	// NB: When removing this feature flag, remove references in ui/src/utils/featureFlags.js
	IQTAnalystNotesUI = registerFeature("Enable IQT Analyst Notes UI", "ROX_IQT_ANALYST_NOTES_UI", true)
)
