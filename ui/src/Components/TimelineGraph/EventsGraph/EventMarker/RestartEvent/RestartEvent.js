import React from 'react';
import PropTypes from 'prop-types';

import { eventTypes } from 'constants/timelineTypes';
import EventTooltip from 'Components/TimelineGraph/EventsGraph/EventTooltip';
import RestartIcon from './RestartIcon';

const RestartEvent = ({ name, type, timestamp, width, height }) => {
    const elementHeight = height || width;
    return (
        // We wrap the tooltip within the specific event Components because the Tooltip Component
        // doesn't seem to work when wrapping it around the rendered html one level above. I suspect
        // it doesn't work because the D3Anchor renders a <g> while this renders an svg element
        <EventTooltip name={name} type={type} timestamp={timestamp}>
            <RestartIcon height={elementHeight} width={width} />
        </EventTooltip>
    );
};

RestartEvent.propTypes = {
    name: PropTypes.string.isRequired,
    type: PropTypes.oneOf(Object.values(eventTypes)).isRequired,
    timestamp: PropTypes.string.isRequired,
    width: PropTypes.number.isRequired,
    height: PropTypes.number,
};

RestartEvent.defaultProps = {
    height: null,
};

export default RestartEvent;
