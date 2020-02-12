import React, { useState } from 'react';
import pluralize from 'pluralize';
import startCase from 'lodash/startCase';

import { searchParams, sortParams, pagingParams } from 'constants/searchParams';
import PageHeader from 'Components/PageHeader';
import ExportButton from 'Components/ExportButton';
import entityLabels from 'messages/entity';
import useCaseLabels from 'messages/useCase';
import getSidePanelEntity from 'utils/getSidePanelEntity';
import parseURL from 'modules/URLParser';
import workflowStateContext from 'Containers/workflowStateContext';
import { WorkflowState } from 'modules/WorkflowState';
import { exportCvesAsCsv } from 'services/VulnerabilitiesService';

import WorkflowSidePanel from './WorkflowSidePanel';
import {
    EntityComponentMap,
    ListComponentMap,
    NavHeaderComponentMap
} from './UseCaseComponentMaps';

const WorkflowListPageLayout = ({ location }) => {
    const workflowState = parseURL(location);
    const { useCase, search, sort, paging, stateStack } = workflowState;
    const pageState = new WorkflowState(
        useCase,
        workflowState.getPageStack(),
        search,
        sort,
        paging
    );

    // set up cache-busting system that either the list or sidepanel can use to trigger list refresh
    const [refreshTrigger, setRefreshTrigger] = useState(0);

    function customCsvExportHandler(fileName) {
        return exportCvesAsCsv(fileName, workflowState);
    }

    // Get the list / entity / nav header components
    const ListComponent = ListComponentMap[useCase];
    const EntityComponent = EntityComponentMap[useCase];
    const NavHeaderComponent = NavHeaderComponentMap[useCase];

    // Page props
    const pageListType = workflowState.getBaseEntity().entityType;
    const pageSearch = workflowState.search[searchParams.page];
    const pageSort = workflowState.sort[sortParams.page];
    const pagePaging = workflowState.paging[pagingParams.page];

    // Sidepanel props
    const { sidePanelEntityId, sidePanelEntityType, sidePanelListType } = getSidePanelEntity(
        workflowState
    );
    const sidePanelSearch = workflowState.search[searchParams.sidePanel];
    const sidePanelSort = workflowState.sort[sortParams.sidePanel];
    const sidePanelPaging = workflowState.paging[pagingParams.sidePanel];
    const selectedRow = workflowState.getSelectedTableRow();

    const header = pluralize(entityLabels[pageListType]);
    const exportFilename = `${useCaseLabels[useCase]} ${pluralize(startCase(header))} Report`;
    const entityContext = {};

    if (selectedRow && stateStack.length > 2) {
        const { entityType, entityId } = selectedRow;
        entityContext[entityType] = entityId;
    }

    return (
        <workflowStateContext.Provider value={pageState}>
            <div className="flex flex-col relative min-h-full">
                <PageHeader header={header} subHeader="Entity List" classes="pr-0">
                    <div className="flex flex-1 justify-end h-10 pr-2">
                        <div className="flex items-center">
                            <ExportButton
                                fileName={exportFilename}
                                type={pageListType}
                                page={useCase}
                                disabled={!!sidePanelEntityId}
                                pdfId="capture-list"
                                customCsvExportHandler={customCsvExportHandler}
                            />
                        </div>
                        <NavHeaderComponent />
                    </div>
                </PageHeader>
                <div className="h-full bg-base-100 relative z-0 min-h-0" id="capture-list">
                    <ListComponent
                        entityListType={pageListType}
                        selectedRowId={selectedRow && selectedRow.entityId}
                        search={pageSearch}
                        sort={pageSort}
                        page={pagePaging}
                        refreshTrigger={refreshTrigger}
                        setRefreshTrigger={setRefreshTrigger}
                    />
                    <WorkflowSidePanel isOpen={!!sidePanelEntityId}>
                        {sidePanelEntityId ? (
                            <EntityComponent
                                entityId={sidePanelEntityId}
                                entityType={sidePanelEntityType}
                                entityListType={sidePanelListType}
                                search={sidePanelSearch}
                                sort={sidePanelSort}
                                page={sidePanelPaging}
                                entityContext={entityContext}
                                setRefreshTrigger={setRefreshTrigger}
                            />
                        ) : (
                            <span />
                        )}
                    </WorkflowSidePanel>
                </div>
            </div>
        </workflowStateContext.Provider>
    );
};

export default WorkflowListPageLayout;
