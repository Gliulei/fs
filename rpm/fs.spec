#软件包的名字
Name: fs
Group: Development/Tools
BuildRoot: %_topdir/BUILDROOT

#软件包的内容介绍
%description
upload or download file tool

#BUILD字段，将通过直接调用源码目录中自动构建工具完成源码编译操作
%build
cd nginx-1.2.1

#调用源码目录中的configure命令
./configure --prefix=/usr/local/nginx

#在源码目录中执行自动构建命令make
make

#安装字段
%install

#调用源码中安装执行脚本
make install
%preun
if [ -z "`ps aux | grep nginx | grep -v grep`" ];then
killall nginx >/dev/null
exit 0
fi

#文件说明字段，声明多余或者缺少都将可能出错
%files
#声明/usr/local/nginx将出现在软件包中/usr/local/nginx