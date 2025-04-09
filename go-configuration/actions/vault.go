package actions

import (
	"context"
	"fmt"
	"strings"

	appConf "github.com/adeo/lmru--devops--argocd-apps/go-configuration/config"

	vault "github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/approle"
)

func GetSecretKVv2WithAppRole(secretPathName, secretVaultName, vaultNamespace string) (map[string]interface{}, error) {
	config := vault.DefaultConfig()

	var returnSecret map[string]interface{}

	client, err := vault.NewClient(config)
	if err != nil {
		return returnSecret, fmt.Errorf("unable to initialize Vault client: %w", err)
	}

	if appConf.RoleID == "" {
		return returnSecret, fmt.Errorf("no role ID was provided in VAULT_APPROLE_ROLE_ID env var")
	}

	client.SetNamespace(vaultNamespace)

	appRoleAuth, err := auth.NewAppRoleAuth(
		appConf.RoleID,
		appConf.SecretID,
		auth.WithMountPath("approle"),
	)
	if err != nil {
		return returnSecret, fmt.Errorf("unable to initialize AppRole auth method: %w", err)
	}

	authInfo, err := client.Auth().Login(context.Background(), appRoleAuth)
	if err != nil {
		return returnSecret, fmt.Errorf("unable to login to AppRole auth method: %w", err)
	}
	if authInfo == nil {
		return returnSecret, fmt.Errorf("no auth info was returned after login")
	}

	getSecret, err := client.KVv2(secretPathName).Get(context.TODO(), secretVaultName)
	if err != nil {
		return returnSecret, fmt.Errorf("%w\n", err)
	}

	return getSecret.Data, nil
}

func PutSecretWithAppRole(secretPathName, secretVaultName, vaultNamespace, vaultRoleID, vaultSecretID, method string, putSecret map[string]interface{}) error {
	config := vault.DefaultConfig()

	client, err := vault.NewClient(config)
	if err != nil {
		return fmt.Errorf("unable to initialize Vault client: %w", err)
	}

	client.SetNamespace(vaultNamespace)

	appRoleAuth, err := auth.NewAppRoleAuth(
		vaultRoleID,
		&auth.SecretID{FromString: vaultSecretID},
		auth.WithMountPath("approle"),
	)
	if err != nil {
		return fmt.Errorf("unable to initialize AppRole auth method: %w", err)
	}

	authInfo, err := client.Auth().Login(context.Background(), appRoleAuth)
	if err != nil {
		return fmt.Errorf("FUNC PutSecretWithAppRole - unable to login to AppRole auth method: %w", err)
	}
	if authInfo == nil {
		return fmt.Errorf("no auth info was returned after login")
	}

	// secretPath := secretPathName + "/data/" + secretVaultName

	getSecreet, err := client.KVv2(secretPathName).Get(context.Background(), secretVaultName)

	if err != nil {
		if !strings.Contains(fmt.Sprint(err), "secret not found") {
			return fmt.Errorf("Secret %s is Found in Vault. Check it.\n", secretVaultName)
		}
	}

	switch method {
	case "put":
		if getSecreet != nil {
			return fmt.Errorf("Secret %s is Found in Namespace %s. Check it.\n", secretVaultName, vaultNamespace)
		}
		_, err = client.KVv2(secretPathName).Put(context.Background(), secretVaultName, putSecret)
		// _, err = client.Logical().Write(secretPath, putSecret)
		if err != nil {
			return fmt.Errorf("unable to PUT secret: %w", err)
		}
	case "update":
		_, err = client.KVv2(secretPathName).Patch(context.Background(), secretVaultName, putSecret)
		if err != nil {
			return fmt.Errorf("unable to UPDATE secret: %w", err)
		}
	}

	return nil
}
