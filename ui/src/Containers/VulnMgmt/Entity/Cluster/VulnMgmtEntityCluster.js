import React from 'react';
import gql from 'graphql-tag';

import useCases from 'constants/useCaseTypes';
import { entityComponentPropTypes, entityComponentDefaultProps } from 'constants/entityPageProps';
import entityTypes from 'constants/entityTypes';
import { defaultCountKeyMap } from 'constants/workflowPages.constants';
import queryService from 'modules/queryService';
import { VULN_CVE_LIST_FRAGMENT } from 'Containers/VulnMgmt/VulnMgmt.fragments';
import WorkflowEntityPage from 'Containers/Workflow/WorkflowEntityPage';
import {
    getPolicyQueryVar,
    tryUpdateQueryWithVulMgmtPolicyClause
} from '../VulnMgmtPolicyQueryUtil';
import VulnMgmtClusterOverview from './VulnMgmtClusterOverview';
import EntityList from '../../List/VulnMgmtList';

const VulmMgmtEntityCluster = ({ entityId, entityListType, search, sort, page, entityContext }) => {
    const overviewQuery = gql`
        query getCluster($id: ID!, $policyQuery: String, $query: String) {
            result: cluster(id: $id) {
                id
                name
                priority
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
                #createdAt
                status {
                    orchestratorMetadata {
                        version
                    }
                }
                #istioEnabled
                policyCount(query: $policyQuery)
                namespaceCount
                deploymentCount
                imageCount
                componentCount
                vulnCount
                vulnerabilities: vulns(query: $query) {
                    ...cveFields
                }
            }
        }
        ${VULN_CVE_LIST_FRAGMENT}
    `;

    function getListQuery(listFieldName, fragmentName, fragment) {
        // @TODO: remove this hack after we are able to search for k8s vulns
        const parsedListFieldName =
            search && search['Vulnerability Type'] ? 'vulns: k8sVulns' : listFieldName;

        return gql`
        query getCluster_${entityListType}($id: ID!, $pagination: Pagination, $query: String${getPolicyQueryVar(
            entityListType
        )}) {
            result: cluster(id: $id) {
                id
                ${defaultCountKeyMap[entityListType]}(query: $query)
                ${parsedListFieldName}(query: $query, pagination: $pagination) { ...${fragmentName} }
            }
        }
        ${fragment}
    `;
    }

    const queryOptions = {
        variables: {
            id: entityId,
            query: tryUpdateQueryWithVulMgmtPolicyClause(entityListType, search, entityContext)
        }
    };

    // @TODO: improve the handling of query vs. policyQuery in the UI Workflow
    //   more info: because Vuln Mgmt requires Policy sub-queries to be filtered by "Category": "Vulnerability Management"
    //   we added a separate query var for just policies, where needed
    //   but it doesn't play well with the initial implementation of backend pagination,
    //   hence, this hack below to limit to only the overview and entities where it is needed
    if (
        !entityListType ||
        entityListType === entityTypes.DEPLOYMENT ||
        entityListType === entityTypes.NAMESPACE
    ) {
        queryOptions.variables.policyQuery = queryService.objectToWhereClause({
            Category: 'Vulnerability Management'
        });
    }

    return (
        <WorkflowEntityPage
            entityId={entityId}
            entityType={entityTypes.CLUSTER}
            entityListType={entityListType}
            useCase={useCases.VULN_MANAGEMENT}
            ListComponent={EntityList}
            OverviewComponent={VulnMgmtClusterOverview}
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

VulmMgmtEntityCluster.propTypes = entityComponentPropTypes;
VulmMgmtEntityCluster.defaultProps = entityComponentDefaultProps;

export default VulmMgmtEntityCluster;
