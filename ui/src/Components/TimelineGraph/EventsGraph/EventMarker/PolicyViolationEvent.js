/* eslint-disable react/display-name */
import React, { forwardRef } from 'react';
import PropTypes from 'prop-types';

const PolicyViolationEvent = forwardRef(({ size }, ref) => {
    return (
        <svg
            data-testid="policy-violation-event"
            width={size}
            height={size}
            viewBox="0 0 15 15"
            version="1.1"
            xmlns="http://www.w3.org/2000/svg"
            ref={ref}
        >
            <g id="Page-1" stroke="none" strokeWidth="1" fill="none" fillRule="evenodd">
                <g id="iqt-timeline-popover" transform="translate(-908.000000, -673.000000)">
                    <g id="Group-73-Copy-3" transform="translate(908.000000, 673.000000)">
                        <g id="Group-71-Copy" transform="translate(0.000000, 0.100000)">
                            <rect
                                id="Rectangle-Copy-38"
                                fill="#FF5782"
                                x="0"
                                y="0"
                                width="14.5799992"
                                height="14.5799992"
                                rx="2.42999987"
                            />
                            <path
                                d="M8.07143356,8.37042581 L6.61628124,8.37042581 L6.29069921,2.78569543 L8.41030465,2.78569543 L8.07143356,8.37042581 Z M6.23754296,10.5896992 C6.23754296,10.2441819 6.32945805,9.97618846 6.513291,9.78571095 C6.69712395,9.59523344 6.97065441,9.49999611 7.3338906,9.49999611 C7.69712679,9.49999611 7.97176466,9.59412603 8.15781246,9.78238868 C8.34386027,9.97065134 8.43688277,10.2397522 8.43688277,10.5896992 C8.43688277,10.9352166 8.34053804,11.20321 8.14784567,11.3936875 C7.9551533,11.584165 7.68383766,11.6794023 7.3338906,11.6794023 C6.97508412,11.6794023 6.70266107,11.5830576 6.51661327,11.3903652 C6.33056547,11.1976729 6.23754296,10.9307869 6.23754296,10.5896992 Z"
                                id="!"
                                fill="#FFFFFF"
                                fillRule="nonzero"
                            />
                        </g>
                    </g>
                </g>
            </g>
        </svg>
    );
});

PolicyViolationEvent.propTypes = {
    size: PropTypes.number.isRequired,
};

export default PolicyViolationEvent;
