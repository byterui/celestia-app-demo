#!/bin/sh

VALIDATOR_NAME=validator1
CHAIN_ID=celestia-app-demod
KEY_NAME=ballman
CHAINFLAG="--chain-id ${CHAIN_ID}"
TOKEN_AMOUNT="10000000000000000000000000demo"
STAKING_AMOUNT="1000000000demo"
NODEIP="--node http://127.0.0.1:26657"

 NAMESPACE_ID=$(echo $RANDOM | md5sum | head -c 16; echo;)
# NAMESPACE_ID=0a5ec754e3041feb
echo $NAMESPACE_ID
DA_BLOCK_HEIGHT=$(curl http://172.19.0.2:26657/block | jq -r '.result.block.header.height')
# DA_BLOCK_HEIGHT=440
echo $DA_BLOCK_HEIGHT

./celestia-app-demod tendermint unsafe-reset-all
./celestia-app-demod init $VALIDATOR_NAME --chain-id $CHAIN_ID

# ./celestia-app-demod keys add $KEY_NAME --keyring-backend test
./celestia-app-demod add-genesis-account $KEY_NAME $TOKEN_AMOUNT --keyring-backend test
./celestia-app-demod gentx $KEY_NAME $STAKING_AMOUNT --chain-id $CHAIN_ID --keyring-backend test

./celestia-app-demod gentx ballman 1000000000demo --chain-id  celestia-app-demod --keyring-backend test

./celestia-app-demod collect-gentxs
./celestia-app-demod start --rpc.laddr tcp://127.0.0.1:46658 --rollkit.aggregator \
true --rollkit.da_layer celestia --rollkit.da_config='{"base_url":"http://172.19.0.4:26659","timeout":60000000000,"fee":6000,"gas_limit":6000000}' --rollkit.namespace_id $NAMESPACE_ID 