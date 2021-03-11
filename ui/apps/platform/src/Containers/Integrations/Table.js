/* eslint-disable react/jsx-no-bind */

import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { selectors } from 'reducers';
import { createStructuredSelector } from 'reselect';

import CheckboxTable from 'Components/CheckboxTable';
import { rtTrActionsClassName } from 'Components/Table';
import { PanelNew, PanelBody, PanelHead, PanelHeadEnd, PanelTitle } from 'Components/Panel';
import PanelButton from 'Components/PanelButton';
import RowActionButton from 'Components/RowActionButton';
import * as Icon from 'react-feather';

import tableColumnDescriptor from 'Containers/Integrations/tableColumnDescriptor';
import NoResultsMessage from 'Components/NoResultsMessage';

class IntegrationTable extends Component {
    static propTypes = {
        integrations: PropTypes.arrayOf(PropTypes.object).isRequired,

        source: PropTypes.oneOf([
            'authPlugins',
            'backups',
            'imageIntegrations',
            'notifiers',
            'authProviders',
            'logIntegrations',
        ]).isRequired,
        type: PropTypes.string.isRequired,
        label: PropTypes.string.isRequired,

        buttonsEnabled: PropTypes.bool.isRequired,
        showAutogenerated: PropTypes.bool.isRequired,

        onAutogeneratedClick: PropTypes.func.isRequired,
        onRowClick: PropTypes.func.isRequired,
        onActivate: PropTypes.func.isRequired,
        onAdd: PropTypes.func.isRequired,
        onDelete: PropTypes.func.isRequired,

        setTable: PropTypes.func.isRequired,
        selectedIntegrationId: PropTypes.string,
        toggleRow: PropTypes.func.isRequired,
        toggleSelectAll: PropTypes.func.isRequired,
        selection: PropTypes.arrayOf(PropTypes.string).isRequired,
    };

    static defaultProps = {
        selectedIntegrationId: null,
    };

    onDeleteHandler = (integration) => (e) => {
        e.stopPropagation();
        this.props.onDelete(integration);
    };

    onActivateHandler = (integration) => (e) => {
        e.stopPropagation();
        this.props.onActivate(integration);
    };

    getPanelButtons = () => {
        const {
            selection,
            onDelete,
            onAutogeneratedClick,
            integrations,
            buttonsEnabled,
            onAdd,
        } = this.props;
        const selectionCount = selection.length;
        const integrationsCount = integrations.length;
        const autogeneratedCount = integrations.filter((x) => x.autogenerated).length;
        const showHideText = `${
            this.props.showAutogenerated ? 'Hide' : 'Show'
        } Autogenerated (${autogeneratedCount})`;

        const addHandler =
            this.props.type === 'awsSecurityHub' && integrationsCount > 0
                ? () => onAdd({ onlyOneIntegrationAllowed: true })
                : onAdd;

        return (
            <>
                {this.props.type === 'docker' && (
                    <PanelButton
                        icon={
                            this.props.showAutogenerated ? (
                                <Icon.EyeOff className="h-4 w-4 ml-1" />
                            ) : (
                                <Icon.Eye className="h-4 w-4 ml-1" />
                            )
                        }
                        className="btn btn-base mr-2 w-56"
                        onClick={onAutogeneratedClick}
                        tooltip={showHideText}
                    >
                        {showHideText}
                    </PanelButton>
                )}
                {selectionCount !== 0 && (
                    <PanelButton
                        icon={<Icon.Trash2 className="h-4 w-4 ml-1" />}
                        className="btn btn-alert mr-3"
                        onClick={onDelete}
                        disabled={integrationsCount === 0 || !buttonsEnabled}
                        tooltip={`Delete (${selectionCount})`}
                    >
                        {`Delete (${selectionCount})`}
                    </PanelButton>
                )}
                {selectionCount === 0 && (
                    <PanelButton
                        icon={<Icon.Plus className="h-4 w-4 ml-1" />}
                        className="btn btn-base mr-3"
                        onClick={addHandler}
                        disabled={!buttonsEnabled}
                        tooltip="New Integration"
                    >
                        New Integration
                    </PanelButton>
                )}
            </>
        );
    };

