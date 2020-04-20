import React from 'react';
import PropTypes from 'prop-types';

import EventMarker from './EventMarker';

const EventsRow = ({
    name,
    events,
    isOdd,
    height,
    width,
    translateX,
    translateY,
    minTimeRange,
    maxTimeRange,
    margin,
}) => {
    const eventMarkerSize = Math.max(0, height / 3);
    const eventMarkerOffsetY = Math.max(0, height / 2);
    return (
        <g
            data-testid="timeline-events-row"
            key={name}
            transform={`translate(${translateX}, ${translateY})`}
        >
            <rect
                fill={isOdd ? 'var(--tertiary-200)' : 'var(--base-100)'}
                stroke="var(--base-300)"
                height={height}
                width={width}
            />
            {events.map(
                ({ id, type, uid, reason, whitelisted, differenceInMilliseconds, timestamp }) => (
                    <EventMarker
                        key={id}
                        name={name}
                        uid={uid}
                        reason={reason}
                        type={type}
                        timestamp={timestamp}
                        whitelisted={whitelisted}
                        differenceInMilliseconds={differenceInMilliseconds}
                        translateX={translateX}
                        translateY={eventMarkerOffsetY}
                        size={eventMarkerSize}
                        minTimeRange={minTimeRange}
                        maxTimeRange={maxTimeRange}
                        margin={margin}
                    />
                )
            )}
        </g>
    );
};

EventsRow.propTypes = {
    minTimeRange: PropTypes.number.isRequired,
    maxTimeRange: PropTypes.number.isRequired,
    margin: PropTypes.number,
    height: PropTypes.number.isRequired,
    width: PropTypes.number.isRequired,
    translateX: PropTypes.number,
    translateY: PropTypes.number,
    name: PropTypes.string.isRequired,
    events: PropTypes.arrayOf(PropTypes.object),
    isOdd: PropTypes.bool,
};

EventsRow.defaultProps = {
    margin: 0,
    translateX: 0,
    translateY: 0,
    events: [],
    isOdd: false,
};

export default EventsRow;
