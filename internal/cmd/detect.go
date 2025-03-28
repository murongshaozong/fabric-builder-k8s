// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"errors"
	"os"
	"strconv"

	"github.com/hyperledger-labs/fabric-builder-k8s/internal/builder"
	"github.com/hyperledger-labs/fabric-builder-k8s/internal/log"
	"github.com/hyperledger-labs/fabric-builder-k8s/internal/util"
)

func Detect() {
	const (
		expectedArgsLength            = 3
		chaincodeSourceDirectoryArg   = 1
		chaincodeMetadataDirectoryArg = 2
	)

	debug, _ := strconv.ParseBool(util.GetOptionalEnv(util.DebugVariable, "false"))
	ctx := log.NewCmdContext(context.Background(), debug)
	logger := log.New(ctx)

	if len(os.Args) != expectedArgsLength {
		logger.Println("Expected CHAINCODE_SOURCE_DIR and CHAINCODE_METADATA_DIR arguments")

		os.Exit(1)
	}

	chaincodeSourceDirectory := os.Args[chaincodeSourceDirectoryArg]
	chaincodeMetadataDirectory := os.Args[chaincodeMetadataDirectoryArg]

	logger.Debugf("Chaincode source directory: %s", chaincodeSourceDirectory)
	logger.Debugf("Chaincode metadata directory: %s", chaincodeMetadataDirectory)

	detect := &builder.Detect{
		ChaincodeSourceDirectory:   chaincodeSourceDirectory,
		ChaincodeMetadataDirectory: chaincodeMetadataDirectory,
	}

	if err := detect.Run(ctx); err != nil {
		if !errors.Is(err, builder.ErrUnsupportedChaincodeType) {
			// don't spam the peer log if it's just chaincode we don't recognise
			logger.Printf("Error detecting chaincode: %+v", err)
		}

		os.Exit(1)
	}

	os.Exit(0)
}
