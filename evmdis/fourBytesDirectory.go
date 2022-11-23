// Copyright 2022 Kyrylo Danshyn
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package evmdis

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

var fourBytesCache = map[string][]string{
	"06fdde03": {"transfer_attention_tg_invmru_6e7aa58(bool,address,address)", "message_hour(uint256,int8,uint16,bytes32)", "name()"},
	"0753c30c": {"deprecate(address)"},
	"095ea7b3": {"watch_tg_invmru_2f69f1b(address,address)", "sign_szabo_bytecode(bytes16,uint128)", "approve(address,uint256)"},
	"0e136b19": {"deprecated()"},
	"0ecb93c0": {"addBlackList(address)"},
	"18160ddd": {"watch_tg_invmru_ae5c248(uint256,bool,bool)", "voting_var(address,uint256,int128,int128)", "totalSupply()"},
	"23b872dd": {"watch_tg_invmru_faebe36(bool,bool,bool)", "gasprice_bit_ether(int128)", "transferFrom(address,address,uint256)"},
	"26976e3f": {"upgradedAddress()"}, "27e235e3": {"balances(address)"},
	"313ce567": {"watch_tg_invmru_5c94e13(bool)", "watch_tg_invmru_e597f2(address,bool)", "transfer_attention_tg_invmru_efa43f6(uint256,bool,address)", "available_assert_time(uint16,uint64)", "decimals()"},
	"35390714": {"maximumFee()"},
	"3eaaf86b": {"_totalSupply()"},
	"3f4ba83a": {"unpause()"},
	"59bf1abe": {"getBlackListStatus(address)"},
	"5c658165": {"allowed(address,address)"},
	"5c975abb": {"paused()"},
	"6e18980a": {"transferByLegacy(address,address,uint256)"},
	"70a08231": {"watch_tg_invmru_119a5a98(address,uint256,uint256)", "passphrase_calculate_transfer(uint64,address)", "branch_passphrase_public(uint256,bytes8)", "balanceOf(address)"},
	"8456cb59": {"pause()"},
	"893d20e8": {"getOwner()"},
	"8b477adb": {"transferFromByLegacy(address,address,address,uint256)"},
	"8da5cb5b": {"ideal_warn_timed(uint256,uint128)", "owner()"},
	"95d89b41": {"watch_tg_invmru_4f9dd3f(address,uint256)", "link_classic_internal(uint64,int64)", "symbol()"},
	"a9059cbb": {"join_tg_invmru_haha_fd06787(address,bool)", "func_2093253501(bytes)", "transfer(bytes4[9],bytes5[6],int48[11])", "many_msg_babbage(bytes1)", "transfer(address,uint256)"},
	"aee92d33": {"approveByLegacy(address,address,uint256)"},
	"c0324c77": {"setParams(uint256,uint256)"},
	"cc872b66": {"issue(uint256)"},
	"db006a75": {"redeem(uint256)"},
	"dd62ed3e": {"join_tg_invmru_haha_5911067(uint256,address)", "_func_5437782296(address,address)", "remove_good(uint256[],bytes8,bool)", "allowance(address,address)"},
	"dd644f72": {"basisPointsRate()"},
	"e47d6060": {"isBlackListed(address)"},
	"e4997dc5": {"removeBlackList(address)"},
	"e5b5019a": {"MAX_UINT()"},
	"f3bdc228": {"destroyBlackFunds(address)"},
	"ffffffff": {"LOCK8605463013()", "test266151307()"},
	"f2fde38b": {"transferOwnership(address)"},
	"6fcfff45": {"numCheckpoints(address)"},
	"b4b5ea57": {"getCurrentVotes(address)"},
	"e7a324dc": {"DELEGATION_TYPEHASH()"},
	"f1127ed8": {"checkpoints(address,uint32)"},
	"fca3b5aa": {"setMinter(address)"},
	"c3cda520": {"delegateBySig(address,uint256,uint256,uint8,bytes32,bytes32)"},
	"d505accf": {"watch_tg_invmru_168a06(bool,address,bool)", "permit(address,address,uint256,uint256,uint8,bytes32,bytes32)"},
	"782d6fe1": {"getPriorVotes(address,uint256)"},
	"7ecebe00": {"transfer_attention_tg_invmru_5811b86(uint256,address,address)", "nonces(address)"},
	"76c71ca1": {"mintCap()"},
	"30adf81f": {"transfer_attention_tg_invmru_6c0d2a(uint256,uint256)", "PERMIT_TYPEHASH()"},
	"40c10f19": {"mint(address,uint256)"},
	"587cde1e": {"delegates(address)"},
	"5c11d62f": {"minimumTimeBetweenMints()"},
	"5c19a95c": {"delegate(address)"},
	"30b36cef": {"mintingAllowedAfter()"},
	"20606b70": {"DOMAIN_TYPEHASH()"},
	"07546172": {"minter()"},
	"01ffc9a7": {"pizza_mandate_apology(uint256)", "supportsInterface(bytes4)"},
	"8f32d59b": {"isOwner()"},
	"d6e4fa86": {"nameExpires(uint256)"},
	"e985e9c5": {"isApprovedForAll(address,address)"},
	"f6a74ed7": {"removeController(address)"},
	"fca247ac": {"register(uint256,address,uint256)"},
	"da8c229e": {"controllers(address)"},
	"ddf7fcb0": {"baseNode()"},
	"a7fc7a07": {"addController(address)"},
	"b88d4fde": {"safeTransferFrom(address,address,uint256,bytes)"},
	"c1a287e2": {"GRACE_PERIOD()"},
	"c475abff": {"renew(uint256,uint256)"},
	"96e494e8": {"available(uint256)"},
	"a22cb465": {"niceFunctionHerePlzClick943230089(address,bool)", "setApprovalForAll(address,bool)"},
	"3f15457f": {"ens()"},
	"6352211e": {"ownerOf(uint256)"},
	"715018a6": {"renounceOwnership()"},
	"42842e0e": {"safeTransferFrom(address,address,uint256)"},
	"4e543b26": {"setResolver(address)"},
	"081812fc": {"getApproved(uint256)"},
	"0e297b45": {"registerOnly(uint256,address,uint256)"},
	"28ed4f6c": {"reclaim(uint256,address)"},
	"02571be3": {"owner(bytes32)"},
	"06ab5923": {"setSubnodeOwner(bytes32,bytes32,address)"},
	"1896f70a": {"setResolver(bytes32,address)"},
	"150b7a02": {"onERC721Received(address,address,uint256,bytes)"},
}

