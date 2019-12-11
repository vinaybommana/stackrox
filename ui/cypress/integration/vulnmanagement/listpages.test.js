import { url, selectors } from '../../constants/VulnManagementPage';
import withAuth from '../../helpers/basicAuth';
import checkFeatureFlag from '../../helpers/features';

const hasExpectedHeaderColumns = colNames => {
    colNames.forEach(col => {
        cy.get(`${selectors.tableColumn}:contains('${col}')`);
    });
};

function validateDataInEntityListPage(entityCountAndName, entityURL) {
    cy.get(selectors.entityRowHeader)
        .invoke('text')
        .then(entityCountFromHeader => {
            if (entityCountAndName.includes('CVE')) {
                const numEntitiesListPage = parseInt(entityCountFromHeader, 10);
                const numEntitiesParentPage = parseInt(entityCountAndName, 10);
                expect(numEntitiesListPage - numEntitiesParentPage).to.be.lessThan(6);
            } else {
                expect(entityCountFromHeader).contains(
                    parseInt(entityCountAndName, 10),
                    `expected entity count ${entityCountAndName} found in the related entity list page`
                );
            }
        });
    cy.visit(entityURL);
}
function validateClickableLinks(colLinks, parentUrl) {
    colLinks.forEach(col => {
        if (col !== 'Policies') {
            cy.get(`${selectors.tableColumnLinks}:contains('${col.toLowerCase()}')`)
                .invoke('text')
                .then(value => {
                    cy.get(
                        `${selectors.tableColumnLinks}:contains('${col.toLowerCase()}')`
                    ).click();
                    validateDataInEntityListPage(value, parentUrl);
                });
        }
        if (col === 'Policies') {
            cy.get(`${selectors.tableColumnLinks}:contains('${col.toLowerCase()}')`)
                .invoke('text')
                .then(value => {
                    expect(value).contains('policies' || 'policy', 'expected text displayed');
                    cy.get(`${selectors.tableColumnLinks}:contains('${value}')`).click();
                    validateDataInEntityListPage(value, parentUrl);
                });
        }
    });
}

function validateClickableLinksEntityListPage(colLinks, parentUrl) {
    colLinks.forEach(col => {
        if (col !== 'Policies') {
            cy.get(`${selectors.tableColumnLinks}:contains('${col.toLowerCase()}')`)
                .eq(0)
                .invoke('text')
                .then(value => {
                    cy.get(`${selectors.tableColumnLinks}:contains('${col.toLowerCase()}')`)
                        .eq(0)
                        .click({ force: true });
                    validateDataInEntityListPage(value, parentUrl);
                });
        }
        if (col === 'Policies') {
            // handle "1 policy", "No failing policies", or "X policies"
            cy.get(`${selectors.tableColumnLinks}:contains('polic')`)
                .eq(0)
                .invoke('text')
                .then(value => {
                    expect(value).contains('polic', 'expected text displayed');
                    cy.get(`${selectors.tableColumnLinks}:contains('${value}')`)
                        .eq(0)
                        .click({ force: true });
                    validateDataInEntityListPage(`${parseInt(value, 10)} polic`, parentUrl);
                });
        }
    });
}

function validateAllCVELinks(prevUrl) {
    cy.get(`${selectors.allCVEColumnLink}`)
        .eq(0)
        .invoke('text')
        .then(value => {
            cy.get(`${selectors.allCVEColumnLink}`)
                .eq(0)
                .click({ force: true });
            validateDataInEntityListPage(value.toUpperCase(), prevUrl);
        });
}

function validateFixableCVELinks(urlBack) {
    cy.get(`${selectors.fixableCVELink}`)
        .eq(0)
        .invoke('text')
        .then(value => {
            cy.get(`${selectors.fixableCVELink}`)
                .eq(0)
                .click({ force: true });
            if (parseInt(value, 10) === 1)
                validateDataInEntityListPage(`${parseInt(value, 10)} CVE`, urlBack);
            if (parseInt(value, 10) > 1)
                validateDataInEntityListPage(`${parseInt(value, 10)} CVES`, urlBack);
        });
}

function validateSort(selector) {
    let current;
    let prev;
    prev = -1000;
    cy.get(selector).each($el => {
        current = parseInt($el.text(), 10);
        const sortOrderStatus = current >= prev;
        expect(sortOrderStatus).to.equals(true, 'sort order is as expected');
        prev = current;
    });
}

