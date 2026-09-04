"""Build a local stcons test genesis from the repo genesis: our dev validator in extraData,
stcStakingTime = now + OFFSET seconds, genesis timestamp = now - 60."""
import json, os, sys, time
REPO = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", "..", ".."))
src = json.load(open(sys.argv[2] if len(sys.argv) > 2 else os.path.join(REPO, "genesis.json")))
offset = int(sys.argv[1]) if len(sys.argv) > 1 else 120
val = "f39fd6e51aad88f6f4ce6ab8827279cfffb92266"
now = int(time.time())
extra = "00"*32 + "01" + val + "00"*48 + "01" + "00"*65          # vanity | count | addr+BLS | turnLength | seal
g = dict(src); g["config"] = dict(src["config"]); g["config"]["stcStakingTime"] = now + offset
for k in ("pascalTime","pragueTime","lorentzTime","maxwellTime","fermiTime","osakaTime","mendelTime","pasteurTime"):
    g["config"].setdefault(k, 0)   # mirror live configs: all upstream forks active from genesis
# STC uses the Cancun blob config for Prague/Osaka too (DefaultPragueBlobConfigSTC = DefaultCancunBlobConfig)
bs = dict(g["config"].get("blobSchedule") or {}); cancun = bs.get("cancun") or {"target": 3, "max": 6, "baseFeeUpdateFraction": 3338477}
bs.setdefault("prague", dict(cancun)); bs.setdefault("osaka", dict(cancun)); g["config"]["blobSchedule"] = bs
g["extraData"] = "0x" + extra; g["timestamp"] = hex(now - 60)
g["alloc"] = dict(src["alloc"]); g["alloc"][val] = {"balance": "0xD3C21BCECCEDA1000000"}
json.dump(g, open("genesis.json", "w"), indent=1)
print("genesis written: stcStakingTime =", now + offset, "(in", offset, "s) | genesis ts =", now - 60, "| chainId", g["config"]["chainId"], "| system contracts:", len([a for a,v in g["alloc"].items() if v.get("code")]))
