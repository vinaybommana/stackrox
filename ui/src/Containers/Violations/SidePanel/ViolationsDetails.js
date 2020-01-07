import React from 'react';
import PropTypes from 'prop-types';

import { knownBackendFlags } from 'utils/featureFlags';

import FeatureEnabled from 'Containers/FeatureEnabled';
import ViolationComments from './ViolationComments';
import DeploytimeMessages from './DeploytimeMessages';
import RuntimeMessages from './RuntimeMessages';

function ViolationsDetails({ violations, processViolation }) {
    return (
        <div className="w-full px-3 pb-5 mt-5">
            <div className="mb-4">
                <FeatureEnabled featureFlag={knownBackendFlags.ROX_IQT_ANALYST_NOTES_UI}>
                    <ViolationComments />
                </FeatureEnabled>
            </div>
            <RuntimeMessages processViolation={processViolation} />
            <DeploytimeMessages violations={violations} />
        </div>
    );
}

ViolationsDetails.propTypes = {
    violations: PropTypes.arrayOf(
        PropTypes.shape({
            message: PropTypes.string.isRequired
        })
    ),
    processViolation: PropTypes.shape({
        message: PropTypes.string.isRequired,
        processes: PropTypes.array.isRequired
    })
};

ViolationsDetails.defaultProps = {
    violations: [],
    processViolation: null
};

export default ViolationsDetails;