function validateSortForCVE(selector) {
    let current;
    let prev;
    let sortOrderStatus = false;
    prev = 1000;
    cy.get(selector).each($el => {
        current = parseFloat($el.text(), 10.0);
        // eslint-disable-next-line no-restricted-globals
        if (!isNaN(prev) && !isNaN(current)) {
            sortOrderStatus = current <= prev;
            expect(sortOrderStatus).to.equals(true, 'sort order is as expected');
            prev = current;
        }
    });
}

function validateTileLinksSidePanelEntityPage(colSelector, relatedEntitiesList, parentUrl) {
    relatedEntitiesList.forEach(col => {
        if (col !== 'CVEs' && col !== 'Fixable' && col !== 'Policies') {
            cy.get(`${selectors.tableColumnLinks}:contains('${col.toLowerCase()}')`)
                .invoke('text')
                .then(value => {
                    cy.get(colSelector)
                        .eq(0)
                        .click({ force: true });
                    cy.get(selectors.getTileLink(col.toUpperCase()))
                        .find(selectors.tileLinkText)
                        .contains(parseInt(value, 10));
                    cy.get(selectors.getTileLink(col.toUpperCase()))
                        .find(selectors.tileLinkValue)
                        .contains(col.toUpperCase());
                    cy.visit(parentUrl);
                });
        }
        if (col === 'CVEs') {
            cy.get(`${selectors.allCVEColumnLink}`)
                .eq(0)
                .invoke('text')
                .then(value => {
                    cy.get(colSelector)
                        .eq(0)
                        .click({ force: true });
                    if (parseInt(value, 10) === 1) {
                        cy.get(selectors.getTileLink('CVE'))
                            .find(selectors.tileLinkValue)
                            .contains('CVE');
                    }
                    if (parseInt(value, 10) > 1) {
                        cy.get(selectors.getTileLink('CVE'))
                            .find(selectors.tileLinkValue)
                            .contains('CVES');
                    }
                    cy.visit(parentUrl);
                });
        }
        if (col === 'Fixable') {
            cy.get(`${selectors.fixableCVELink}`)
                .invoke('text')
                .then(value => {
                    cy.get(colSelector)
                        .eq(0)
                        .click({ force: true });
                    cy.get(selectors.tabButton)
                        .contains('Fixable CVEs')
                        .click();
                    cy.get(selectors.getSidePanelTabHeader('fixable')).contains(
                        parseInt(value, 10)
                    );
                    cy.visit(parentUrl);
                });
        }
        if (col === 'Policies') {
            cy.get(`${selectors.tableColumnLinks}`)
                .contains(/(?:policies|policy)/)
                .invoke('text')
                .then(value => {
                    if (
                        (value.includes('policies') || value.includes('policy')) &&
                        value !== 'No failing policies'
                    ) {
                        cy.get(colSelector)
                            .eq(0)
                            .click({ force: true });
                        let colText = '';
                        if (parseInt(value, 10) > 1) colText = 'POLICIES';
                        if (parseInt(value, 10) === 1) colText = 'POLICY';
                        expect(
                            cy
                                .get(selectors.getTileLink(colText))
                                .find(selectors.tileLinkText)
                                .contains(parseInt(value, 10)),
                            'policy count displayed on tile is valid'
                        );
                        expect(
                            cy
                                .get(selectors.getTileLink(colText))
                                .find(selectors.tileLinkValue)
                                .contains(colText.toUpperCase()),
                            'policy text displayed is valid'
                        );
                    }
                });
        }
    });
}

