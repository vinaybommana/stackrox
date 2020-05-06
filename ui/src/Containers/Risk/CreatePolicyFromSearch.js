import React, { useContext } from 'react';
import PropTypes from 'prop-types';
import ReactRouterPropTypes from 'react-router-prop-types';
import { withRouter } from 'react-router-dom';
import { connect } from 'react-redux';
import { Plus } from 'react-feather';

import wizardStages from 'Containers/Policies/Wizard/wizardStages';
import workflowStateContext from 'Containers/workflowStateContext';
import PanelButton from 'Components/PanelButton';
import { actions as notificationActions } from 'reducers/notifications';
import { actions as pageActions } from 'reducers/policies/page';
import { actions as wizardActions } from 'reducers/policies/wizard';
import { generatePolicyFromSearch } from 'services/PoliciesService';
import searchOptionsToQuery from 'services/searchOptionsToQuery';
import { convertToRestSearch } from './riskPageUtils';

function CreatePolicyFromSearch({
    history,
    openWizard,
    setWizardPolicy,
    setWizardStage,
    addToast,
    removeToast,
}) {
    const workflowState = useContext(workflowStateContext);

    // this utility filters out incomplete search pairs
    const currentSearch = workflowState.getCurrentSearchState();
    const policySearchOptions = convertToRestSearch(currentSearch);

    function createPolicyFromSearch() {
        const queryString = searchOptionsToQuery(policySearchOptions);

        generatePolicyFromSearch(queryString)
            .then((response) => {
                history.push({
                    pathname: `/main/policies`,
                });

                const newPolicy = { ...response?.policy, severity: null };

                setWizardPolicy(newPolicy);
                setWizardStage(wizardStages.edit);
                openWizard();
            })
            .catch((err) => {
                addToast(`Could not create a policy from this search: ${err.message}`);
                setTimeout(removeToast, 5000);
            });
    }

    const isPolicyBtnDisabled = !policySearchOptions?.length;

    return (
        <PanelButton
            icon={<Plus className="h-4 w-4" />}
            className="btn-icon btn-tertiary whitespace-no-wrap h-10 ml-4"
            onClick={createPolicyFromSearch}
            disabled={isPolicyBtnDisabled}
            tooltip="Create Policy from Current Search"
            dataTestId="panel-button-create-policy-from-search"
        >
            Create Policy
        </PanelButton>
    );
}

CreatePolicyFromSearch.propTypes = {
    openWizard: PropTypes.func.isRequired,
    setWizardStage: PropTypes.func.isRequired,
    setWizardPolicy: PropTypes.func.isRequired,
    addToast: PropTypes.func.isRequired,
    removeToast: PropTypes.func.isRequired,
    history: ReactRouterPropTypes.history.isRequired,
};

const mapDispatchToProps = {
    openWizard: pageActions.openWizard,
    setWizardStage: wizardActions.setWizardStage,
    setWizardPolicy: wizardActions.setWizardPolicy,
    addToast: notificationActions.addNotification,
    removeToast: notificationActions.removeOldestNotification,
};

export default withRouter(connect(null, mapDispatchToProps)(CreatePolicyFromSearch));
