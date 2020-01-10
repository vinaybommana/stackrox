import React, { useState } from 'react';
import PropTypes from 'prop-types';
import pluralize from 'pluralize';
import sortBy from 'lodash/sortBy';

import CollapsibleCard from 'Components/CollapsibleCard';
import Button from 'Components/Button';
import Comment from './Comment';

const CommentThread = ({ type, comments, onSave, onDelete, defaultLimit, defaultOpen }) => {
    const [limit, setLimit] = useState(defaultLimit);

    const sortedComments = sortBy(comments, ['createdTime', 'email']);
    const { length } = sortedComments;
    const hasMoreComments = limit < length;

    function showMoreComments() {
        setLimit(limit + defaultLimit);
    }

    return (
        <CollapsibleCard
            title={`${length} ${type} ${pluralize('Comment', length)}`}
            open={defaultOpen}
        >
            <div className="p-3">
                {sortedComments.slice(0, limit).map((comment, i) => (
                    <div key={comment.id} className={i === 0 ? 'mt-0' : 'mt-3'}>
                        <Comment comment={comment} onSave={onSave} onDelete={onDelete} />
                    </div>
                ))}
                {hasMoreComments && (
                    <div className="flex flex-1 justify-center mt-3">
                        <Button
                            className="bg-primary-200 border border-primary-800 hover:bg-primary-300 p-1 rounded-full rounded-sm text-sm text-success-900 uppercase"
                            text="Load More Comments"
                            onClick={showMoreComments}
                        />
                    </div>
                )}
            </div>
        </CollapsibleCard>
    );
};

CommentThread.propTypes = {
    type: PropTypes.string.isRequired,
    comments: PropTypes.arrayOf(
        PropTypes.shape({
            id: PropTypes.string,
            message: PropTypes.string,
            email: PropTypes.string,
            createdTime: PropTypes.string,
            updatedTime: PropTypes.string,
            canModify: PropTypes.bool
        })
    ),
    onSave: PropTypes.func.isRequired,
    onDelete: PropTypes.func.isRequired,
    defaultLimit: PropTypes.number,
    defaultOpen: PropTypes.bool
};

CommentThread.defaultProps = {
    comments: [],
    defaultLimit: 3,
    defaultOpen: false
};

export default CommentThread;
