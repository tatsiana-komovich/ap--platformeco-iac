package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/adeo/lmru--devops--argocd-apps/go-configuration/actions"
	appConf "github.com/adeo/lmru--devops--argocd-apps/go-configuration/config"
)

func createCommonFromTemplate(
	clusterInfo appConf.ClusterInfo, pathName, resource string) string {

	var createFilePath, templatepath string

	// парсим файл (шаблон) для нужной сущности в арго
	switch resource {
	case "Nodes":
		templatepath = fmt.Sprintf("templates/nodes/%s", pathName)
	case "Monitoring":
		templatepath = fmt.Sprintf("templates/monitoring/%s", pathName)
	case "Istio":
		templatepath = fmt.Sprintf("templates/istio/%s", pathName)
	case "Alertmanager":
		templatepath = fmt.Sprintf("templates/observability/alertmanager/%s", pathName)
	case "DeckhouseModules":
		templatepath = fmt.Sprintf("templates/deckhouse-modules/%s", pathName)
	case "NamespaceConfigurator":
		templatepath = fmt.Sprintf("templates/deckhouse-modules/%s", pathName)
	default:
		templatepath = fmt.Sprintf("templates/%s", pathName)
	}

	// читаем шаблон файла
	t, err := template.ParseFiles(templatepath)
	if err != nil {
		return fmt.Sprintf("Файл templates/%s не смог быть прочитан\n %s", pathName, err)
	}

	// нужно проверить для какой сущности генерируем конфиг, чтобы настроить путь
	switch {
	case resource == "Nodes":
		clusterPrefix := strings.Split(clusterInfo.Cluster, "-")
		createFilePath = fmt.Sprintf("../clusters/%s/%s/%s/nodes/%s.yaml", clusterPrefix[0], clusterInfo.Environment, clusterPrefix[1], pathName)
	case resource == "Monitoring":
		createFilePath = fmt.Sprintf("../common/observability/flant/%s/values/%s-%s-cluster.yaml",
			pathName, clusterInfo.Cluster, clusterInfo.Environment)
	case resource == "Istio":
		folderPrefix := "p"
		if clusterInfo.Environment == "stage" {
			folderPrefix = "s"
		}
		createFilePath = fmt.Sprintf("../common/istio/%s/%s/%s/templates/%s-%s-cluster.yaml",
			pathName, clusterInfo.Environment, fmt.Sprintf("istio-%s-mesh", folderPrefix), clusterInfo.Cluster, clusterInfo.Environment)
	case resource == "Alertmanager":
		createFilePath = fmt.Sprintf("../common/observability/lm/alertmanager/values/%s/%s-%s-cluster/%s.yaml",
			clusterInfo.Environment, clusterInfo.Cluster, clusterInfo.Environment, pathName)
		customString := `{{ define "generator.url" }} {{ (index .Alerts 0).GeneratorURL | reReplaceAll "&g0.tab=1$" "&g0.tab=0" }} {{ end }}`
		clusterInfo.CustomString = customString
	case resource == "DeckhouseModules":
		createFilePath = fmt.Sprintf("../common/%s/values/%s-%s-cluster.yaml",
			pathName, clusterInfo.Cluster, clusterInfo.Environment)
	case resource == "NamespaceConfigurator":
		createFilePath = fmt.Sprintf("../common/deckhouse-modules/%s/values/%s-%s-cluster.yaml",
			pathName, clusterInfo.Cluster, clusterInfo.Environment)
	case pathName == "clusters-permissions":
		createFilePath = fmt.Sprintf("../common/%s/%s/values/%s-%s-cluster.yaml", pathName, clusterInfo.Environment, clusterInfo.Cluster, clusterInfo.Environment)
	default:
		createFilePath = fmt.Sprintf("../common/%s/values/%s-%s-cluster.yaml", pathName, clusterInfo.Cluster, clusterInfo.Environment)
	}

	// проверяем файл на существование
	err = actions.ExistFile(createFilePath)
	if err != nil {
		return fmt.Sprintf("%s", err)
	}

	// создаем файл по нужному пути
	f, err := os.Create(createFilePath)
	if err != nil {
		return fmt.Sprintf("Не получилось создать файл по пути %s -- %s", createFilePath, err)
	}

	defer f.Close()

	// записываем инфомарцию в файл согласно шаблону
	// переменные передаются в шаблон
	err = t.Execute(f, clusterInfo)
	if err != nil {
		return fmt.Sprintf("%s", err)
	}

	// возвращаем, что файл успешно создан
	return fmt.Sprintf("Файл %s создан", createFilePath)
}

