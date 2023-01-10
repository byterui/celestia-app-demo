#!/bin/sh

VALIDATOR_NAME=validator1
CHAIN_ID=celestia-app-demod
KEY_NAME=ballman
CHAINFLAG="--chain-id ${CHAIN_ID}"
TOKEN_AMOUNT="10000000000000000000000000stake"
STAKING_AMOUNT="1000000000stake"
NODEIP="--node http://127.0.0.1:26657"

NAMESPACE_ID=$(echo $RANDOM | md5sum | head -c 16; echo;)
echo $NAMESPACE_ID
DA_BLOCK_HEIGHT=$(curl https://rpc.limani.celestia-devops.dev/block | jq -r '.result.block.header.height')
echo $DA_BLOCK_HEIGHT

./celestia-app-demod tendermint unsafe-reset-all
./celestia-app-demod init $VALIDATOR_NAME --chain-id $CHAIN_ID

./celestia-app-demod keys add $KEY_NAME --keyring-backend test
./celestia-app-demod add-genesis-account $KEY_NAME $TOKEN_AMOUNT --keyring-backend test
./celestia-app-demod gentx $KEY_NAME $STAKING_AMOUNT --chain-id $CHAIN_ID --keyring-backend test
./celestia-app-demod collect-gentxs
./celestia-app-demod start --rollmint.aggregator true --rollmint.da_layer celestia --rollmint.da_config='{"base_url":"http://localhost:26658","timeout":60000000000,"gas_limit":6000000}' --rollmint.namespace_id $NAMESPACE_ID --rollmint.da_start_height $DA_BLOCK_HEIGHT 