import React from 'react';
import { withRouter } from 'react-router-dom';
import ReactRouterPropTypes from 'react-router-prop-types';
import URLService from 'modules/URLService';
import { resourceTypes, standardTypes } from 'constants/entityTypes';

import StandardsByEntity from 'Containers/Compliance2/widgets/StandardsByEntity';
import StandardsAcrossEntity from 'Containers/Compliance2/widgets/StandardsAcrossEntity';
import ComplianceByStandard from 'Containers/Compliance2/widgets/ComplianceByStandard';
import WaveBackground from 'images/wave-bg.svg';
import WaveBackground2 from 'images/wave-bg-2.svg';
import DashboardHeader from './Header';

const ComplianceDashboardPage = ({ match, location }) => {
    const params = URLService.getParams(match, location);

    return (
        <section className="flex flex-col relative">
            <DashboardHeader
                classes="bg-gradient-horizontal z-10 sticky pin-t text-base-100"
                bgStyle={{
                    boxShadow: 'hsla(230, 75%, 63%, 0.62) 0px 5px 30px 0px',
                    '--start': 'hsl(226, 70%, 60%)',
                    '--end': 'hsl(226, 64%, 56%)'
                }}
                params={params}
            />
            <img
                className="absolute pin-l pointer-events-none z-10 w-full"
                src={WaveBackground2}
                style={{ mixBlendMode: 'lighten', top: '-60px' }}
                alt="Waves"
            />
            <div
                className="flex-1 relative bg-gradient-diagonal p-6"
                style={{ '--start': '#F5F2FF', '--end': '#F0F6FF' }}
            >
                <img
                    className="absolute pin-l pointer-events-none w-full"
                    src={WaveBackground}
                    style={{ top: '-130px' }}
                    alt="Wave"
                />
                <div className="grid grid-gap-6 md:grid-auto-fit md:grid-dense">
                    <StandardsAcrossEntity type={resourceTypes.CLUSTER} params={params} />
                    <StandardsByEntity type={resourceTypes.CLUSTER} params={params} />
                    <StandardsAcrossEntity type={resourceTypes.NAMESPACE} params={params} />
                    <StandardsAcrossEntity type={resourceTypes.NODE} params={params} />
                    <ComplianceByStandard type={standardTypes.PCI_DSS_3_2} params={params} />
                    <ComplianceByStandard type={standardTypes.NIST_800_190} params={params} />
                    <ComplianceByStandard type={standardTypes.HIPAA_164} params={params} />
                    <ComplianceByStandard type={standardTypes.CIS_DOCKER_V1_1_0} params={params} />
                    <ComplianceByStandard
                        type={standardTypes.CIS_KUBERENETES_V1_2_0}
                        params={params}
                    />
                </div>
            </div>
        </section>
    );
};

ComplianceDashboardPage.propTypes = {
    match: ReactRouterPropTypes.match.isRequired,
    location: ReactRouterPropTypes.location.isRequired
};

export default withRouter(ComplianceDashboardPage);
