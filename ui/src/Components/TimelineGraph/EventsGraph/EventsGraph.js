import React from 'react';
import PropTypes from 'prop-types';

import EventsRow from './EventsRow';

const MAX_ROW_HEIGHT = 48;
const MIN_ROW_HEIGHT = 0;

const EventsGraph = ({
    data,
    translateX,
    translateY,
    minTimeRange,
    maxTimeRange,
    height,
    width,
    numRows,
    margin,
    isHeightAdjustable,
}) => {
    const rowHeight = isHeightAdjustable
        ? Math.min(Math.max(MIN_ROW_HEIGHT, Math.floor(height / numRows) - 1), MAX_ROW_HEIGHT)
        : MAX_ROW_HEIGHT;
    return (
        <g
            data-testid="timeline-events-graph"
            transform={`translate(${translateX}, ${translateY})`}
        >
            {data.map((datum, index) => {
                const { id, name, events } = datum;
                const isOddRow = index % 2 !== 0;
                return (
                    <EventsRow
                        key={id}
                        name={name}
                        events={events}
                        isOdd={isOddRow}
                        height={rowHeight}
                        width={width}
                        translateX={0}
                        translateY={index * rowHeight}
                        minTimeRange={minTimeRange}
                        maxTimeRange={maxTimeRange}
                        margin={margin}
                    />
                );
            })}
        </g>
    );
};

EventsGraph.propTypes = {
    minTimeRange: PropTypes.number.isRequired,
    maxTimeRange: PropTypes.number.isRequired,
    data: PropTypes.arrayOf(PropTypes.object).isRequired,
    numRows: PropTypes.number.isRequired,
    margin: PropTypes.number,
    height: PropTypes.number.isRequired,
    width: PropTypes.number.isRequired,
    translateX: PropTypes.number,
    translateY: PropTypes.number,
    isHeightAdjustable: PropTypes.bool,
};

EventsGraph.defaultProps = {
    margin: 0,
    translateX: 0,
    translateY: 0,
    isHeightAdjustable: false,
};

export default EventsGraph;
