package actions

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"

	appConf "github.com/adeo/lmru--devops--argocd-apps/go-configuration/config"

	"gopkg.in/yaml.v2"
)

func OpenFoldersYaml() (appConf.Folder, error) {
	var openFolder appConf.Folder
	envFileOpen, err := ioutil.ReadFile("folders.yml")
	err = yaml.Unmarshal(envFileOpen, &openFolder)
	if err != nil {
		return openFolder, err
	}
	return openFolder, nil
}

func OpenClusterInfoYaml(filePath string) (appConf.ClusterInfo, error) {
	var clusterInfo appConf.ClusterInfo
	envFileOpen, err := ioutil.ReadFile(filePath)
	if err != nil {
		return clusterInfo, err
	}
	err = yaml.Unmarshal(envFileOpen, &clusterInfo)
	if err != nil {
		return clusterInfo, err
	}
	return clusterInfo, nil
}

// проверяем существует ли файл
func ExistFile(pathFile string) error {
	// ращбиваем путь до файла
	folderSplitPath := strings.Split(pathFile, "/")
	// слеиваем путь папок
	folderPath := strings.Join(folderSplitPath[:len(folderSplitPath)-1], "/")
	// проверяем есть ли такой файл
	if _, err := os.Stat(pathFile); errors.Is(err, os.ErrNotExist) {
		// если файла нет (или не может быть проверен), то проверяем, есть ли такой фолдер
		if _, err := os.Stat(folderPath); errors.Is(err, os.ErrNotExist) {
			// создаем папку / папки, куда будет помещен файл
			os.MkdirAll(folderPath, 0765)
		}
		return nil
	}
	return errors.New(fmt.Sprintf("Файл %s существует", pathFile))
}

// проверяем папку на существование и удаляем, если есть
func DeleteDir(path string) error {
	// проверяем есть ли такая папка
	if stat, err := os.Stat(path); err == nil && stat.IsDir() {
		err := os.RemoveAll(path)
		if err != nil {
			return errors.New(fmt.Sprintf("%s", err))
		}
		appConf.InfoLog.Printf("Папка удалена - %s", path)
		return nil
	}
	return errors.New(fmt.Sprintf("Функция DeleteDir - %s", path))
}

// Функция генерирует рандомную строку в 16 символов,
// которая содержит числа и буквы
func GeneratePass() string {
	rand.Seed(time.Now().UnixNano())
	length := 16
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(appConf.Chars[rand.Intn(len(appConf.Chars))])
	}
	str := b.String()
	return str
}

func CreateBackup(method string, clusterInfo appConf.ClusterInfo) {
	// Имя секрета в Vault в NS Flant
	secretVaultName := fmt.Sprintf("%s-%s-cluster", clusterInfo.Cluster, clusterInfo.Environment)

	if clusterInfo.Environment == "stage" {
		clusterInfo.Cluster = fmt.Sprintf("%s-%s", clusterInfo.Cluster, clusterInfo.Environment)
	}

	// Пробуем достать из Vault креды от S3 хранилища
	getS3SecretFromVault, err := GetSecretKVv2WithAppRole(
		appConf.VaultS3SecretPath,
		appConf.VaultS3SecretName,
		appConf.DevopsVaultNamespace,
	)
	if err != nil {
		appConf.ErrorLog.Fatalf("%s\n", err)
	}

	// Пробуем создать S3 Bucket и фолдер в нем
	err = MainS3(
		fmt.Sprintf("%s", getS3SecretFromVault[appConf.S3SecretIDName]),
		fmt.Sprintf("%s", getS3SecretFromVault[appConf.S3SecretKeyName]),
		clusterInfo.Cluster, // название бакета в S3
	)
	if err != nil {
		appConf.ErrorLog.Printf("\t%s", err)
		return
	}

	// Формируем карту - это будет записано как секрет в Vault
	// putSecret := make(map[string]interface{})
	putSecret := map[string]interface{}{
		"AWS_ACCESS_KEY_ID":     getS3SecretFromVault[appConf.S3SecretIDName],
		"AWS_BUCKET":            fmt.Sprintf("%s/%s", clusterInfo.Cluster, appConf.BucketFolderName),
		"AWS_SECRET_ACCESS_KEY": getS3SecretFromVault[appConf.S3SecretKeyName],
		"RESTIC_PASSWORD":       GeneratePass(),
		"RESTIC_REPOSITORY":     fmt.Sprintf("s3://storage.yandexcloud.net/%s/%s", clusterInfo.Cluster, appConf.BucketFolderName),
	}

	UpdateVaultSecret(
		appConf.PutVaultSecret,
		appConf.VaultApproleSecretPath,         // getApprolePath
		appConf.VaultApproleFlantResticBackups, // getApproleName
		appConf.DevopsVaultNamespace,           // getNamespace
		appConf.BackupVaultSecretPath,          // putSecretPath
		secretVaultName,                        // putSecretName
		appConf.BackupVaultNamespace,           //putSecretNamespace
		putSecret)

}

func UpdateVaultSecret(
	method, getApprolePath, getApproleName, getNamespace, putSecretPath, putSecretName, putSecretNamespace string,
	putSecret map[string]interface{}) {

	// putSecret - что будем класть
	// getApprolePath - куда сходить за кредами
	// getApproleName - как называется секрет
	// getNamespace - куда сходить за кредами
	// putSecretPath - куда положить
	// putSecretName - как будет называться новый секрет
	// putSecretNamespace - куда положить новый секрет

	// Пробуем достать из Vault креды NS Flant в Vault
	getFlantApproleFromVault, err := GetSecretKVv2WithAppRole(
		getApprolePath,
		getApproleName,
		getNamespace,
	)
	if err != nil {
		appConf.ErrorLog.Fatalf("%s\n", err)
	}

	// Кладем в Vault конфиг для рестика
	err = PutSecretWithAppRole(
		putSecretPath,      // путь до секрета в vault
		putSecretName,      // название секрета в vault
		putSecretNamespace, // название NS Vault, в который надо сходить
		fmt.Sprintf("%s", getFlantApproleFromVault[appConf.FlantApproleRoleID]),   // role_id от vault approle
		fmt.Sprintf("%s", getFlantApproleFromVault[appConf.FlantApproleSecretID]), // secret_id от vault approle
		method,    // указываем что будем делать с секретом - put / update / delete
		putSecret, // что будет записано в сам секрет
	)
	if err != nil {
		appConf.ErrorLog.Printf("Could not %s creds in Vault \t%s", method, err)
		return
	}
	appConf.InfoLog.Printf("Запись в Vault успешно создана - namespace: %s - secret-name: %s", putSecretNamespace, putSecretName)
}