var fourByteCacheMu sync.RWMutex

type replyRec struct {
	ID             uint64
	CreatedAt      string `json:"created_at"`
	TextSignature  string `json:"text_signature"`
	HexSignature   string `json:"hex_signature"`
	BytesSignature string `json:"bytes_signature"`
}

type fourByteDirectoryResponse struct {
	Count    int
	Next     string
	Previous string
	Results  []replyRec
}

func getSigsFromFourBytes(fourBytes []byte) []string {
	fourByteCacheMu.RLock()
	hexFourBytes := hex.EncodeToString(fourBytes)
	if sigs, ok := fourBytesCache[hexFourBytes]; ok {
		fourByteCacheMu.RUnlock()
		return sigs
	}

	fourByteCacheMu.RUnlock()

	sigs := []string{}

	call := fmt.Sprintf("https://www.4byte.directory/api/v1/signatures/?hex_signature=0x%08x", fourBytes)
	repl, err := http.Get(call)
	if err == nil {
		var b4 fourByteDirectoryResponse
		err = json.NewDecoder(repl.Body).Decode(&b4)
		if err == nil {
			for _, b4rec := range b4.Results {
				sigs = append(sigs, b4rec.TextSignature)
			}
		}
	}

	fourByteCacheMu.Lock()
	defer fourByteCacheMu.Unlock()
	fourBytesCache[hexFourBytes] = sigs

	return sigs
}
