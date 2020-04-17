import React, { forwardRef } from 'react';
import PropTypes from 'prop-types';

const PolicyViolationIcon = forwardRef(({ height, width }, ref) => {
    const iconHeight = height || width;
    return (
        <svg
            data-testid="policy-violation-event"
            width={width}
            height={iconHeight}
            viewBox="0 0 15 15"
            version="1.1"
            xmlns="http://www.w3.org/2000/svg"
            ref={ref}
        >
            <g id="Singles" stroke="none" strokeWidth="1" fill="none" fillRule="evenodd">
                <g id="iqt-timeline" transform="translate(-689.000000, -673.000000)" fill="#5677DD">
                    <rect
                        id="Rectangle-Copy-38"
                        x="689"
                        y="673.1"
                        width="14.5799992"
                        height="14.5799992"
                        rx="2.42999987"
                    />
                </g>
            </g>
        </svg>
    );
});

PolicyViolationIcon.propTypes = {
    width: PropTypes.number.isRequired,
    height: PropTypes.number
};

PolicyViolationIcon.defaultProps = {
    height: null
};

export default PolicyViolationIcon;
