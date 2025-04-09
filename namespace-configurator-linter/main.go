package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	appConf "github.com/adeo/lmru--devops--argocd-apps/namespace-configurator-linter/config"
	"gopkg.in/yaml.v3"
)

type MyStruct appConf.Folder

var (
	ErrInvalidValue = errors.New("invalid value for field 'value'")
	ErrNoValue      = errors.New("no value for field 'value'")
)

func openYamlFile(pathFile string) (error, appConf.Folder) {
	var openFolder appConf.Folder

	envFileOpen, err := os.ReadFile(pathFile)
	if err != nil {
		return err, openFolder
	}

	err = yaml.Unmarshal(envFileOpen, &openFolder)

	if err != nil {
		return err, openFolder
	}
	return err, openFolder
}

func validateAsteriskUsage(pattern string) error {
	if strings.Contains(pattern, "-*") || strings.Contains(pattern, "-+") {
		return errors.New("недопустимый символ сразу после дефиса")
	}

	return nil
}

func fileLint(filePath string) (error, string) {
	var exclude bool
	var excludeSecond bool
	returnText := ""
	err, readFile := openYamlFile(filePath)
	if err != nil {
		return errors.New(fmt.Sprintf("не удалось открыть файл %s\n%s", filePath, err)), returnText
	}
	// Проверяем валидность namespace_pattern
	for _, product := range readFile.Products {
		for _, pattern := range product.NamespacePattern {
			// Проверяем regex синтаксис
			_, err := regexp.Compile(pattern)
			if err != nil {
				returnText += fmt.Sprintf("\nInvalid regex syntax in ProductID '%s' in file %s:\n\tPattern: %s\n\tError: %v",
					product.ProductID, filePath, pattern, err)
				continue
			}
			// Проверяем наличие '-*' или '-+'
			if err := validateAsteriskUsage(pattern); err != nil {
				returnText += fmt.Sprintf("\nInvalid regex pattern in ProductID '%s' in file %s:\n\tPattern: %s\n\tError: %v",
					product.ProductID, filePath, pattern, err)
			}
		}
	}
	// Запускаем ряд циклов
	// здесь мы проверяем все NamespacePattern друг с другом
	// во всех ProductID одного файла
	for i := 0; i <= len(readFile.Products)-2; i++ {
		if len(readFile.Products[i].ProductID) == 0 {
			returnText += fmt.Sprintf(
				"\nProductID must be product-id not %s in file %s",
				readFile.Products[i], filePath)
		}
		// берем NamespacePattern из ProductID по порядку (сверху)
		for _, valueFirst := range readFile.Products[i].NamespacePattern {
			// запускаем второй цикл, в котором мы переберем все NamespacePattern
			// от ProductID. отличного от того, что взят в первом цикле
			for k := i + 1; k <= len(readFile.Products)-1; k++ {
				for _, valueSecond := range readFile.Products[k].NamespacePattern {
					// сделаем переменные по размеру длины текста в названии NamespacePattern
					minLenWord := valueFirst
					maxWord := valueSecond

					// на всякий случай проверим, что длина NamespacePattern из первого цикла
					// больше или меньше блины NamespacePattern из второго цикла
					// так проще находить вхождения одного в другом
					if len(valueFirst) > len(valueSecond) {
						minLenWord = valueSecond
						maxWord = valueFirst
					}

					// сравниваем два значения между друг другом
					// например ^billing-resource будет входить в ^billing-resource$
					// но ^billing$ не входит в ^billing-resource$
					if strings.Contains(maxWord, minLenWord) {
						// если у нас есть вхождения, то нам надо дополнить текст
						// который вернется после прохождения всех файлов
						returnText += fmt.Sprintf(
							"\nDetected similars in file %s:\n\tProductID --%s-- and --%s-- have similar NamespacePattern: %s and %s",
							filePath, readFile.Products[i].ProductID, readFile.Products[k].ProductID, valueFirst, valueSecond)
					}

					str := regexp.MustCompile(`[0-9A-Za-z-]+`)
					pattern, err := regexp.MatchString(minLenWord, str.FindString(maxWord))
					if err != nil {
						returnText += fmt.Sprintf(

							"\nDetected similars in file %s:\n\tProductID --%s-- and --%s-- have similar NamespacePattern: %s and %s",
							filePath, readFile.Products[i].ProductID, readFile.Products[k].ProductID, valueFirst, valueSecond)
					}
					for _, excludeValue := range readFile.Products[i].Exclusions {

						exclude, err = regexp.MatchString(minLenWord, str.FindString(excludeValue))
						// если у нас есть вхождения, то нам надо дополнить текст
						// который вернется после прохождения всех файлов
						if err != nil {
							returnText += fmt.Sprintf(

								"\nDetected similars in file %s:\n\tProductID --%s-- and --%s-- have similar NamespacePattern: %s and %s",
								filePath, readFile.Products[i].ProductID, readFile.Products[k].ProductID, valueFirst, valueSecond)
						}
					}
					for _, excludeValueSecond := range readFile.Products[k].Exclusions {

						excludeSecond, err = regexp.MatchString(minLenWord, str.FindString(excludeValueSecond))
						// если у нас есть вхождения, то нам надо дополнить текст
						// который вернется после прохождения всех файлов
						if err != nil {
							returnText += fmt.Sprintf(

								"\nDetected similars in file %s:\n\tProductID --%s-- and --%s-- have similar NamespacePattern: %s and %s",
								filePath, readFile.Products[i].ProductID, readFile.Products[k].ProductID, valueFirst, valueSecond)
						}
					}
					if pattern && !exclude && !excludeSecond {
						returnText += fmt.Sprintf(
							"\nDetected similars in file %s:\n\tProductID --%s-- and --%s-- have similar NamespacePattern: %s and %s %s",
							filePath, readFile.Products[i].ProductID, readFile.Products[k].ProductID, valueFirst, valueSecond, readFile.Products[i].NamespacePattern)
					}
				}
			}

		}

	}

	return err, returnText
}

func main() {
	// определяем флаги для приложеньки
	pathFolderValues := flag.String("path", "../common/deckhouse-modules/namespace-configurator/values", "Path to folder with conf files")

	flag.Parse()
	files, err := os.ReadDir(*pathFolderValues)
	if err != nil {
		appConf.ErrorLog.Fatalf("Не удалось прочитать содержимое папки %s\n%s\n", *pathFolderValues, err)
	}

	lintText := ""
	for _, file := range files {
		if file.IsDir() {
			continue // если вдруг в папке есть папка, то ее надо пропустить
		}
		filePath := fmt.Sprintf("%s/%s", *pathFolderValues, file.Name())
		err, similars := fileLint(filePath)
		if len(similars) > 0 {
			lintText += similars
		}
		if err != nil {
			appConf.ErrorLog.Println(err)
			continue
		}
	}
	if len(lintText) > 0 {
		appConf.ErrorLog.Fatal(lintText)
	}
}
