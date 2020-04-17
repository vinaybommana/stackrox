import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { createStructuredSelector } from 'reselect';
import cloneDeep from 'lodash/cloneDeep';
import { Copy, Download, Edit } from 'react-feather';

import { selectors } from 'reducers';
import { actions as notificationActions } from 'reducers/notifications';
import { actions as wizardActions } from 'reducers/policies/wizard';
import wizardStages from 'Containers/Policies/Wizard/wizardStages';
import PanelButton from 'Components/PanelButton';
import { exportPolicies } from 'services/PoliciesService';
import { knownBackendFlags } from 'utils/featureFlags';
import FeatureEnabled from 'Containers/FeatureEnabled';

function DetailsButtons({ wizardPolicy, setWizardStage, setWizardPolicy, addToast, removeToast }) {
    function goToEdit() {
        setWizardStage(wizardStages.edit);
    }

    function exportOnePolicy() {
        const policiesToExport = [wizardPolicy.id];
        exportPolicies(policiesToExport).catch(err => {
            addToast(`Could not delete all of the selected policies: ${err.message}`);
            setTimeout(removeToast, 5000);
        });
    }

    function onPolicyClone() {
        const newPolicy = cloneDeep(wizardPolicy);
        newPolicy.id = '';
        newPolicy.name += ' (COPY)';
        setWizardPolicy(newPolicy);
        setWizardStage(wizardStages.edit);
    }

    return (
        <React.Fragment>
            <PanelButton
                icon={<Copy className="h-4 w-4" />}
                className="btn btn-base mr-2"
                onClick={onPolicyClone}
                tooltip="Clone policy"
            >
                Clone
            </PanelButton>
            <FeatureEnabled featureFlag={knownBackendFlags.ROX_POLICY_IMPORT_EXPORT}>
                <PanelButton
                    icon={<Download className="h-4 w-4" />}
                    className="btn btn-base mr-2"
                    onClick={exportOnePolicy}
                    tooltip="Export policy"
                >
                    Export
                </PanelButton>
            </FeatureEnabled>
            <PanelButton
                icon={<Edit className="h-4 w-4" />}
                className="btn btn-base mr-2"
                onClick={goToEdit}
                tooltip="Edit policy"
            >
                Edit
            </PanelButton>
        </React.Fragment>
    );
}

DetailsButtons.propTypes = {
    wizardPolicy: PropTypes.shape({
        id: PropTypes.string
    }).isRequired,
    setWizardStage: PropTypes.func.isRequired,
    setWizardPolicy: PropTypes.func.isRequired,
    addToast: PropTypes.func.isRequired,
    removeToast: PropTypes.func.isRequired
};

const mapStateToProps = createStructuredSelector({
    wizardPolicy: selectors.getWizardPolicy
});

const mapDispatchToProps = {
    setWizardStage: wizardActions.setWizardStage,
    setWizardPolicy: wizardActions.setWizardPolicy,
    addToast: notificationActions.addNotification,
    removeToast: notificationActions.removeOldestNotification
};

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(DetailsButtons);
