import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { createStructuredSelector } from 'reselect';

import { useTheme } from 'Containers/ThemeProvider';
import Message from 'Components/Message';
import { selectors } from 'reducers';
import AppWrapper from '../AppWrapper';

function closeThisWindow() {
    window.close();
}

function TestLoginResultsPage({ authProviderTestResults }) {
    const { isDarkMode } = useTheme();

    if (!authProviderTestResults) {
        closeThisWindow();
    }

    const userAttributes = Object.entries(authProviderTestResults.userAttributes);
    const displayedAttributes = userAttributes.map(([key, value]) => `${key}: ${value}`).join(', ');

    return (
        <AppWrapper>
            <section
                className={`flex flex-col items-center justify-center h-full py-5 ${
                    isDarkMode ? 'bg-base-0' : 'bg-primary-800'
                } `}
            >
                <div className="flex flex-col items-center bg-base-100 w-4/5 relative">
                    <Message
                        type="info"
                        message={
                            <div>
                                <h3>Authentication successful</h3>
                                <p>User ID: {authProviderTestResults?.userID}</p>
                                <p>User attributes: {displayedAttributes}</p>
                            </div>
                        }
                    />
                    <div className="flex p-4">
                        <button
                            type="button"
                            className="btn-icon btn-tertiary whitespace-no-wrap h-10 ml-4"
                            onClick={closeThisWindow}
                            dataTestId="button-close-window"
                        >
                            Close Window
                        </button>
                    </div>
                </div>
            </section>
        </AppWrapper>
    );
}

TestLoginResultsPage.propTypes = {
    authProviderTestResults: PropTypes.shape({
        userID: PropTypes.string,
        userAttributes: PropTypes.shape({}),
    }),
};

TestLoginResultsPage.defaultProps = {
    authProviderTestResults: null,
};

const mapStateToProps = createStructuredSelector({
    authProviderTestResults: selectors.getLoginAuthProviderTestResults,
});

export default connect(mapStateToProps, null)(TestLoginResultsPage);
