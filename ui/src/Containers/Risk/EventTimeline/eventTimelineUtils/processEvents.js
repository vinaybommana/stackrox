import getDifferenceInMilliseconds from './getDifferenceInMilliseconds';

/**
 * Processes the events data and returns a new array of events with modified data
 * @param {Object[]} events - The events data returned by the API call
 * @param {string} entityStartTime - The timestamp for the entity's start time
 * @returns {Object[]} - The processed events data
 */
function processEvents(events, entityStartTime) {
    return events.map((event) => ({
        ...event,
        differenceInMilliseconds: getDifferenceInMilliseconds(event.timestamp, entityStartTime),
    }));
}

export default processEvents;
