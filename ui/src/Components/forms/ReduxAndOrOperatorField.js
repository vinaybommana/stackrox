import React from 'react';
import PropTypes from 'prop-types';
import { Field } from 'redux-form';

import AndOrOperator from 'Components/AndOrOperator';
import BOOLEAN_LOGIC_VALUES from 'constants/booleanLogicValues';

function ReduxAndOrOperator({ input, disabled }) {
    const { value, onChange } = input;
    function onToggle() {
        const newValue =
            value === BOOLEAN_LOGIC_VALUES.AND ? BOOLEAN_LOGIC_VALUES.OR : BOOLEAN_LOGIC_VALUES.AND;
        onChange(newValue);
    }
    return <AndOrOperator value={value} onToggle={onToggle} disabled={disabled} />;
}

function ReduxAndOrOperatorField({ name, disabled }) {
    return (
        <Field
            key={name}
            name={name}
            id={name}
            component={ReduxAndOrOperator}
            disabled={disabled}
        />
    );
}

ReduxAndOrOperatorField.propTypes = {
    name: PropTypes.string.isRequired,
    disabled: PropTypes.bool.isRequired,
};

export default ReduxAndOrOperatorField;
