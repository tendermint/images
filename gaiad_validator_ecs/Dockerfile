FROM tendermint/gaia:develop

COPY validator_setup /usr/bin/validator_setup
COPY gaiad_cp_conf_start /usr/bin/gaiad_cp_conf_start

RUN chmod 755 /usr/bin/validator_setup && validator_setup
RUN chmod 755 /usr/bin/gaiad_cp_conf_start

EXPOSE 26656
EXPOSE 26657

ENTRYPOINT ["gaiad_cp_conf_start"]
