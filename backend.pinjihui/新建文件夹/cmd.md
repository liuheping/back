#### 设置代理等常用命令
    export PATH=~/go/bin:$PATH
    export http_proxy=http://127.0.0.1:8118
    export https_proxy=http://127.0.0.1:8118
    go generate ./schema

#### 包管理 
    dep ensure

#### Linux64位环境编译
    GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -o backend.pinjihui

#### 服务器登陆
    ssh root@39.108.97.201

#### 常用拷贝
    scp backend.pinjihui root@39.108.97.201:/home/www/go/src/pinjihui.com/backend.pinjihui
    scp Config.toml root@39.108.97.201:/home/www/go/src/pinjihui.com/backend.pinjihui

#### 程序路径
    /home/www/go/src/pinjihui.com/backend.pinjihui

#### 执行SQL脚本
    psql -h localhost -d pinjihui -U postgres -f test_data_v47.sql

#### 七牛云帐号密码
    2117390@qq.com
    pinjihui

#### 日志查看
    tailf nohup.out

#### 服务器仓库
     git clone git@test.pinjihui.com:/home/git/backend.pinjihui.git

#### 数据库版本管理
    ./cmd/goose --dir=./data/1.0/migrations-sql postgres "host=39.108.97.201 user=postgres dbname=pinjihui password=pinjihui\!@#123 sslmode=disable" up

    ./cmd/goose --dir=./data/1.0/migrations-sql postgres "host=39.108.97.201 user=postgres dbname=pinjihui_online password=pinjihui\!@#123 sslmode=disable" up

#### 启动ss
    sslocal -c  /home/hp/下载/shadowsocks.json