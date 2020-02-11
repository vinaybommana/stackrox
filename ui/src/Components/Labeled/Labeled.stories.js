import React from 'react';

import Labeled from './Labeled';

export default {
    title: 'Labeled',
    component: Labeled
};

export const withTextLabelAndValue = () => <Labeled label="Label">Value</Labeled>;

export const withTextLabelAndInput = () => (
    <Labeled label="Enter value">
        <input className="border-2" />
    </Labeled>
);

export const withRenderPropLabelAndInput = () => {
    function Label() {
        return (
            <p>
                Enter value <i className="text-alert-700">(required)</i>
            </p>
        );
    }

    return (
        <Labeled label={Label}>
            <input className="border-2" />
        </Labeled>
    );
};
