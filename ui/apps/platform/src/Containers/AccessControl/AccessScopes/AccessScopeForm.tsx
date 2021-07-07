import React, { ReactElement, useEffect, useState } from 'react';
import { useFormik } from 'formik';
import * as yup from 'yup';
import {
    Alert,
    AlertVariant,
    Button,
    Flex,
    FlexItem,
    Form,
    FormGroup,
    TextInput,
    Title,
    Toolbar,
    ToolbarContent,
    ToolbarGroup,
    ToolbarItem,
    Tooltip,
} from '@patternfly/react-core';
import { OutlinedQuestionCircleIcon } from '@patternfly/react-icons';

import {
    AccessScope,
    EffectiveAccessScopeCluster,
    LabelSelector,
    LabelSelectorsKey,
    computeEffectiveAccessScopeClusters,
    getIsValidRules,
} from 'services/RolesService';

import { AccessControlQueryAction } from '../accessControlPaths';

import { LabelSelectorsEditingState } from './accessScopes.utils';
import EffectiveAccessScopeTable from './EffectiveAccessScopeTable';
import LabelInclusion from './LabelInclusion';

const labelIconEffectiveAccessScope = (
    <Tooltip
        content={
            <div>
                Computed <strong>state</strong> of clusters and namespaces
                <br />
                from <strong>manual</strong> inclusion, <strong>label</strong> inclusion,
                <br />
                or <strong>hierarchical</strong> inclusion:
                <br />
                included cluster: therefore all of its namespaces
                <br />
                included namespace: therefore its cluster
            </div>
        }
        isContentLeftAligned
        maxWidth="24em"
    >
        <div className="pf-c-button pf-m-plain pf-m-small">
            <OutlinedQuestionCircleIcon />
        </div>
    </Tooltip>
);

const labelIconLabelInclusion = (
    <Tooltip
        content={
            <div>
                A label inclusion tab has label <strong>selector</strong> cards
                <br />
                At least one selector must be satisfied (s1 or s2 or s3)
            </div>
        }
        isContentLeftAligned
        maxWidth="24em"
    >
        <div className="pf-c-button pf-m-plain pf-m-small">
            <OutlinedQuestionCircleIcon />
        </div>
    </Tooltip>
);

export type AccessScopeFormProps = {
    isActionable: boolean;
    action?: AccessControlQueryAction;
    accessScope: AccessScope;
    accessScopes: AccessScope[];
    handleCancel: () => void;
    handleEdit: () => void;
    handleSubmit: (values: AccessScope) => Promise<null>; // because the form has only catch and finally
};

