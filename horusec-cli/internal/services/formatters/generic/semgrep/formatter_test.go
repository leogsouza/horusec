// Copyright 2020 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package semgrep

import (
	"errors"
	"testing"

	"github.com/ZupIT/horusec/development-kit/pkg/entities/horusec"
	cliConfig "github.com/ZupIT/horusec/horusec-cli/config"
	"github.com/ZupIT/horusec/horusec-cli/internal/entities/workdir"
	"github.com/ZupIT/horusec/horusec-cli/internal/services/docker"
	"github.com/ZupIT/horusec/horusec-cli/internal/services/formatters"
	"github.com/stretchr/testify/assert"
)

func TestParseOutput(t *testing.T) {
	t.Run("Should return 1 vulnerabilities with no errors", func(t *testing.T) {
		dockerAPIControllerMock := &docker.Mock{}
		dockerAPIControllerMock.On("SetAnalysisID")
		analysis := &horusec.Analysis{}
		config := &cliConfig.Config{}
		config.SetWorkDir(&workdir.WorkDir{})

		output := "{ \"results\":[ { \"check_id\":\"python.lang.correctness.useless-comparison.no-strings-as-booleans\"," +
			" \"path\":\"bad/vulpy.py\", \"start\":{ \"line\":36, \"col\":1 }, \"end\":{ \"line\":37, \"col\":23 }, " +
			"\"extra\":{ \"message\":\"Using strings as booleans in Python has unexpected results.\\n`\\\"one\\\" and " +
			"\\\"two\\\"` will return \\\"two\\\".\\n`\\\"one\\\" or \\\"two\\\"` will return \\\"one\\\".\\n In Python" +
			", strings are truthy, evaluating to True.\\n\", \"metavars\":{ }, \"metadata\":{ }, \"severity\":\"ERROR\"" +
			", \"lines\":\"if csp:\\n    print('CSP:', csp)\" } } ] }"

		dockerAPIControllerMock.On("CreateLanguageAnalysisContainer").Return(output, nil)

		service := formatters.NewFormatterService(analysis, dockerAPIControllerMock, config, &horusec.Monitor{})
		formatter := NewFormatter(service)

		formatter.StartAnalysis("")
		assert.Len(t, analysis.AnalysisVulnerabilities, 1)
	})

	t.Run("Should return 1 vulnerabilities with no errors", func(t *testing.T) {
		dockerAPIControllerMock := &docker.Mock{}
		dockerAPIControllerMock.On("SetAnalysisID")
		analysis := &horusec.Analysis{}
		config := &cliConfig.Config{}
		config.SetWorkDir(&workdir.WorkDir{})

		output := "{ \"results\":[ { \"check_id\":\"python.lang.correctness.useless-comparison.no-strings-as-booleans\"," +
			" \"path\":\"bad/vulpy.py\", \"start\":{ \"line\":36, \"col\":1 }, \"end\":{ \"line\":37, \"col\":23 }, " +
			"\"extra\":{ \"message\":\"Using strings as booleans in Python has unexpected results.\\n`\\\"one\\\" and " +
			"\\\"two\\\"` will return \\\"two\\\".\\n`\\\"one\\\" or \\\"two\\\"` will return \\\"one\\\".\\n In Python" +
			", strings are truthy, evaluating to True.\\n\", \"metavars\":{ }, \"metadata\":{ }, \"severity\":\"WARNING\"" +
			", \"lines\":\"if csp:\\n    print('CSP:', csp)\" } } ] }"

		dockerAPIControllerMock.On("CreateLanguageAnalysisContainer").Return(output, nil)

		service := formatters.NewFormatterService(analysis, dockerAPIControllerMock, config, &horusec.Monitor{})
		formatter := NewFormatter(service)

		formatter.StartAnalysis("")
		assert.Len(t, analysis.AnalysisVulnerabilities, 1)
	})

	t.Run("Should return 1 vulnerabilities with no errors", func(t *testing.T) {
		dockerAPIControllerMock := &docker.Mock{}
		dockerAPIControllerMock.On("SetAnalysisID")
		analysis := &horusec.Analysis{}
		config := &cliConfig.Config{}
		config.SetWorkDir(&workdir.WorkDir{})

		output := "{ \"results\":[ { \"check_id\":\"python.lang.correctness.useless-comparison.no-strings-as-booleans\"," +
			" \"path\":\"bad\", \"start\":{ \"line\":36, \"col\":1 }, \"end\":{ \"line\":37, \"col\":23 }, " +
			"\"extra\":{ \"message\":\"Using strings as booleans in Python has unexpected results.\\n`\\\"one\\\" and " +
			"\\\"two\\\"` will return \\\"two\\\".\\n`\\\"one\\\" or \\\"two\\\"` will return \\\"one\\\".\\n In Python" +
			", strings are truthy, evaluating to True.\\n\", \"metavars\":{ }, \"metadata\":{ }, \"severity\":\"test\"" +
			", \"lines\":\"if csp:\\n    print('CSP:', csp)\" } } ] }"

		dockerAPIControllerMock.On("CreateLanguageAnalysisContainer").Return(output, nil)

		service := formatters.NewFormatterService(analysis, dockerAPIControllerMock, config, &horusec.Monitor{})
		formatter := NewFormatter(service)

		formatter.StartAnalysis("")
		assert.Len(t, analysis.AnalysisVulnerabilities, 1)
	})

	t.Run("Should return error when invalid output", func(t *testing.T) {
		dockerAPIControllerMock := &docker.Mock{}
		dockerAPIControllerMock.On("SetAnalysisID")
		analysis := &horusec.Analysis{}
		config := &cliConfig.Config{}
		config.SetWorkDir(&workdir.WorkDir{})

		output := "!!"

		dockerAPIControllerMock.On("CreateLanguageAnalysisContainer").Return(output, nil)

		service := formatters.NewFormatterService(analysis, dockerAPIControllerMock, config, &horusec.Monitor{})
		formatter := NewFormatter(service)

		formatter.StartAnalysis("")
		assert.NotEmpty(t, analysis.Errors)
	})

	t.Run("Should return error when executing container", func(t *testing.T) {
		dockerAPIControllerMock := &docker.Mock{}
		dockerAPIControllerMock.On("SetAnalysisID")
		analysis := &horusec.Analysis{}
		config := &cliConfig.Config{}
		config.SetWorkDir(&workdir.WorkDir{})

		dockerAPIControllerMock.On("CreateLanguageAnalysisContainer").Return("", errors.New("test"))

		service := formatters.NewFormatterService(analysis, dockerAPIControllerMock, config, &horusec.Monitor{})
		formatter := NewFormatter(service)

		formatter.StartAnalysis("")
		assert.NotEmpty(t, analysis.Errors)
	})

	t.Run("Should not execute tool because it's ignored", func(t *testing.T) {
		analysis := &horusec.Analysis{}
		dockerAPIControllerMock := &docker.Mock{}
		config := &cliConfig.Config{}
		config.SetWorkDir(&workdir.WorkDir{})
		config.SetToolsToIgnore([]string{"semgrep"})

		service := formatters.NewFormatterService(analysis, dockerAPIControllerMock, config, &horusec.Monitor{})
		formatter := NewFormatter(service)

		formatter.StartAnalysis("")
	})
}
