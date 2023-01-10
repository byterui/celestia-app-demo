# celestiaappdemo
**celestiaappdemo** is a blockchain built using Cosmos SDK and Tendermint and created with [Ignite CLI](https://ignite.com/cli).

## Get started   
step1: run celestia light client(use mocha testnet) (https://docs.celestia.org/nodes/light-node)
```
create wallet command:
1. cel-key add ballman  --keyring-backend test --node.type light
2. require test coin from faucet:(https://discord.com/channels/638338779505229824/1018935905991536710)

run light node command:
celestia light init

celestia light start --core.ip https://rpc-mocha.pops.one:9090 --keyring.accname ballman

```  
step2: 
```
ignite chain build -o ./
./init.sh
````

## use example
### Query Balance
get native balance command
```
./celestia-app-demod query bank balances cosmos1jckkk484daz7dvqlp5laa7m733erwd5847mpqs
```
response
```
balances:
- amount: "9999999999999999000000000"
  denom: stake
pagination:
  next_key: null
  total: "0"

```
### Send coin
```
create new account: ./celestia-app-demod keys add goldman --keyring-backend test

./celestia-app-demod tx bank send cosmos1jckkk484daz7dvqlp5laa7m733erwd5847mpqs cosmos199uds665u99ux3f5lmzlk9sslepkxq5e40mqaw  1899999999999999000000000stake --keyring-backend test --chain-id celestia-app-demod
```