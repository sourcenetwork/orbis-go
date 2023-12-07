#!/bin/sh

MONIKER=demo
VALIDATOR=validator

sourcehubd init $MONIKER --chain-id sourcehub

sourcehubd keys add ${VALIDATOR} --keyring-backend test

VALIDATOR_ADDRESS=$(sourcehubd keys show ${VALIDATOR} --address --keyring-backend test)

sourcehubd genesis add-genesis-account $VALIDATOR_ADDRESS 100000000stake

sourcehubd genesis gentx ${VALIDATOR} 70000000stake --chain-id sourcehub --keyring-backend test

sourcehubd genesis collect-gentxs
