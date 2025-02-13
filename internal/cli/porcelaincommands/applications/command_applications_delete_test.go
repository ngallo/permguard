// Copyright 2024 Nitro Agility S.r.l.
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
//
// SPDX-License-Identifier: Apache-2.0

package applications

import (
	"fmt"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"

	aziclicommon "github.com/permguard/permguard/internal/cli/common"
	aztestutils "github.com/permguard/permguard/internal/cli/porcelaincommands/testutils"
	azmocks "github.com/permguard/permguard/internal/cli/porcelaincommands/testutils/mocks"
	azconfigs "github.com/permguard/permguard/pkg/cli/options"
	azerrors "github.com/permguard/permguard/pkg/core/errors"
	azmodelaap "github.com/permguard/permguard/pkg/transport/models/aap"
)

// TestDeleteCommandForApplicationsDelete tests the deleteCommandForApplicationsDelete function.
func TestDeleteCommandForApplicationsDelete(t *testing.T) {
	args := []string{"-h"}
	outputs := []string{"The official Permguard Command Line Interface", "Copyright © 2022 Nitro Agility S.r.l.", "This command deletes a remote application."}
	aztestutils.BaseCommandTest(t, createCommandForApplicationDelete, args, false, outputs)
}

// TestCliApplicationsDeleteWithError tests the command for creating an application with an error.
func TestCliApplicationsDeleteWithError(t *testing.T) {
	tests := []struct {
		OutputType string
		HasError   bool
	}{
		{
			OutputType: "terminal",
			HasError:   true,
		},
		{
			OutputType: "json",
			HasError:   false,
		},
	}
	for _, test := range tests {
		args := []string{"applications", "delete", "--appid", "581616507495", "--output", test.OutputType}
		outputs := []string{""}

		v := viper.New()
		v.Set(azconfigs.FlagName(aziclicommon.FlagPrefixAAP, aziclicommon.FlagSuffixAAPTarget), "localhost:9092")

		depsMocks := azmocks.NewCliDependenciesMock()
		cmd := createCommandForApplicationDelete(depsMocks, v)
		cmd.PersistentFlags().StringP(aziclicommon.FlagWorkingDirectory, aziclicommon.FlagWorkingDirectoryShort, ".", "work directory")
		cmd.PersistentFlags().StringP(aziclicommon.FlagOutput, aziclicommon.FlagOutputShort, test.OutputType, "output format")
		cmd.PersistentFlags().BoolP(aziclicommon.FlagVerbose, aziclicommon.FlagVerboseShort, true, "true for verbose output")

		aapClient := azmocks.NewGrpcAAPClientMock()
		aapClient.On("DeleteApplication", mock.Anything).Return(nil, azerrors.ErrClientParameter)

		printerMock := azmocks.NewPrinterMock()
		printerMock.On("Println", mock.Anything).Return()
		printerMock.On("PrintlnMap", mock.Anything).Return()
		printerMock.On("Error", azerrors.ErrClientParameter).Return()

		depsMocks.On("CreatePrinter", mock.Anything, mock.Anything).Return(printerMock, nil)
		depsMocks.On("CreateGrpcAAPClient", mock.Anything).Return(aapClient, nil)

		aztestutils.BaseCommandWithParamsTest(t, v, cmd, args, true, outputs)
		if test.HasError {
			printerMock.AssertCalled(t, "Error", azerrors.ErrClientParameter)
		} else {
			printerMock.AssertNotCalled(t, "Error", azerrors.ErrClientParameter)
		}
	}
}

// TestCliApplicationsDeleteWithSuccess tests the command for creating an application with an error.
func TestCliApplicationsDeleteWithSuccess(t *testing.T) {
	tests := []string{
		"terminal",
		"json",
	}
	for _, outputType := range tests {
		args := []string{"applications", "delete", "--appid", "581616507495", "--output", outputType}
		outputs := []string{""}

		v := viper.New()
		v.Set("output", outputType)
		v.Set(azconfigs.FlagName(aziclicommon.FlagPrefixAAP, aziclicommon.FlagSuffixAAPTarget), "localhost:9092")

		depsMocks := azmocks.NewCliDependenciesMock()
		cmd := createCommandForApplicationDelete(depsMocks, v)
		cmd.PersistentFlags().StringP(aziclicommon.FlagWorkingDirectory, aziclicommon.FlagWorkingDirectoryShort, ".", "work directory")
		cmd.PersistentFlags().StringP(aziclicommon.FlagOutput, aziclicommon.FlagOutputShort, outputType, "output format")
		cmd.PersistentFlags().BoolP(aziclicommon.FlagVerbose, aziclicommon.FlagVerboseShort, true, "true for verbose output")

		aapClient := azmocks.NewGrpcAAPClientMock()
		application := &azmodelaap.Application{
			ApplicationID: 581616507495,
			Name:          "mycorporate",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		aapClient.On("DeleteApplication", mock.Anything).Return(application, nil)

		printerMock := azmocks.NewPrinterMock()
		outputPrinter := map[string]any{}

		if outputType == "terminal" {
			applicationID := fmt.Sprintf("%d", application.ApplicationID)
			outputPrinter[applicationID] = application.Name
		} else {
			outputPrinter["applications"] = []*azmodelaap.Application{application}
		}
		printerMock.On("PrintMap", outputPrinter).Return()
		printerMock.On("PrintlnMap", outputPrinter).Return()

		depsMocks.On("CreatePrinter", mock.Anything, mock.Anything).Return(printerMock, nil)
		depsMocks.On("CreateGrpcAAPClient", mock.Anything).Return(aapClient, nil)

		aztestutils.BaseCommandWithParamsTest(t, v, cmd, args, false, outputs)
		printerMock.AssertCalled(t, "PrintlnMap", outputPrinter)
	}
}
