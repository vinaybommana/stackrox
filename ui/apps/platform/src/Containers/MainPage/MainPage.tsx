import React, { ReactElement } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useHistory } from 'react-router-dom';
import { Page, Button } from '@patternfly/react-core';
import { OutlinedCommentsIcon } from '@patternfly/react-icons';
import { gql, useQuery } from '@apollo/client';

import LoadingSection from 'Components/PatternFly/LoadingSection';
import useFeatureFlags from 'hooks/useFeatureFlags';
import usePermissions from 'hooks/usePermissions';
import { selectors } from 'reducers';
import { actions } from 'reducers/feedback';
import { clustersBasePath } from 'routePaths';

import AnnouncementBanner from './Banners/AnnouncementBanner';
import CredentialExpiryBanner from './Banners/CredentialExpiryBanner';
import DatabaseStatusBanner from './Banners/DatabaseStatusBanner';
import OutdatedVersionBanner from './Banners/OutdatedVersionBanner';
import ServerStatusBanner from './Banners/ServerStatusBanner';

import Masthead from './Header/Masthead';

import PublicConfigFooter from './PublicConfig/PublicConfigFooter';
import PublicConfigHeader from './PublicConfig/PublicConfigHeader';

import NavigationSidebar from './Sidebar/NavigationSidebar';

import Body from './Body';
import Notifications from './Notifications';
import AcsFeedbackModal from './AcsFeedbackModal';

type ClusterCountResponse = {
    clusterCount: number;
};

const CLUSTER_COUNT = gql`
    query summary_counts {
        clusterCount
    }
`;

function MainPage(): ReactElement {
    const history = useHistory();
    const dispatch = useDispatch();

    const { isFeatureFlagEnabled, isLoadingFeatureFlags } = useFeatureFlags();
    const { hasReadAccess, hasReadWriteAccess, isLoadingPermissions } = usePermissions();
    const isLoadingPublicConfig = useSelector(selectors.isLoadingPublicConfigSelector);
    const isLoadingCentralCapabilities = useSelector(selectors.getIsLoadingCentralCapabilities);

    // Check for clusters under management
    // if none, and user can admin Clusters, redirect to clusters section
    // (only applicable in Cloud Services version)
    const hasClusterWritePermission = hasReadWriteAccess('Cluster');

    useQuery<ClusterCountResponse>(CLUSTER_COUNT, {
        onCompleted: (data) => {
            if (hasClusterWritePermission && data?.clusterCount < 1) {
                history.push(clustersBasePath);
            }
        },
    });

    // Prerequisites from initial requests for conditional rendering that affects all authenticated routes:
    // feature flags: for NavigationSidebar and Body
    // permissions: for NavigationSidebar and Body
    // public config: for PublicConfigHeader and PublicConfigFooter and analytics
    if (
        isLoadingFeatureFlags ||
        isLoadingPermissions ||
        isLoadingPublicConfig ||
        isLoadingCentralCapabilities
    ) {
        return <LoadingSection message="Loading..." />;
    }

    const hasAdministrationWritePermission = hasReadWriteAccess('Administration');

    return (
        <>
            <Notifications />
            <PublicConfigHeader />
            <AnnouncementBanner />
            <CredentialExpiryBanner
                component="CENTRAL"
                hasAdministrationWritePermission={hasAdministrationWritePermission}
            />
            <CredentialExpiryBanner
                component="SCANNER"
                hasAdministrationWritePermission={hasAdministrationWritePermission}
            />
            <OutdatedVersionBanner />
            <DatabaseStatusBanner />
            <ServerStatusBanner />
            <div id="PageParent">
                <Button
                    style={{
                        bottom: 'calc(var(--pf-global--spacer--lg) * 6)',
                        position: 'absolute',
                        right: '0',
                        transform: 'rotate(270deg)',
                        transformOrigin: 'bottom right',
                        zIndex: 20000,
                    }}
                    icon={<OutlinedCommentsIcon />}
                    iconPosition="left"
                    variant="danger"
                    id="feedback-trigger-button"
                    onClick={() => {
                        dispatch(actions.setFeedbackModalVisibility(true));
                    }}
                >
                    Feedback
                </Button>
                <AcsFeedbackModal />
                <Page
                    mainContainerId="main-page-container"
                    header={<Masthead />}
                    isManagedSidebar
                    sidebar={
                        <NavigationSidebar
                            hasReadAccess={hasReadAccess}
                            isFeatureFlagEnabled={isFeatureFlagEnabled}
                        />
                    }
                >
                    <Body
                        hasReadAccess={hasReadAccess}
                        isFeatureFlagEnabled={isFeatureFlagEnabled}
                    />
                </Page>
            </div>
            <PublicConfigFooter />
        </>
    );
}

export default MainPage;
