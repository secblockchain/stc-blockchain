## Smart Technology Chain


## Building the source


Building `geth` requires both a Go (version 1.23 or later) and a C compiler. You can install
them using your favourite package manager. Once the dependencies are installed, run

```shell
make geth
```

or, to build the full suite of utilities:

```shell
make all
```

## Executables

The go-ethereum project comes with several wrappers/executables found in the `cmd`
directory.

|  Command   | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        |
| :--------: | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **`geth`** | Our main Ethereum CLI client. It is the entry point into the Ethereum network (main-, test- or private net), capable of running as a full node (default), archive node (retaining all historical state) or a light node (retrieving data live). It can be used by other processes as a gateway into the Ethereum network via JSON RPC endpoints exposed on top of HTTP, WebSocket and/or IPC transports. `geth --help` and the [CLI page](https://geth.ethereum.org/docs/fundamentals/command-line-options) for command line options. |
|  `devp2p`  | Utilities to interact with nodes on the networking layer, without running a full blockchain.                                                                                                                                                                                                                                                                                                                                                                                                                                       |
|  `abigen`  | Source code generator to convert Ethereum contract definitions into easy-to-use, compile-time type-safe Go packages. It operates on plain [Ethereum contract ABIs](https://docs.soliditylang.org/en/develop/abi-spec.html) with expanded functionality if the contract bytecode is also available. However, it also accepts Solidity source files, making development much more streamlined. Please see our [Native DApps](https://geth.ethereum.org/docs/developers/dapp-developer/native-bindings) page for details.                                  |
|   `evm`    | Developer utility version of the EVM (Ethereum Virtual Machine) that is capable of running bytecode snippets within a configurable environment and execution mode. Its purpose is to allow isolated, fine-grained debugging of EVM opcodes (e.g. `evm --code 60ff60ff --debug run`).                                                                                                                                                                                                                                               |
| `rlpdump`  | Developer utility tool to convert binary RLP ([Recursive Length Prefix](https://ethereum.org/en/developers/docs/data-structures-and-encoding/rlp)) dumps (data encoding used by the Ethereum protocol both network as well as consensus wise) to user-friendlier hierarchical representation (e.g. `rlpdump --hex CE0183FFFFFFC4C304050583616263`).                                                                                                                                                                                |

## Running `geth`

Going through all the possible command line flags is out of scope here (please consult our
[CLI Wiki page](https://geth.ethereum.org/docs/fundamentals/command-line-options)),
but we've enumerated a few common parameter combos to get you up to speed quickly
on how you can run your own `geth` instance.

### Hardware Requirements

Minimum:

* CPU with 4+ cores
* 8GB RAM
* 1TB free storage space to sync the Mainnet
* 8 MBit/sec download Internet service

Recommended:

* Fast CPU with 8+ cores
* 16GB+ RAM
* High-performance SSD with at least 1TB of free space
* 25+ MBit/sec download Internet service



### Configuration

As an alternative to passing the numerous flags to the `geth` binary, you can also pass a
configuration file via:

```shell
$ geth --config /path/to/your_config.toml
```

To get an idea of how the file should look like you can use the `dumpconfig` subcommand to
export your existing configuration:

```shell
$ geth --your-favourite-flags dumpconfig
```


### Programmatically interfacing `geth` nodes


HTTP based JSON-RPC API options:

  * `--http` Enable the HTTP-RPC server
  * `--http.addr` HTTP-RPC server listening interface (default: `localhost`)
  * `--http.port` HTTP-RPC server listening port (default: `8545`)
  * `--http.api` API's offered over the HTTP-RPC interface (default: `eth,net,web3`)
  * `--http.corsdomain` Comma separated list of domains from which to accept cross-origin requests (browser enforced)
  * `--ws` Enable the WS-RPC server
  * `--ws.addr` WS-RPC server listening interface (default: `localhost`)
  * `--ws.port` WS-RPC server listening port (default: `8546`)
  * `--ws.api` API's offered over the WS-RPC interface (default: `eth,net,web3`)
  * `--ws.origins` Origins from which to accept WebSocket requests
  * `--ipcdisable` Disable the IPC-RPC server
  * `--ipcpath` Filename for IPC socket/pipe within the datadir (explicit paths escape it)

You'll need to use your own programming environments' capabilities (libraries, tools, etc) to
connect via HTTP, WS or IPC to a `geth` node configured with the above flags and you'll
need to speak [JSON-RPC](https://www.jsonrpc.org/specification) on all transports. You
can reuse the same connection for multiple requests!

**Note: Please understand the security implications of opening up an HTTP/WS based
transport before doing so! Hackers on the internet are actively trying to subvert
Ethereum nodes with exposed APIs! Further, all browser tabs can access locally
running web servers, so malicious web pages could try to subvert locally available
APIs!**

