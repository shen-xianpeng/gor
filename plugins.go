package main

import (
	"io"
)

type InOutPlugins struct {
	Inputs  []io.Reader
	Outputs []io.Writer
}

var Plugins *InOutPlugins = new(InOutPlugins)

func InitPlugins() {
	for _, options := range Settings.inputDummy {
		Plugins.Inputs = append(Plugins.Inputs, NewDummyInput(options))
	}

	for _, options := range Settings.outputDummy {
		Plugins.Outputs = append(Plugins.Outputs, NewDummyOutput(options))
	}

	for _, options := range Settings.inputRAW {
		Plugins.Inputs = append(Plugins.Inputs, NewRAWInput(options))
	}

	for _, options := range Settings.inputTCP {
		Plugins.Inputs = append(Plugins.Inputs, NewTCPInput(options))
	}

	for _, options := range Settings.outputTCP {
		Plugins.Outputs = append(Plugins.Outputs, NewTCPOutput(options))
	}

	for _, options := range Settings.inputData {
		Plugins.Inputs = append(Plugins.Inputs, NewDataInput(options))
	}

	for _, options := range Settings.inputFile {
		Plugins.Inputs = append(Plugins.Inputs, NewFileInput(options))
	}

	for _, options := range Settings.outputFile {
		Plugins.Outputs = append(Plugins.Outputs, NewFileOutput(options))
	}

	for _, options := range Settings.outputHTTP {
		Plugins.Outputs = append(Plugins.Outputs, NewHTTPOutput(options, Settings.outputHTTPHeaders, Settings.outputHTTPMethods, Settings.outputHTTPUrlRegexp, Settings.outputHTTPHeaderFilters, Settings.outputHTTPHeaderHashFilters, Settings.outputHTTPElasticSearch, Settings.outputHTTPUrlRewrite))
	}
}
