#!/bin/bash -x

source /etc/profile.d/set_env.sh

cd ${GOPATH}/src/github.com/cosmos/cosmos-sdk
mkdir ${HOME}/sim-logs

sim() {
  seed=$1
  file="${HOME}/sim-logs/gaia-simulation-seed-${seed}.stdout"
  go test -mod=readonly ./cmd/gaia/app -run TestFullGaiaSimulation \
                                       -SimulationEnabled=true \
                                       -SimulationNumBlocks=${BLOCKS} \
                                       -SimulationVerbose=true \
                                       -SimulationCommit=true \
                                       -SimulationSeed=$seed \
                                       -SimulationPeriod=${PERIOD} \
                                       -v -timeout 24h > ${file}
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
    go run notify_slack.go 1 ${seed} ${SLACK_MSG_TS} > "${HOME}/sim-logs/slack-$seed"
  fi
  i=$(($i+1))
done

go run notify_slack.go 2 "${args[*]}" ${SLACK_MSG_TS} > "${HOME}/sim-logs/slack-done"

sudo shutdown -h now
