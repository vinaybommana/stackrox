import React from 'react';
import gql from 'graphql-tag';
import WorkflowEntityPage from 'Containers/Workflow/WorkflowEntityPage';
import queryService from 'modules/queryService';
import entityTypes from 'constants/entityTypes';
import { defaultCountKeyMap } from 'constants/workflowPages.constants';
import useCases from 'constants/useCaseTypes';
import { VULN_CVE_LIST_FRAGMENT } from 'Containers/VulnMgmt/VulnMgmt.fragments';
import EntityList from '../../List/VulnMgmtList';
import VulnMgmtComponentOverview from './VulnMgmtComponentOverview';
import {
    getPolicyQueryVar,
    tryUpdateQueryWithVulMgmtPolicyClause
} from '../VulnMgmtPolicyQueryUtil';

const VulnMgmtComponent = ({ entityId, entityListType, search, entityContext, sort, page }) => {
    const overviewQuery = gql`
        query getComponent($id: ID!, $query: String) {
            result: component(id: $id) {
                id
                name
                version
                priority
                vulnCount(query: $query)
                deploymentCount(query: $query)
                imageCount(query: $query)
                topVuln {
                    cvss
                    scoreVersion
                }
                fixableCVEs: vulns(query: "Fixed By:r/.*") {
                    ...cveFields
                }
            }
        }
        ${VULN_CVE_LIST_FRAGMENT}
    `;

    function getListQuery(listFieldName, fragmentName, fragment) {
        return gql`
        query getComponentSubEntity${entityListType}($id: ID!, $pagination: Pagination, $query: String${getPolicyQueryVar(
            entityListType
        )}) {
            result: component(id: $id) {
                id
                ${defaultCountKeyMap[entityListType]}(query: $query)
                ${listFieldName}(query: $query, pagination: $pagination) { ...${fragmentName} }
            }
        }
        ${fragment}
    `;
    }

    const queryOptions = {
        variables: {
            id: entityId,
            query: tryUpdateQueryWithVulMgmtPolicyClause(entityListType, search, entityContext),
            policyQuery: queryService.objectToWhereClause({ Category: 'Vulnerability Management' })
        }
    };

    return (
        <WorkflowEntityPage
            entityId={entityId}
            entityType={entityTypes.COMPONENT}
            entityListType={entityListType}
            useCase={useCases.VULN_MANAGEMENT}
            ListComponent={EntityList}
            OverviewComponent={VulnMgmtComponentOverview}
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

export default VulnMgmtComponent;
