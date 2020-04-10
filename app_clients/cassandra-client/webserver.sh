# Init wait time in seconds
I_W_S=${INIT_WAIT_SECONDS:-10}
# Port of the liveness service
L_S_P=${LIVENESS_SVC_PORT:-8088}

# init sleep before business logic updates status
# typically INIT_WAIT_SECONDS
sleep $I_W_S

# start webserver
while true 
do 
   { printf 'HTTP/1.0 200 OK\r\nContent-Length: %d\r\n\r\n' "$(wc -c < /var/tmp/action.status)"; cat /var/tmp/action.status; } | nc -l $L_S_P
done 
