#!/bin/bash

# Chaos toolkit litmus local package installation
declare -a chaos_litmus_packages=("chaos")
for chaos_litmus_package in "${chaos_litmus_packages[@]}"
do
  pwd
#  mkdir /app/"$chaos_litmus_package"/
#  cd "$chaos_litmus_package"
#  ls -lrt
#  cp -rf . /app/"$chaos_litmus_package"/
  cd /app/"$chaos_litmus_package"/
  pwd
  ls -ltr
  python setup.py develop
  pip install -U .
done

# Preserve order for chaostest and lib in the beginning, thats the core
declare -a chaosexperiments=("chaostoolkit" "chaostoolkit-lib" "chaostoolkit-kubernetes" "chaostoolkit-reporting")
for chaosexperiment in "${chaosexperiments[@]}"
do
  ls -lrt
  mkdir /app/"$chaosexperiment"/
  cd /app/"$chaosexperiment"
  pip install --no-cache-dir -U "$chaosexperiment"
  ls -ltr
done

# For json path and other custom packages you can use the below
declare -a packages=("jsonpath2")
for package in "${packages[@]}"
do
  pip install --no-cache-dir -U "$package"
  ls -ltr
done



rm -rf /tmp/* /root/.cache




#nohup kubectl proxy --port=8080 &>/dev/null &
#wait

