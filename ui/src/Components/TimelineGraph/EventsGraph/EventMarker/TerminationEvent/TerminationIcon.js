import React, { forwardRef } from 'react';
import PropTypes from 'prop-types';

const TerminationIcon = forwardRef(({ height, width }, ref) => {
    const iconHeight = height || width;
    return (
        <svg
            data-testid="termination-event"
            width={width}
            height={iconHeight}
            viewBox="0 0 16 16"
            version="1.1"
            xmlns="http://www.w3.org/2000/svg"
            ref={ref}
        >
            <g id="Singles" stroke="none" strokeWidth="1" fill="none" fillRule="evenodd">
                <g id="iqt-timeline" transform="translate(-650.000000, -843.000000)" fill="#FF9064">
                    <path
                        d="M659.202024,844.183144 L665.813796,855.451253 C666.197113,856.10452 665.969636,856.939849 665.305712,857.317013 C665.094692,857.43689 664.855321,857.5 664.611656,857.5 L651.388112,857.5 C650.621479,857.5 650,856.888496 650,856.134169 C650,855.894415 650.064139,855.658885 650.185972,855.451253 L656.797744,844.183144 C657.18106,843.529877 658.030016,843.306051 658.69394,843.683215 C658.90496,843.803092 659.080192,843.975511 659.202024,844.183144 Z"
                        id="Polygon-Copy-5"
                        transform="translate(658.000000, 850.500000) scale(1, -1) translate(-658.000000, -850.500000) "
                    />
                </g>
            </g>
        </svg>
    );
});

TerminationIcon.propTypes = {
    width: PropTypes.number.isRequired,
    height: PropTypes.number,
};

TerminationIcon.defaultProps = {
    height: null,
};

export default TerminationIcon;
