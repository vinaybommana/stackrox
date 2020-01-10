import React, { Component } from 'react';
import PropTypes from 'prop-types';
import Collapsible from 'react-collapsible';
import * as Icon from 'react-feather';

class CollapsibleCard extends Component {
    static propTypes = {
        title: PropTypes.string.isRequired,
        children: PropTypes.node.isRequired,
        open: PropTypes.bool,
        titleClassName: PropTypes.string,
        renderWhenOpened: PropTypes.func,
        renderWhenClosed: PropTypes.func,
        cardClassName: PropTypes.string,
        headerComponents: PropTypes.element
    };

    static defaultProps = {
        open: true,
        titleClassName:
            'border-b border-base-300 leading-normal cursor-pointer flex justify-end items-center hover:bg-primary-100 hover:border-primary-300',
        renderWhenOpened: null,
        renderWhenClosed: null,
        cardClassName: 'border border-base-400',
        headerComponents: null
    };

    renderTriggerElement = cardState => {
        const icons = {
            opened: <Icon.ChevronUp className="h-4 w-4" />,
            closed: <Icon.ChevronDown className="h-4 w-4" />
        };
        const { title, titleClassName, headerComponents } = this.props;
        return (
            <div className={titleClassName}>
                <h1 className="flex flex-1 p-3 pb-2 text-base-600 font-700 text-lg">{title}</h1>
                {headerComponents && <div>{headerComponents}</div>}
                <div className="flex px-3">{icons[cardState]}</div>
            </div>
        );
    };

    renderWhenOpened = () => this.renderTriggerElement('opened');

    renderWhenClosed = () => this.renderTriggerElement('closed');

    render() {
        const renderWhenOpened = this.props.renderWhenOpened
            ? this.props.renderWhenOpened
            : this.renderWhenOpened;
        const renderWhenClosed = this.props.renderWhenClosed
            ? this.props.renderWhenClosed
            : this.renderWhenClosed;
        return (
            <div className={`bg-base-100 text-base-600 rounded ${this.props.cardClassName}`}>
                <Collapsible
                    open={this.props.open}
                    trigger={renderWhenClosed()}
                    triggerWhenOpen={renderWhenOpened()}
                    transitionTime={100}
                    lazyRender
                >
                    {this.props.children}
                </Collapsible>
            </div>
        );
    }
}

export default CollapsibleCard;
