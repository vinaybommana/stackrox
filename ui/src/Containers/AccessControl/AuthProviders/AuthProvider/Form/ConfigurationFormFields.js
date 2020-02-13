import React from 'react';
import PropTypes from 'prop-types';
import { formValues } from 'redux-form';

import { knownBackendFlags as featureFlags } from 'utils/featureFlags';
import FeatureEnabled from 'Containers/FeatureEnabled';
import Labeled from 'Components/Labeled';
import FormFieldLabel from 'Components/forms/FormFieldLabel';
import ReduxTextField from 'Components/forms/ReduxTextField';
import ReduxPasswordField from 'Components/forms/ReduxPasswordField';
import ReduxSelectField from 'Components/forms/ReduxSelectField';
import ReduxTextAreaField from 'Components/forms/ReduxTextAreaField';
import ReduxCheckboxField from 'Components/forms/ReduxCheckboxField';

const baseURL = `${window.location.protocol}//${window.location.host}`;
const oidcFragmentCallbackURL = `${baseURL}/auth/response/oidc`;
const oidcPostCallbackURL = `${baseURL}/sso/providers/oidc/callback`;
const samlACSURL = `${baseURL}/sso/providers/saml/acs`;

const CommonFields = ({ disabled }) => (
    <>
        <Labeled label={<FormFieldLabel text="Integration Name" required />}>
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

const OidcFormFields = ({ disabled, configValues, change }) => {
    function onModeChange(event, newValue) {
        // client secret is supported only with HTTP POST
        change('config.do_not_use_client_secret', newValue !== 'post');
        change('config.client_secret', '');
    }

    function onDoNotUseClientSecretChange(event, newValue) {
        if (newValue) change('config.client_secret', '');
    }

    const clientSecretSupported = configValues.mode === 'post';

    // use client secret placeholder as an explanation text
    let clientSecretPlaceholder = 'Client Secret provided by your IdP';
    if (!clientSecretSupported) {
        clientSecretPlaceholder = 'Client Secret is only supported with HTTP POST callback mode';
    } else if (configValues.clientOnly?.clientSecretStored) {
        clientSecretPlaceholder = configValues.do_not_use_client_secret
            ? 'Disabled, the currently stored secret will be removed'
            : 'Leave this field empty to keep the currently stored secret';
    } else if (configValues.do_not_use_client_secret) {
        clientSecretPlaceholder = 'Disabled';
    }

    // user is expected to enter the value unless opted out or ok to leave empty to preserve the old value
    const clientSecretRequired =
        !configValues.do_not_use_client_secret && !configValues.clientOnly?.clientSecretStored;
    const clientSecretLabel = (
        <FormFieldLabel text="Client Secret" required={clientSecretRequired} />
    );

    const doNotUseClientSecretDisabled = disabled || !clientSecretSupported;
    const clientSecret = (
        <>
            <Labeled label={clientSecretLabel}>
                <ReduxPasswordField
                    name="config.client_secret"
                    disabled={disabled || configValues.do_not_use_client_secret}
                    placeholder={clientSecretPlaceholder}
                />
                <div className="mt-2">
                    <ReduxCheckboxField
                        name="config.do_not_use_client_secret"
                        id="do-not-use-client-secret-checkbox"
                        disabled={doNotUseClientSecretDisabled}
                        onChange={onDoNotUseClientSecretChange}
                    />
                    <label
                        className={`ml-2 ${doNotUseClientSecretDisabled && 'text-base-500'}`}
                        htmlFor="do-not-use-client-secret-checkbox"
                    >
                        Do not use Client Secret (not recommended)
                    </label>
                </div>
            </Labeled>
        </>
    );
    return (
        <>
            <CommonFields disabled={disabled} />
            <Labeled label={<FormFieldLabel text="Callback Mode" required />}>
                <ReduxSelectField
                    name="config.mode"
                    options={[
                        { value: 'fragment', label: 'Fragment' },
                        { value: 'post', label: 'HTTP POST' }
                    ]}
                    disabled={disabled}
                    onChange={onModeChange}
                />
            </Labeled>
            <Labeled label={<FormFieldLabel text="Issuer" required />}>
                <ReduxTextField
                    name="config.issuer"
                    placeholder="tenant.auth-provider.com"
                    disabled={disabled}
                />
            </Labeled>
            <Labeled label={<FormFieldLabel text="Client ID" required />}>
                <ReduxTextField name="config.client_id" disabled={disabled} />
            </Labeled>
            <FeatureEnabled featureFlag={featureFlags.ROX_REFRESH_TOKENS}>
                {clientSecret}
            </FeatureEnabled>
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
};

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

const ConfigurationFormFields = ({ providerType, disabled, configValues, change }) => {
    const FormFieldsComponent = formFieldsComponents[providerType];
    if (!FormFieldsComponent)
        throw new Error(`Unknown auth provider type passed to the form component: ${providerType}`);

    return <FormFieldsComponent disabled={disabled} configValues={configValues} change={change} />;
};
ConfigurationFormFields.propTypes = {
    providerType: PropTypes.oneOf(Object.keys(formFieldsComponents)).isRequired,
    disabled: PropTypes.bool.isRequired,
    configValues: PropTypes.shape({}),
    change: PropTypes.func.isRequired
};

ConfigurationFormFields.defaultProps = {
    configValues: {}
};

export default formValues({ configValues: 'config' })(ConfigurationFormFields);