    getColumns = () => {
        const { source, type } = this.props;
        const columns = [...tableColumnDescriptor[source][type]];
        columns.push({
            Header: '',
            accessor: '',
            headerClassName: 'hidden',
            className: rtTrActionsClassName,
            Cell: ({ original }) => this.renderRowActionButtons(original),
        });
        return columns;
    };

    renderRowActionButtons = (integration) => {
        const { source } = this.props;
        let activateBtn = null;
        if (source === 'authProviders') {
            const enableTooltip = `${!integration.validated ? 'Enable' : 'Disable'} auth provider`;
            const enableIconColor = integration.disabled ? 'text-primary-600' : 'text-success-600';
            const enableIconHoverColor = integration.disabled
                ? 'text-primary-700'
                : 'text-success-700';
            activateBtn = (
                <RowActionButton
                    text={enableTooltip}
                    onClick={this.onActivateHandler(integration)}
                    className={`hover:bg-primary-200 ${enableIconColor} hover:${enableIconHoverColor}`}
                    icon={<Icon.Power className="my-1 h-4 w-4" />}
                />
            );
        }
        return (
            <div className="border-2 border-r-2 border-base-400 bg-base-100 flex">
                {activateBtn}
                <RowActionButton
                    text="Delete integration"
                    onClick={this.onDeleteHandler(integration)}
                    border={`${source === 'authProviders' ? 'border-l-2 border-base-400' : ''}`}
                    icon={<Icon.Trash2 className="my-1 h-4 w-4" />}
                />
            </div>
        );
    };

    renderTableContent = () => {
        let rows = this.props.integrations;

        let { label } = this.props;
        if (label === undefined) {
            label = this.props.type;
        }

        let defaultSorted = [];
        if (this.props.type === 'docker') {
            defaultSorted = [{ id: 'autogenerated' }];
            if (!this.props.showAutogenerated) {
                rows = rows.filter((x) => !x.autogenerated);
            }
        }

        if (!rows.length) {
            return <NoResultsMessage message={`No ${label} integrations`} />;
        }
        return (
            <CheckboxTable
                ref={this.props.setTable}
                rows={rows}
                columns={this.getColumns()}
                onRowClick={this.props.onRowClick}
                toggleRow={this.props.toggleRow}
                toggleSelectAll={this.props.toggleSelectAll}
                selection={this.props.selection}
                selectedRowId={this.props.selectedIntegrationId}
                noDataText={`No ${label} integrations`}
                minRows={20}
                defaultSorted={defaultSorted}
            />
        );
    };

    render() {
        const { type, selection, integrations } = this.props;
        let { label } = this.props;
        if (label === undefined) {
            label = type;
        }
        const selectionCount = selection.length;
        const integrationsCount = integrations.length;
        const headerText =
            selectionCount !== 0
                ? `${selectionCount} ${label} Integration${
                      selectionCount === 1 ? '' : 's'
                  } selected`
                : `${integrationsCount} ${label} Integration${integrationsCount === 1 ? '' : 's'}`;

        return (
            <div className="bg-base-100 flex-shrink-1 overflow-hidden w-full">
                <PanelNew testid="panel">
                    <PanelHead>
                        <PanelTitle isUpperCase testid="panel-header" text={headerText} />
                        <PanelHeadEnd>{this.getPanelButtons()}</PanelHeadEnd>
                    </PanelHead>
                    <PanelBody>{this.renderTableContent()}</PanelBody>
                </PanelNew>
            </div>
        );
    }
}

const mapStateToProps = createStructuredSelector({
    clusters: selectors.getClusters,
});

export default connect(mapStateToProps)(IntegrationTable);
