import React from 'react';
import PropTypes from 'prop-types';
import ReactRouterPropTypes from 'react-router-prop-types';
import URLService from 'modules/URLService';
import pageTypes from 'constants/pageTypes';

import Panel from 'Components/Panel';
import ComplianceEntityPage from 'Containers/Compliance2/Entity/Page';
import AppLink from 'Components/AppLink';
import { standardBaseTypes } from 'constants/entityTypes';

const ComplianceListSidePanel = ({ match, location, selectedRow, clearSelectedRow }) => {
    const { name } = selectedRow;
    const { context, query, entityType } = URLService.getParams(match, location);

    const pageParams = {
        context,
        pageType: pageTypes.ENTITY,
        entityType,
        entityId: selectedRow.id
    };

    const linkParams = {
        query,
        entityId: selectedRow.id,
        entityType
    };

    const headerTextComponent = (
        <div className="w-full flex items-center">
            <div>
                <AppLink
                    context={context}
                    externalLink
                    pageType={pageTypes.ENTITY}
                    entityType={entityType}
                    params={linkParams}
                    className="w-full flex text-primary-700 hover:text-primary-800 focus:text-primary-700"
                >
                    <div
                        className="flex flex-1 uppercase items-center tracking-wide pl-4 leading-normal font-700"
                        data-test-id="panel-header"
                    >
                        {standardBaseTypes[entityType]
                            ? `${standardBaseTypes[entityType]} ${name}`
                            : name}
                    </div>
                </AppLink>
            </div>
        </div>
    );

    return (
        <Panel
            className="bg-primary-200 z-40 w-full h-full absolute pin-r pin-t md:w-1/2 min-w-72 md:relative"
            headerTextComponent={headerTextComponent}
            onClose={clearSelectedRow}
        >
            <ComplianceEntityPage params={pageParams} sidePanelMode />
        </Panel>
    );
};

ComplianceListSidePanel.propTypes = {
    match: ReactRouterPropTypes.match.isRequired,
    location: ReactRouterPropTypes.location.isRequired,
    selectedRow: PropTypes.shape({}),
    clearSelectedRow: PropTypes.func.isRequired
};

ComplianceListSidePanel.defaultProps = {
    selectedRow: null
};

export default ComplianceListSidePanel;
