import { hasFeatureFlag } from '../../helpers/features';
import withAuth from '../../helpers/basicAuth';
import {
    assertSortedItems,
    callbackForPairOfAscendingNumberValuesFromElements,
    callbackForPairOfDescendingNumberValuesFromElements,
} from '../../helpers/sort';
import {
    hasTableColumnHeadings,
    interactAndWaitForVulnerabilityManagementEntities,
    verifyFilteredSecondaryEntitiesLink,
    verifySecondaryEntities,
    visitVulnerabilityManagementEntities,
} from '../../helpers/vulnmanagement/entities';

const entitiesKey = 'deployments';

export function getCountAndNounFromCVEsLinkResults([, count]) {
    return {
        panelHeaderText: `${count} ${count === '1' ? 'CVE' : 'CVES'}`,
        relatedEntitiesCount: count,
        relatedEntitiesNoun: count === '1' ? 'CVE' : 'CVES',
    };
}

describe('Vulnerability Management Deployments', () => {
    withAuth();

    before(function beforeHook() {
        if (hasFeatureFlag('ROX_POSTGRES_DATASTORE')) {
            this.skip();
        }
    });
    it('should display table columns', () => {
        visitVulnerabilityManagementEntities(entitiesKey);

        hasTableColumnHeadings([
            '', // hidden
            'Deployment',
            'CVEs',
            'Latest Violation',
            'Policy Status',
            'Cluster',
            'Namespace',
            'Images',
            'Risk Priority',
        ]);
    });

    it('should sort the Risk Priority column', () => {
        visitVulnerabilityManagementEntities(entitiesKey);

        const thSelector = '.rt-th:contains("Risk Priority")';
        const tdSelector = '.rt-td:nth-child(9)';

        // 0. Initial table state indicates that the column is sorted ascending.
        cy.get(thSelector).should('have.class', '-sort-asc');
        cy.get(tdSelector).then((items) => {
            assertSortedItems(items, callbackForPairOfAscendingNumberValuesFromElements);
        });

        // 1. Sort descending by the column.
        interactAndWaitForVulnerabilityManagementEntities(() => {
            cy.get(thSelector).click();
        }, entitiesKey);
        cy.location('search').should(
            'eq',
            '?sort[0][id]=Deployment%20Risk%20Priority&sort[0][desc]=true'
        );

        cy.get(thSelector).should('have.class', '-sort-desc');
        cy.get(tdSelector).then((items) => {
            assertSortedItems(items, callbackForPairOfDescendingNumberValuesFromElements);
        });

        // 2. Sort ascending by the column.
        cy.get(thSelector).click(); // no request because initial response has been cached
        cy.location('search').should(
            'eq',
            '?sort[0][id]=Deployment%20Risk%20Priority&sort[0][desc]=false'
        );

        cy.get(thSelector).should('have.class', '-sort-asc');
        // Do not assert because of potential timing problem: get td elements before table re-renders.
    });

    // Argument 3 in verify functions is one-based index of column which has the links.

    // Some tests might fail in local deployment.

    it('should display links for all CVEs', () => {
        verifySecondaryEntities(
            entitiesKey,
            'cves',
            2,
            /^\d+ CVEs?$/,
            getCountAndNounFromCVEsLinkResults
        );
    });

    it('should display links for fixable CVEs', () => {
        verifyFilteredSecondaryEntitiesLink(
            entitiesKey,
            'cves',
            2,
            /^\d+ Fixable$/,
            getCountAndNounFromCVEsLinkResults
        );
    });

    it('should display links for images', () => {
        verifySecondaryEntities(entitiesKey, 'images', 7, /^\d+ images?$/);
    });
});