function validateTabsInSidePanelWithTileLinks(parentUrl, colSelector, relatedEntitiesList) {
    relatedEntitiesList.forEach(col => {
        if (
            col !== 'CVEs' &&
            col !== 'Fixable' &&
            col !== 'Policies' &&
            col !== 'Failing Policies'
        ) {
            cy.get(`${selectors.tableColumnLinks}:contains('${col.toLowerCase()}')`)
                .invoke('text')
                .then(value => {
                    cy.get(colSelector)
                        .eq(0)
                        .click({ force: true });
                    cy.get(selectors.sidePanelExpandButton).click({ force: true });
                    cy.get(selectors.getSidePanelTabLink(col.toLowerCase())).click({ force: true });
                    expect(cy.get(selectors.tabHeader).contains(parseInt(value, 10)));
                    cy.visit(parentUrl);
                });
        }
        if (col === 'CVEs') {
            cy.get(`${selectors.allCVEColumnLink}`)
                .eq(0)
                .invoke('text')
                .then(value => {
                    cy.get(colSelector)
                        .eq(0)
                        .click({ force: true });
                    cy.get(selectors.sidePanelExpandButton).click({ force: true });
                    cy.get(selectors.getSidePanelTabLink(col.toUpperCase())).click({ force: true });
                    expect(cy.get(selectors.tabHeader).contains(parseInt(value, 10)));
                    cy.visit(parentUrl);
                });
        }
        if (col === 'Policies' || col === 'Failing Policies') {
            cy.get(`${selectors.tableColumnLinks}`)
                .contains(/(?:policies|policy)/)
                .invoke('text')
                .then(value => {
                    if (
                        (value.includes('policies') || value.includes('policy')) &&
                        value !== 'No failing policies'
                    ) {
                        cy.get(colSelector)
                            .eq(0)
                            .click({ force: true });
                        let colText = '';
                        if (parseInt(value, 10) > 1) colText = 'POLICIES';
                        if (parseInt(value, 10) === 1) colText = 'POLICY';
                        if (col === 'Failing Policies') colText = 'POLICIES';

                        cy.get(selectors.sidePanelExpandButton).click({ force: true });
                        cy.get(selectors.getSidePanelTabLink(colText.toLowerCase())).click({
                            force: true
                        });

                        expect(cy.get(selectors.tabHeader).contains(parseInt(value, 10)));
                        cy.visit(parentUrl);
                    }
                });
        }
    });
}

