#!/bin/bash -x

source /etc/profile.d/set_env.sh

cd ${GOPATH}/src/github.com/cosmos/cosmos-sdk
mkdir ${HOME}/sim-logs

#seeds=(1 2 4 7 9 20 32 123 124 582 1893 2989 3012 4728 37827 981928 87821 891823782 989182 89182391 \
#11 22 44 77 99 2020 3232 123123 124124 582582 18931893 29892989 30123012 47284728 37827)

sim() {
  seed=$1
  file="${HOME}/sim-logs/gaia-simulation-seed-${seed}-date-$(date -u +"%Y-%m-%dT%H:%M:%S+00:00").stdout"
  go test ./cmd/gaia/app -run TestFullGaiaSimulation \
                         -SimulationEnabled=true \
                         -SimulationNumBlocks=${BLOCKS} \
                         -SimulationVerbose=true \
                         -SimulationCommit=true \
                         -SimulationSeed=$seed \
                         -SimulationPeriod=${PERIOD} \
                         -v -timeout 24h > $file
}

i=0
pids=()
for seed in ${SEEDS[@]}; do
  sim $seed &
  pids[${i}]=$!
  i=$(($i+1))
  sleep 10
done

code=0
i=0
for pid in ${pids[@]}; do
  wait $pid
  last=$?
  seed=${SEEDS[${i}]}
  if [[ $last -ne 0 ]]
  then
    go run notify_slack.go 1 $seed > "${HOME}/sim-logs/slack-$seed"
  else
    echo "WAT"
    go run notify_slack.go 0 $seed > "${HOME}/sim-logs/slack-$seed"
  fi
  i=$(($i+1))
done

sudo shutdown -h now
