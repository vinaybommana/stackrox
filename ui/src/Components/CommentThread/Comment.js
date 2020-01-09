import React, { useState } from 'react';
import { useForm } from 'react-hook-form';
import PropTypes from 'prop-types';
import { format } from 'date-fns';
import { Edit, Trash2, XCircle } from 'react-feather';

import dateTimeFormat from 'constants/dateTimeFormat';

import TextArea from 'Components/forms/TextArea';
import CustomDialogue from 'Components/CustomDialogue';

const ActionButtons = ({ isEditing, canModify, onToggleEdit, onDelete }) => {
    if (isEditing) {
        return (
            <div>
                <XCircle
                    className="h-4 w-4 ml-2 text-success-800 cursor-pointer hover:text-success-500"
                    onClick={onToggleEdit}
                />
            </div>
        );
    }
    return (
        <div className={`${!canModify && 'invisible'}`}>
            <Edit
                className="h-4 w-4 mx-2 text-primary-800 cursor-pointer hover:text-primary-500"
                onClick={onToggleEdit}
            />
            <Trash2
                className="h-4 w-4 text-primary-800 cursor-pointer hover:text-primary-500"
                onClick={onDelete}
            />
        </div>
    );
};

const InputForm = ({ value, onSubmit }) => {
    const { register, handleSubmit, errors } = useForm();
    return (
        <form onSubmit={handleSubmit(onSubmit)}>
            <TextArea
                name="message"
                required
                register={register}
                errors={errors}
                rows="5"
                cols="33"
                defaultValue={value}
                placeholder="Write a comment here..."
            />
            <div className="flex justify-end">
                <input
                    className="bg-success-300 border border-success-800 p-1 rounded-sm text-sm text-success-900 uppercase hover:bg-success-400 cursor-pointer"
                    type="submit"
                    value="Save"
                />
            </div>
        </form>
    );
};

const Comment = ({ comment, onDelete, onSave, defaultEdit }) => {
    const [isEditing, setEdit] = useState(defaultEdit);
    const [isDialogueOpen, setIsDialogueOpen] = useState(false);
    const { email, createdTime, updatedTime, message, canModify } = comment;

    const isCommentUpdated = updatedTime && createdTime !== updatedTime;

    function onToggleEdit() {
        setEdit(!isEditing);
    }

    function onSubmit(data) {
        onToggleEdit();
        onSave(comment, data.message);
    }

    function onDeleteHandler() {
        setIsDialogueOpen(true);
    }

    function cancelDeletion() {
        setIsDialogueOpen(false);
    }

    function confirmDeletion() {
        onDelete(comment);
        setIsDialogueOpen(false);
    }

    return (
        <div
            className={`${
                isEditing
                    ? 'bg-success-200 border-success-500'
                    : 'bg-primary-100 border-primary-300'
            } border rounded-lg p-2`}
        >
            <div className="flex flex-1">
                <div className="text-primary-800 flex flex-1">{email}</div>
                <ActionButtons
                    isEditing={isEditing}
                    canModify={canModify}
                    onToggleEdit={onToggleEdit}
                    onDelete={onDeleteHandler}
                />
            </div>
            <div className="text-base-500 text-xs mt-1">
                {format(createdTime, dateTimeFormat)} {isCommentUpdated && '(edited)'}
            </div>
            <div className="mt-2 text-primary-800 leading-normal">
                {isEditing ? <InputForm value={message} onSubmit={onSubmit} /> : message}
            </div>
            {isDialogueOpen && (
                <CustomDialogue
                    title="Delete Comment?"
                    onConfirm={confirmDeletion}
                    confirmText="Yes"
                    onCancel={cancelDeletion}
                />
            )}
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
    onDelete: PropTypes.func.isRequired,
    onSave: PropTypes.func.isRequired,
    defaultEdit: PropTypes.bool
};

Comment.defaultProps = {
    defaultEdit: false
};

export default Comment;
