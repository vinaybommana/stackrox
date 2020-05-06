import React from 'react';
import PropTypes from 'prop-types';

import Fields from 'Containers/Policies/Wizard/Details/Fields';
import ConfigurationFields from 'Containers/Policies/Wizard/Details/ConfigurationFields';
import FeatureEnabled from 'Containers/FeatureEnabled';
import BooleanPolicySection from 'Containers/Policies/Wizard/Form/BooleanPolicySection';
import { knownBackendFlags } from 'utils/featureFlags';

function PolicyDetails({ initialValues }) {
    if (!initialValues) return null;

    return (
        <div className="w-full h-full">
            <div className="flex flex-col w-full overflow-auto pb-5">
                <Fields policy={initialValues} />
                <ConfigurationFields policy={initialValues} />
                <FeatureEnabled featureFlag={knownBackendFlags.ROX_BOOLEAN_POLICY_LOGIC}>
                    <BooleanPolicySection readOnly initialValues={initialValues} />
                </FeatureEnabled>
            </div>
        </div>
    );
}

PolicyDetails.propTypes = {
    initialValues: PropTypes.shape({
        name: PropTypes.string,
    }).isRequired,
};

export default PolicyDetails;
