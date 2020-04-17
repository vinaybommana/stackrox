import React, { forwardRef } from 'react';
import PropTypes from 'prop-types';

const WhitelistedProcessActivityIcon = forwardRef(({ height, width }, ref) => {
    const iconHeight = height || width;
    return (
        <svg
            data-testid="whitelisted-process-activity-event"
            width={width}
            height={iconHeight}
            viewBox="0 0 15 15"
            version="1.1"
            xmlns="http://www.w3.org/2000/svg"
            ref={ref}
        >
            <g id="Page-1" stroke="none" strokeWidth="1" fill="none" fillRule="evenodd">
                <g id="iqt-timeline-popover" transform="translate(-689.000000, -673.000000)">
                    <g id="Group-70" transform="translate(689.000000, 673.100000)">
                        <rect
                            id="Rectangle-Copy-38"
                            fill="#56DDB2"
                            x="2.14939178e-13"
                            y="-1.24344979e-14"
                            width="14.5799992"
                            height="14.5799992"
                            rx="2.42999987"
                        />
                        <path
                            d="M4.45866889,6.7675062 C4.1519101,6.40536439 3.62404701,6.37328328 3.27965367,6.69585103 C2.93526034,7.01841878 2.90475151,7.57348555 3.2115103,7.93562737 L5.89698061,11.1059394 C6.24340685,11.5149103 6.85744163,11.4942771 7.17832948,11.0628827 L11.8226518,4.8191603 C12.1067965,4.43716259 12.0426478,3.88527656 11.6793717,3.58648835 C11.3160956,3.28770014 10.7912573,3.35515476 10.5071127,3.73715248 L6.4789459,9.1525292 L4.45866889,6.7675062 Z"
                            id="Path-2"
                            fill="#FFFFFF"
                            fillRule="nonzero"
                        />
                    </g>
                </g>
            </g>
        </svg>
    );
});

WhitelistedProcessActivityIcon.propTypes = {
    width: PropTypes.number.isRequired,
    height: PropTypes.number
};

WhitelistedProcessActivityIcon.defaultProps = {
    height: null
};

export default WhitelistedProcessActivityIcon;
