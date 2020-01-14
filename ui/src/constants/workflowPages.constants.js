import entityTypes from 'constants/entityTypes';

export const overviewLimit = 5;

export const dashboardLimit = 8;

export const LIST_PAGE_SIZE = 50;

export const defaultCountKeyMap = {
    [entityTypes.COMPONENT]: 'componentCount',
    [entityTypes.CVE]: 'vulnCount',
    [entityTypes.DEPLOYMENT]: 'deploymentCount',
    [entityTypes.NAMESPACE]: 'namespaceCount',
    [entityTypes.IMAGE]: 'imageCount',
    [entityTypes.POLICY]: 'failingPolicyCount'
};
