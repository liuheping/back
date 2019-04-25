#!/bin/sh
cd /home/www/go/src/pinjihui.com/backend.pinjihui/code/backend.pinjihui
git pull
kill -9 $(pidof backend.pinjihui)
cp -f /home/www/go/src/pinjihui.com/backend.pinjihui/code/backend.pinjihui/Config.toml /home/www/go/src/pinjihui.com/backend.pinjihui
cp -f /home/www/go/src/pinjihui.com/backend.pinjihui/code/backend.pinjihui/backend.pinjihui /home/www/go/src/pinjihui.com/backend.pinjihui
cd /home/www/go/src/pinjihui.com/backend.pinjihui
nohup ./backend.pinjihui > nohup.out 2>&1&

cd /home/www/go/src/pinjihui.com/backend.pinjihui/code/backend.pinjihui
kill -9 $(pidof backend.pinjihui_online)
cp -f /home/www/go/src/pinjihui.com/backend.pinjihui/code/backend.pinjihui/Config_online.toml /home/www/go/src/pinjihui.com/backend.pinjihui_online
cp -f /home/www/go/src/pinjihui.com/backend.pinjihui/code/backend.pinjihui/backend.pinjihui /home/www/go/src/pinjihui.com/backend.pinjihui_online
cd /home/www/go/src/pinjihui.com/backend.pinjihui_online
mv Config_online.toml Config.toml
mv backend.pinjihui backend.pinjihui_online
nohup ./backend.pinjihui_online > nohup.out 2>&1&