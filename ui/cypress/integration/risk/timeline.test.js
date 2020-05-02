import { selectors, url } from '../../constants/RiskPage';

import * as api from '../../constants/apiEndpoints';

import withAuth from '../../helpers/basicAuth';

function setRoutes() {
    cy.server();
    cy.route('GET', api.risks.riskyDeployments).as('deployments');
    cy.route('GET', api.risks.getDeploymentWithRisk).as('getDeployment');
    cy.route('POST', api.graphql(api.risks.graphqlOps.getDeploymentEventTimeline)).as(
        'getDeploymentEventTimeline'
    );
    cy.route('POST', api.graphql(api.risks.graphqlOps.getPodEventTimeline)).as(
        'getPodEventTimeline'
    );
}

function openDeployment(deploymentName) {
    cy.visit(url);
    cy.wait('@deployments');

    cy.get(`${selectors.table.rows}:contains(${deploymentName})`).click();
    cy.wait('@getDeployment');
}

function openEventTimeline() {
    openDeployment('collector');
    // open the process discovery tab
    cy.get(selectors.sidePanel.processDiscoveryTab).click();
    cy.get(selectors.eventTimelineOverviewButton).click();
}

describe('Risk Page Event Timeline - Filtering Events By Type', () => {
    withAuth();

    const FILTER_OPTIONS = {
        SHOW_ALL: 0,
        POLICY_VIOLATIONS: 1,
        PROCESS_ACTIVITIES: 2,
        RESTARTS: 3,
        TERMINATIONS: 4,
    };

    it('should filter policy violation events', () => {
        setRoutes();
        // mocking data to thoroughly test the filtering
        cy.route(
            'POST',
            api.graphql(api.risks.graphqlOps.getDeploymentEventTimeline),
            'fixture:risks/eventTimeline/deploymentEventTimeline.json'
        ).as('getDeploymentEventTimeline');

        openEventTimeline();

        cy.wait('@getDeploymentEventTimeline');

        // the policy violation event should be visible
        cy.get(selectors.eventTimeline.select.value).should('contain', 'Show All');
        cy.get(selectors.eventTimeline.timeline.mainView.event.policyViolation);

        // filter by something else
        cy.get(selectors.eventTimeline.select.input).click();
        cy.get(
            `${selectors.eventTimeline.select.options}:eq(${FILTER_OPTIONS.PROCESS_ACTIVITIES})`
        ).click({ force: true });

        // the policy violation event should not be visible
        cy.get(selectors.eventTimeline.timeline.mainView.event.policyViolation).should('not.exist');
    });

    it('should filter process activity events and whitelisted process activity events', () => {
        setRoutes();
        // mocking data to thoroughly test the filtering
        cy.route(
            'POST',
            api.graphql(api.risks.graphqlOps.getDeploymentEventTimeline),
            'fixture:risks/eventTimeline/deploymentEventTimeline.json'
        ).as('getDeploymentEventTimeline');

        openEventTimeline();

        cy.wait('@getDeploymentEventTimeline');

        // the process activity + whitelisted process activity event should be visible
        cy.get(selectors.eventTimeline.select.value).should('contain', 'Show All');
        cy.get(selectors.eventTimeline.timeline.mainView.event.processActivity);
        cy.get(selectors.eventTimeline.timeline.mainView.event.whitelistedProcessActivity);

        // filter by something else
        cy.get(selectors.eventTimeline.select.input).click();
        cy.get(
            `${selectors.eventTimeline.select.options}:eq(${FILTER_OPTIONS.POLICY_VIOLATIONS})`
        ).click({ force: true });

        // the process activity + whitelisted process activity event should not be visible
        cy.get(selectors.eventTimeline.timeline.mainView.event.processActivity).should('not.exist');
        cy.get(selectors.eventTimeline.timeline.mainView.event.whitelistedProcessActivity).should(
            'not.exist'
        );
    });

    it('should filter container restart events', () => {
        setRoutes();
        // mocking data to thoroughly test the filtering
        cy.route(
            'POST',
            api.graphql(api.risks.graphqlOps.getDeploymentEventTimeline),
            'fixture:risks/eventTimeline/deploymentEventTimeline.json'
        ).as('getDeploymentEventTimeline');

        openEventTimeline();

        cy.wait('@getDeploymentEventTimeline');

        // the container restart event should be visible
        cy.get(selectors.eventTimeline.select.value).should('contain', 'Show All');
        cy.get(selectors.eventTimeline.timeline.mainView.event.restart);

        // filter by something else
        cy.get(selectors.eventTimeline.select.input).click();
        cy.get(
            `${selectors.eventTimeline.select.options}:eq(${FILTER_OPTIONS.POLICY_VIOLATIONS})`
        ).click({ force: true });

        // thecontainer restart event should not be visible
        cy.get(selectors.eventTimeline.timeline.mainView.event.restart).should('not.exist');
    });

    it('should filter container termination events', () => {
        setRoutes();
        // mocking data to thoroughly test the filtering
        cy.route(
            'POST',
            api.graphql(api.risks.graphqlOps.getDeploymentEventTimeline),
            'fixture:risks/eventTimeline/deploymentEventTimeline.json'
        ).as('getDeploymentEventTimeline');

        openEventTimeline();

        cy.wait('@getDeploymentEventTimeline');

        // the container termination event should be visible
        cy.get(selectors.eventTimeline.select.value).should('contain', 'Show All');
        cy.get(selectors.eventTimeline.timeline.mainView.event.termination);

        // filter by something else
        cy.get(selectors.eventTimeline.select.input).click();
        cy.get(
            `${selectors.eventTimeline.select.options}:eq(${FILTER_OPTIONS.POLICY_VIOLATIONS})`
        ).click({ force: true });

        // the container termination event should not be visible
        cy.get(selectors.eventTimeline.timeline.mainView.event.containerTermination).should(
            'not.exist'
        );
    });
});

