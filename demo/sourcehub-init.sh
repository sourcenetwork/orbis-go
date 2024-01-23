#!/bin/sh

MONIKER=demo
VALIDATOR1=validator1
VALIDATOR2=validator2
VALIDATOR3=validator3

sourcehubd init $MONIKER --chain-id sourcehub

sourcehubd keys add ${VALIDATOR1} --keyring-backend test
sourcehubd keys add ${VALIDATOR2} --keyring-backend test
sourcehubd keys add ${VALIDATOR3} --keyring-backend test

VALIDATOR1_ADDRESS=$(sourcehubd keys show ${VALIDATOR1} --address --keyring-backend test)
VALIDATOR2_ADDRESS=$(sourcehubd keys show ${VALIDATOR2} --address --keyring-backend test)
VALIDATOR3_ADDRESS=$(sourcehubd keys show ${VALIDATOR3} --address --keyring-backend test)

sourcehubd genesis add-genesis-account $VALIDATOR1_ADDRESS 100000000stake
sourcehubd genesis add-genesis-account $VALIDATOR2_ADDRESS 100000000stake
sourcehubd genesis add-genesis-account $VALIDATOR3_ADDRESS 100000000stake

sourcehubd genesis gentx ${VALIDATOR1} 70000000stake --chain-id sourcehub --keyring-backend test
sourcehubd genesis gentx ${VALIDATOR2} 70000000stake --chain-id sourcehub --keyring-backend test
sourcehubd genesis gentx ${VALIDATOR3} 70000000stake --chain-id sourcehub --keyring-backend test

sourcehubd genesis collect-gentxs

sed -i -e 's/timeout_commit = "5s"/timeout_commit = "1s"/g' /root/.sourcehub/config/config.toml
cat /root/.sourcehub/config/config.toml
