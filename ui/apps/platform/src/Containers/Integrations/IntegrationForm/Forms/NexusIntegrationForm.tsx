import React, { ReactElement } from 'react';
import { TextInput, PageSection, Form, Checkbox } from '@patternfly/react-core';
import * as yup from 'yup';

import { ImageIntegrationBase } from 'services/ImageIntegrationsService';

import usePageState from 'Containers/Integrations/hooks/usePageState';
import FormMessage from 'Components/PatternFly/FormMessage';
import FormTestButton from 'Components/PatternFly/FormTestButton';
import FormSaveButton from 'Components/PatternFly/FormSaveButton';
import FormCancelButton from 'Components/PatternFly/FormCancelButton';
import useIntegrationForm from '../useIntegrationForm';
import { IntegrationFormProps } from '../integrationFormTypes';

import IntegrationFormActions from '../IntegrationFormActions';
import FormLabelGroup from '../FormLabelGroup';

export type NexusIntegration = {
    categories: 'REGISTRY'[];
    docker: {
        endpoint: string;
        username: string;
        password: string;
        insecure: boolean;
    };
    type: 'nexus';
} & ImageIntegrationBase;

export type NexusIntegrationFormValues = {
    config: NexusIntegration;
    updatePassword: boolean;
};

export const validationSchema = yup.object().shape({
    config: yup.object().shape({
        name: yup.string().trim().required('An integration name is required'),
        categories: yup
            .array()
            .of(yup.string().trim().oneOf(['REGISTRY']))
            .min(1, 'Must have at least one type selected')
            .required('A category is required'),
        docker: yup.object().shape({
            endpoint: yup.string().trim().required('An endpoint is required'),
            username: yup
                .string()
                .test(
                    'username-test',
                    'A username is required if the integration has a password',
                    (value, context: yup.TestContext) => {
                        const hasPassword = !!context.parent.password;
                        if (!hasPassword) {
                            return true;
                        }
                        const trimmedValue = value?.trim();
                        return !!trimmedValue;
                    }
                ),
            password: yup
                .string()
                .test(
                    'password-test',
                    'A password is required if the integration has a username',
                    (value, context: yup.TestContext) => {
                        const requirePasswordField =
                            // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                            // @ts-ignore
                            context?.from[2]?.value?.updatePassword || false;
                        const hasUsername = !!context.parent.username;

                        if (!requirePasswordField || !hasUsername) {
                            return true;
                        }

                        const trimmedValue = value?.trim();
                        return !!trimmedValue;
                    }
                ),
            insecure: yup.bool(),
        }),
        skipTestIntegration: yup.bool(),
        type: yup.string().matches(/nexus/),
    }),
    updatePassword: yup.bool(),
});

export const defaultValues: NexusIntegrationFormValues = {
    config: {
        id: '',
        name: '',
        categories: ['REGISTRY'],
        docker: {
            endpoint: '',
            username: '',
            password: '',
            insecure: false,
        },
        autogenerated: false,
        clusterId: '',
        clusters: [],
        skipTestIntegration: false,
        type: 'nexus',
    },
    updatePassword: true,
};

