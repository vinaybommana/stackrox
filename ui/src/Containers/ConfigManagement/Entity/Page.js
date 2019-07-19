import React, { useContext } from 'react';
import ReactRouterPropTypes from 'react-router-prop-types';
import { withRouter } from 'react-router-dom';
import URLService from 'modules/URLService';
import { entityQueryKeys } from 'constants/entityTypes';

import SidePanelAnimation from 'Components/animations/SidePanelAnimation';

import pluralize from 'pluralize';
import ExportButton from 'Components/ExportButton';
import searchContext from 'Containers/searchContext';
import searchContexts from 'constants/searchContexts';
import PageHeader from './EntityPageHeader';
import Tabs from './EntityTabs';
import Overview from '../Entity';
import List from '../EntityList';
import SidePanel from '../SidePanel/SidePanel';

const EntityPage = ({ match, location, history }) => {
    const params = URLService.getParams(match, location);
    const {
        pageEntityType,
        pageEntityId,
        entityListType1,
        entityId1,
        entityType2,
        entityListType2,
        entityId2,
        query
    } = params;
    const searchParam = useContext(searchContext);

    function onTabClick({ value }) {
        const urlBuilder = URLService.getURL(match, location).base(pageEntityType, pageEntityId);
        const url = value !== null ? urlBuilder.push(value).url() : urlBuilder.url();
        history.push(url);
    }

    function onRowClick(entityId) {
        const urlBuilder = URLService.getURL(match, location).push(entityId);
        history.push(urlBuilder.url());
    }

    function onRelatedEntityClick(entityType, entityId) {
        const urlBuilder = URLService.getURL(match, location).base(entityType, entityId);
        history.push(urlBuilder.url());
    }

    function onRelatedEntityListClick(entityListType) {
        const urlBuilder = URLService.getURL(match, location).push(entityListType);
        history.push(urlBuilder.url());
    }

    function onSidePanelClose() {
        const urlBuilder = URLService.getURL(match, location).clearSidePanelParams();
        history.replace(urlBuilder.url());
    }
    let component;
    if (!entityListType1) {
        component = (
            <Overview
                entityType={pageEntityType}
                entityId={pageEntityId}
                onRelatedEntityClick={onRelatedEntityClick}
                onRelatedEntityListClick={onRelatedEntityListClick}
            />
        );
    } else {
        const listQuery = {
            [`${entityQueryKeys[pageEntityType]}`]: pageEntityId,
            ...query[searchParam]
        };

        component = (
            <div className="flex flex-1 w-full h-full bg-base-100 relative">
                <List
                    className={entityId1 ? 'overlay' : ''}
                    entityListType={entityListType1}
                    entityId={entityId1}
                    onRowClick={onRowClick}
                    query={listQuery}
                />
                <searchContext.Provider value={searchContexts.sidePanel}>
                    <SidePanelAnimation className="w-3/4" condition={!!entityId1}>
                        <SidePanel
                            className="w-full h-full bg-base-100 border-l-2 border-base-300"
                            entityType1={entityListType1}
                            entityId1={entityId1}
                            entityType2={entityType2}
                            entityListType2={entityListType2}
                            entityId2={entityId2}
                            onClose={onSidePanelClose}
                            query={query}
                        />
                    </SidePanelAnimation>
                </searchContext.Provider>
            </div>
        );
    }

    const exportFilename = `${pluralize(pageEntityType)}`;
    const { urlParams } = URLService.getURL(match, location);
    let pdfId = 'capture-dashboard-stretch';
    if (urlParams.entityListType1) {
        pdfId = 'capture-list';
    }

    return (
        <div className="flex flex-1 flex-col bg-base-200">
            <PageHeader entityType={pageEntityType} entityId={pageEntityId}>
                <div className="flex flex-1 justify-end">
                    <div className="flex">
                        <div className="flex items-center">
                            <ExportButton
                                fileName={exportFilename}
                                type={pageEntityType}
                                page="configManagement"
                                pdfId={pdfId}
                            />
                        </div>
                    </div>
                </div>
            </PageHeader>
            <Tabs
                pageEntityId={pageEntityId}
                entityType={pageEntityType}
                entityListType={entityListType1}
                onClick={onTabClick}
            />
            {component}
        </div>
    );
};

EntityPage.propTypes = {
    match: ReactRouterPropTypes.match.isRequired,
    location: ReactRouterPropTypes.location.isRequired,
    history: ReactRouterPropTypes.history.isRequired
};

export default withRouter(EntityPage);
