import React from 'react';
import PropTypes from 'prop-types';
import { Field } from 'redux-form';

import RadioButtonGroup from 'Components/RadioButtonGroup';

function ReduxRadioButtonGroup({ input, buttons, groupClassName }) {
    const { value, onChange } = input;
    return (
        <RadioButtonGroup
            buttons={buttons}
            onClick={onChange}
            selected={value}
            groupClassName={groupClassName}
        />
    );
}

function ReduxRadioButtonGroupField({ name, buttons, groupClassName }) {
    return (
        <Field
            key={name}
            name={name}
            id={name}
            component={ReduxRadioButtonGroup}
            buttons={buttons}
            groupClassName={groupClassName}
        />
    );
}

ReduxRadioButtonGroupField.propTypes = {
    name: PropTypes.string.isRequired,
    buttons: PropTypes.arrayOf(
        PropTypes.shape({
            text: PropTypes.string.isRequired,
            value: PropTypes.bool,
        })
    ).isRequired,
    groupClassName: PropTypes.string,
};

ReduxRadioButtonGroupField.defaultProps = {
    groupClassName: '',
};

export default ReduxRadioButtonGroupField;
