# eddb2sqlite

Downloads and import to the SQLite important info from eddb.io to porvide possibility to get best prices of any commodity.

Usage:
- download
- go get (gcc/mingw required in PATH to build sqlite lib)
- to update full DB: `rm -rf ./data; rm -rf ./db/eddb.sqlite`
- to update only prices: `rm -rf ./data/listings*; rm -rf ./db/eddb.sqlite`
- go run main.go
- open result DB in any SQLite DB Browser

Sample SQL (not optimal but enough):
```
SELECT cmd.name as cmd, lst.supply, lst.buy_price, st.name as st, st.max_landing_pad_size as stpad, sm.name as smname, datetime(lst.collected_at, 'unixepoch') as collected_dt
FROM listings as lst
LEFT JOIN commodities as cmd ON cmd.eddb_id = lst.commodity_id
LEFT JOIN stations as st ON st.eddb_id = lst.station_id
LEFT JOIN systems as sm ON sm.eddb_id = st.system_id
WHERE commodity_id IN (SELECT eddb_id FROM commodities WHERE name = 'Indite' or name = 'Gallite' or name = 'Bertrandite' or name = 'Gold' or name = 'Silver')
AND lst.supply > 0
AND lst.buy_price < 6000
AND st.type_id <> 24
AND st.max_landing_pad_size = 'L'
ORDER BY supply DESC
```

Sample result:
