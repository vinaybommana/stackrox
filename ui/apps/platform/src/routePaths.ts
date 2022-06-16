/**
 * Application route paths constants.
 */

import { resourceTypes, standardEntityTypes, rbacConfigTypes } from 'constants/entityTypes';
import useFeatureFlags, { IsFeatureFlagEnabled } from 'hooks/useFeatureFlags';
import usePermissions, { HasReadAccess } from 'hooks/usePermissions';
import { FeatureFlagEnvVar } from 'types/featureFlag';
import { ResourceName } from 'types/roleResources';

export const mainPath = '/main';
export const loginPath = '/login';
export const testLoginResultsPath = '/test-login-results';
export const authResponsePrefix = '/auth/response/';

export const dashboardPath = `${mainPath}/dashboard`;
export const dashboardPathPF = `${mainPath}/dashboard-pf`;
export const networkBasePath = `${mainPath}/network`;
export const networkPath = `${networkBasePath}/:deploymentId?/:externalType?`;
export const violationsBasePath = `${mainPath}/violations`;
export const violationsPath = `${violationsBasePath}/:alertId?`;
export const clustersBasePath = `${mainPath}/clusters`;
export const clustersPathWithParam = `${clustersBasePath}/:clusterId?`;
export const clustersListPath = `${mainPath}/clusters-pf`;
export const integrationsPath = `${mainPath}/integrations`;
export const integrationsListPath = `${integrationsPath}/:source/:type`;
export const integrationCreatePath = `${integrationsPath}/:source/:type/create`;
export const integrationDetailsPath = `${integrationsPath}/:source/:type/view/:id`;
export const integrationEditPath = `${integrationsPath}/:source/:type/edit/:id`;
export const policyManagementBasePath = `${mainPath}/policy-management`;
export const policiesBasePath = `${policyManagementBasePath}/policies`;
export const policiesPath = `${policiesBasePath}/:policyId?/:command?`;
export const deprecatedPoliciesBasePath = `${mainPath}/policies`;
export const deprecatedPoliciesPath = `${deprecatedPoliciesBasePath}/:policyId?/:command?`;
export const riskBasePath = `${mainPath}/risk`;
export const riskPath = `${riskBasePath}/:deploymentId?`;
export const secretsPath = `${mainPath}/configmanagement/secrets/:secretId?`;
export const apidocsPath = `${mainPath}/apidocs`;
export const accessControlPath = `${mainPath}/access`;
export const accessControlBasePathV2 = `${mainPath}/access-control`;
export const accessControlPathV2 = `${accessControlBasePathV2}/:entitySegment?/:entityId?`;
export const userBasePath = `${mainPath}/user`;
export const userRolePath = `${userBasePath}/roles/:roleName`;
export const systemConfigPath = `${mainPath}/systemconfig`;
export const complianceBasePath = `${mainPath}/compliance`;
export const compliancePath = `${mainPath}/:context(compliance)`;
export const configManagementPath = `${mainPath}/configmanagement`;
export const dataRetentionPath = `${mainPath}/retention`;
export const systemHealthPath = `${mainPath}/system-health`;
export const systemHealthPathPF = `${mainPath}/system-health-pf`;
export const productDocsPath = '/docs/product';

// Vuln Management Paths

export const vulnManagementPath = `${mainPath}/vulnerability-management`;
export const vulnManagementPoliciesPath = `${vulnManagementPath}/policies`;
export const vulnManagementCVEsPath = `${vulnManagementPath}/cves`;
export const vulnManagementImageCVEsPath = `${vulnManagementPath}/image-cves`;
export const vulnManagementNodeCVEsPath = `${vulnManagementPath}/node-cves`;
export const vulnManagementPlatformCVEsPath = `${vulnManagementPath}/cluster-cves`;
export const vulnManagementClustersPath = `${vulnManagementPath}/clusters`;
export const vulnManagementNamespacesPath = `${vulnManagementPath}/namespaces`;
export const vulnManagementDeploymentsPath = `${vulnManagementPath}/deployments`;
export const vulnManagementImagesPath = `${vulnManagementPath}/images`;
export const vulnManagementComponentsPath = `${vulnManagementPath}/components`;
export const vulnManagementNodesPath = `${vulnManagementPath}/nodes`;

// The following paths are not part of the infinite nesting Workflow in Vuln Management
export const vulnManagementReportsPath = `${vulnManagementPath}/reports`;
export const vulnManagementReportsPathWithParam = `${vulnManagementPath}/reports/:reportId`;

export const vulnManagementRiskAcceptancePath = `${vulnManagementPath}/risk-acceptance`;
export const vulnManagementPendingApprovalsPath = `${vulnManagementRiskAcceptancePath}/pending-approvals`;
export const vulnManagementApprovedDeferralsPath = `${vulnManagementRiskAcceptancePath}/approved-deferrals`;
export const vulnManagementApprovedFalsePositivesPath = `${vulnManagementRiskAcceptancePath}/approved-false-positives`;

/**
 * New Framwork-related route paths
 */

