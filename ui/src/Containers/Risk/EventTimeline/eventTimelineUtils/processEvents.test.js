import processEvents from './processEvents';

describe('processEvents', () => {
    it('should add a field to the event object to show the difference in milliseconds between the event time and entity start time', () => {
        const startTime = '2020-04-20T15:00:00Z';
        const events = [
            {
                timestamp: '2020-04-20T16:00:00Z'
            },
            {
                timestamp: '2020-04-20T20:00:00Z'
            },
            {
                timestamp: '2020-04-21T15:00:00Z'
            }
        ];

        const value = processEvents(events, startTime);

        expect(value).toEqual([
            {
                timestamp: '2020-04-20T16:00:00Z',
                differenceInMilliseconds: 3600000
            },
            {
                timestamp: '2020-04-20T20:00:00Z',
                differenceInMilliseconds: 18000000
            },
            {
                timestamp: '2020-04-21T15:00:00Z',
                differenceInMilliseconds: 86400000
            }
        ]);
    });
});
