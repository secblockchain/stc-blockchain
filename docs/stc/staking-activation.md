# Staking activation hard fork (`stcStakingTime`)

## Why this exists

On both STC networks the staking & governance system contracts — StakeHub `0x…1005`,
GovToken `0x…1008`, Governor `0x…1007`, Timelock `0x…1009` — are deployed from genesis but were
**never initialised**: `StakeHub.minSelfDelegationSEP()` is 0, no validator is registered,
`Governor.votingPeriod()` is 0. Block production runs on the genesis validator set in
ValidatorSet `0x…1000` (which *is* initialised by `initContract` at block 1).

Their `initialize()` functions are system-transaction-only (`onlyCoinbase` / `onlyZeroGasPrice`),
so no wallet or governance vote can call them. BSC solves this with the engine hook
`initializeFeynmanContract` at the Feynman fork transition; this fork had no such hook, so
nothing could ever initialise them. Consequence: no staking, no delegation, no governance, and
**no way to add a new validator** (`createValidator` requires an initialised StakeHub).

## What this change adds

| Piece | File | Effect |
|---|---|---|
| Fork switch `STCStakingTime` | `params/config.go` | new timestamp fork (`IsSTCStaking`, `IsOnSTCStaking`), included in fork-order and compatibility checks; JSON key `stcStakingTime` |
| Init hook | `consensus/stcons/stcons.go` → `initStakingContracts` | at the transition block, system txs call `initialize()` on StakeHub → GovToken → Governor → Timelock (import and mining paths) |
| Bytecode slot | `core/systemcontracts/upgrade.go` | at the same block, optionally replaces StakeHub bytecode with the parameterised build (`stakeHubStakingCode`); a hard guard makes an empty slot a no-op instead of erasing the contract |

The 5 genesis validators keep producing throughout: `validatorset.go` always unions genesis
validators into the active set. After activation they (and any newcomer) register via
`StakeHub.createValidator`.

## How to ship it (no regenesis — chain data, balances and contracts are preserved)

1. **Parameters.** Build the StakeHub whose `initialize()` carries STC's economics
   (minimum self-stake 10,000,000 SEP, the 100,000,000 SEP maximum-stake cap — a custom check the
   stock contract does not have — unbonding period, jail times, max elected validators). Paste its
   runtime bytecode into `stakeHubStakingCode` in `core/systemcontracts/upgrade.go`.
   *If you skip this, the StakeHub bytecode already on chain is initialised with its built-in defaults.*
2. **Schedule.** Set `STCStakingTime` in `params/config.go` for the testnet config first
   (`newUint64(<unix ts>)`, ~48–72 h ahead), build, roll out, observe; then the mainnet config.
   The chain config is code-embedded (`STCGenesisHash` → `STCChainConfig`), so **do not edit
   `genesis.json`** on the servers (`start-node.sh` wipes chaindata when genesis changes).
3. **Roll out to every node before the timestamp** — all validators, RPC nodes, the explorer node.
   Verify with `geth version`. A node left on the old binary rejects the fork block and drops off.
4. **Watch the fork block**: logs show `initialize staking contract` ×4, then
   ```
   cast call 0x0000000000000000000000000000000000001005 'minSelfDelegationSEP()(uint256)' --rpc-url <rpc>   # non-zero
   cast call 0x0000000000000000000000000000000000001007 'name()(string)' --rpc-url <rpc>                     # non-empty
   ```
5. **Register validators**: each operator calls `createValidator` (the operations console's
   "Become a validator" builds and pre-flights the transaction for the operator's wallet).
   Parameters remain tunable afterwards through Governor proposals.

## Local rehearsal

```
# genesis: your validator in extraData, stcStakingTime = genesis timestamp + 60
geth --datadir d init genesis.json
geth --datadir d --networkid <id> --mine --miner.etherbase <val> --unlock <val> --password pw --allow-insecure-unlock --http
# after the first block with time >= stcStakingTime:
cast call 0x…1005 'minSelfDelegationSEP()(uint256)' --rpc-url http://127.0.0.1:8545
```

## Rehearsal result (2026-09-04, local single-validator chain built from this branch)

Harness: `docs/stc/rehearsal/make_genesis.py` (repo genesis + our validator in `extraData` +
`stcStakingTime = now + N seconds`) and `docs/stc/rehearsal/verify.sh` (post-fork reads).
The validator key is Hardhat account #0 (`0xf39F…2266`, well-known test key), never a real key.

```bash
python3 docs/stc/rehearsal/make_genesis.py 120           # fork 120 s after start
geth init --datadir ./data genesis.json
geth --datadir ./data account import --password password.txt key.hex
geth --datadir ./data --networkid 19517 --nodiscover --mine --miner.etherbase 0xf39F…2266 \
     --unlock 0xf39F…2266 --password password.txt --allow-insecure-unlock \
     --http --http.port 9645 --http.api eth,net,web3
docs/stc/rehearsal/verify.sh http://127.0.0.1:9645       # run after the fork timestamp
```

Observed:

- Before the fork every StakeHub/Governor/GovToken/Timelock read returned `0` / empty — the same
  uninitialised state mainnet and testnet are in today.
- Block 39 was the first block with `time >= stcStakingTime`. The node logged
  `Apply upgrade stc-staking-activation at height 39`, then
  `initialize staking contract` for `0x…1005`, `0x…1008`, `0x…1007`, `0x…1009` (StakeHub,
  GovToken, Governor, Timelock) — all in that one block, as system transactions.
- After the fork: `minSelfDelegationSEP = 1e26` (100,000,000 SEP), `maxElectedValidators = 9`,
  GovToken `STC Governance Token / govSEP`, Governor `votingPeriod = 28800`, Timelock `minDelay = 21600`.
- The chain kept sealing normally after the fork (head 69 at shutdown, zero errors in the log).

Important for the 10M–100M SEP band: the StakeHub bytecode that is in genesis today hard-codes
the minimum self-delegation at **100,000,000 SEP** and has **no maximum**. Because no replacement
bytecode is registered yet (`stakeHubStakingCode` is empty), the fork logs
`No bytecode registered for system-contract upgrade, keeping existing code` and initialises the
existing contract with that 100M default. The developer must compile a StakeHub with
`minimum = 10_000_000 SEP`, a `maximum = 100_000_000 SEP` check in `createValidator`/`delegate`,
and paste its runtime bytecode into `stakeHubStakingCode` before the fork timestamp is scheduled.
