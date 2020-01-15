import React from 'react';
import PropTypes from 'prop-types';
import Tooltip from 'rc-tooltip';
import 'rc-tooltip/assets/bootstrap.css';

const PanelButton = ({
    children,
    icon,
    text,
    onClick,
    className,
    disabled,
    tooltip,
    alwaysVisibleText
}) => {
    const tooltipContent = tooltip || text;
    const tooltipClassName = !tooltip ? 'visible xl:invisible' : '';
    return (
        <Tooltip
            placement="top"
            mouseLeaveDelay={0}
            overlay={<div>{tooltipContent}</div>}
            overlayClassName={tooltipClassName}
        >
            <button
                type="button"
                className={className}
                onClick={onClick}
                disabled={disabled}
                data-test-id={`${text.toLowerCase()}-button`}
            >
                {icon && <span className="flex items-center">{icon}</span>}
                {children && (
                    <span className={`mx-2 items-center ${alwaysVisibleText ? 'flex' : 'hidden'}`}>
                        {children}
                    </span>
                )}
            </button>
        </Tooltip>
    );
};

PanelButton.propTypes = {
    children: PropTypes.node,
    icon: PropTypes.node,
    text: PropTypes.string,
    onClick: PropTypes.func.isRequired,
    className: PropTypes.string,
    disabled: PropTypes.bool,
    tooltip: PropTypes.node,
    alwaysVisibleText: PropTypes.bool
};

PanelButton.defaultProps = {
    children: null,
    icon: null,
    text: '',
    className: '',
    disabled: false,
    tooltip: '',
    alwaysVisibleText: false
};

export default PanelButton;
