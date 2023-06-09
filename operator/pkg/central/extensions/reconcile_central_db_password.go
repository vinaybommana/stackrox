package extensions

import (
	"context"
	"strings"

	"github.com/go-logr/logr"
	"github.com/operator-framework/helm-operator-plugins/pkg/extensions"
	"github.com/pkg/errors"
	platform "github.com/stackrox/rox/operator/apis/platform/v1alpha1"
	"github.com/stackrox/rox/pkg/renderer"
	coreV1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrlClient "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	centralDBPasswordKey = `password`

	// canonicalCentralDBPasswordSecretName is the name of the secret that is mounted into Central (and Central DB).
	// This is not configurable; if a user specifies a different password secret, the password from that needs to be
	// mirrored into the canonical password secret.
	canonicalCentralDBPasswordSecretName = `central-db-password`
)

// ReconcileCentralDBPasswordExtension returns an extension that takes care of reconciling the central-db-password secret.
func ReconcileCentralDBPasswordExtension(client ctrlClient.Client) extensions.ReconcileExtension {
	return wrapExtension(func(ctx context.Context, central *platform.Central, client ctrlClient.Client, statusUpdater func(statusFunc updateStatusFunc), log logr.Logger) error {
		return reconcileCentralDBPassword(ctx, central, client)
	}, client)
}

func reconcileCentralDBPassword(ctx context.Context, c *platform.Central, client ctrlClient.Client) error {

	var (
		err                  error
		hasReferencedSecret  bool
		referencedSecretName string
		isExternalDB         bool
		password             = renderer.CreatePassword()
	)

	if c.Spec.Central != nil {
		isExternalDB = c.Spec.Central.IsExternalDB()
	}

	if c.Spec.Central != nil && c.Spec.Central.DB != nil && c.Spec.Central.DB.PasswordSecret != nil {
		hasReferencedSecret = true
		referencedSecretName = c.Spec.Central.DB.PasswordSecret.Name
	}

	if !hasReferencedSecret && isExternalDB {
		return errors.New("spec.central.db.passwordSecret must be set when using an external database")
	}

	if hasReferencedSecret {
		if len(referencedSecretName) == 0 {
			return errors.New("central.db.passwordSecret.name must be set")
		}
		var referencedSecret coreV1.Secret
		if err := client.Get(ctx, ctrlClient.ObjectKey{Namespace: c.GetNamespace(), Name: referencedSecretName}, &referencedSecret); err != nil {
			return errors.Wrapf(err, "failed to get spec.central.db.passwordSecret %q", referencedSecretName)
		}

		password, err = validateCentralDBPassword(referencedSecret)
		if err != nil {
			return errors.Wrapf(err, "reading central db password from secret %s", canonicalCentralDBPasswordSecretName)
		}

		// if the referenced secret name == canonical secret name, we don't need to do anything.
		if hasReferencedSecret && referencedSecretName == canonicalCentralDBPasswordSecretName {
			return nil
		}
	}

	// we need to create or update the canonical secret
	var secret coreV1.Secret
	if err := ctrlClient.IgnoreNotFound(client.Get(ctx, ctrlClient.ObjectKey{Namespace: c.GetNamespace(), Name: canonicalCentralDBPasswordSecretName}, &secret)); err != nil {
		return errors.Wrapf(err, "failed to get central-db-password secret %q", canonicalCentralDBPasswordSecretName)
	}
	if secret.Name == "" {
		// secret doesn't exist, so we need to create it
		// we do not set the owner reference, because this password is bound to the lifetime of the PVC which we might
		// not be managing. For security, we do not want to delete the password when the Central instance is deleted.
		secret = coreV1.Secret{
			ObjectMeta: v1.ObjectMeta{
				Name:      canonicalCentralDBPasswordSecretName,
				Namespace: c.GetNamespace(),
			},
			Data: map[string][]byte{
				centralDBPasswordKey: []byte(password),
			},
		}
		if err := client.Create(ctx, &secret); err != nil {
			return errors.Wrapf(err, "failed to create central-db-password secret %q", canonicalCentralDBPasswordSecretName)
		}

	} else {

		// make sure that the secret owner reference is unset. This is to ensure that PVCs which are not deleted
		// when Centrals are deleted do not have their passwords deleted.
		shouldUpdateOwnerReference := false
		centralOwnerRef := v1.NewControllerRef(c, c.GroupVersionKind())
		centralGK := c.GroupVersionKind().GroupKind()
		for i := len(secret.OwnerReferences) - 1; i >= 0; i-- {
			ownerRef := secret.OwnerReferences[i]
			ownerGV, err := schema.ParseGroupVersion(ownerRef.APIVersion)
			if err != nil {
				continue
			}
			ownerGK := ownerGV.WithKind(ownerRef.Kind).GroupKind()
			if ownerRef.UID == centralOwnerRef.UID &&
				ownerRef.Name == centralOwnerRef.Name &&
				ownerGK == centralGK {
				secret.OwnerReferences = append(secret.OwnerReferences[:i], secret.OwnerReferences[i+1:]...)
				shouldUpdateOwnerReference = true
			}
		}

		// check if the password needs to be updated
		shouldUpdatePassword := false
		passwordIsEmpty := secret.Data == nil || len(secret.Data[centralDBPasswordKey]) == 0
		passwordsAreDifferent := secret.Data == nil || string(secret.Data[centralDBPasswordKey]) != password

		if passwordIsEmpty || hasReferencedSecret && passwordsAreDifferent {
			shouldUpdatePassword = true
		}

		if shouldUpdatePassword {
			// update the password
			secret.Data = map[string][]byte{
				centralDBPasswordKey: []byte(password),
			}
		}

		if shouldUpdateOwnerReference || shouldUpdatePassword {
			if err := client.Update(ctx, &secret); err != nil {
				return errors.Wrapf(err, "failed to update central-db-password secret %q", canonicalCentralDBPasswordSecretName)
			}
		}
	}

	return nil
}

func validateCentralDBPassword(secret coreV1.Secret) (string, error) {
	if secret.Data == nil {
		return "", errors.Errorf("secret %q does not contain a %q entry", secret.Name, centralDBPasswordKey)
	}
	passwordBytes, ok := secret.Data[centralDBPasswordKey]
	if !ok {
		return "", errors.Errorf("secret %q does not contain a %q entry", secret.Name, centralDBPasswordKey)
	}
	password := strings.TrimSpace(string(passwordBytes))
	if len(password) == 0 {
		return "", errors.Errorf("secret %q contains an empty %q entry", secret.Name, centralDBPasswordKey)
	}
	if strings.ContainsAny(password, "\r\n") {
		return "", errors.Errorf("secret %q contains a multi-line %q entry", secret.Name, centralDBPasswordKey)
	}
	return password, nil
}
