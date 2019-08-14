import React, { useContext } from 'react';
import entityTypes from 'constants/entityTypes';
import dateTimeFormat from 'constants/dateTimeFormat';
import { format } from 'date-fns';
import Query from 'Components/ThrowingQuery';
import Loader from 'Components/Loader';
import PageNotFound from 'Components/PageNotFound';
import CollapsibleSection from 'Components/CollapsibleSection';
import RelatedEntity from 'Containers/ConfigManagement/Entity/widgets/RelatedEntity';
import RelatedEntityListCount from 'Containers/ConfigManagement/Entity/widgets/RelatedEntityListCount';
import Metadata from 'Containers/ConfigManagement/Entity/widgets/Metadata';
import Rules from 'Containers/ConfigManagement/Entity/widgets/Rules';
import RulePermissions from 'Containers/ConfigManagement/Entity/widgets/RulePermissions';
import gql from 'graphql-tag';
import queryService from 'modules/queryService';
import searchContext from 'Containers/searchContext';
import { entityComponentPropTypes, entityComponentDefaultProps } from 'constants/entityPageProps';
import { SUBJECT_WITH_CLUSTER_FRAGMENT } from 'queries/subject';
import { SERVICE_ACCOUNT_FRAGMENT } from 'queries/serviceAccount';
import getSubListFromEntity from '../List/utilities/getSubListFromEntity';
import EntityList from '../List/EntityList';

const Role = ({ id, entityListType, query }) => {
    const searchParam = useContext(searchContext);

    const variables = {
        id,
        where: queryService.objectToWhereClause(query[searchParam])
    };

    const QUERY = gql`
        query k8sRole($id: ID!) {
            clusters {
                id
                k8srole(role: $id) {
                    id
                    name
                    type
                    verbs
                    createdAt
                    roleNamespace {
                        metadata {
                            id
                            name
                        }
                    }
                    ${
                        entityListType === entityTypes.SERVICE_ACCOUNT
                            ? 'serviceAccounts {...serviceAccountFields}'
                            : 'serviceAccountCount'
                    }
                    ${
                        entityListType === entityTypes.SUBJECT
                            ? 'subjects {...subjectWithClusterFields}'
                            : 'subjectCount'
                    }
                    rules {
                        apiGroups
                        nonResourceUrls
                        resourceNames
                        resources
                        verbs
                    }
                    clusterName
                    clusterId
                }
            }
        }

    ${entityListType === entityTypes.SUBJECT ? SUBJECT_WITH_CLUSTER_FRAGMENT : ''}
    ${entityListType === entityTypes.SERVICE_ACCOUNT ? SERVICE_ACCOUNT_FRAGMENT : ''}


    `;
    return (
        <Query query={QUERY} variables={variables}>
            {({ loading, data }) => {
                if (loading) return <Loader />;
                const { clusters } = data;
                if (!clusters || !clusters.length)
                    return <PageNotFound resourceType={entityTypes.ROLE} />;

                const { k8srole: entity } = clusters[0];

                if (entityListType) {
                    return (
                        <EntityList
                            entityListType={entityListType}
                            data={getSubListFromEntity(entity, entityListType)}
                            query={query}
                        />
                    );
                }

                const {
                    type,
                    createdAt,
                    roleNamespace,
                    serviceAccountCount,
                    subjectCount,
                    labels = [],
                    annotations = [],
                    rules,
                    clusterName,
                    clusterId
                } = entity;
                const { name: namespaceName, id: namespaceId } = roleNamespace
                    ? roleNamespace.metadata
                    : {};

                const metadataKeyValuePairs = [
                    { key: 'Role Type', value: type },
                    {
                        key: 'Created',
                        value: createdAt ? format(createdAt, dateTimeFormat) : 'N/A'
                    }
                ];

                return (
                    <div className="bg-primary-100 w-full">
                        <CollapsibleSection title="Role Details">
                            <div className="flex mb-4 flex-wrap">
                                <Metadata
                                    className="mx-4 bg-base-100 h-48 mb-4"
                                    keyValuePairs={metadataKeyValuePairs}
                                    labels={labels}
                                    annotations={annotations}
                                />
                                <RelatedEntity
                                    className="mx-4 min-w-48 h-48 mb-4"
                                    entityType={entityTypes.CLUSTER}
                                    name="Cluster"
                                    value={clusterName}
                                    entityId={clusterId}
                                />
                                {roleNamespace && (
                                    <RelatedEntity
                                        className="mx-4 min-w-48 h-48 mb-4"
                                        entityType={entityTypes.NAMESPACE}
                                        name="Namespace Scope"
                                        value={namespaceName}
                                        entityId={namespaceId}
                                    />
                                )}
                                <RelatedEntityListCount
                                    className="mx-4 min-w-48 h-48 mb-4"
                                    name="Users & Groups"
                                    value={subjectCount}
                                    entityType={entityTypes.SUBJECT}
                                />
                                <RelatedEntityListCount
                                    className="mx-4 min-w-48 h-48 mb-4"
                                    name="Service Accounts"
                                    value={serviceAccountCount}
                                    entityType={entityTypes.SERVICE_ACCOUNT}
                                />
                            </div>
                        </CollapsibleSection>
                        <CollapsibleSection title="Role Permissions And Rules">
                            <div className="flex mb-4">
                                <RulePermissions rules={rules} className="mx-4 bg-base-100" />
                                <Rules rules={rules} className="mx-4 bg-base-100" />
                            </div>
                        </CollapsibleSection>
                    </div>
                );
            }}
        </Query>
    );
};

Role.propTypes = entityComponentPropTypes;
Role.defaultProps = entityComponentDefaultProps;

export default Role;
