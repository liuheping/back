-- +goose Up
-- 分词
-- 全文索引的实现要靠 PgSQL 的 gin 索引。分词功能 PgSQL 内置了英文、西班牙文等，但中文分词需要借助开源插件 zhparser；
--
-- SCWS
-- 要使用 zhparser，我们首先要安装 SCWS 分词库，SCWS 是 Simple Chinese Word Segmentation 的首字母缩写（即：简易中文分词系统），其 GitHub 项目地址为 https://github.com/hightman/scws，我们下载之后可以直接安装。
--
-- 安装完后，就可以在命令行中使用 scws 命令进行测试分词了， 其参数主要有：
--
-- -c utf8 指定字符集
-- -d dict 指定字典 可以是 xdb 或 txt 格式
-- -M 复合分词的级别， 1~15，按位异或的 1|2|4|8 依次表示 短词|二元|主要字|全部字，默认不复合分词，这个参数可以帮助调整到最想要的分词效果。
-- zhpaser
-- 下载 zhparser 源码 git clone https:github.com/amutu/zhparser.git；
-- 安装前需要先配置环境变量：export PATH=$PATH:/path/to/pgsql；
-- make && make install编译 zhparser；
-- centos 报错：
-- Makefile:18: /usr/pgsql-10/lib/pgxs/src/makefiles/pgxs.mk: 没有那个文件或目录
-- make: *** 没有规则可以创建目标“/usr/pgsql-10/lib/pgxs/src/makefiles/pgxs.mk”。 停止
-- 解决方法：
-- yum install postgresql10-devel
-- 登陆 PgSQL 使用 CREATE EXTENSION zhparser; 启用插件；
-- 添加分词配置
--
-- CREATE TEXT SEARCH CONFIGURATION parser_name (PARSER = zhparser); // 添加配置
-- ALTER TEXT SEARCH CONFIGURATION parser_name ADD MAPPING FOR n,v,a,i,e,l,j WITH simple; // 设置分词规则 （n 名词 v 动词等，详情阅读下面的文档）
-- 给某一列的分词结果添加 gin 索引 create index idx_name on table using gin(to_tsvector('parser_name', field));
-- 启用插件；
CREATE EXTENSION zhparser;
-- 添加分词配置
CREATE TEXT SEARCH CONFIGURATION zhparser (PARSER = zhparser);
-- 设置分词规则
ALTER TEXT SEARCH CONFIGURATION zhparser ADD MAPPING FOR n,v,a,i,e,l,j WITH simple;
-- +goose Down
DROP EXTENSION zhparser;
DROP TEXT SEARCH CONFIGURATION zhparser;