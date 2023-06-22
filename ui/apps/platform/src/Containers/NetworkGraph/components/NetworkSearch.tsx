import React, { useEffect, useState } from 'react';

import SearchFilterInput from 'Components/SearchFilterInput';
import useURLSearch from 'hooks/useURLSearch';
import searchOptionsToQuery from 'services/searchOptionsToQuery';
import { getSearchOptionsForCategory } from 'services/SearchService';
import { orchestratorComponentsOption } from 'utils/orchestratorComponents';
import { NetworkScopeHierarchy } from '../utils/hierarchyUtils';

import './NetworkSearch.css';

const searchCategory = 'DEPLOYMENTS';
const searchOptionExclusions = [
    'Cluster',
    'Deployment',
    'Namespace',
    'Namespace ID',
    'Orchestrator Component',
];

type NetworkSearchProps = {
    scopeHierarchy: NetworkScopeHierarchy | null;
    isDisabled: boolean;
};

function NetworkSearch({ scopeHierarchy, isDisabled }: NetworkSearchProps) {
    const [searchOptions, setSearchOptions] = useState<string[]>([]);
    const { searchFilter, setSearchFilter } = useURLSearch();

    useEffect(() => {
        const { request, cancel } = getSearchOptionsForCategory(searchCategory);
        request
            .then((options) => {
                const filteredOptions = options.filter((o) => !searchOptionExclusions.includes(o));
                setSearchOptions(filteredOptions);
            })
            .catch(() => {
                // A request error will disable the search filter.
            });

        return cancel;
    }, [setSearchOptions]);

    function onSearch(options) {
        const newOptions = { ...options };
        newOptions.Cluster = scopeHierarchy?.cluster?.name ?? '';
        newOptions.Namespace = scopeHierarchy?.namespaces ?? [];
        newOptions.Deployment = scopeHierarchy?.deployments ?? [];

        setSearchFilter(newOptions);
    }

    const prependAutocompleteQuery = [...orchestratorComponentsOption];

    return (
        <SearchFilterInput
            className="pf-u-w-100 theme-light pf-search-shim"
            placeholder="Filter deployments"
            searchFilter={searchFilter}
            searchCategory="DEPLOYMENTS"
            searchOptions={searchOptions}
            autocompleteQueryPrefix={searchOptionsToQuery(prependAutocompleteQuery)}
            handleChangeSearchFilter={onSearch}
            isDisabled={isDisabled}
        />
    );
}

export default NetworkSearch;
