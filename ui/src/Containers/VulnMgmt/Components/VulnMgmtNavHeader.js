import React from 'react';

import entityTypes from 'constants/entityTypes';
import EntitiesMenu from 'Components/workflow/EntitiesMenu';
import PoliciesCountTile from './PoliciesCountTile';
import CvesCountTile from './CvesCountTile';

const entityMenuTypes = [
    entityTypes.CLUSTER,
    entityTypes.NAMESPACE,
    entityTypes.DEPLOYMENT,
    entityTypes.IMAGE,
    entityTypes.COMPONENT
];

const VulnMgmtNavHeader = () => (
    <div className="flex h-full ml-3 pl-3 border-l border-base-400">
        <PoliciesCountTile />
        <CvesCountTile />
        <div className="flex w-32">
            <EntitiesMenu text="Application & Infrastructure" options={entityMenuTypes} />
        </div>
    </div>
);

export default VulnMgmtNavHeader;