function AccessScopeForm({
    isActionable,
    action,
    accessScope,
    accessScopes,
    handleCancel,
    handleEdit,
    handleSubmit,
}: AccessScopeFormProps): ReactElement {
    const [counterComputing, setCounterComputing] = useState(0);
    const [alertCompute, setAlertCompute] = useState<ReactElement | null>(null);
    const [clusters, setClusters] = useState<EffectiveAccessScopeCluster[]>([]);

    const [isSubmitting, setIsSubmitting] = useState(false);
    const [alertSubmit, setAlertSubmit] = useState<ReactElement | null>(null);

    /*
     * Disable Submit button while editing a label selector.
     * Prevent simultaneous editing on both tabs.
     */
    const [
        labelSelectorsEditingState,
        setLabelSelectorsEditingState,
    ] = useState<LabelSelectorsEditingState | null>(null);

    const { dirty, errors, handleChange, isValid, resetForm, setFieldValue, values } = useFormik({
        initialValues: accessScope,
        onSubmit: () => {},
        validationSchema: yup.object({
            name: yup
                .string()
                .required()
                .test(
                    'non-unique-name',
                    'Another access scope already has this name',
                    // Return true if current input name is initial name
                    // or no other access scope already has this name.
                    (nameInput) =>
                        nameInput === accessScope.name ||
                        accessScopes.every(({ name }) => nameInput !== name)
                ),
            description: yup.string(),
        }),
    });

    /*
     * A label selector or set requirement is temporarily invalid when it is added,
     * before its first requirement or value has been added.
     */
    const isValidRules = getIsValidRules(values.rules);

    useEffect(() => {
        if (isValidRules) {
            setCounterComputing((counterPrev) => counterPrev + 1);
            computeEffectiveAccessScopeClusters(values.rules)
                .then((clustersComputed) => {
                    setClusters(clustersComputed);
                    setAlertCompute(null);
                })
                .catch((error) => {
                    setAlertCompute(
                        <Alert
                            title="Compute effective access scope failed"
                            variant={AlertVariant.danger}
                            isInline
                        >
                            {error.message}
                        </Alert>
                    );
                })
                .finally(() => {
                    setCounterComputing((counterPrev) => counterPrev - 1);
                });
        }
    }, [isValidRules, values.rules]);

    function onChange(_value, event) {
        handleChange(event);
    }

    function handleIncludedClustersChange(clusterName: string, isChecked: boolean) {
        const { includedClusters } = values.rules;
        return setFieldValue(
            'rules.includedClusters',
            isChecked
                ? [...includedClusters, clusterName]
                : includedClusters.filter(
                      (includedClusterName) => includedClusterName !== clusterName
                  )
        );
    }

    function handleIncludedNamespacesChange(
        clusterName: string,
        namespaceName: string,
        isChecked: boolean
    ) {
        const { includedNamespaces } = values.rules;
        return setFieldValue(
            'rules.includedNamespaces',
            isChecked
                ? [...includedNamespaces, { clusterName, namespaceName }]
                : includedNamespaces.filter(
                      ({
                          clusterName: includedClusterName,
                          namespaceName: includedNamespaceName,
                      }) =>
                          includedClusterName !== clusterName ||
                          includedNamespaceName !== namespaceName
                  )
        );
    }

    function handleLabelSelectorsChange(
        labelSelectorsKey: LabelSelectorsKey,
        labelSelectors: LabelSelector[]
    ) {
        return setFieldValue(`rules.${labelSelectorsKey}`, labelSelectors);
    }

    function onClickSubmit() {
        // TODO submit through Formik, especially to update its initialValue.
        // For example, to make a change, submit, and then make the opposite change.
        setIsSubmitting(true);
        setAlertSubmit(null);
        handleSubmit(values)
            .catch((error) => {
                setAlertSubmit(
                    <Alert
                        title="Failed to save access scope"
                        variant={AlertVariant.danger}
                        isInline
                    >
                        {error.message}
                    </Alert>
                );
            })
            .finally(() => {
                setIsSubmitting(false);
            });
    }

    function onClickCancel() {
        resetForm();
        handleCancel(); // close form if action=create but not if action=update
    }

    const hasAction = Boolean(action);
    const isViewing = !hasAction;

    const nameErrorMessage = values.name.length !== 0 && errors.name ? errors.name : '';
    const nameValidatedState = nameErrorMessage ? 'error' : 'default';

    return (
        <Form id="access-scope-form">
            <Toolbar inset={{ default: 'insetNone' }}>
                <ToolbarContent>
                    <ToolbarItem>
                        <Title headingLevel="h2">
                            {action === 'create' ? 'Add access scope' : accessScope.name}
                        </Title>
                    </ToolbarItem>
                    {isActionable && action !== 'create' && (
                        <ToolbarGroup variant="button-group" alignment={{ default: 'alignRight' }}>
                            <ToolbarItem>
                                <Button
                                    variant="primary"
                                    onClick={handleEdit}
                                    isDisabled={action === 'update'}
                                    isSmall
                                >
                                    Edit access scope
                                </Button>
                            </ToolbarItem>
                        </ToolbarGroup>
                    )}
                </ToolbarContent>
            </Toolbar>
            {alertSubmit}
            <FormGroup
                label="Name"
                fieldId="name"
                isRequired
                validated={nameValidatedState}
                helperTextInvalid={nameErrorMessage}
                className="pf-m-horizontal"
            >
                <TextInput
                    type="text"
                    id="name"
                    value={values.name}
                    validated={nameValidatedState}
                    onChange={onChange}
                    isDisabled={isViewing}
                    isRequired
                    className="pf-m-limit-width"
                />
            </FormGroup>
            <FormGroup label="Description" fieldId="description" className="pf-m-horizontal">
                <TextInput
                    type="text"
                    id="description"
                    value={values.description}
                    onChange={onChange}
                    isDisabled={isViewing}
                />
            </FormGroup>
            {alertCompute}
            <Flex
                direction={{ default: 'column', xl: 'row' }}
                spaceItems={{ default: 'spaceItemsNone', xl: 'spaceItemsLg' }}
                className="pf-u-pt-lg"
            >
                <FlexItem className="pf-u-flex-basis-0" flex={{ default: 'flex_1' }}>
                    <FormGroup
                        label="Allowed resources"
                        fieldId="effectiveAccessScope"
                        labelIcon={labelIconEffectiveAccessScope}
                        className="pf-u-pb-lg"
                    >
                        <EffectiveAccessScopeTable
                            counterComputing={counterComputing}
                            clusters={clusters}
                            includedClusters={values.rules.includedClusters}
                            includedNamespaces={values.rules.includedNamespaces}
                            handleIncludedClustersChange={handleIncludedClustersChange}
                            handleIncludedNamespacesChange={handleIncludedNamespacesChange}
                            hasAction={hasAction}
                        />
                    </FormGroup>
                </FlexItem>
                <FlexItem className="pf-u-flex-basis-0" flex={{ default: 'flex_1' }}>
                    <FormGroup
                        label="Label inclusion"
                        fieldId="labelInclusion"
                        labelIcon={labelIconLabelInclusion}
                        className="pf-u-pb-lg"
                    >
                        <LabelInclusion
                            clusterLabelSelectors={values.rules.clusterLabelSelectors}
                            namespaceLabelSelectors={values.rules.namespaceLabelSelectors}
                            hasAction={hasAction}
                            labelSelectorsEditingState={labelSelectorsEditingState}
                            setLabelSelectorsEditingState={setLabelSelectorsEditingState}
                            handleLabelSelectorsChange={handleLabelSelectorsChange}
                        />
                    </FormGroup>
                </FlexItem>
            </Flex>
            {hasAction && (
                <Toolbar inset={{ default: 'insetNone' }}>
                    <ToolbarContent>
                        <ToolbarGroup variant="button-group">
                            <ToolbarItem>
                                <Button
                                    variant="primary"
                                    onClick={onClickSubmit}
                                    isDisabled={
                                        !dirty ||
                                        !isValid ||
                                        !isValidRules ||
                                        Boolean(labelSelectorsEditingState) ||
                                        isSubmitting
                                    }
                                    isLoading={isSubmitting}
                                    isSmall
                                >
                                    Save
                                </Button>
                            </ToolbarItem>
                            <ToolbarItem>
                                <Button variant="tertiary" onClick={onClickCancel} isSmall>
                                    Cancel
                                </Button>
                            </ToolbarItem>
                        </ToolbarGroup>
                    </ToolbarContent>
                </Toolbar>
            )}
        </Form>
    );
}

export default AccessScopeForm;
