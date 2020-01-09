import { saveFile } from 'services/DownloadService';
import entityTypes from 'constants/entityTypes';
import axios from './instance';

const baseUrl = '/v1/cves';
const csvUrl = '/api/vm/export/csv';

/**
 * Send request to suppress / unsuppress CVE with a given ID.
 *
 * @param {!string} CVE unique identifier
 * @param {!boolean} true if CVE should be suppressed, false for unsuppress
 * @returns {Promise<AxiosResponse, Error>} fulfilled in case of success or rejected with an error
 */
export function updateCveSuppressedState(cve, suppressed = false) {
    return axios.patch(`${baseUrl}/${cve}`, { suppressed });
}

export function getCvesInCsvFormat(fileName, searchParamsList) {
    const searchString = searchParamsList
        .map(item => {
            return `${item.key}=${item.value}`;
        })
        .join('&');

    const url = searchString ? `${csvUrl}?${searchString}` : csvUrl;

    return saveFile({
        method: 'get',
        url,
        data: null,
        name: `${fileName}.csv`
    });
}

const searchFields = {
    [entityTypes.CLUSTER]: 'Cluster+ID',
    [entityTypes.COMPONENT]: 'Component+ID',
    [entityTypes.DEPLOYMENT]: 'Deployment+ID',
    [entityTypes.NAMESPACE]: 'Namespace+ID',
    [entityTypes.IMAGE]: 'Image+ID'
};

export function exportCvesAsCsv(fileName, workflowState) {
    const pageStack = workflowState.getPageStack();

    const searchParamsList = [];
    if (pageStack.length > 1) {
        const parentEntity = pageStack[pageStack.length - 2];
        const searchEntity = parentEntity.t;
        const id = parentEntity.i;
        searchParamsList.push({ key: searchFields[searchEntity], value: id });
    }
    // @TODO: add other search params from the workflow state, if present

    return getCvesInCsvFormat(fileName, searchParamsList);
}
