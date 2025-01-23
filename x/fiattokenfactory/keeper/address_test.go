// Copyright 2025 Circle Internet Group, Inc.  All rights reserved.
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

package keeper_test

import (
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/btcutil/bech32"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/keeper"
	"github.com/stretchr/testify/require"
)

func TestDecodeNoLimitToBase256(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		ibcChain     string
		ibcChainAddr string
		ibcChainHrp  string
		nobleAddr    string
		bech32m      bool
	}{
		// The tests below ensure we can correctly decode addresses from IBC chains (both bech32/Segwit and bech32m/Taproot addresses).
		"cosmos": {
			ibcChainAddr: "cosmos1hjz2rjqfn7yhaawqgfk6j6hv5dtf9nau70fusm",
			ibcChainHrp:  "cosmos",
			nobleAddr:    "noble1hjz2rjqfn7yhaawqgfk6j6hv5dtf9naukvu5g4",
		},
		"osmosis": {
			ibcChainAddr: "osmo1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3aq6l09",
			ibcChainHrp:  "osmo",
			nobleAddr:    "noble1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3acu8pe",
		},
		"dydx": {
			ibcChainAddr: "dydx18vgsfaarveyg7xy585657ak8a9jvut9z8yuzmv",
			ibcChainHrp:  "dydx",
			nobleAddr:    "noble18vgsfaarveyg7xy585657ak8a9jvut9zx78wr4",
		},
		"namada": {
			// Generated by importing the following mnemonic into namada cli and running `namadaw derive --pre-genesis`.
			// mule inform liberty tray polar planet marble ketchup tiny brush hedgehog kiss project cable thank dismiss island fortune snake rice vicious feed intact canvas
			ibcChainAddr: "tpknam1qzdjad7ta2246ms4z82dz8zhv2trhw7w4fpnpuj56ekjakwcc3xqwvzr6ak",
			ibcChainHrp:  "tpknam",
			nobleAddr:    "noble1qzdjad7ta2246ms4z82dz8zhv2trhw7w4fpnpuj56ekjakwcc3xqwmvmf5j",
			bech32m:      true,
		},
		"penumbra": {
			// Generated by importing the following mnemonic into pcli (https://guide.penumbra.zone/pcli) and running `pcli view address`.
			// approve bracket canyon yard such jungle patch decade monster scissors burden gold stone essay shield scatter net dynamic salad umbrella play trophy lake blossom
			ibcChainAddr: "penumbra1ld2kghffzgwq4597ejpgmnwxa7ju0cndytuxtsjh8qhjyfuwq0rwd5flnw4a3fgclw7m5puh50nskn2c88flhne2hzchnpxru609d5wgmqqvhdf0sy2tktqfcm2p2tmxceqwvv",
			ibcChainHrp:  "penumbra",
			nobleAddr:    "noble1ld2kghffzgwq4597ejpgmnwxa7ju0cndytuxtsjh8qhjyfuwq0rwd5flnw4a3fgclw7m5puh50nskn2c88flhne2hzchnpxru609d5wgmqqvhdf0sy2tktqfcm2p2tmxq2k7my",
			bech32m:      true,
		},
		"penumbra compatible": {
			// Generated with the same mnemonic and running `pcli view address --compat`.
			ibcChainAddr: "penumbracompat11ld2kghffzgwq4597ejpgmnwxa7ju0cndytuxtsjh8qhjyfuwq0rwd5flnw4a3fgclw7m5puh50nskn2c88flhne2hzchnpxru609d5wgmqqvhdf0sy2tktqfcm2p2tmxeuc86n",
			ibcChainHrp:  "penumbracompat1",
			nobleAddr:    "noble1ld2kghffzgwq4597ejpgmnwxa7ju0cndytuxtsjh8qhjyfuwq0rwd5flnw4a3fgclw7m5puh50nskn2c88flhne2hzchnpxru609d5wgmqqvhdf0sy2tktqfcm2p2tmxq2k7my",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// Ensure we can decode the IBC chain address to noble address
			hrp, bz, err := keeper.DecodeNoLimitToBase256(tc.ibcChainAddr)
			require.NoError(t, err, fmt.Sprintf("failed to decode %s address", name))
			require.Equal(t, tc.ibcChainHrp, hrp)

			encoded, err := convertAndEncodeBase256("noble", bz, false)
			require.NoError(t, err, "failed to encode bytes to noble address")
			require.Equal(t, tc.nobleAddr, encoded)

			hrp, bz, err = keeper.DecodeNoLimitToBase256(tc.nobleAddr)
			require.NoError(t, err, "failed to decode noble address")
			require.Equal(t, "noble", hrp)

			encoded, err = convertAndEncodeBase256(tc.ibcChainHrp, bz, tc.bech32m)
			require.NoError(t, err, fmt.Sprintf("failed to encode bytes to %s address", tc.ibcChain))
			require.Equal(t, tc.ibcChainAddr, encoded)
		})
	}
}

// Function that reverts `DecodeNoLimitToBase256` and converts
func convertAndEncodeBase256(hrp string, data []byte, bech32m bool) (string, error) {
	converted, _ := bech32.ConvertBits(data, 8, 5, true)
	var encoded string
	if bech32m {
		encoded, _ = bech32.EncodeM(hrp, converted)
	} else {
		encoded, _ = bech32.Encode(hrp, converted)
	}
	return encoded, nil
}
