// Copyright 2026 shing1211
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

package trade

import (
	"fmt"
	"strings"
)

// HKDerivativesPartyIDSource is the party ID source code required for HK
// derivatives orders. It matches the BCAN party ID format used on HK equity
// orders.
const HKDerivativesPartyIDSource = "D"

// HKDerivativesPartyRole is the party role code required for HK derivatives
// orders. Role "3" corresponds to the BCAN party ID format used in HK equity
// regulatory reporting.
const HKDerivativesPartyRole = "3"

// ValidateHKDerivativesOrder checks that an order targeting the HK market has
// the required party ID information when it is a futures or options order. HK
// derivatives orders must carry BCAN party identifiers for regulatory
// compliance, mirroring the requirement on HK equity orders.
func ValidateHKDerivativesOrder(r OrderRequest) error {
	if r.Market != MarketHK {
		return nil
	}
	if r.InstrumentType != InstrumentTypeFutures && r.InstrumentType != InstrumentTypeOption {
		return nil
	}
	if len(r.NoPartyIDs) == 0 {
		return fmt.Errorf("no_party_ids is required for HK %s orders", r.InstrumentType)
	}
	for i, p := range r.NoPartyIDs {
		switch {
		case strings.TrimSpace(p.PartyID) == "":
			return fmt.Errorf("no_party_ids[%d]: party_id is required", i)
		case p.PartyIDSource != HKDerivativesPartyIDSource:
			return fmt.Errorf("no_party_ids[%d]: party_id_source must be %q, got %q", i, HKDerivativesPartyIDSource, p.PartyIDSource)
		case p.PartyRole != HKDerivativesPartyRole:
			return fmt.Errorf("no_party_ids[%d]: party_role must be %q, got %q", i, HKDerivativesPartyRole, p.PartyRole)
		}
	}
	return nil
}

// BuildHKDerivativesPartyIDs constructs a single-element NoPartyIDs slice
// populated with the required BCAN party ID values for HK derivatives orders.
// The supplied partyID should be the broker-assigned client identifier, for
// example "ABC123.2568".
func BuildHKDerivativesPartyIDs(partyID string) []PartyID {
	return []PartyID{
		{
			PartyID:       partyID,
			PartyIDSource: HKDerivativesPartyIDSource,
			PartyRole:     HKDerivativesPartyRole,
		},
	}
}
