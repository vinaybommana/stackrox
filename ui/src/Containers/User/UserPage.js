import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { createStructuredSelector } from 'reselect';

import { selectors } from 'reducers';
import PageHeader from 'Components/PageHeader';
import getUserAttributeMap from 'utils/userDataUtils';

const UserPage = ({ userData }) => {
    const { userAttributes } = userData;
    const userAttributeMap = getUserAttributeMap(userAttributes);
    const header = userAttributeMap.name;
    const subHeader = userAttributeMap.email;
    return (
        <section className="flex flex-1 h-full w-full">
            <div className="flex flex-1 flex-col w-full">
                <PageHeader header={header} subHeader={subHeader} capitalize={false} />
                <div className="flex-1 relative p-6 xxxl:p-8">
                    <div
                        className="grid grid-gap-6 xxxl:grid-gap-8 md:grid-auto-fit xxl:grid-auto-fit-wide md:grid-dense"
                        style={{ '--min-tile-height': '160px' }}
                    />
                </div>
            </div>
        </section>
    );
};

UserPage.propTypes = {
    userData: PropTypes.shape({
        userAttributes: PropTypes.arrayOf(PropTypes.shape({})),
    }).isRequired,
};

const mapStateToProps = createStructuredSelector({
    userData: selectors.getCurrentUser,
});

export default connect(mapStateToProps, null)(UserPage);
