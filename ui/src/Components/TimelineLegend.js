import React from 'react';

import Tooltip from 'Components/Tooltip';
import TooltipOverlay from 'Components/TooltipOverlay';
import Button from 'Components/Button';
import PolicyViolationIcon from 'Components/TimelineGraph/EventsGraph/EventMarker/PolicyViolationEvent/PolicyViolationIcon';
import ProcessActivityIcon from 'Components/TimelineGraph/EventsGraph/EventMarker/ProcessActivityEvent/ProcessActivityIcon';
import WhitelistedProcessActivityIcon from 'Components/TimelineGraph/EventsGraph/EventMarker/ProcessActivityEvent/WhitelistedProcessActivityIcon';
import RestartIcon from 'Components/TimelineGraph/EventsGraph/EventMarker/RestartEvent/RestartIcon';
import TerminationIcon from 'Components/TimelineGraph/EventsGraph/EventMarker/TerminationEvent/TerminationIcon';

const ICON_SIZE = 15;

const TimelineLegend = () => {
    const content = (
        <TooltipOverlay>
            <div className="flex items-center mb-2">
                <PolicyViolationIcon width={ICON_SIZE} />
                <span className="ml-2">Policy Violation</span>
            </div>
            <div className="flex items-center mb-2">
                <ProcessActivityIcon width={ICON_SIZE} />
                <span className="ml-2">Process Activity</span>
            </div>
            <div className="flex items-center mb-2">
                <WhitelistedProcessActivityIcon width={ICON_SIZE} />
                <span className="ml-2">Whitelisted Process Activity</span>
            </div>
            <div className="flex items-center mb-2">
                <RestartIcon width={ICON_SIZE} />
                <span className="ml-2">Container Restart</span>
            </div>
            <div className="flex items-center">
                <TerminationIcon width={ICON_SIZE} />
                <span className="ml-2">Container Termination</span>
            </div>
        </TooltipOverlay>
    );
    return (
        <Tooltip trigger="click" position="right" content={content}>
            <div>
                <Button className="btn btn-base" dataTestId="timeline-legend" text="Show Legend" />
            </div>
        </Tooltip>
    );
};

export default TimelineLegend;
