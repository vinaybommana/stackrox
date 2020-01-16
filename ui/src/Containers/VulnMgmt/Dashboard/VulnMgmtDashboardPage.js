import React, { useContext } from 'react';
import { withRouter } from 'react-router-dom';
import ReactRouterPropTypes from 'react-router-prop-types';

import entityTypes from 'constants/entityTypes';
import DashboardLayout from 'Components/DashboardLayout';
import ExportButton from 'Components/ExportButton';
import RadioButtonGroup from 'Components/RadioButtonGroup';
import workflowStateContext from 'Containers/workflowStateContext';
import { dashboardLimit } from 'constants/workflowPages.constants';
import TopRiskyEntitiesByVulnerabilities from '../widgets/TopRiskyEntitiesByVulnerabilities';
import TopRiskiestImagesAndComponents from '../widgets/TopRiskiestImagesAndComponents';
import FrequentlyViolatedPolicies from '../widgets/FrequentlyViolatedPolicies';
import RecentlyDetectedVulnerabilities from '../widgets/RecentlyDetectedVulnerabilities';
import MostCommonVulnerabilities from '../widgets/MostCommonVulnerabilities';
import DeploymentsWithMostSeverePolicyViolations from '../widgets/DeploymentsWithMostSeverePolicyViolations';
import ClustersWithMostK8sIstioVulnerabilities from '../widgets/ClustersWithMostK8sIstioVulnerabilities';
import VulnMgmtNavHeader from '../Components/VulnMgmtNavHeader';

// layout-specific graph widget counts

const VulnDashboardPage = ({ history }) => {
    const workflowState = useContext(workflowStateContext);
    const searchState = workflowState.getCurrentSearchState();

    const cveFilterButtons = [
        {
            text: 'Fixable'
        },
        {
            text: 'All'
        }
    ];

    function handleCveFilterToggle(value) {
        const selectedOption = cveFilterButtons.find(button => button.text === value);
        const newValue = selectedOption.text || 'All';

        let targetUrl;
        if (newValue === 'Fixable') {
            targetUrl = workflowState
                .setSearch({
                    IsFixable: 'true'
                })
                .toUrl();
        } else {
            const allSearch = { ...searchState };
            delete allSearch.IsFixable;

            targetUrl = workflowState.setSearch(allSearch).toUrl();
        }

        history.push(targetUrl);
    }

    const cveFilter = searchState.IsFixable ? 'Fixable' : 'All';

    const headerComponents = (
        <>
            <div className="flex items-center">
                <RadioButtonGroup
                    buttons={cveFilterButtons}
                    headerText="Filter CVEs"
                    onClick={handleCveFilterToggle}
                    selected={cveFilter}
                />
                <ExportButton
                    fileName="Vulnerability Management Dashboard Report"
                    page={workflowState.useCase}
                    pdfId="capture-dashboard"
                />
            </div>
            <VulnMgmtNavHeader />
        </>
    );
    return (
        <DashboardLayout headerText="Vulnerability Management" headerComponents={headerComponents}>
            <div className="sx-4 sy-2">
                <TopRiskyEntitiesByVulnerabilities
                    defaultSelection={entityTypes.DEPLOYMENT}
                    cveFilter={cveFilter}
                />
            </div>
            <div className="s-2">
                <TopRiskiestImagesAndComponents limit={dashboardLimit} />
            </div>
            <div className="s-2">
                <FrequentlyViolatedPolicies />
            </div>
            <div className="s-2">
                <RecentlyDetectedVulnerabilities search={searchState} limit={dashboardLimit} />
            </div>
            <div className="sx-2 sy-4">
                <MostCommonVulnerabilities search={searchState} />
            </div>
            <div className="s-2">
                <DeploymentsWithMostSeverePolicyViolations limit={dashboardLimit} />
            </div>
            <div className="s-2">
                <ClustersWithMostK8sIstioVulnerabilities />
            </div>
        </DashboardLayout>
    );
};

VulnDashboardPage.propTypes = {
    history: ReactRouterPropTypes.history.isRequired
};

export default withRouter(VulnDashboardPage);
