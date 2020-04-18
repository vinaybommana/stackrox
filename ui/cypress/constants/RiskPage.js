import tableSelectors from '../selectors/table';
import { processTagsSelectors } from '../selectors/tags';
import { processCommentsSelectors, commentsDialogSelectors } from '../selectors/comments';
import scopeSelectors from '../helpers/scopeSelectors';

export const url = '/main/risk';

export const errorMessages = {
    deploymentNotFound: 'Deployment not found',
    riskNotFound: 'Risk not found',
    processNotFound: 'No processes discovered'
};

const sidePanelSelectors = scopeSelectors('[data-testid="panel"]:eq(1)', {
    firstProcessCard: scopeSelectors('[data-testid="process-discovery-card"]:first', {
        header: '[data-testid="process"]',
        tags: processTagsSelectors,
        comments: processCommentsSelectors
    }),

    riskIndicatorsTab: 'button[data-testid="tab"]:contains("Risk Indicators")',
    deploymentDetailsTab: 'button[data-testid="tab"]:contains("Deployment Details")',
    processDiscoveryTab: 'button[data-testid="tab"]:contains("Process Discovery")',

    cancelButton: 'button[data-testid="cancel"]'
});

export const selectors = {
    risk: 'nav.left-navigation li:contains("Risk") a',
    errMgBox: 'div.error-message',
    panelTabs: {
        riskIndicators: 'button[data-testid="tab"]:contains("Risk Indicators")',
        deploymentDetails: 'button[data-testid="tab"]:contains("Deployment Details")',
        processDiscovery: 'button[data-testid="tab"]:contains("Process Discovery")'
    },
    cancelButton: 'button[data-testid="cancel"]',
    search: {
        searchLabels: '.react-select__multi-value__label',
        // selectors for legacy tests
        searchModifier: '.react-select__multi-value__label:first',
        searchWord: '.react-select__multi-value__label:eq(1)'
    },
    mounts: {
        label: 'div:contains("Mounts"):last',
        items: 'div:contains("Mounts"):last + ul li div'
    },
    imageLink: 'div:contains("Image Name") + a',
    table: scopeSelectors('[data-testid="panel"]:first', tableSelectors),
    collapsible: {
        card: '.Collapsible',
        header: '.Collapsible__trigger',
        body: '.Collapsible__contentInner'
    },
    suspiciousProcesses: "[data-testid='suspicious-process']",
    networkNodeLink: '[data-testid="network-node-link"]',
    sidePanel: sidePanelSelectors,
    commentsDialog: commentsDialogSelectors
};
