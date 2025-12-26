#! /bin/bash
# import data from a Kostal LogDaten.dat file
set -ex

# curl -u <user>:<password> 192.168.2.125/LogDaten.dat > ~/LogDaten.dat
cat ~/LogDaten.dat | tail -n+8 | tr '\t' ',' | tr -d ' ' | awk -F, -v OFS=, -v ORS=: '{NF=42}1' | tr -d '\r' | tr ':' '\n' > ~/LogDaten.csv
docker cp ~/LogDaten.csv solar_db:/var/lib/postgresql/data/LogDaten.csv
docker exec -it solar_db psql -U raskob -d postgres -c "\copy pvdata FROM '/var/lib/postgresql/data/LogDaten.csv' DELIMITER ',' CSV HEADER;"
