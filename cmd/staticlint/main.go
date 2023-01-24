// Custom static analyzer.
// To run use flag ./cmd/..
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gostaticanalysis/nilerr"
	"github.com/rs/zerolog/log"
	"github.com/timakin/bodyclose/passes/bodyclose"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/tests"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"

	"github.com/gtgaleevtimur/reduction-url-service/cmd/staticlint/analyzer"
)

// Config структура описывающая проверки кода.
type Config struct {
	StaticCheck []string
	StyleCheck  []string
}

func main() {
	file, err := os.Executable()
	if err != nil {
		log.Fatal().Err(err)
	}
	data, err := os.ReadFile(filepath.Join(filepath.Dir(file), `config.json`))
	if err != nil {
		log.Fatal().Err(err)
	}

	var config Config
	if err = json.Unmarshal(data, &config); err != nil {
		log.Fatal().Err(err)
	}

	analyzers := []*analysis.Analyzer{
		analyzer.ExitAnalyzer,
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		assign.Analyzer,
		bools.Analyzer,
		copylock.Analyzer,
		errorsas.Analyzer,
		framepointer.Analyzer,
		httpresponse.Analyzer,
		nilfunc.Analyzer,
		tests.Analyzer,
		nilerr.Analyzer,
		bodyclose.Analyzer,
	}

	for _, v := range staticcheck.Analyzers {
		for _, y := range config.StaticCheck {
			if strings.HasPrefix(v.Name, y) {
				analyzers = append(analyzers, v)
			}
		}
	}
	for _, v := range stylecheck.Analyzers {
		for _, y := range config.StyleCheck {
			if strings.HasPrefix(v.Name, y) {
				analyzers = append(analyzers, v)
			}
		}
	}

	fmt.Println("Go run static checks:\n", analyzers)

	multichecker.Main(analyzers...)
}
