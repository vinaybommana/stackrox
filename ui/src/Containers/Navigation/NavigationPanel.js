import React, { Component } from 'react';
import { NavLink as Link, withRouter } from 'react-router-dom';
import PropTypes from 'prop-types';
import { createStructuredSelector } from 'reselect';
import { connect } from 'react-redux';
import downloadDiagnostics from 'services/DebugService';
import { isBackendFeatureFlagEnabled, knownBackendFlags } from 'utils/featureFlags';

import { selectors } from 'reducers';
import {
    clustersPath,
    policiesListPath,
    integrationsPath,
    accessControlPath,
    systemConfigPath
} from 'routePaths';
import * as Icon from 'react-feather';
import PanelButton from 'Components/PanelButton';
import { filterLinksByFeatureFlag } from './navHelpers';

const navLinks = [
    {
        text: 'Clusters',
        to: clustersPath.replace('/:clusterId', '')
    },
    {
        text: 'System Policies',
        to: policiesListPath
    },
    {
        text: 'Integrations',
        to: integrationsPath
    },
    {
        text: 'Access Control',
        to: accessControlPath
    },
    {
        text: 'System Configuration',
        to: systemConfigPath,
        data: 'system-config'
    }
];

class NavigationPanel extends Component {
    static propTypes = {
        panelType: PropTypes.string.isRequired,
        onClose: PropTypes.func.isRequired,
        featureFlags: PropTypes.arrayOf(
            PropTypes.shape({
                envVar: PropTypes.string.isRequired,
                enabled: PropTypes.bool.isRequired
            })
        ).isRequired
    };

    constructor(props) {
        super(props);
        this.panels = {
            configure: this.renderConfigurePanel
        };
    }

    renderConfigurePanel = () => (
        <ul className="flex flex-col overflow-auto list-reset uppercase tracking-wide bg-primary-800 border-r border-l border-primary-900">
            <li className="border-b-2 border-primary-500 px-1 py-5 pl-2 pr-2 text-base-100 font-700">
                Configure StackRox Settings
            </li>
            {filterLinksByFeatureFlag(this.props.featureFlags, navLinks).map(navLink => (
                <li key={navLink.text} className="text-sm">
                    <Link
                        to={navLink.to}
                        onClick={this.props.onClose(true, 'configure')}
                        className="block no-underline text-base-100 px-1 font-700 border-b py-5 border-primary-900 pl-2 pr-2 hover:bg-base-700"
                        data-test-id={navLink.data || navLink.text}
                    >
                        {navLink.text}
                    </Link>
                </li>
            ))}
            {isBackendFeatureFlagEnabled(
                this.props.featureFlags,
                knownBackendFlags.ROX_DIAGNOSTIC_BUNDLE,
                false
            ) && (
                <li key="Download Diagnostic Data" className="text-sm border-b border-primary-900">
                    <PanelButton
                        icon={<Icon.Download className="h-4 w-4 ml-1 text-primary-400" />}
                        className="flex leading-normal font-700 text-sm text-base-100 no-underline py-5 px-1 items-center uppercase"
                        onClick={this.downloadDiagnostics}
                        alwaysVisibleText
                        tooltip={
                            <div className="w-auto">
                                <h2 className="mb-2 font-700 text-lg uppercase">
                                    What we collect:
                                </h2>
                                <p className="mb-2">
                                    The diagnostic bundle contains information pertaining to the
                                    system health of the StackRox deployments
                                    <br /> in the central cluster as well as all currently connected
                                    secured clusters.
                                </p>
                                <p className="mb-1">It includes:</p>
                                <ul className="mb-1 w-full list-disc">
                                    <li>Heap profile of Central</li>
                                    <li>
                                        Database storage information (database size, free space on
                                        volume)
                                    </li>
                                    <li>
                                        Component health information for all StackRox components
                                        (version, memory usage, error conditions)
                                    </li>
                                    <li>
                                        Coarse-grained usage statistics (API endpoint invocation
                                        counts)
                                    </li>
                                    <li>
                                        Logs of all StackRox components from the last 20 minutes
                                    </li>
                                    <li>
                                        Logs of recently crashed StackRox components from up to 20
                                        minutes before the last crash
                                    </li>
                                    <li>
                                        Kubernetes YAML definitions of StackRox components
                                        (excluding Kubernetes secrets)
                                    </li>
                                    <li>Kubernetes events of objects in the StackRox namespaces</li>
                                    <li>
                                        Information about nodes in each secured cluster (kernel and
                                        OS versions, resource pressure, taints)
                                    </li>
                                    <li>
                                        Environment information about each secured cluster
                                        (Kubernetes version, if applicable cloud provider)
                                    </li>
                                </ul>
                            </div>
                        }
                    >
                        Download Diagnostic Data{' '}
                        <Icon.HelpCircle className="h-4 w-4 text-primary-400 ml-2" />
                    </PanelButton>
                </li>
            )}
        </ul>
    );

    downloadDiagnostics = () => {
        downloadDiagnostics();
    };

    render() {
        return (
            <div
                className="navigation-panel w-full flex theme-light"
                data-test-id="configure-subnav"
            >
                {this.panels[this.props.panelType]()}
                <button
                    aria-label="Close Configure sub-navigation menu"
                    type="button"
                    className="flex-1 opacity-50 bg-primary-700"
                    onClick={this.props.onClose(true)}
                />
            </div>
        );
    }
}

const mapStateToProps = createStructuredSelector({
    featureFlags: selectors.getFeatureFlags
});

export default withRouter(connect(mapStateToProps)(NavigationPanel));
