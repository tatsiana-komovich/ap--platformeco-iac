package actions

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	appConf "github.com/adeo/lmru--devops--argocd-apps/go-configuration/config"
	"gopkg.in/yaml.v3"
)

func GetClustersInfo(clusterInfo appConf.ClusterInfo) (error, string) {
	// Пробуем достать из Vault креды от S3 хранилища
	getTokenFromVault, err := GetSecretKVv2WithAppRole(
		appConf.CommanderVaultSecretPath,
		appConf.CommanderVaultSecretName,
		appConf.DevopsVaultNamespace,
	)
	if err != nil {
		return fmt.Errorf("can't get Token from Vault: %w", err), ""
	}

	// устанавливаем время подключения
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}

	// собираем новый запрос в commander
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/clusters", appConf.CommanderAPI), nil)
	if err != nil {
		return fmt.Errorf("can't make a new GET request to commander: %w", err), ""
	}

	// добавляем Хедер для авторизации
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", fmt.Sprintf("%s", getTokenFromVault[appConf.CommanderBasicAuth])))

	// выполняем запрос собранный запрос
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("can't make a new GET request to commander: %w", err), ""
	}

	// парсим полученный ответ в структуру
	var clustersInfo []appConf.CommanderClustersInfo
	responseBody, err := ioutil.ReadAll(resp.Body)
	err = yaml.Unmarshal(responseBody, &clustersInfo)

	// теперь нам необходимо распарсить все, что попало в Config из commandera
	// распарссить и преобразовать в yaml
	// чтобы можно было распарсить в структуру
	for _, values := range clustersInfo {
		// здесь все, что попало в Config из commandera
		dec := yaml.NewDecoder(bytes.NewReader([]byte(values.Config)))

		// структура для выделения ID сетки из YC
		var newConfig appConf.ClusterInfoConfig

		clusterNameDefault := clusterInfo.Cluster
		if clusterInfo.Environment == "stage" {
			clusterNameDefault = fmt.Sprintf("%s-%s", clusterInfo.Cluster, clusterInfo.Environment)
		}
		// начинаем парсить полученный ответ от коммандера в yaml (структуру)
		for dec.Decode(&newConfig) == nil {
			// если длина ответа больше 0 (если есть хоть что-то)
			if len(newConfig.WithNATInstance.ExternalSubnetID) > 0 {
				// если нужный нам кластер
				if values.Name == clusterNameDefault {
					return nil, newConfig.WithNATInstance.ExternalSubnetID
				}
			}
		}
	}
	return fmt.Errorf("empty body from function GetClustersInfo"), ""
}
