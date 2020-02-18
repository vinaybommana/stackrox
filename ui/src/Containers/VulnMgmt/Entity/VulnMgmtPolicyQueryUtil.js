import entityTypes from 'constants/entityTypes';
import queryService from 'modules/queryService';

const entitiesWithPolicyField = [entityTypes.CLUSTER, entityTypes.NAMESPACE, entityTypes.POLICY];

export function getPolicyQueryVar(entityType) {
    return entitiesWithPolicyField.includes(entityType) ? ', $policyQuery: String' : '';
}

export function tryUpdateQueryWithVulMgmtPolicyClause(entityType, search, entityContext) {
    const whereObj = { ...search, ...queryService.entityContextToQueryObject(entityContext) };
    return entityType === entityTypes.POLICY
        ? queryService.objectToWhereClause({ ...whereObj, Category: 'Vulnerability Management' })
        : queryService.objectToWhereClause(whereObj);
}

// returns `policyQuery` if the subquery is for a policy; else returns regular `query`
export function getQueryVar(entityType) {
    return entityType === entityTypes.POLICY ? '$policyQuery' : '$query';
}
