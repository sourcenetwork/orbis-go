#!/bin/sh

SOURCEHUBD=/ko-app/sourcehubd
MONIKER=demo
VALIDATOR=validator

sourcehubd init $MONIKER --chain-id sourcehub

sourcehubd keys add ${VALIDATOR} --keyring-backend test

VALIDATOR_ADDRESS=$(sourcehubd keys show ${VALIDATOR} --address --keyring-backend test)

sourcehubd add-genesis-account $VALIDATOR_ADDRESS 100000000stake

sourcehubd gentx ${VALIDATOR} 70000000stake --chain-id sourcehub --keyring-backend test

sourcehubd collect-gentxs
