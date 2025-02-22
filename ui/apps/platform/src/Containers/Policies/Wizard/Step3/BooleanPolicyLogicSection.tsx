import React from 'react';
import { Flex, GridItem } from '@patternfly/react-core';
import { useFormikContext } from 'formik';

import { Policy } from 'types/policy.proto';
import useFeatureFlags from 'hooks/useFeatureFlags';
import { getCriteriaAllowedByLifecycle } from 'Containers/Policies/policies.utils';
import { policyConfigurationDescriptor, auditLogDescriptor } from './policyCriteriaDescriptors';
import PolicySection from './PolicySection';

import './BooleanPolicyLogicSection.css';

type BooleanPolicyLogicSectionProps = {
    readOnly?: boolean;
};

function BooleanPolicyLogicSection({ readOnly = false }: BooleanPolicyLogicSectionProps) {
    const { values } = useFormikContext<Policy>();
    const { isFeatureFlagEnabled } = useFeatureFlags();

    const unfilteredDescriptors =
        values.eventSource === 'AUDIT_LOG_EVENT'
            ? auditLogDescriptor
            : policyConfigurationDescriptor;
    const descriptors = unfilteredDescriptors.filter((unfilteredDescriptor) => {
        if (typeof unfilteredDescriptor.featureFlagDependency === 'string') {
            return isFeatureFlagEnabled(unfilteredDescriptor.featureFlagDependency);
        }
        return true;
    });
    const descriptorsFilteredByLifecycle = getCriteriaAllowedByLifecycle(
        descriptors,
        values.lifecycleStages
    );

    return (
        <>
            {values.policySections?.map((_, sectionIndex) =>
                readOnly ? (
                    // eslint-disable-next-line react/no-array-index-key
                    <React.Fragment key={sectionIndex}>
                        {/* this grid item takes up the default 5 columns specified in the Grid component in PolicyDetailContent */}
                        <GridItem>
                            <PolicySection
                                sectionIndex={sectionIndex}
                                descriptors={descriptorsFilteredByLifecycle}
                                readOnly={readOnly}
                            />
                        </GridItem>
                        {/* this grid item takes up 1 column specified here so that two policy sections & OR dividers can fit in one row */}
                        <GridItem lg={1}>
                            {sectionIndex !== values.policySections.length - 1 && (
                                <Flex
                                    alignSelf={{ default: 'alignSelfCenter' }}
                                    alignItems={{ default: 'alignItemsCenter' }}
                                    direction={{ default: 'row', lg: 'column' }}
                                    flexWrap={{ default: 'nowrap' }}
                                    spaceItems={{ default: 'spaceItemsSm' }}
                                    className="or-divider-container"
                                >
                                    <div className="or-divider" />
                                    <div className="pf-u-align-self-center">OR</div>
                                    <div className="or-divider" />
                                </Flex>
                            )}
                        </GridItem>
                    </React.Fragment>
                ) : (
                    // eslint-disable-next-line react/no-array-index-key
                    <React.Fragment key={sectionIndex}>
                        <PolicySection
                            sectionIndex={sectionIndex}
                            descriptors={descriptorsFilteredByLifecycle}
                            readOnly={readOnly}
                        />
                        {sectionIndex !== values.policySections.length - 1 && (
                            <Flex
                                alignSelf={{ default: 'alignSelfCenter' }}
                                alignItems={{ default: 'alignItemsCenter' }}
                                direction={{ default: 'row', lg: 'column' }}
                                flexWrap={{ default: 'nowrap' }}
                                spaceItems={{ default: 'spaceItemsSm' }}
                                className="or-divider-container"
                            >
                                <div className="or-divider" />
                                <div className="pf-u-align-self-center">OR</div>
                                <div className="or-divider" />
                            </Flex>
                        )}
                    </React.Fragment>
                )
            )}
        </>
    );
}

export default BooleanPolicyLogicSection;
