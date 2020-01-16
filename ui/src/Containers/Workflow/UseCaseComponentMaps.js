import VulnMgmtList from 'Containers/VulnMgmt/List/VulnMgmtList';
import VulnMgmtEntity from 'Containers/VulnMgmt/Entity/VulnMgmtEntity';
import VulnMgmtNavHeader from 'Containers/VulnMgmt/Components/VulnMgmtNavHeader';
import useCaseTypes from 'constants/useCaseTypes';

export const NavHeaderComponentMap = {
    [useCaseTypes.VULN_MANAGEMENT]: VulnMgmtNavHeader
};

export const ListComponentMap = {
    [useCaseTypes.VULN_MANAGEMENT]: VulnMgmtList
};

export const EntityComponentMap = {
    [useCaseTypes.VULN_MANAGEMENT]: VulnMgmtEntity
};