describe('Risk Page Event Timeline - Drilling Down To Container Events', () => {
    withAuth();

    it("should drill down on a pod to see that pod's containers", () => {
        setRoutes();
        openEventTimeline();

        cy.wait('@getDeploymentEventTimeline');

        Cypress.Promise.all([
            cy.get(selectors.eventTimeline.timeline.namesList.firstListedName),
            cy.get(selectors.eventTimeline.timeline.mainView.eventsInFirstRow),
        ]).then(([firstListedName, events]) => {
            const firstPodName = firstListedName.text();
            const numEventsInFirstPod = events.length;

            // click the button and drill down to see containers
            cy.get(selectors.eventTimeline.timeline.namesList.drillDownButtonInFirstRow).click();

            cy.wait('@getPodEventTimeline');

            // the back button should be visible
            cy.get(selectors.eventTimeline.backButton);
            // the pod name should be shown in the panel header
            cy.get(selectors.eventTimeline.panelHeader.header).should('contain', firstPodName);
            // we should see the same number of events across containers
            cy.get(selectors.eventTimeline.timeline.mainView.allEvents).should(
                'have.length',
                numEventsInFirstPod
            );
        });
    });
});

describe('Risk Page Event Timeline - Event Details', () => {
    withAuth();

    /**
     * Finds the single event in the mock data for a specified type and returns the timestamp
     * @param {Object[]} json - the mock data in JSON
     * @param {string} type - the event type
     * @param {bool} isWhitelisted - indicates if we're looking for a whitelisted process activity
     * @returns {string} - the timestamp of the event
     */
    function getEventTimeByType(json, type, isWhitelisted = false) {
        return json.data.pods[0].events.filter(
            (event) => event.type === type && isWhitelisted === !!event.whitelisted
        )[0].timestamp;
    }

    /**
     * Finds an event based on the event type and returns the formatted timestamp
     * @param {string} type - the event type
     * @param {bool} isWhitelisted - indicates if we're looking for a whitelisted process activity
     * @returns {Promise<string>} - a promise that, once resolved, will return the formatted timestamp of an event for the specified event typee
     */
    function getTimeByType(type, isWhitelisted = false) {
        return cy.fixture('risks/eventTimeline/deploymentEventTimeline.json').then((json) => {
            const eventTime = getEventTimeByType(json, type, isWhitelisted);
            const formattedEventTime = Cypress.moment(eventTime).format(
                '[Event time:] MM/DD/YYYY | h:mm:ssA'
            );
            return formattedEventTime;
        });
    }

    it('shows the policy violation event details', () => {
        setRoutes();
        // mocking data to thoroughly test the event details
        cy.route(
            'POST',
            api.graphql(api.risks.graphqlOps.getDeploymentEventTimeline),
            'fixture:risks/eventTimeline/deploymentEventTimeline.json'
        ).as('getDeploymentEventTimeline');

        openEventTimeline();

        cy.wait('@getDeploymentEventTimeline');

        // trigger the tooltip
        cy.get(selectors.eventTimeline.timeline.mainView.event.policyViolation).trigger(
            'mouseenter'
        );

        // the header should include the event name
        cy.get(selectors.tooltip.title).contains('Ubuntu Package Manager Execution');
        // the body should include the following
        cy.get(selectors.tooltip.body).contains('Type: Policy Violation');
        // since the displayed time depends on the time zone, we don't want to check against a  hardcoded value
        getTimeByType('PolicyViolationEvent').then((formattedEventTime) => {
            cy.get(selectors.tooltip.body).contains(formattedEventTime);
        });
    });

    it('shows the process activity event details', () => {
        setRoutes();
        // mocking data to thoroughly test the event details
        cy.route(
            'POST',
            api.graphql(api.risks.graphqlOps.getDeploymentEventTimeline),
            'fixture:risks/eventTimeline/deploymentEventTimeline.json'
        ).as('getDeploymentEventTimeline');

        openEventTimeline();

        cy.wait('@getDeploymentEventTimeline');

        // trigger the tooltip
        cy.get(selectors.eventTimeline.timeline.mainView.event.processActivity).trigger(
            'mouseenter'
        );

        // the header should include the event name
        cy.get(selectors.tooltip.title).contains('/bin/bash');
        // the body should include the following
        cy.get(selectors.tooltip.body).contains('Type: Process Activity');
        cy.get(selectors.tooltip.body).contains('Arguments: None');
        cy.get(selectors.tooltip.body).contains('UID: 0');
        // since the displayed time depends on the time zone, we don't want to check against a  hardcoded value
        getTimeByType('ProcessActivityEvent').then((formattedEventTime) => {
            cy.get(selectors.tooltip.body).contains(formattedEventTime);
        });
    });

    it('shows the whitelisted process activity event details', () => {
        setRoutes();
        // mocking data to thoroughly test the event details
        cy.route(
            'POST',
            api.graphql(api.risks.graphqlOps.getDeploymentEventTimeline),
            'fixture:risks/eventTimeline/deploymentEventTimeline.json'
        ).as('getDeploymentEventTimeline');

        openEventTimeline();

        cy.wait('@getDeploymentEventTimeline');

        // trigger the tooltip
        cy.get(selectors.eventTimeline.timeline.mainView.event.whitelistedProcessActivity).trigger(
            'mouseenter'
        );

        // the header should include the event name
        cy.get(selectors.tooltip.title).contains('/usr/sbin/nginx');
        // the body should include the following
        cy.get(selectors.tooltip.body).contains('Type: Process Activity');
        cy.get(selectors.tooltip.body).contains('Arguments: -g daemon off;');
        cy.get(selectors.tooltip.body).contains('UID: 0');
        // since the displayed time depends on the time zone, we don't want to check against a  hardcoded value
        getTimeByType('ProcessActivityEvent', true).then((formattedEventTime) => {
            cy.get(selectors.tooltip.body).contains(formattedEventTime);
        });
    });

    it('shows the container restart event details', () => {
        setRoutes();
        // mocking data to thoroughly test the event details
        cy.route(
            'POST',
            api.graphql(api.risks.graphqlOps.getDeploymentEventTimeline),
            'fixture:risks/eventTimeline/deploymentEventTimeline.json'
        ).as('getDeploymentEventTimeline');

        openEventTimeline();

        cy.wait('@getDeploymentEventTimeline');

        // trigger the tooltip
        cy.get(selectors.eventTimeline.timeline.mainView.event.restart).trigger('mouseenter');

        // the header should include the event name
        cy.get(selectors.tooltip.title).contains('nginx');
        // the body should include the following
        cy.get(selectors.tooltip.body).contains('Type: Container Restart');
        // since the displayed time depends on the time zone, we don't want to check against a  hardcoded value
        getTimeByType('ContainerRestartEvent').then((formattedEventTime) => {
            cy.get(selectors.tooltip.body).contains(formattedEventTime);
        });
    });

    it('shows the container restart event details', () => {
        setRoutes();
        // mocking data to thoroughly test the event details
        cy.route(
            'POST',
            api.graphql(api.risks.graphqlOps.getDeploymentEventTimeline),
            'fixture:risks/eventTimeline/deploymentEventTimeline.json'
        ).as('getDeploymentEventTimeline');

        openEventTimeline();

        cy.wait('@getDeploymentEventTimeline');

        // trigger the tooltip
        cy.get(selectors.eventTimeline.timeline.mainView.event.termination).trigger('mouseenter');

        // the header should include the event name
        cy.get(selectors.tooltip.title).contains('nginx');
        // the body should include the following
        cy.get(selectors.tooltip.body).contains('Type: Container Termination');
        cy.get(selectors.tooltip.body).contains('Reason: OOMKilled');
        // since the displayed time depends on the time zone, we don't want to check against a  hardcoded value
        getTimeByType('ContainerTerminationEvent').then((formattedEventTime) => {
            cy.get(selectors.tooltip.body).contains(formattedEventTime);
        });
    });
});

