import React from 'react';
import PropTypes from 'prop-types';
import { event } from 'd3-selection';
import { brushX } from 'd3-brush';

import D3Anchor from 'Components/D3Anchor';

const Brush = ({ translateX, translateY, width, height, onSelectionChange, margin }) => {
    function brushEnded() {
        if (!event.sourceEvent) return; // Only transition after input.
        if (!event.selection) {
            onSelectionChange(null);
            return;
        }
        const selection = {
            start: event.selection[0],
            end: event.selection[1],
        };
        onSelectionChange(selection);
    }

    // the "container" argument is a reference to the container for the D3-related element
    function onUpdate(container) {
        const minHorizontalExtent = margin;
        const maxHorizontalExtent = width - margin;
        const brush = container.call(
            brushX()
                .extent([
                    [minHorizontalExtent, 0],
                    [maxHorizontalExtent, height],
                ])
                .on('end', brushEnded)
        );
        brush
            .select('rect.selection')
            .style('fill', 'var(--accent-500)')
            .style('stroke', 'var(--accent-500)');
    }

    return (
        <D3Anchor
            dataTestId="timeline-brush"
            translateX={translateX}
            translateY={translateY}
            onUpdate={onUpdate}
        />
    );
};

Brush.propTypes = {
    margin: PropTypes.number,
    height: PropTypes.number.isRequired,
    width: PropTypes.number.isRequired,
    translateX: PropTypes.number,
    translateY: PropTypes.number,
    onSelectionChange: PropTypes.func.isRequired,
};

Brush.defaultProps = {
    margin: 0,
    translateX: 0,
    translateY: 0,
};

export default Brush;