export const urlEntityListTypes = {
    [resourceTypes.NAMESPACE]: 'namespaces',
    [resourceTypes.CLUSTER]: 'clusters',
    [resourceTypes.NODE]: 'nodes',
    [resourceTypes.DEPLOYMENT]: 'deployments',
    [resourceTypes.IMAGE]: 'images',
    [resourceTypes.SECRET]: 'secrets',
    [resourceTypes.POLICY]: 'policies',
    [resourceTypes.CVE]: 'cves',
    [resourceTypes.IMAGE_CVE]: 'image-cves',
    [resourceTypes.NODE_CVE]: 'node-cves',
    [resourceTypes.CLUSTER_CVE]: 'cluster-cves',
    [resourceTypes.COMPONENT]: 'components',
    [standardEntityTypes.CONTROL]: 'controls',
    [rbacConfigTypes.SERVICE_ACCOUNT]: 'serviceaccounts',
    [rbacConfigTypes.SUBJECT]: 'subjects',
    [rbacConfigTypes.ROLE]: 'roles',
};

export const urlEntityTypes = {
    [resourceTypes.NAMESPACE]: 'namespace',
    [resourceTypes.CLUSTER]: 'cluster',
    [resourceTypes.NODE]: 'node',
    [resourceTypes.DEPLOYMENT]: 'deployment',
    [resourceTypes.IMAGE]: 'image',
    [resourceTypes.SECRET]: 'secret',
    [resourceTypes.POLICY]: 'policy',
    [resourceTypes.CVE]: 'cve',
    [resourceTypes.IMAGE_CVE]: 'image-cve',
    [resourceTypes.NODE_CVE]: 'node-cve',
    [resourceTypes.CLUSTER_CVE]: 'cluster-cve',
    [resourceTypes.COMPONENT]: 'component',
    [standardEntityTypes.CONTROL]: 'control',
    [standardEntityTypes.STANDARD]: 'standard',
    [rbacConfigTypes.SERVICE_ACCOUNT]: 'serviceaccount',
    [rbacConfigTypes.SUBJECT]: 'subject',
    [rbacConfigTypes.ROLE]: 'role',
};

/*
 * Declare feature flags and resource requirements for route paths.
 */

export type RouteDescriptor = {
    featureFlagDependency?: FeatureFlagEnvVar;
    readAccessPredicate: ReadAccessPredicate;
};

// Evaluate resource requirements given user permissions via hasReadAccess function.
type ReadAccessPredicate = (hasReadAccess: HasReadAccess) => boolean;

/*
 * Simplified resources like Access might remove need for the following:
 *
 * Evaluate routes like access-control that render alternative sub-routes,
 * like auth-providers or roles, which have different resource requirements.
 * Export the sub-route predicates so container can render permitted subset of sub-routes.
 */
/*
function readAccessAlternatives(readAccessPredicates: ReadAccessPredicate[]): ReadAccessPredicate {
    return (hasReadAccess) =>
        readAccessPredicates.some((readAccessPredicate) => readAccessPredicate(hasReadAccess));
}
*/

/*
 * Call with 0 resource names to render route like dashboard unconditionally.
 * Call with 1 or more resource names to require read access for all of them.
 */
function readAccessResourceNames(resourceNames: ResourceName[]): ReadAccessPredicate {
    return (hasReadAccess) => resourceNames.every(hasReadAccess);
}

const renderUnconditionally = readAccessResourceNames([]);

/*
 * Map key is base path of route like violationsBasePath = '/main/violations'
 * not including parameters in path prop of some Route elements like path="/main/violations/:alertId?"
 *
 * Specify only resource requirements for primary requests of containers.
 * Rendered containers are responsible for conditional rendering of elements:
 * data might depend on resource requirements for secondary requests
 * buttons might depend on hasReadWriteAccess instead of hasReadAccess
 */
