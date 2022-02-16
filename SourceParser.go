/*
SourceParser interface defines the methods that a parser plugin needs to implement
*/

package main

type SourceParser interface {
	ParseConfigFile(configPath string)
	BuildDependencyGraph(filePath string, depth uint8)
	GetDependencyListForFile(filePath string) []string
	GetDependencyGraphForFile(filePath string)
}

type jsSourceParser struct {
	dependencyGraph string
}

func (parser jsSourceParser) ParseConfigFile(configFile string) {

}

func (parser jsSourceParser) BuildDependencyGraph(startFile string, depth uint8) {

}

func (parser jsSourceParser) GetDependencyListForFile(filePath string) []string {
	return make([]string, 0)
}

func (parser jsSourceParser) GetDependencyGraphForFile(filePath string) {

}
