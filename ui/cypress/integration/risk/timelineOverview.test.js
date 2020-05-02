import { selectors, url } from '../../constants/RiskPage';

import * as api from '../../constants/apiEndpoints';

import withAuth from '../../helpers/basicAuth';

function setRoutes() {
    cy.server();
    cy.route('GET', api.risks.riskyDeployments).as('deployments');
    cy.route('GET', api.risks.getDeploymentWithRisk).as('getDeployment');
}

function openDeployment(deploymentName) {
    cy.visit(url);
    cy.wait('@deployments');

    cy.get(`${selectors.table.rows}:contains(${deploymentName})`).click();
    cy.wait('@getDeployment');
}

describe('Risk Page Event Timeline - Timeline Overview', () => {
    withAuth();

    it('should show the correct number of events in the timeline overview', () => {
        setRoutes();
        // select a deployment to open the side panel
        openDeployment('collector');
        // open the process discovery tab
        cy.get(selectors.sidePanel.processDiscoveryTab).click();

        let sumOfEvents = 0;

        cy.get(selectors.eventTimelineOverview.eventCounts)
            // go through each tile and count up the events
            .each((element) => {
                sumOfEvents += parseInt(element.text(), 10);
            })
            // compare the summed up number of events with the total number of events
            .then(() => {
                cy.get(selectors.eventTimelineOverview.totalNumEventsText).contains(
                    `${sumOfEvents} EVENTS`
                );
            });
    });

    it('should show the timeline graph when the overview is clicked', () => {
        setRoutes();
        // select a deployment to open the side panel
        openDeployment('collector');
        // open the process discovery tab
        cy.get(selectors.sidePanel.processDiscoveryTab).click();
        // click the overview button
        cy.get(selectors.eventTimelineOverviewButton).click();
        // the event timeline graph should show up
        cy.get(selectors.eventTimeline.timeline);
    });
});
