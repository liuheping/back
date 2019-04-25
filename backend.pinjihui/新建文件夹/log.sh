#!/bin/sh
cp -f /home/www/go/src/pinjihui.com/backend.pinjihui/nohup.out /home/www/go/src/pinjihui.com/backend.pinjihui/log/$(date +%Y-%m-%d).out
cat /dev/null > nohup.out
