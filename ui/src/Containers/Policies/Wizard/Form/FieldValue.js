import React from 'react';
import PropTypes from 'prop-types';

import ReduxAndOrOperatorField from 'Components/forms/ReduxAndOrOperatorField';
import FormFieldRemoveButton from 'Components/FormFieldRemoveButton';
import Field from './Field';

function FieldValue({ name, length, booleanOperatorName, fieldKey, removeValueHandler, index }) {
    return (
        <>
            <div className="flex">
                <Field key={name} field={fieldKey} name={name} />
                {/* only show remove button if there is more than one value */}
                {length > 1 && (
                    <FormFieldRemoveButton
                        field={name}
                        onClick={removeValueHandler}
                        className="border-base-300 hover:border-base-400 hover:text-base-600 rounded-r text-base-100 text-base-500"
                    />
                )}
            </div>
            {/* only show and/or operator if not at end of array */}
            {index + 1 !== length && (
                <ReduxAndOrOperatorField
                    name={booleanOperatorName}
                    disabled={!fieldKey.canBooleanLogic}
                />
            )}
        </>
    );
}

FieldValue.propTypes = {
    name: PropTypes.string.isRequired,
    length: PropTypes.number.isRequired,
    fieldKey: PropTypes.shape({
        canBooleanLogic: PropTypes.bool,
    }).isRequired,
    booleanOperatorName: PropTypes.string.isRequired,
    removeValueHandler: PropTypes.func.isRequired,
    index: PropTypes.number.isRequired,
};

export default FieldValue;
