package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

func getDependencyListInFile() []string {
	return []string{}
}

func getFilesInDirectory() []string {
	return []string{}
}

func getDirectories() []string {
	return []string{}
}

func main() {
	sourceDirArg := flag.String("source", "", "")
	entryPointFileNameArg := flag.String("entry", "", "")

	flag.Parse()
	sourceDir := string(*sourceDirArg)
	entryPointFileName := string(*entryPointFileNameArg)

	// 1. read cli flags langType, sourceDir, buildConfig, depth (0 = infinite, also default)
	// 2. get type of project from langType flag (JS/Python/Ruby)
	// 3. load module to parse based on language
	parser := new(jsSourceParser)
	parser.ParseConfigFile("hellopath")

	if sourceDir == "" {
		fmt.Print("Please provide root directory of your JS project")
		return
	}
	fmt.Printf("JS Source: %v\n", sourceDir)
	fmt.Printf("Entrypoint JS/TS File: %v\n\n", entryPointFileName)

	files, err := ioutil.ReadDir(sourceDir)
	if err != nil {
		log.Fatalln(err)
		return
	}

	for _, file := range files {
		if entryPointFileName == file.Name() && !file.IsDir() {
			hasSlash := sourceDir[len(sourceDir)-1] == '/'
			filePath := sourceDir
			if hasSlash {
				filePath += entryPointFileName
			} else {
				filePath += "/" + entryPointFileName
			}
			imports, err := getImportsInFile(filePath)
			if err != nil {
				return
			}
			for _, imp := range imports {
				fmt.Println(imp)
			}
		}
	}
}

func getImportsInFile(absoluteFilePath string) ([]string, error) {
	importFilePaths := make([]string, 0)
	file, err := os.Open(absoluteFilePath)
	basePathTokens := strings.Split(absoluteFilePath, "/")
	basePath := strings.Join(basePathTokens[0:len(basePathTokens)-1], "/")

	/**
	* regex to filter out the following patterns
	* 1. require('path/to/script')
	* 2. import * from 'path/to/script'
	* 3. import(path/to/script)
	**/
	importRegex := regexp.MustCompile(`require\(["'](?P<importPath>.+)["']\)|import .* from ['"](?P<importPath>.+)['"];|import\((/\*.+\*/[ ])*["'](?P<importPath>.+)["']\)`)

	if err != nil {
		return importFilePaths, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		currentLine := scanner.Text()
		match := importRegex.FindStringSubmatch(currentLine)
		if len(match) > 0 {
			path := ""
			isLibrary := false
			if match[1] != "" {
				path, isLibrary = convertToAbsolutePath(basePath, match[1])
			} else if match[2] != "" {
				path, isLibrary = convertToAbsolutePath(basePath, match[2])
			} else if match[4] != "" {
				path, isLibrary = convertToAbsolutePath(basePath, match[4])
			}
			if path != "" && !isLibrary {
				importFilePaths = append(importFilePaths, path)
			}
		}

	}

	return importFilePaths, nil
}

func convertToAbsolutePath(basePath string, relativePath string) (string, bool) {
	basePathLength := len(basePath)
	hasEndingSlash := basePath[basePathLength-1] == '/'
	trimmedBasePath := basePath

	if hasEndingSlash {
		trimmedBasePath += basePath[0 : basePathLength-1]
	}
	relativePathRegex := regexp.MustCompile(`^\./(.+/)*.+|(\.\./)+(.+/)*.+`)
	regexMatch := relativePathRegex.FindStringSubmatch(relativePath)

	if regexMatch == nil || regexMatch[0] == "" {
		return relativePath, true
	}

	basePathSectionList := strings.Split(basePath, "/")
	relativePathSectionList := strings.Split(relativePath, "/")
	result := make([]string, 0)
	for _, relativePathSection := range relativePathSectionList {
		if relativePathSection == ".." {
			basePathSectionList = basePathSectionList[0 : len(basePathSectionList)-1]
		} else if relativePathSection != "." {
			result = append(result, relativePathSection)
		}
	}
	result = append(basePathSectionList, result...)

	return strings.Join(result, "/"), false
}
