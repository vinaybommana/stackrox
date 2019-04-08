export const url = '/main/network';

export const selectors = {
    network: 'nav.left-navigation li:contains("Network") a',
    simulatorSuccessMessage: 'div:contains("Policies processed")',
    panels: {
        creatorPanel: '[data-test-id="network-creator-panel"]',
        uploadPanel: '[data-test-id="upload-yaml-panel"]'
    },
    legend: {
        items: '[data-test-id=legend] > div:last > div > div'
    },
    namespaces: {
        all: 'g.container > rect',
        getNamespace: namespace => `g.container > rect.namespace-${namespace}`
    },
    services: {
        all: 'g.namespace .node',
        getServicesForNamespace: namespace => `g.namespace-${namespace} .node`
    },
    links: {
        all: '.link',
        bidirectional: '.link[marker-start="url(#start)"]',
        namespaces: '.link.namespace',
        services: '.link.service'
    },
    buttons: {
        simulatorButtonOn: '[data-test-id="simulator-button-on"]',
        simulatorButtonOff: '[data-test-id="simulator-button-off"]'
    }
};
