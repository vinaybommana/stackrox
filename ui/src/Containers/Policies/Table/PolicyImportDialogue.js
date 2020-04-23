import React, { useState, useRef, useCallback } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { useDropzone } from 'react-dropzone';
import { Upload } from 'react-feather';
import pluralize from 'pluralize';

import CustomDialogue from 'Components/CustomDialogue';
import Message from 'Components/Message';
import { fileUploadColors } from 'constants/visuals/colors';
import { actions as pageActions } from 'reducers/policies/page';
import { importPolicies } from 'services/PoliciesService';

const PolicyImportDialogue = ({ closeAction, importPolicySuccess }) => {
    const [messageObj, setMessageObj] = useState(null);
    const [policies, setPolicies] = useState([]);
    const dialogueRef = useRef(null);

    const onDrop = useCallback((acceptedFiles) => {
        setMessageObj(null);

        acceptedFiles.forEach((file) => {
            // check file type.
            if (file && !file.name.includes('.json')) {
                setMessageObj({ type: 'warn', message: 'Only JSON files are supported.' });
                return;
            }

            const reader = new FileReader();
            reader.onload = () => {
                const fileContent = reader.result;
                try {
                    const jsonObj = JSON.parse(fileContent);
                    setPolicies(jsonObj.policies);
                } catch (err) {
                    setMessageObj({ type: 'error', message: err.message });
                }
            };
            reader.onerror = (e) => {
                reader.abort();
                setMessageObj({ type: 'error', message: e.message });
            };
            reader.readAsText(file);
        });
    }, []);

    const { getRootProps, getInputProps } = useDropzone({ onDrop });

    function startImport() {
        importPolicies(policies)
            .then(() => {
                // TODO: handle responses that indicate a conflict in policy name or ID
                //  and offer the user the option
                //  - to rename the policy being imported,
                //  - or overwrite the existing policy with the policy being imported

                const importedPolicyId = policies[0]?.id; // API always returns a list
                if (importedPolicyId) {
                    setMessageObj({ type: 'info', message: 'Policy successfully imported' });
                    setTimeout(handleClose, 3000);

                    importPolicySuccess(importedPolicyId);
                }
            })
            .catch((err) => {
                setMessageObj({
                    type: 'error',
                    message: `A network error occurred: ${err.message}`,
                });
            });
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
            confirmDisabled={policies.length < 1}
            onCancel={handleClose}
        >
            <div
                className="overflow-auto p-4"
                ref={dialogueRef}
                data-testid="policy-import-modal-content"
            >
                <>
                    {messageObj && <Message type={messageObj.type} message={messageObj.message} />}
                    <div className="flex flex-col bg-base-100 rounded-sm shadow flex-grow flex-shrink-0 mb-4">
                        <div className="my-3 px-3 font-600 text-lg leading-loose text-base-600">
                            Upload a policy that has been exported from StackRox system.
                        </div>
                        <div
                            {...getRootProps()}
                            className="bg-warning-100 border border-dashed border-warning-500 cursor-pointer flex flex-col h-full hover:bg-warning-200 justify-center min-h-32 mt-3 outline-none py-3 self-center uppercase w-full"
                        >
                            <input {...getInputProps()} />
                            <div className="flex flex-shrink-0 flex-col">
                                <div
                                    className="mt-3 h-18 w-18 self-center rounded-full flex items-center justify-center flex-shrink-0"
                                    style={{
                                        background: fileUploadColors.BACKGROUND_COLOR,
                                        color: fileUploadColors.ICON_COLOR,
                                    }}
                                >
                                    <Upload
                                        className="h-8 w-8"
                                        strokeWidth="1.5px"
                                        data-testid="upload-icon"
                                    />
                                </div>
                                <span className="font-700 mt-3 text-center text-warning-800">
                                    Choose a policy file in JSON format
                                </span>
                            </div>
                        </div>
                    </div>
                    {policies.length > 0 && (
                        <div className="flex flex-col bg-base-100 flex-grow flex-shrink-0 mb-4">
                            <h3 className="b-2 font-700 text-lg">
                                The following {`${pluralize('policy', policies.length)}`} will be
                                imported:
                            </h3>
                            {policies.map((policy) => (
                                <li
                                    key={policy.id}
                                    className="p-2 text-sm text-primary-800 font-600 w-full"
                                >
                                    {policy.name}
                                </li>
                            ))}
                        </div>
                    )}
                </>
            </div>
        </CustomDialogue>
    );
};

PolicyImportDialogue.propTypes = {
    closeAction: PropTypes.func.isRequired,
    importPolicySuccess: PropTypes.func.isRequired,
};

const mapDispatchToProps = {
    importPolicySuccess: pageActions.importPolicySuccess,
};

export default connect(null, mapDispatchToProps)(PolicyImportDialogue);