function NexusIntegrationForm({
    initialValues = null,
    isEditable = false,
}: IntegrationFormProps<NexusIntegration>): ReactElement {
    const formInitialValues = { ...defaultValues, ...initialValues };
    if (initialValues) {
        formInitialValues.config = { ...formInitialValues.config, ...initialValues };
        // We want to clear the password because backend returns '******' to represent that there
        // are currently stored credentials
        formInitialValues.config.docker.password = '';

        // Don't assume user wants to change password; that has caused confusing UX.
        formInitialValues.updatePassword = false;
    }
    const {
        values,
        touched,
        errors,
        dirty,
        isValid,
        setFieldValue,
        handleBlur,
        isSubmitting,
        isTesting,
        onSave,
        onTest,
        onCancel,
        message,
    } = useIntegrationForm<NexusIntegrationFormValues>({
        initialValues: formInitialValues,
        validationSchema,
    });
    const { isCreating } = usePageState();
    const areCredentialsRequired =
        values.updatePassword || !!values.config.docker.password || !!values.config.docker.username;

    function onChange(value, event) {
        return setFieldValue(event.target.id, value);
    }

    function onUpdateCredentialsChange(value, event) {
        setFieldValue('config.docker.password', '');
        return setFieldValue(event.target.id, value);
    }

    return (
        <>
            <PageSection variant="light" isFilled hasOverflowScroll>
                <FormMessage message={message} />
                <Form isWidthLimited>
                    <FormLabelGroup
                        label="Integration name"
                        isRequired
                        fieldId="config.name"
                        touched={touched}
                        errors={errors}
                    >
                        <TextInput
                            isRequired
                            type="text"
                            id="config.name"
                            value={values.config.name}
                            onChange={onChange}
                            onBlur={handleBlur}
                            isDisabled={!isEditable}
                        />
                    </FormLabelGroup>
                    <FormLabelGroup
                        label="Endpoint"
                        isRequired
                        fieldId="config.docker.endpoint"
                        touched={touched}
                        errors={errors}
                    >
                        <TextInput
                            isRequired
                            type="text"
                            id="config.docker.endpoint"
                            value={values.config.docker.endpoint}
                            onChange={onChange}
                            onBlur={handleBlur}
                            isDisabled={!isEditable}
                        />
                    </FormLabelGroup>
                    <FormLabelGroup
                        label="Username"
                        isRequired={areCredentialsRequired}
                        fieldId="config.docker.username"
                        touched={touched}
                        errors={errors}
                        helperText="A username is required if the integration has a password"
                    >
                        <TextInput
                            isRequired={areCredentialsRequired}
                            type="text"
                            id="config.docker.username"
                            value={values.config.docker.username}
                            onChange={onChange}
                            onBlur={handleBlur}
                            isDisabled={!isEditable || !values.updatePassword}
                        />
                    </FormLabelGroup>
                    {!isCreating && isEditable && (
                        <FormLabelGroup
                            fieldId="updatePassword"
                            helperText="Enable this option to replace currently stored credentials (if any)"
                            errors={errors}
                        >
                            <Checkbox
                                label="Update stored credentials"
                                id="updatePassword"
                                isChecked={values.updatePassword}
                                onChange={onUpdateCredentialsChange}
                                onBlur={handleBlur}
                                isDisabled={!isEditable}
                            />
                        </FormLabelGroup>
                    )}
                    <FormLabelGroup
                        label="Password"
                        isRequired={areCredentialsRequired}
                        fieldId="config.docker.password"
                        touched={touched}
                        errors={errors}
                        helperText="A password is required if the integration has a username"
                    >
                        <TextInput
                            isRequired={areCredentialsRequired}
                            type="password"
                            id="config.docker.password"
                            value={values.config.docker.password}
                            onChange={onChange}
                            onBlur={handleBlur}
                            isDisabled={!isEditable || !values.updatePassword}
                            placeholder={
                                values.updatePassword
                                    ? ''
                                    : 'Currently-stored password will be used.'
                            }
                        />
                    </FormLabelGroup>
                    <FormLabelGroup fieldId="config.docker.insecure" errors={errors}>
                        <Checkbox
                            label="Disable TLS certificate validation (insecure)"
                            id="config.docker.insecure"
                            isChecked={values.config.docker.insecure}
                            onChange={onChange}
                            onBlur={handleBlur}
                            isDisabled={!isEditable}
                        />
                    </FormLabelGroup>
                    <FormLabelGroup fieldId="config.skipTestIntegration" errors={errors}>
                        <Checkbox
                            label="Create integration without testing"
                            id="config.skipTestIntegration"
                            isChecked={values.config.skipTestIntegration}
                            onChange={onChange}
                            onBlur={handleBlur}
                            isDisabled={!isEditable}
                        />
                    </FormLabelGroup>
                </Form>
            </PageSection>
            {isEditable && (
                <IntegrationFormActions>
                    <FormSaveButton
                        onSave={onSave}
                        isSubmitting={isSubmitting}
                        isTesting={isTesting}
                        isDisabled={!dirty || !isValid}
                    >
                        Save
                    </FormSaveButton>
                    <FormTestButton
                        onTest={onTest}
                        isSubmitting={isSubmitting}
                        isTesting={isTesting}
                        isDisabled={!isValid}
                    >
                        Test
                    </FormTestButton>
                    <FormCancelButton onCancel={onCancel}>Cancel</FormCancelButton>
                </IntegrationFormActions>
            )}
        </>
    );
}

export default NexusIntegrationForm;
