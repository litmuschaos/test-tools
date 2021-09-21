#! /bin/sh
echo "Jmeter Query Load Generation For PostgreSQL Application"
echo $DATA_SOURCE_URI

sh ./jmeter -JpasswordU=$PASSWORD -Jusername=$USERNAME -JdsURI=$DATA_SOURCE_URI -n -t postgres-load.jmx -l result.jtl