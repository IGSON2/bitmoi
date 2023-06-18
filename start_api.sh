#!/bin/sh

migrate -path backend/db/migrate -database "mysql://root:123@tcp(db:3306)/bitmoi" -verbose up

if [ $? -eq 0 ]; then
    echo migration has done successfully
else
    echo magration failed
fi

echo fetch candles..

/bitmoi/bitmoi store --interval=1h --pairs=BTC,ETH,ADA --timestamp=1682899200000

/bitmoi/bitmoi

if [ $? -eq 0 ]; then
    echo run bitmoi api server PID : $!
else
    echo failed to start bitmoi api server.
fi