import React from 'react';
import gql from 'graphql-tag';

import useCases from 'constants/useCaseTypes';
import { workflowEntityPropTypes, workflowEntityDefaultProps } from 'constants/entityPageProps';
import queryService from 'modules/queryService';
import entityTypes from 'constants/entityTypes';
import { defaultCountKeyMap } from 'constants/workflowPages.constants';
import WorkflowEntityPage from 'Containers/Workflow/WorkflowEntityPage';
import { VULN_CVE_LIST_FRAGMENT } from 'Containers/VulnMgmt/VulnMgmt.fragments';
import VulnMgmtNamespaceOverview from './VulnMgmtNamespaceOverview';
import EntityList from '../../List/VulnMgmtList';
import { getPolicyQueryVar, getQueryVar } from '../VulnMgmtPolicyQueryUtil';

const VulnMgmtNamespace = ({ entityId, entityListType, search, sort, page, entityContext }) => {
    const overviewQuery = gql`
        query getNamespace($id: ID!, $policyQuery: String, $query: String) {
            result: namespace(id: $id) {
                metadata {
                    priority
                    name
                    ${entityContext[entityTypes.CLUSTER] ? '' : 'clusterName clusterId'}
                    id
                    labels {
                        key
                        value
                    }
                }
                policyStatus(query: $policyQuery) {
                    status
                    failingPolicies {
                        id
                        name
                        description
                        policyStatus
                        latestViolation
                        severity
                        deploymentCount(query: $query)
                        lifecycleStages
                        enforcementActions
                    }
                }
                policyCount(query: $policyQuery)
                vulnCount
                deploymentCount: numDeployments 
                imageCount 
                componentCount(query: $query)
                vulnerabilities: vulns(query: $query) {
                    ...cveFields
                }
            }
        }
        ${VULN_CVE_LIST_FRAGMENT}
    `;

    function getListQuery(listFieldName, fragmentName, fragment) {
        return gql`
        query getNamespace${entityListType}($id: ID!, $pagination: Pagination, $query: String${getPolicyQueryVar(
            entityListType
        )}) {
            result: namespace(id: $id) {
                metadata {
                    id
                }
                ${defaultCountKeyMap[entityListType]}(query: ${getQueryVar(entityListType)})
                ${listFieldName}(query: ${getQueryVar(
            entityListType
        )}, pagination: $pagination) { ...${fragmentName} }
            }
        }
        ${fragment}
    `;
    }
    const newEntityContext = { ...entityContext, [entityTypes.NAMESPACE]: entityId };

    const entityContextObj = queryService.entityContextToQueryObject(newEntityContext);

    const queryOptions = {
        variables: {
            id: entityId,
            query: queryService.objectToWhereClause({ ...search, ...entityContextObj }),
            policyQuery: queryService.objectToWhereClause({ Category: 'Vulnerability Management' })
        }
    };

    return (
        <WorkflowEntityPage
            entityId={entityId}
            entityType={entityTypes.NAMESPACE}
            entityListType={entityListType}
            useCase={useCases.VULN_MANAGEMENT}
            ListComponent={EntityList}
            OverviewComponent={VulnMgmtNamespaceOverview}
            overviewQuery={overviewQuery}
            getListQuery={getListQuery}
            search={search}
            sort={sort}
            page={page}
            queryOptions={queryOptions}
            entityContext={entityContext}
        />
    );
};

VulnMgmtNamespace.propTypes = workflowEntityPropTypes;
VulnMgmtNamespace.defaultProps = workflowEntityDefaultProps;

export default VulnMgmtNamespace;
