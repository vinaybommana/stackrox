import React from 'react';
import PropTypes from 'prop-types';
import { format } from 'date-fns';
import { Edit, Trash2 } from 'react-feather';

import dateTimeFormat from 'constants/dateTimeFormat';

const Comment = ({ comment, onDelete }) => {
    const { email, createdTime, updatedTime, message, canModify } = comment;
    function onEdit() {
        // do nothing for now
    }
    function onDeleteHandler() {
        onDelete(comment);
    }
    const isCommentUpdated = updatedTime && createdTime !== updatedTime;
    return (
        <div className="bg-primary-100 border border-primary-300 rounded-lg p-2">
            <div className="flex flex-1">
                <div className="text-primary-800 flex flex-1">{email}</div>
                <div className={`${!canModify && 'invisible'}`}>
                    <Edit
                        className="h-4 w-4 mx-2 text-primary-800 cursor-pointer hover:text-primary-500"
                        onClick={onEdit}
                    />
                    <Trash2
                        className="h-4 w-4 text-primary-800 cursor-pointer hover:text-primary-500"
                        onClick={onDeleteHandler}
                    />
                </div>
            </div>
            <div className="text-base-500 text-xs mt-1">
                {format(createdTime, dateTimeFormat)} {isCommentUpdated && '(edited)'}
            </div>
            <div className="mt-2 text-primary-800 leading-normal">{message}</div>
        </div>
    );
};

Comment.propTypes = {
    comment: PropTypes.shape({
        id: PropTypes.string,
        message: PropTypes.string,
        email: PropTypes.string,
        createdTime: PropTypes.string,
        updatedTime: PropTypes.string,
        canModify: PropTypes.bool
    }).isRequired,
    onDelete: PropTypes.func.isRequired
};

export default Comment;
