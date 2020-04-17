import React, { useState, useRef } from 'react';
import PropTypes from 'prop-types';

import CustomDialogue from 'Components/CustomDialogue';
import Message from 'Components/Message';

const PolicyImportDialogue = ({ closeAction }) => {
    const [messageObj, setMessageObj] = useState(null);
    const dialogueRef = useRef(null);

    function startImport() {
        setMessageObj({ type: 'info', message: 'Policy successfully imported' });
    }

    function handleClose() {
        closeAction();
    }

    return (
        <CustomDialogue
            className="max-w-3/4 md:max-w-2/3 lg:max-w-1/2"
            title="Import a Policy"
            onConfirm={startImport}
            confirmText="Begin Import"
            confirmDisabled={false}
            onCancel={handleClose}
        >
            <div className="overflow-auto p-4" ref={dialogueRef}>
                <>
                    {messageObj && <Message type={messageObj.type} message={messageObj.message} />}
                    <div>File upload goes here</div>
                </>
            </div>
        </CustomDialogue>
    );
};

PolicyImportDialogue.propTypes = {
    closeAction: PropTypes.func.isRequired
};

export default PolicyImportDialogue;
