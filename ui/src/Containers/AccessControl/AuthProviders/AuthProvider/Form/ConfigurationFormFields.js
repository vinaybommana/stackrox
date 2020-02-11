import React from 'react';
import PropTypes from 'prop-types';

import Labeled from 'Components/Labeled';
import ReduxTextField from 'Components/forms/ReduxTextField';
import ReduxSelectField from 'Components/forms/ReduxSelectField';
import ReduxTextAreaField from 'Components/forms/ReduxTextAreaField';

const baseURL = `${window.location.protocol}//${window.location.host}`;
const oidcFragmentCallbackURL = `${baseURL}/auth/response/oidc`;
const oidcPostCallbackURL = `${baseURL}/sso/providers/oidc/callback`;
const samlACSURL = `${baseURL}/sso/providers/saml/acs`;

const CommonFields = ({ disabled }) => (
    <>
        <Labeled label="Integration Name">
            <ReduxTextField
                name="name"
                placeholder="Name for this integration"
                disabled={disabled}
            />
        </Labeled>
    </>
);

const Note = ({ header, children }) => (
    <div className="w-full mb-5">
        <div className="text-tertiary-800 bg-tertiary-200 p-3 pb-2 rounded border-2 border-tertiary-300 ">
            <p className="border-b-2 border-tertiary-300 pb-3">
                <strong>Note: </strong> {header}
            </p>
            {children}
        </div>
    </div>
);

const OidcFormFields = ({ disabled }) => (
    <>
        <CommonFields disabled={disabled} />
        <Labeled label="Callback Mode">
            <ReduxSelectField
                name="config.mode"
                options={[
                    { value: 'fragment', label: 'Fragment' },
                    { value: 'post', label: 'HTTP POST' }
                ]}
                defaultValue="post"
                disabled={disabled}
            />
        </Labeled>
        <Labeled label="Issuer">
            <ReduxTextField
                name="config.issuer"
                placeholder="tenant.auth-provider.com"
                disabled={disabled}
            />
        </Labeled>
        <Labeled label="Client ID">
            <ReduxTextField name="config.client_id" disabled={disabled} />
        </Labeled>
        <Note header="if required by your IdP, use the following callback URLs:">
            <ul className="pl-4 mt-2 leading-loose">
                <li>
                    For <span className="font-700">Fragment</span> mode:{' '}
                    <a
                        className="text-tertiary-800 hover:text-tertiary-900"
                        href={oidcFragmentCallbackURL}
                    >
                        {oidcFragmentCallbackURL}
                    </a>
                </li>
                <li>
                    For <span className="font-700">HTTP POST</span> mode:{' '}
                    <a
                        className="text-tertiary-800 hover:text-tertiary-900"
                        href={oidcPostCallbackURL}
                    >
                        {oidcPostCallbackURL}
                    </a>
                </li>
            </ul>
        </Note>
    </>
);

const Auth0FormFields = ({ disabled }) => (
    <>
        <CommonFields />
        <Labeled label="Auth0 Tenant">
            <ReduxTextField
                name="config.issuer"
                placeholder="your-tenant.auth0.com"
                disabled={disabled}
            />
        </Labeled>
        <Labeled label="Client ID">
            <ReduxTextField name="config.client_id" disabled={disabled} />
        </Labeled>
        <Note header="if required by your IdP, use the following callback URL:">
            <ul className="pl-4 mt-2 leading-loose">
                <li>
                    <a
                        className="text-tertiary-800 hover:text-tertiary-900"
                        href={oidcFragmentCallbackURL}
                    >
                        {oidcFragmentCallbackURL}
                    </a>
                </li>
            </ul>
        </Note>
    </>
);

const SamlFormFields = ({ disabled }) => (
    <>
        <CommonFields />
        <Labeled label="ServiceProvider Issuer">
            <ReduxTextField
                name="config.sp_issuer"
                placeholder="https://prevent.stackrox.io/"
                disabled={disabled}
            />
        </Labeled>
        <div className="w-full mb-5">
            <div className="border-b border-base-400 border-dotted flex pb-2">
                Option 1: Dynamic Configuration
            </div>
        </div>
        <Labeled label="IdP Metadata URL">
            <ReduxTextField
                name="config.idp_metadata_url"
                placeholder="https://idp.example.com/metadata"
                disabled={disabled}
            />
        </Labeled>
        <div className="w-full mb-5">
            <div className="border-b border-base-400 border-dotted flex pb-2">
                Option 2: Static Configuration
            </div>
        </div>
        <Labeled label="IdP Issuer">
            <ReduxTextField
                name="config.idp_issuer"
                placeholder="https://idp.example.com/"
                disabled={disabled}
            />
        </Labeled>
        <Labeled label="IdP SSO URL">
            <ReduxTextField
                name="config.idp_sso_url"
                placeholder="https://idp.example.com/login"
                disabled={disabled}
            />
        </Labeled>
        <Labeled label="Name/ID Format">
            <ReduxTextField
                name="config.idp_nameid_format"
                placeholder="urn:oasis:names:tc:SAML:1.1:nameid-format:persistent"
                disabled={disabled}
            />
        </Labeled>
        <Labeled label="IdP Certificate (PEM)">
            <ReduxTextAreaField
                name="config.idp_cert_pem"
                placeholder={
                    '-----BEGIN CERTIFICATE-----\nYour certificate data\n-----END CERTIFICATE-----'
                }
                disabled={disabled}
            />
        </Labeled>
        <Note header="if required by your IdP, use the following Assertion Consumer Service (ACS) URL:">
            <ul className="pl-4 mt-2 leading-loose">
                <li>
                    <a className="text-tertiary-800 hover:text-tertiary-900" href={samlACSURL}>
                        {samlACSURL}
                    </a>
                </li>
            </ul>
        </Note>
    </>
);

const UserPkiFormFields = ({ disabled }) => (
    <>
        <CommonFields />
        <Labeled label="CA Certificates (PEM)">
            <ReduxTextAreaField
                name="config.keys"
                placeholder={
                    '-----BEGIN CERTIFICATE-----\nAuthority certificate data\n-----END CERTIFICATE-----'
                }
                disabled={disabled}
            />
        </Labeled>
    </>
);

const formFieldsComponents = {
    oidc: OidcFormFields,
    auth0: Auth0FormFields,
    saml: SamlFormFields,
    userpki: UserPkiFormFields
};

const ConfigurationFormFields = ({ providerType, disabled }) => {
    const FormFieldsComponent = formFieldsComponents[providerType];
    if (!FormFieldsComponent)
        throw new Error(`Unknown auth provider type passed to the form component: ${providerType}`);

    return <FormFieldsComponent disabled={disabled} />;
};
ConfigurationFormFields.propTypes = {
    providerType: PropTypes.oneOf(Object.keys(formFieldsComponents)).isRequired,
    disabled: PropTypes.bool.isRequired
};

export default ConfigurationFormFields;
