#!/bin/bash
# Post-fork verification for the staking-activation rehearsal chain.
R=${1:-http://127.0.0.1:9645}
HUB=0x0000000000000000000000000000000000001005; GOV=0x0000000000000000000000000000000000001007; TOK=0x0000000000000000000000000000000000001008; TL=0x0000000000000000000000000000000000001009
echo "head: $(cast block-number --rpc-url $R)"
echo "StakeHub.minSelfDelegationSEP : $(cast call $HUB 'minSelfDelegationSEP()(uint256)' --rpc-url $R)"
echo "StakeHub.unbondPeriod         : $(cast call $HUB 'unbondPeriod()(uint256)' --rpc-url $R)"
echo "StakeHub.maxElectedValidators : $(cast call $HUB 'maxElectedValidators()(uint256)' --rpc-url $R)"
echo "StakeHub.getProtector         : $(cast call $HUB 'getProtector()(address)' --rpc-url $R)"
echo "Governor.name                 : $(cast call $GOV 'name()(string)' --rpc-url $R)"
echo "Governor.votingPeriod         : $(cast call $GOV 'votingPeriod()(uint256)' --rpc-url $R)"
echo "GovToken.symbol/totalSupply   : $(cast call $TOK 'symbol()(string)' --rpc-url $R) / $(cast call $TOK 'totalSupply()(uint256)' --rpc-url $R)"
echo "Timelock.getMinDelay          : $(cast call $TL 'getMinDelay()(uint256)' --rpc-url $R 2>&1 | head -1)"
echo "--- fork block: system txs (gasPrice 0) to the four contracts ---"
HEAD=$(cast block-number --rpc-url $R)
for ((n=HEAD; n>0 && n>HEAD-120; n--)); do
  TX=$(cast block $n --json --rpc-url $R | python3 -c "import json,sys; b=json.load(sys.stdin); txs=b.get('transactions',[]); print(len(txs), b['timestamp'])")
  set -- $TX; if [ "$1" -gt 0 ]; then echo "block $n (time $2): $1 txs"; cast block $n --json --rpc-url $R | python3 -c "
import json,sys; b=json.load(sys.stdin)
for h in b['transactions']: print('   tx', h)" ; for h in $(cast block $n --json --rpc-url $R | python3 -c "import json,sys;[print(h) for h in json.load(sys.stdin)['transactions']]"); do cast tx $h --json --rpc-url $R | python3 -c "import json,sys;t=json.load(sys.stdin);print('     → to', t['to'], '| gasPrice', int(t['gasPrice'],16), '| input', t['input'][:10])"; done; break; fi
done
