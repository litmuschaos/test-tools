# Pause

Pause is a cli tool that can be used to begin a commands in paused state. It helps to perform some operations before the actual command 
execution like loading `cgroup`. 

## Usage

```sh
mv pause /usr/local/bin
pause nsutil -p -n -t 39590 -- stress-ng -c 2 -t 60s
```
