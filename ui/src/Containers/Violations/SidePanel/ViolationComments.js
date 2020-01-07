import React from 'react';

import CommentThread from 'Components/CommentThread';

const ViolationComments = () => {
    const comments = [
        {
            id: '1',
            message:
                'Completely unrelated, but check out this subreddit https://www.reddit.com/r/aww/',
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
            id: '1',
            message: 'Also, do you want to hear a joke?',
            email: 'sc@stackrox.com',
            createdTime: '2019-12-30T22:21:31.218853651Z',
            updatedTime: '2019-12-30T22:21:31.218853651Z',
            canModify: true
        },
        {
            id: '2',
            message: 'No',
            email: 'ls@stackrox.com',
            createdTime: '2019-12-31T21:21:31.218853651Z',
            updatedTime: '2019-12-31T21:21:31.218853651Z',
            canModify: false
        }
    ];

    function onSave() {}

    function onDelete() {}

    return <CommentThread comments={comments} onSave={onSave} onDelete={onDelete} />;
};

export default ViolationComments;
