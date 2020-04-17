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
            <g transform="matrix(1,0,0,1,-650,-843)">
                <g id="Singles">
                    <g id="iqt-timeline" fill="#FF9064">
                        <path
                            id="Polygon-Copy-5"
                            d="M659.202,844.183L665.814,855.451C666.197,856.105 665.97,856.94 665.306,857.317C665.095,857.437 664.855,857.5 664.612,857.5L651.388,857.5C650.621,857.5 650,856.888 650,856.134C650,855.894 650.064,855.659 650.186,855.451L656.798,844.183C657.181,843.53 658.03,843.306 658.694,843.683C658.905,843.803 659.08,843.976 659.202,844.183Z"
                        />
                    </g>
                </g>
            </g>
        </svg>
    );
});

TerminationIcon.propTypes = {
    width: PropTypes.number.isRequired,
    height: PropTypes.number
};

TerminationIcon.defaultProps = {
    height: null
};

export default TerminationIcon;
