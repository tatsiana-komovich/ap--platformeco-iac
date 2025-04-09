package config

import (
	"log"
	"os"

	auth "github.com/hashicorp/vault/api/auth/approle"
)

var (
	ErrorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	InfoLog  = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	Chars    = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZÅÄÖ" +
		"abcdefghijklmnopqrstuvwxyzåäö" +
		"0123456789")

	RoleID                  = os.Getenv("VAULT_RESTIC_ROLE_ID")
	SecretID                = &auth.SecretID{FromEnv: "VAULT_RESTIC_SECRET_ID"}
	BackupVaultNamespace    = "flant"
	BackupVaultSecretPath   = "restic_backups"
	ArgocdSAVaultSecretPath = "services"
	ArgocdSAVaultSecretName = "argocd-secrets"
	DevopsVaultNamespace    = "it-devops"
	// Commander
	CommanderAPI             = "https://api-commander.devops.lmru.tech"
	CommanderVaultSecretPath = "service-accounts/data"
	CommanderVaultSecretName = "commander"
	CommanderBasicAuth       = "basic-auth"
	// S3 config
	S3Region          = "ru-central1"
	S3URL             = "https://storage.yandexcloud.net"
	S3PartitionID     = "yc"
	BucketFolderName  = "k8s-backup/"
	S3SecretIDName    = "secret_id"
	S3SecretKeyName   = "secret_key"
	VaultS3SecretPath = "service-accounts"
	VaultS3SecretName = "/yandex-cloud/s3/kubernetes-backup"
	// Vault
	UpdateVaultSecret                = "update"
	PutVaultSecret                   = "put"
	FlantApproleRoleID               = "role_id"
	FlantApproleSecretID             = "secret_id"
	VaultApproleSecretPath           = "service-accounts"
	VaultApproleFlantResticBackups   = "/approles/flant-restic-backups-api-rw"
	VaultApproleDevopsServicesArgocd = "/approles/it-devops-services-argocd-rw"
	// Folder with files for Delete function
	PathToCommonFiles = "../common"
)

type Folder struct {
	Folders struct {
		Common                []string
		Nodes                 []string
		Monitoring            []string
		Alertmanager          []string
		Istio                 []string
		DeckhouseModules      []string
		NamespaceConfigurator []string
	}
}

type ClusterInfo struct {
	Cluster           string
	Domain            string
	Datacenter        string
	Environment       string
	Zone              string
	CustomString      string
	ArgoSAToken       string
	UseCustomCrowdApp bool
	Restic            struct {
		AwsAccessKeyID     string
		AwsSecretAccessKey string
	}
	Argo struct {
		CrowdPass string
		KubeToken string
	}
}

// Struct for Commander response "GET Clusters"
type CommanderClustersInfo struct {
	ID        int32
	Name      string
	Resources string
	Config    string
}

// Struct for Commander CONFIG (from response Clusters)
type ClusterInfoConfig struct {
	WithNATInstance struct {
		ExternalSubnetID string `yaml:"externalSubnetID"`
	} `yaml:"withNATInstance"`
}
