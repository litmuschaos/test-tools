#!/bin/bash

# Preserve order for chaostoolkit and lib in the beginning, thats the core
declare -a chaosexperiments=("chaostoolkit" "chaostoolkit-lib" "chaostoolkit-kubernetes" "chaostoolkit-reporting")
for chaosexperiment in "${chaosexperiments[@]}"
do
  ls -lrt
  mkdir /app/"$chaosexperiment"/
  cd /app/"$chaosexperiment"
  pip install --no-cache-dir -U "$chaosexperiment"
  ls -ltr
done
rm -rf /tmp/* /root/.cache

#nohup kubectl proxy --port=8080 &>/dev/null &
#wait