describe('Risk Page Event Timeline - Pagination', () => {
    withAuth();

    it('should be able to page between sets of pods when there are 10+', () => {
        setRoutes();
        // mocking data to thoroughly test the pagination
        cy.route(
            'POST',
            api.graphql(api.risks.graphqlOps.getDeploymentEventTimeline),
            'fixture:risks/eventTimeline/deploymentEventTimelineForFirstSetOfPods.json'
        ).as('getFirstHalfOfDeploymentEventTimeline');

        openEventTimeline();

        cy.wait('@getFirstHalfOfDeploymentEventTimeline');

        // we should see the first 10 pods out of a total of 15
        cy.get(selectors.eventTimeline.timeline.namesList.listOfNames).should('have.length', 10);

        cy.route(
            'POST',
            api.graphql(api.risks.graphqlOps.getDeploymentEventTimeline),
            'fixture:risks/eventTimeline/deploymentEventTimelineForSecondSetOfPods.json'
        ).as('getSecondHalfOfDeploymentEventTimeline');

        // go to the next page
        cy.get(selectors.eventTimeline.timeline.pagination.nextPage).click({ force: true });
        cy.wait('@getSecondHalfOfDeploymentEventTimeline');

        // we should see the last 5 pods out of the total of 15s
        cy.get(selectors.eventTimeline.timeline.namesList.listOfNames).should('have.length', 5);
    });
});

