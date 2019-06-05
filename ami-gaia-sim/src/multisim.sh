#!/bin/bash -x

source /etc/profile.d/set_env.sh

cd ${GOPATH}/src/github.com/cosmos/cosmos-sdk
mkdir ${HOME}/sim-logs

sim() {
  seed=$1
  file="${HOME}/sim-logs/gaia-sim-seed-${seed}.stdout"
  go test github.com/cosmos/cosmos-sdk/simapp -run TestFullAppSimulation \
                                              -SimulationEnabled=true \
                                              -SimulationNumBlocks=${BLOCKS} \
                                              -SimulationVerbose=true \
                                              -SimulationCommit=true \
                                              -SimulationSeed=$seed \
                                              -SimulationPeriod=${PERIOD} \
                                              -v -timeout 5h > ${file}
}

args=("$@")

i=0
pids=()
for seed in ${args[@]}; do
  sim ${seed} &
  pids[$i]=$!
  i=$(($i+1))
done

code=0
i=0
for pid in ${pids[@]}; do
  wait ${pid}
  last=$?
  seed=${args[$i]}
  if [[ ${last} -ne 0 ]]; then
    ./notify_slack -failed=true \
                   -num-blocks=${BLOCKS} \
                   -seed-num=${seed} \
                   -period=${PERIOD} \
                   -slack-token=${SLACK_TOKEN} \
                   -channel-id=${SLACK_CHANNEL_ID} > "${HOME}/sim-logs/slack-$seed"
  fi
  i=$(($i+1))
done

./notify_slack -ts=${SLACK_MSG_TS} \
               -slack-token=${SLACK_TOKEN} \
               -channel-id=${SLACK_CHANNEL_ID} \
               -seeds="${args[*]}" 2>&1 "${HOME}/sim-logs/slack-done"

sudo shutdown -h now
