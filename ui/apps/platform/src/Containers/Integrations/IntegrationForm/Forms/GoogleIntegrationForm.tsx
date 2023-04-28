import React, { ReactElement } from 'react';
import {
    Checkbox,
    Form,
    PageSection,
    TextArea,
    TextInput,
    ToggleGroup,
    ToggleGroupItem,
} from '@patternfly/react-core';
import * as yup from 'yup';

import { GoogleImageIntegration } from 'types/imageIntegration.proto';

import usePageState from 'Containers/Integrations/hooks/usePageState';
import FormMessage from 'Components/PatternFly/FormMessage';
import FormTestButton from 'Components/PatternFly/FormTestButton';
import FormSaveButton from 'Components/PatternFly/FormSaveButton';
import FormCancelButton from 'Components/PatternFly/FormCancelButton';
import useIntegrationForm from '../useIntegrationForm';
import { IntegrationFormProps } from '../integrationFormTypes';

import IntegrationFormActions from '../IntegrationFormActions';
import FormLabelGroup from '../FormLabelGroup';

import { categoriesUtilsForRegistryScanner } from '../../utils/integrationUtils';

const { categoriesAlternatives, getCategoriesText, matchCategoriesAlternative, validCategories } =
    categoriesUtilsForRegistryScanner;

export type GoogleIntegrationFormValues = {
    config: GoogleImageIntegration;
    updatePassword: boolean;
};

export const validationSchema = yup.object().shape({
    config: yup.object().shape({
        name: yup.string().trim().required('An integration name is required'),
        categories: yup
            .array()
            .of(yup.string().trim().oneOf(validCategories))
            .min(1, 'Must have at least one type selected')
            .required('A category is required'),
        google: yup.object().shape({
            endpoint: yup.string().trim().required('An endpoint is required'),
            project: yup.string().trim().required('A project is required'),
            serviceAccount: yup
                .string()
                .test(
                    'serviceAccount-test',
                    'A service account key is required',
                    (value, context: yup.TestContext) => {
                        const requirePasswordField =
                            // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                            // @ts-ignore
                            context?.from[2]?.value?.updatePassword || false;

                        if (!requirePasswordField) {
                            return true;
                        }

                        const trimmedValue = value?.trim();
                        return !!trimmedValue;
                    }
                ),
        }),
        skipTestIntegration: yup.bool(),
        type: yup.string().matches(/google/),
    }),
    updatePassword: yup.bool(),
});

export const defaultValues: GoogleIntegrationFormValues = {
    config: {
        id: '',
        name: '',
        categories: ['REGISTRY'],
        google: {
            endpoint: '',
            project: '',
            serviceAccount: '',
        },
        autogenerated: false,
        clusterId: '',
        skipTestIntegration: false,
        type: 'google',
    },
    updatePassword: true,
};

function GoogleIntegrationForm({
    initialValues = null,
    isEditable = false,
}: IntegrationFormProps<GoogleImageIntegration>): ReactElement {
    const formInitialValues = { ...defaultValues, ...initialValues };
    if (initialValues) {
        formInitialValues.config = { ...formInitialValues.config, ...initialValues };
        // We want to clear the password because backend returns '******' to represent that there
        // are currently stored credentials
        formInitialValues.config.google.serviceAccount = '';

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
    } = useIntegrationForm<GoogleIntegrationFormValues>({
        initialValues: formInitialValues,
        validationSchema,
    });

    const { isCreating } = usePageState();

    function onChange(value, event) {
        return setFieldValue(event.target.id, value);
    }

    function onUpdateCredentialsChange(value, event) {
        setFieldValue('config.google.serviceAccount', '');
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
                            type="text"
                            id="config.name"
                            placeholder="(ex. Google Registry and Scanner)"
                            value={values.config.name}
                            onChange={onChange}
                            onBlur={handleBlur}
                            isDisabled={!isEditable}
                        />
                    </FormLabelGroup>
                    <FormLabelGroup
                        label="Type"
                        isRequired
                        fieldId="config.categories"
                        touched={touched}
                        errors={errors}
                    >
                        <ToggleGroup id="config.categories" areAllGroupsDisabled={!isEditable}>
                            {categoriesAlternatives.map((categoriesAlternative) => {
                                const [categoriesAlternativeItem0] = categoriesAlternative;
                                const text = getCategoriesText(categoriesAlternativeItem0);
                                const isSelected = matchCategoriesAlternative(
                                    categoriesAlternative,
                                    values.config.categories
                                );
                                return (
                                    <ToggleGroupItem
                                        key={text}
                                        text={text}
                                        isSelected={isSelected}
                                        onChange={() =>
                                            setFieldValue(
                                                'config.categories',
                                                categoriesAlternativeItem0
                                            )
                                        }
                                    />
                                );
                            })}
                        </ToggleGroup>
                    </FormLabelGroup>
                    <FormLabelGroup
                        label="Registry endpoint"
                        isRequired
                        fieldId="config.google.endpoint"
                        touched={touched}
                        errors={errors}
                        helperText="example, gcr.io"
                    >
                        <TextInput
                            type="text"
                            id="config.google.endpoint"
                            value={values.config.google.endpoint}
                            onChange={onChange}
                            onBlur={handleBlur}
                            isDisabled={!isEditable}
                        />
                    </FormLabelGroup>
                    <FormLabelGroup
                        label="Project"
                        isRequired
                        fieldId="config.google.project"
                        touched={touched}
                        errors={errors}
                    >
                        <TextInput
                            type="text"
                            id="config.google.project"
                            value={values.config.google.project}
                            onChange={onChange}
                            onBlur={handleBlur}
                            isDisabled={!isEditable}
                        />
                    </FormLabelGroup>
                    {!isCreating && isEditable && (
                        <FormLabelGroup
                            fieldId="updatePassword"
                            helperText="Enable this option to replace currently stored credentials (if any)"
                            errors={errors}
                        >
                            <Checkbox
                                id="updatePassword"
                                label="Update stored credentials"
                                isChecked={values.updatePassword}
                                onChange={onUpdateCredentialsChange}
                                onBlur={handleBlur}
                                isDisabled={!isEditable}
                            />
                        </FormLabelGroup>
                    )}
                    <FormLabelGroup
                        isRequired={values.updatePassword}
                        label="Service account key (JSON)"
                        fieldId="config.google.serviceAccount"
                        touched={touched}
                        errors={errors}
                    >
                        <TextArea
                            isRequired={values.updatePassword}
                            id="config.google.serviceAccount"
                            value={values.config.google.serviceAccount}
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
                    <FormLabelGroup
                        fieldId="config.skipTestIntegration"
                        touched={touched}
                        errors={errors}
                    >
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

export default GoogleIntegrationForm;