func main() {
	// определяем флаги для приложеньки
	create := flag.Bool("create", false, "Create files for cluster.")
	backup := flag.Bool("backup", false, "Create or add backup for cluster")
	delete := flag.Bool("delete", false, "Delete cluster files")
	file := flag.String("file", "/tmp/clusterInfo.yaml", "Path to file in yaml format (default /tmp/clusterInfo.yaml)")

	// переменная для многопотока
	results := make(chan string, 20)
	timeout := time.After(30 * time.Second)

	flag.Parse()

	if (*create && *delete) || (*backup && *delete) {
		appConf.InfoLog.Println("Нельзя использовать создание / бэкапирование и удаление одновременно")
		return
	}

	switch {
	default:
		appConf.InfoLog.Println("Не указано ни одного аргумента.")
	case *delete:
		// открываем файл с настройками кластера
		clusterInfo, err := actions.OpenClusterInfoYaml(*file)
		if err != nil {
			appConf.ErrorLog.Fatalf("Что-то пошло не так с чтением файла %s\t%s\n", *file, err)
		}
		if strings.Contains("prod", clusterInfo.Environment) {
			clusterInfo.Environment = "prod"
		}

		err = filepath.Walk(appConf.PathToCommonFiles,
			func(path string, info os.FileInfo, err error) error {
				fullCLusterName := fmt.Sprintf("%s-%s", clusterInfo.Cluster, clusterInfo.Environment)
				fileClusterName := strings.Split(path, "/")
				if strings.Contains(fileClusterName[len(fileClusterName)-1], fullCLusterName) {
					// проверяем путь это файл или папка
					if info.IsDir() {
						// если папка, то нужно удалить все содержимое внутри и саму папку
						err := os.RemoveAll(path)
						if err != nil {
							appConf.ErrorLog.Fatalf("%s", err)
						}
						appConf.InfoLog.Printf("Папка удалена - %s", path)
						return nil
					}
					// удаляем файл по пути
					err = os.Remove(path)
					if err != nil {
						appConf.ErrorLog.Fatalf("%s", err)
					}
					appConf.InfoLog.Printf("Файл удален - %s", path)
				}
				return nil
			})
		if err != nil {
			appConf.ErrorLog.Fatalf("%s", err)
		}
		// удаляем папку с конфигами NodeGroup
		clusterPrefix := strings.Split(clusterInfo.Cluster, "-")
		clusterNodesConfigsPath := fmt.Sprintf("../clusters/%s/%s/%s", clusterPrefix[0], clusterInfo.Environment, clusterPrefix[1])
		err = actions.DeleteDir(clusterNodesConfigsPath)

		appConf.InfoLog.Println("===================================================================")
		appConf.InfoLog.Println("Не забудьте удалить секреты из Vault - https://vault.lmru.tech/ui/vault/secrets/services/show/argocd-secrets?namespace=it-devops")
		appConf.InfoLog.Println("А также из неймспейса flant - https://vault.lmru.tech/ui/vault/secrets/restic_backups/list?namespace=flant")
		appConf.InfoLog.Println("Бакет в ЯО - https://console.cloud.yandex.ru/folders/b1glvj5h6scj9kpm1pc1/storage/buckets/")
		appConf.InfoLog.Println("Из репозитория ArgoCRD - https://github.lmru.tech/adeo/lmru--devops--argocd/tree/master/argocd-crd")
		appConf.InfoLog.Println("Из репозитория NGINX - https://github.lmru.tech/adeo/lmru--itplatform-app--nginx/blob/master/configuration/prod/")
		appConf.InfoLog.Println("===================================================================")

	case *backup:
		// открываем файл с настройками кластера
		clusterInfo, err := actions.OpenClusterInfoYaml(*file)
		if err != nil {
			appConf.ErrorLog.Fatalf("Что-то пошло не так с чтением файла %s\t%s\n", *file, err)
		}
		if strings.Contains("prod", clusterInfo.Environment) {
			clusterInfo.Environment = "prod"
		}
		backupFileCreate := createCommonFromTemplate(
			clusterInfo,
			"restic-backups",
			"",
		)
		appConf.InfoLog.Println(backupFileCreate)
		actions.CreateBackup(appConf.PutVaultSecret, clusterInfo)

	case *create:
		// открываем файл с настройками кластера
		clusterInfo, err := actions.OpenClusterInfoYaml(*file)
		if err != nil {
			appConf.ErrorLog.Fatalf("Что-то пошло не так с чтением файла %s\t%w\n", *file, err)
		}
		if strings.Contains("prod", clusterInfo.Environment) {
			clusterInfo.Environment = "prod"
		}
		// пробуем открыть файл с перечнем сущностей, для которых надо сделать конфиг
		folders, err := actions.OpenFoldersYaml()
		if err != nil {
			appConf.ErrorLog.Fatalf("Что-то пошло не так с чтением файла folders.yaml\n%s\n", err)
		}

		// делаем переменную, чтобы в нее положить словарь
		// полученным из файла folders.yml
		var dictFolders map[string][]string
		// содержимое, полученно из folders.yml преобразуем в json
		jsonFolders, _ := json.Marshal(&folders.Folders)
		// и преобразуем обратно в словарь + интерфейс
		_ = json.Unmarshal(jsonFolders, &dictFolders)
		// идем по ключам из нашего folders.yml
		for keyName, listNames := range dictFolders {
			// идем по списку файлов
			for _, oneName := range listNames {
				go func(value string) {
					// создаем файл из шаблона
					results <- createCommonFromTemplate(
						clusterInfo,
						oneName,
						keyName,
					)
				}(oneName)
				select {
				case res := <-results:
					appConf.InfoLog.Println(res)
				case <-timeout:
					appConf.ErrorLog.Printf("Timed out!")
					return
				}
			}
		}
		// Кладем в Vault токен для ArgoCD
		putSecret := map[string]interface{}{
			fmt.Sprintf("%s-%s-cluster", clusterInfo.Cluster, clusterInfo.Environment): clusterInfo.ArgoSAToken,
		}

		actions.UpdateVaultSecret(
			appConf.UpdateVaultSecret,
			appConf.VaultApproleSecretPath,           // getApprolePath
			appConf.VaultApproleDevopsServicesArgocd, // getApproleName
			appConf.DevopsVaultNamespace,             // getNamespace
			appConf.ArgocdSAVaultSecretPath,          // putSecretPath
			appConf.ArgocdSAVaultSecretName,          // putSecretName
			appConf.DevopsVaultNamespace,             //putSecretNamespace
			putSecret)

		// создаем запись в Vault с данными для restic
		// создаем Bucket в S3
		actions.CreateBackup(appConf.PutVaultSecret, clusterInfo)
	}
}