export const routeDescriptorMap: Record<string, RouteDescriptor> = {
    // Sidebar Unexpandable1
    [dashboardPath]: {
        readAccessPredicate: renderUnconditionally,
    },
    [dashboardPathPF]: {
        featureFlagDependency: 'ROX_SECURITY_METRICS_PHASE_ONE',
        readAccessPredicate: renderUnconditionally,
    },
    [networkBasePath]: {
        readAccessPredicate: readAccessResourceNames([]), // NetworkGraph, and NetworkPolicy?
    },
    [violationsBasePath]: {
        readAccessPredicate: readAccessResourceNames([]), // Alert
    },
    [complianceBasePath]: {
        readAccessPredicate: readAccessResourceNames([]), // Compliance
    },

    // Sidebar VulnerabilityManagement
    [vulnManagementPath]: {
        readAccessPredicate: readAccessResourceNames([]),
    },
    [vulnManagementRiskAcceptancePath]: {
        readAccessPredicate: readAccessResourceNames([]), // VulnerabilityManagementApprovals and/or VulnerabilityManagementRequests?
    },
    [vulnManagementReportsPath]: {
        readAccessPredicate: readAccessResourceNames(['VulnerabilityReports']),
    },

    // Sidebar Unexpandable2
    [configManagementPath]: {
        readAccessPredicate: readAccessResourceNames([]),
    },
    [riskBasePath]: {
        readAccessPredicate: readAccessResourceNames([]), // Deployment, and DeploymentExtension?
    },

    // Sidebar PlatformConfiguration
    [clustersBasePath]: {
        readAccessPredicate: readAccessResourceNames([]), // Cluster
    },
    /*
    [clustersListPath]: {
        featureFlagDependency: 'ROX_TODO', // replace conditional development rendering with backend feature flag
        readAccessPredicate: readAccessResourceNames([]), // Cluster
    },
    */
    [policiesBasePath]: {
        readAccessPredicate: readAccessResourceNames([]), // Policy
    },
    [integrationsPath]: {
        readAccessPredicate: readAccessResourceNames([]), // AuthPlugin is obsolete? APIToken; BackupPlugins, ImageIntegration, Notifier, SignatureIntegration; superseded by Integration?
    },
    [accessControlBasePathV2]: {
        readAccessPredicate: readAccessResourceNames([]), // Access
    },
    [systemConfigPath]: {
        readAccessPredicate: readAccessResourceNames(['Config']),
    },
    [systemHealthPath]: {
        readAccessPredicate: readAccessResourceNames([]),
    },
    [systemHealthPathPF]: {
        featureFlagDependency: 'ROX_SYSTEM_HEALTH_PF',
        readAccessPredicate: readAccessResourceNames([]),
    },

    // Header
    [apidocsPath]: {
        readAccessPredicate: readAccessResourceNames([]),
    },
    // Help Center is an external link to /docs/product
    [userBasePath]: {
        readAccessPredicate: readAccessResourceNames([]),
    },
};

/*
 * Evaluate feature flags and resource requirements for route paths.
 */

export type IsRenderedRoutePath = (routePath: string) => boolean;

/*
 * Higher-order function if caller needs to have predicate functions in its scope.
 * For example, MainPath because:
 * Body needs both isFeatureFlagEnabled and isRenderedRoutePath.
 * NaviationSidebar needs only isRenderedRoutePath.
 */
export function getIsRenderedRoutePath(
    hasReadAccess: HasReadAccess,
    isFeatureFlagEnabled: IsFeatureFlagEnabled
): IsRenderedRoutePath {
    return (routePath: string) => {
        const routeDescriptor = routeDescriptorMap[routePath];

        // Delete if we replace string with RoutePath string union type.
        if (routeDescriptor === undefined) {
            return true;
        }

        const { featureFlagDependency, readAccessPredicate } = routeDescriptor;

        if (typeof featureFlagDependency === 'string') {
            if (!isFeatureFlagEnabled(featureFlagDependency)) {
                return false;
            }
        }

        return readAccessPredicate(hasReadAccess);
    };
}

/*
 * Hook if caller does not need to have predicate functions in its scope.
 * It might be a moot point if pages need usePermissions and hasReadWriteAccess for conditional rendering of buttons.
 */
export function useRenderedRoutePaths(): IsRenderedRoutePath {
    const { isFeatureFlagEnabled } = useFeatureFlags();
    const { hasReadAccess } = usePermissions();

    return getIsRenderedRoutePath(hasReadAccess, isFeatureFlagEnabled);
}

/*
 * Labels for route paths.
 * Map key is base path like routeDescriptorMap above.
 */

const vulnManagementPathToLabelMap = {
    [vulnManagementPath]: 'Dashboard',
    // TODO: add mapping for Deferrals
    [vulnManagementReportsPath]: 'Reporting',
    [vulnManagementRiskAcceptancePath]: 'Risk Acceptance',
};

export const basePathToLabelMap = {
    [dashboardPath]: 'Dashboard',
    [networkBasePath]: 'Network Graph',
    [violationsBasePath]: 'Violations',
    [complianceBasePath]: 'Compliance',
    ...vulnManagementPathToLabelMap,
    [configManagementPath]: 'Configuration Management',
    [riskBasePath]: 'Risk',
    [apidocsPath]: 'API Reference',
    [productDocsPath]: 'Help Center',
    [clustersBasePath]: 'Clusters',
    [policiesBasePath]: 'Policy Management',
    [integrationsPath]: 'Integrations',
    [accessControlPath]: 'Access Control',
    [accessControlBasePathV2]: 'Access Control',
    [systemConfigPath]: 'System Configuration',
    [systemHealthPath]: 'System Health',
    [loginPath]: 'Log In',
    [userBasePath]: 'User Profile',
};

const entityListTypeMatcher = `(${Object.values(urlEntityListTypes).join('|')})`;
const entityTypeMatcher = `(${Object.values(urlEntityTypes).join('|')})`;

export const workflowPaths = {
    DASHBOARD: `${mainPath}/:context`,
    LIST: `${mainPath}/:context/:pageEntityListType${entityListTypeMatcher}/:entityId1?/:entityType2?/:entityId2?`,
    ENTITY: `${mainPath}/:context/:pageEntityType${entityTypeMatcher}/:pageEntityId?/:entityType1?/:entityId1?/:entityType2?/:entityId2?`,
};
