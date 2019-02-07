FROM tendermint/gaia:develop

COPY validator_setup /usr/bin/validator_setup

ARG pwd

RUN chmod 755 /usr/bin/validator_setup && validator_setup

EXPOSE 26656
EXPOSE 26657

ENTRYPOINT ["gaiad"]
CMD ["start"]