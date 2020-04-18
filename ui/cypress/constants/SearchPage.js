import * as api from './apiEndpoints';

export const selectors = {
    panelHeader: 'div[data-testid="panel"]',
    searchBtn: 'button:contains("Search")',
    pageSearchSuggestions: 'div.Select-menu-outer',
    categoryTabs: '[data-testid="tab"]',
    searchInput: '.search-modal input',
    pageSearchInput: '[data-testid="page-header"] .react-select__input > input',
    searchOptions: '.react-select__option',
    searchResultsHeader: '.bg-base-100.flex-1 > .text-xl',
    viewOnViolationsChip:
        'div.rt-tbody > .rt-tr-group:first-child .rt-tr .rt-td:nth-child(3) ul > li:first-child > button',
    viewOnRiskChip:
        'div.rt-tbody > .rt-tr-group:nth-child(2) .rt-tr .rt-td:nth-child(3) ul > li:first-child > button',
    viewOnPoliciesChip:
        'div.rt-tbody > .rt-tr-group:nth-child(3) .rt-tr .rt-td:nth-child(3) ul > li:first-child > button ',
    viewOnImagesChip:
        'div.rt-tbody > .rt-tr-group:nth-child(4) .rt-tr .rt-td:nth-child(3) ul > li:first-child > button'
};

export const operations = {
    enterPageSearch: (searchObj, inputSelector = selectors.pageSearchInput) => {
        cy.route(api.search.autocomplete).as('searchAutocomplete');
        function selectSearchOption(optionText) {
            // typing is slow, assuming we'll get autocomplete results, select them
            // also, likely it'll mimic better typical user's behavior
            cy.get(inputSelector).type(`${optionText.charAt(0)}`);
            cy.wait('@searchAutocomplete');
            cy.get(selectors.searchOptions)
                .contains(optionText)
                .first()
                .click({ force: true });
        }

        Object.entries(searchObj).forEach(([searchCategory, searchValue]) => {
            selectSearchOption(searchCategory);

            if (Array.isArray(searchValue)) {
                searchValue.forEach(val => selectSearchOption(val));
            } else {
                selectSearchOption(searchValue);
            }
        });
        cy.get(inputSelector).blur(); // remove focus to close the autocomplete popup
    }
};
