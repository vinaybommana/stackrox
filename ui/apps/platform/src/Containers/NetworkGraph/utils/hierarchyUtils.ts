import { SearchFilter } from 'types/search';
import { NamespaceWithDeployments } from 'hooks/useFetchNamespaceDeployments';
import { Cluster } from 'types/cluster.proto';

export type NetworkScopeHierarchy = {
    cluster: {
        id: string;
        name: string;
    };
    namespaces: string[];
    deployments: string[];
    remainingQuery: Omit<SearchFilter, 'Cluster' | 'Namespace' | 'Deployment'>;
};

export function getScopeHierarchyFromSearch(
    searchFilter: SearchFilter,
    clusters: Cluster[]
): NetworkScopeHierarchy | null {
    const urlCluster = searchFilter.Cluster;
    if (!urlCluster || Array.isArray(urlCluster)) {
        return null;
    }

    const cluster = clusters.find((cl) => cl.name === urlCluster);
    if (!cluster) {
        return null;
    }

    const workingQuery = { ...searchFilter };
    const hierarchy: NetworkScopeHierarchy = {
        cluster,
        namespaces: [],
        deployments: [],
        remainingQuery: workingQuery,
    };

    if (searchFilter.Namespace) {
        hierarchy.namespaces = Array.isArray(searchFilter.Namespace)
            ? searchFilter.Namespace
            : [searchFilter.Namespace];
        delete hierarchy.remainingQuery.Namespace;
    }

    if (searchFilter.Deployment) {
        hierarchy.deployments = Array.isArray(searchFilter.Deployment)
            ? searchFilter.Deployment
            : [searchFilter.Deployment];
    }

    return hierarchy;
}

export function getDeploymentLookupMap(
    deploymentsByNamespace: NamespaceWithDeployments[]
): Record<string, string[]> {
    return deploymentsByNamespace.reduce<Record<string, string[]>>((acc, ns) => {
        const deployments = ns.deployments.map((deployment) => deployment.name);
        return { ...acc, [ns.metadata.name]: deployments };
    }, {});
}

export function getDeploymentsAllowedByNamespaces(
    deploymentLookupMap: Record<string, string[]>,
    namespaceSelection: string[]
) {
    const newDeploymentLookup = Object.fromEntries(
        Object.entries(deploymentLookupMap).filter(([key]) => namespaceSelection.includes(key))
    );
    const allowedDeployments = Object.values(newDeploymentLookup).flat(1);

    return allowedDeployments;
}

export default {
    getScopeHierarchyFromSearch,
    getDeploymentLookupMap,
};
