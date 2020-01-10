import React, { useState } from 'react';

import CommentThread from 'Components/CommentThread';

const defaultComments = [
    {
        id: '1',
        message: 'Completely unrelated, but check out this subreddit https://www.reddit.com/r/aww/',
        email: 'sc@stackrox.com',
        createdTime: '2019-12-29T21:21:31.218853651Z',
        updatedTime: '2019-12-30T21:21:31.218853651Z',
        canModify: true
    },
    {
        id: '2',
        message: 'Oh nice! This is the content I like',
        email: 'ls@stackrox.com',
        createdTime: '2019-12-30T21:21:31.218853651Z',
        updatedTime: '2019-12-30T21:21:31.218853651Z',
        canModify: false
    },
    {
        id: '3',
        message: 'Also, do you want to hear a joke?',
        email: 'sc@stackrox.com',
        createdTime: '2019-12-30T22:21:31.218853651Z',
        updatedTime: '2019-12-30T22:21:31.218853651Z',
        canModify: true
    },
    {
        id: '4',
        message: 'No',
        email: 'ls@stackrox.com',
        createdTime: '2019-12-31T21:21:31.218853651Z',
        updatedTime: '2019-12-31T21:21:31.218853651Z',
        canModify: false
    }
];

const ViolationComments = () => {
    const [comments, setComments] = useState(defaultComments);

    function onSave(comment, message) {
        const newComments = comments.filter(datum => datum.id !== comment.id);
        newComments.push({ ...comment, message });
        setComments(newComments);
    }

    function onDelete(comment) {
        setComments(comments.filter(datum => datum.id !== comment.id));
    }

    return <CommentThread comments={comments} onSave={onSave} onDelete={onDelete} />;
};

export default React.memo(ViolationComments);