describe('Risk Page Event Timeline - Legend', () => {
    withAuth();

    it('should show the timeline legend', () => {
        setRoutes();
        openEventTimeline();

        cy.wait('@getDeploymentEventTimeline');

        // show the legend
        cy.get(selectors.eventTimeline.legend).click();

        // make sure the process activity icon and text shows up
        cy.get(`${selectors.tooltip.legendContents}:eq(0):contains("Process Activity")`);
        cy.get(
            `${selectors.tooltip.legendContents}:eq(0) ${selectors.tooltip.legendContent.event.processActivity}`
        );

        // make sure the policy violation icon and text shows up
        cy.get(
            `${selectors.tooltip.legendContents}:eq(1):contains("Process Activity with Violation")`
        );
        cy.get(
            `${selectors.tooltip.legendContents}:eq(1) ${selectors.tooltip.legendContent.event.policyViolation}`
        );

        // make sure the whitelisted process activity icon and text shows up
        cy.get(
            `${selectors.tooltip.legendContents}:eq(2):contains("Whitelisted Process Activity")`
        );
        cy.get(
            `${selectors.tooltip.legendContents}:eq(2) ${selectors.tooltip.legendContent.event.whitelistedProcessActivity}`
        );

        // make sure the container restart icon and text shows up
        cy.get(`${selectors.tooltip.legendContents}:eq(3):contains("Container Restart")`);
        cy.get(
            `${selectors.tooltip.legendContents}:eq(3) ${selectors.tooltip.legendContent.event.restart}`
        );

        // make sure the container termination icon and text shows up
        cy.get(`${selectors.tooltip.legendContents}:eq(4):contains("Container Termination")`);
        cy.get(
            `${selectors.tooltip.legendContents}:eq(4) ${selectors.tooltip.legendContent.event.termination}`
        );
    });
});
