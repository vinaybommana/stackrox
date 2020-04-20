import React from 'react';
import PropTypes from 'prop-types';

import { eventTypes } from 'constants/timelineTypes';
import EventTooltip from 'Components/TimelineGraph/EventsGraph/EventTooltip';
import ProcessActivityIcon from './ProcessActivityIcon';
import WhitelistedProcessActivityIcon from './WhitelistedProcessActivityIcon';

const ProcessActivityEvent = ({ name, type, uid, timestamp, whitelisted, width, height }) => {
    const elementHeight = height || width;
    return (
        // We wrap the tooltip within the specific event Components because the Tooltip Component
        // doesn't seem to work when wrapping it around the rendered html one level above. I suspect
        // it doesn't work because the D3Anchor renders a <g> while this renders an svg element
        <EventTooltip name={name} type={type} uid={uid} timestamp={timestamp}>
            {whitelisted ? (
                <WhitelistedProcessActivityIcon height={elementHeight} width={width} />
            ) : (
                <ProcessActivityIcon height={elementHeight} width={width} />
            )}
        </EventTooltip>
    );
};

ProcessActivityEvent.propTypes = {
    name: PropTypes.string.isRequired,
    type: PropTypes.oneOf(Object.values(eventTypes)).isRequired,
    uid: PropTypes.string,
    timestamp: PropTypes.string.isRequired,
    whitelisted: PropTypes.bool,
    width: PropTypes.number.isRequired,
    height: PropTypes.number,
};

ProcessActivityEvent.defaultProps = {
    uid: null,
    whitelisted: false,
    height: null,
};

export default ProcessActivityEvent;
