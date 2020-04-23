import React from 'react';
import PropTypes from 'prop-types';
import { useDrop } from 'react-dnd';
import { Trash2 } from 'react-feather';
import { Field, FieldArray } from 'redux-form';

import reduxFormPropTypes from 'constants/reduxFormPropTypes';
import Button from 'Components/Button';
import SectionHeaderInput from 'Components/SectionHeaderInput';
import AndOrOperator from 'Components/AndOrOperator';
import PolicyFieldCard from './PolicyFieldCard';
import { policyConfiguration } from './descriptors';
import { getPolicyCriteriaFieldKeys } from './utils';

const getEmptyPolicyFieldCard = (fieldKey) => ({
    field_name: fieldKey.name,
    boolean_operator: 'OR',
    values: [
        {
            value: '',
        },
    ],
    negate: false,
    fieldKey,
});

function PolicySection({ fields, sectionName, removeSectionHandler }) {
    const allFields = fields.getAll();
    const acceptedFields = getPolicyCriteriaFieldKeys(allFields);

    const [{ isOver, canDrop }, drop] = useDrop({
        accept: acceptedFields,
        drop: ({ fieldKey }) => {
            const newPolicyFieldCard = getEmptyPolicyFieldCard(fieldKey);
            fields.push(newPolicyFieldCard);
        },
        canDrop: ({ fieldKey }) => {
            return !allFields.find((field) => field.field_name === fieldKey.name);
        },
        collect: (monitor) => ({
            isOver: monitor.isOver(),
            canDrop: monitor.canDrop(),
        }),
    });

    function removeFieldHandler(index) {
        return () => {
            fields.remove(index);
        };
    }

    const disabledDrop = !canDrop && isOver;

    return (
        <>
            <div className="bg-base-300 border-2 border-base-100 rounded">
                <div className="flex justify-between items-center border-b-2 border-base-400">
                    <Field name={sectionName} component={SectionHeaderInput} />
                    <Button
                        onClick={removeSectionHandler}
                        icon={<Trash2 className="w-5 h-5" />}
                        className="p-2 border-l-2 border-base-400 hover:bg-base-400"
                    />
                </div>
                <div className="p-2">
                    {fields.map((name, i) => {
                        const field = fields.get(i);
                        const { field_name: fieldName } = field;
                        let { fieldKey } = field;
                        if (!fieldKey) {
                            fieldKey = policyConfiguration.descriptor.find(
                                (fieldObj) => fieldObj.name === fieldName
                            );
                        }
                        return (
                            <FieldArray
                                key={name}
                                name={`${name}.values`}
                                component={PolicyFieldCard}
                                booleanOperatorName={`${name}.boolean_operator`}
                                removeFieldHandler={removeFieldHandler(i)}
                                fieldKey={fieldKey}
                                toggleFieldName={`${name}.negate`}
                            />
                        );
                    })}
                    <div
                        ref={drop}
                        className={`${
                            disabledDrop
                                ? 'bg-base-300 border-base-400'
                                : 'bg-base-200 border-base-300'
                        } rounded border-2 border-dashed flex font-700 justify-center p-3 text-base-500 text-sm uppercase`}
                    >
                        Drop a policy field inside
                    </div>
                </div>
            </div>
            <AndOrOperator disabled />
        </>
    );
}

PolicySection.propTypes = {
    ...reduxFormPropTypes,
    sectionName: PropTypes.string.isRequired,
    removeSectionHandler: PropTypes.func.isRequired,
};

export default PolicySection;
