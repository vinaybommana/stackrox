import React from 'react';
import PropTypes from 'prop-types';

import SaveButton from 'Components/SaveButton';

function Button(props) {
    if (!props.isEditing) {
        return (
            <button className="btn btn-base mr-3" type="button" onClick={props.onEdit}>
                Edit
            </button>
        );
    }
    return (
        <div className="flex mr-3">
            <button className="btn btn-base mr-2" type="button" onClick={props.onCancel}>
                Cancel
            </button>
            <SaveButton formName="role-form" />
        </div>
    );
}

Button.propTypes = {
    isEditing: PropTypes.bool.isRequired,
    onEdit: PropTypes.func.isRequired,
    onCancel: PropTypes.func.isRequired,
};

export default Button;