describe('Entities list Page', () => {
    before(function beforeHook() {
        // skip the whole suite if vuln mgmt isn't enabled
        if (checkFeatureFlag('ROX_VULN_MGMT_UI', false)) {
            this.skip();
        }
    });

    withAuth();

    it('should display all the columns and links expected in clusters list page', () => {
        cy.visit(url.list.clusters);
        hasExpectedHeaderColumns([
            'Cluster',
            'CVEs',
            'K8S Version',
            'Namespaces',
            'Deployments',
            'Policies',
            'Policy Status',
            'Latest Violation',
            'Risk Priority'
        ]);
        validateClickableLinks(['Namespace', 'Deployment', 'Policies'], url.list.clusters);
        validateAllCVELinks(url.list.clusters);
        validateFixableCVELinks(url.list.clusters);
        validateSort(selectors.riskScoreCol);
        validateTileLinksSidePanelEntityPage(
            selectors.tableFirstColumn,
            ['Namespace', 'Deployment', 'Policies', 'CVEs', 'Fixable'],
            url.list.clusters
        );

        validateTabsInSidePanelWithTileLinks(url.list.clusters, selectors.tableFirstColumn, [
            'Namespace',
            'Deployment',
            'Policies',
            'CVEs'
        ]);
    });

    it('should display all the columns and links expected in namespaces list page', () => {
        cy.visit(url.list.namespaces);
        hasExpectedHeaderColumns([
            'Cluster',
            'CVEs',
            'Images',
            'Namespace',
            'Deployments',
            'Policies',
            'Policy Status',
            'Latest Violation',
            'Risk Priority'
        ]);
        validateClickableLinksEntityListPage(['image', 'deployment'], url.list.namespaces);
        validateAllCVELinks(url.list.namespaces);
        validateFixableCVELinks(url.list.namespaces);
        validateSort(selectors.riskScoreCol);
        validateTileLinksSidePanelEntityPage(
            selectors.tableFirstColumn,
            ['Deployment', 'Image', 'Policies', 'CVEs', 'Fixable'],
            url.list.namespaces
        );
        validateTabsInSidePanelWithTileLinks(url.list.namespaces, selectors.tableFirstColumn, [
            'Deployment',
            'Image',
            'Policies',
            'CVEs'
        ]);
    });

    it('should display all the columns and links expected in deployments list page', () => {
        cy.visit(url.list.deployments);
        hasExpectedHeaderColumns([
            'Cluster',
            'CVEs',
            'Images',
            'Namespace',
            'Deployment',
            'Policies',
            'Policy Status',
            'Latest Violation',
            'Risk Priority'
        ]);
        cy.get(selectors.tableBodyColumn).each($el => {
            const columnValue = $el.text().toLowerCase();
            if (columnValue !== 'no failing policies' && columnValue.includes('polic')) {
                validateClickableLinksEntityListPage(['Policies'], url.list.deployments);
                validateTileLinksSidePanelEntityPage(
                    selectors.tableFirstColumn,
                    ['Policies'],
                    url.list.deployments
                );
                validateTabsInSidePanelWithTileLinks(
                    url.list.deployments,
                    selectors.tableFirstColumn,
                    ['Failing Policies']
                );
            }
            if (columnValue !== 'no images' && columnValue.includes('image')) {
                validateClickableLinksEntityListPage(['image'], url.list.deployments);
                validateTileLinksSidePanelEntityPage(
                    selectors.tableFirstColumn,
                    ['Image'],
                    url.list.deployments
                );
                validateTabsInSidePanelWithTileLinks(
                    url.list.deployments,
                    selectors.tableFirstColumn,
                    ['Image']
                );
            }
            if (columnValue !== 'no cves' && columnValue.includes('cve')) {
                validateAllCVELinks(url.list.deployments);
                if (columnValue.includes('fixable')) {
                    validateFixableCVELinks(url.list.deployments);
                    validateTileLinksSidePanelEntityPage(
                        selectors.tableFirstColumn,
                        ['Fixable'],
                        url.list.deployments
                    );
                }
                validateTileLinksSidePanelEntityPage(
                    selectors.tableFirstColumn,
                    ['CVEs'],
                    url.list.deployments
                );
                validateTabsInSidePanelWithTileLinks(
                    url.list.deployments,
                    selectors.tableFirstColumn,
                    ['CVEs']
                );
            }
        });
        validateSort(selectors.riskScoreCol);
    });

    it('should display all the columns and links expected in images list page', () => {
        cy.visit(url.list.images);
        hasExpectedHeaderColumns([
            'Image',
            'CVEs',
            'Top CVSS',
            'Created',
            'Scan Time',
            'Image Status',
            'Deployments',
            'Components',
            'Risk Priority'
        ]);
        validateClickableLinksEntityListPage(['deployment', 'component'], url.list.images);
        validateAllCVELinks(url.list.images);
        validateFixableCVELinks(url.list.images);
        validateSort(selectors.riskScoreCol);
        validateTileLinksSidePanelEntityPage(
            selectors.tableFirstColumn,
            ['Deployment', 'Component'],
            url.list.images
        );
        validateTabsInSidePanelWithTileLinks(url.list.images, selectors.tableFirstColumn, [
            'Deployment',
            'Component'
        ]);
    });

    it('should display all the columns expected in components list page', () => {
        cy.visit(url.list.components);
        hasExpectedHeaderColumns([
            'Component',
            'CVEs',
            'Top CVSS',
            'Images',
            'Deployments',
            'Risk Priority'
        ]);
        validateClickableLinksEntityListPage(['deployment', 'image'], url.list.components);
        validateAllCVELinks(url.list.components);
        validateFixableCVELinks(url.list.components);
        validateSort(selectors.componentsRiskScoreCol);
        validateTileLinksSidePanelEntityPage(
            selectors.tableFirstColumn,
            ['Deployment', 'Image', 'CVEs'],
            url.list.components
        );
        validateTabsInSidePanelWithTileLinks(url.list.components, selectors.tableFirstColumn, [
            'Deployment',
            'Image',
            'CVEs'
        ]);
    });

    it('should display all the columns and links expected in cves list page', () => {
        cy.visit(url.list.cves);
        hasExpectedHeaderColumns([
            'CVE',
            'Fixable',
            'CVSS Score',
            'Env. Impact',
            'Impact Score',
            'Scanned',
            'Published',
            'Deployments'
        ]);
        validateClickableLinksEntityListPage(['image', 'deployment', 'component'], url.list.cves);
        validateSortForCVE(selectors.cvesCvssScoreCol);
        validateTileLinksSidePanelEntityPage(
            selectors.tableFirstColumn,
            ['Deployment', 'Component', 'Image'],
            url.list.cves
        );
        validateTabsInSidePanelWithTileLinks(url.list.cves, selectors.tableFirstColumn, [
            'Deployment',
            'Component',
            'Image'
        ]);
    });
});
