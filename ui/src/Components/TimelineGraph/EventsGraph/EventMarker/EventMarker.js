import React from 'react';
import PropTypes from 'prop-types';
import { scaleLinear } from 'd3-scale';

import { getWidth } from 'utils/d3Utils';
import { eventTypes } from 'constants/timelineTypes';
import mainViewSelector from 'Components/TimelineGraph/MainView/selectors';
import D3Anchor from 'Components/D3Anchor';
import PolicyViolationEvent from './PolicyViolationEvent';
import RestartEvent from './RestartEvent';
import ProcessActivityEvent from './ProcessActivityEvent';
import TerminationEvent from './TerminationEvent';

const EventMarker = ({
    name,
    args,
    type,
    uid,
    reason,
    timestamp,
    whitelisted,
    differenceInMilliseconds,
    translateX,
    translateY,
    minTimeRange,
    maxTimeRange,
    size,
    margin,
}) => {
    // the "container" argument is a reference to the container for the D3-related element
    function onUpdate(container) {
        const width = getWidth(mainViewSelector);
        const minRange = margin;
        const maxRange = width - margin;
        const xScale = scaleLinear()
            .domain([minTimeRange, maxTimeRange])
            .range([minRange, maxRange]);
        const x = xScale(differenceInMilliseconds).toFixed(0);

        container.attr(
            'transform',
            `translate(${Number(translateX) + Number(x) - size / 2}, ${
                Number(translateY) - size / 2
            })`
        );
    }

    return (
        <D3Anchor
            dataTestId="timeline-event-marker"
            translateX={translateX}
            translateY={translateY}
            onUpdate={onUpdate}
        >
            {type === eventTypes.POLICY_VIOLATION && (
                <PolicyViolationEvent name={name} type={type} timestamp={timestamp} width={size} />
            )}
            {type === eventTypes.PROCESS_ACTIVITY && (
                <ProcessActivityEvent
                    name={name}
                    args={args}
                    type={type}
                    uid={uid}
                    timestamp={timestamp}
                    whitelisted={whitelisted}
                    width={size}
                />
            )}
            {type === eventTypes.RESTART && (
                <RestartEvent name={name} type={type} timestamp={timestamp} width={size} />
            )}
            {type === eventTypes.TERMINATION && (
                <TerminationEvent
                    name={name}
                    type={type}
                    reason={reason}
                    timestamp={timestamp}
                    width={size}
                />
            )}
        </D3Anchor>
    );
};

EventMarker.propTypes = {
    name: PropTypes.string.isRequired,
    args: PropTypes.string,
    type: PropTypes.string.isRequired,
    uid: PropTypes.number,
    reason: PropTypes.string,
    timestamp: PropTypes.string.isRequired,
    whitelisted: PropTypes.bool,
    differenceInMilliseconds: PropTypes.number.isRequired,
    translateX: PropTypes.number.isRequired,
    translateY: PropTypes.number.isRequired,
    minTimeRange: PropTypes.number.isRequired,
    maxTimeRange: PropTypes.number.isRequired,
    size: PropTypes.number.isRequired,
    margin: PropTypes.number,
};

EventMarker.defaultProps = {
    uid: null,
    args: null,
    reason: null,
    whitelisted: false,
    margin: 0,
};

export default EventMarker;
