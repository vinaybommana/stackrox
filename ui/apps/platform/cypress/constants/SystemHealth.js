const baseURL = '/main/system-health';

export const url = {
    dashboard: baseURL,
};

export const selectors = {
    bundle: {
        generateDiagnosticBundleButton: 'button:contains("Generate Diagnostic Bundle")',
        filterByClusters: '[data-testid="filter-by-clusters"]',
        filterByStartingTime: '#filter-by-starting-time',
        startingTimeMessage: '[data-testid="starting-time-message"]',
        downloadDiagnosticBundleButton: 'button:contains("Download Diagnostic Bundle")',
    },
    clusters: {
        categoryCount: '[data-testid="count"]',
        categoryLabel: '[data-testid="label"]',
        healthyText: '[data-testid="healthy-text"]',
        healthySubtext: '[data-testid="healthy-subtext"]',
        problemCount: '[data-testid="problem-count"]',
        viewAllButton: '[data-testid="cluster-health"] a:contains("View All")',
        widgets: {
            clusterOverview: '[data-testid="cluster-overview"]',
            collectorStatus: '[data-testid="collector-status"]',
            sensorStatus: '[data-testid="sensor-status"]',
            sensorUpgrade: '[data-testid="sensor-upgrade"]',
            credentialExpiration: '[data-testid="credential-expiration"]',
        },
    },
    integrations: {
        errorMessage: '[data-testid="error-message"]',
        healthyText: '[data-testid="healthy-text"]',
        integrationName: '[data-testid="integration-name"]',
        integrationLabel: '[data-testid="integration-label"]',
        lastContact: '[data-testid="last-contact"]',
        viewAllButton: 'a:contains("View All")',
        widgets: {
            imageIntegrations: '[data-testid="image-integrations"]',
            pluginIntegrations: '[data-testid="plugin-integrations"]',
            backupIntegrations: '[data-testid="backup-integrations"]',
        },
    },
};
