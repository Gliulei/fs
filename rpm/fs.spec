Name: fs
Version: %(echo $VERSION)
Release: %(echo $RELEASE)%{?dist}
Summary: fs
Group: Development/Tools
BuildRoot: %_topdir/BUILDROOT
License: GPL

#软件包的内容介绍
%description
upload or download file tool

#安装字段
%install
SRC_DIR=$OLDPWD/..
cd $SRC_DIR/

make build
echo ${RPM_BUILD_ROOT}
mkdir -p ${RPM_BUILD_ROOT}/usr/local/bin
cp -rf $SRC_DIR/fs ${RPM_BUILD_ROOT}/usr/local/bin/fs

# package infomation
%files
# set file attribute here
%defattr(-,root,root,0755)
# need not list every file here, keep it as this
/usr/local/bin/fs

## create an empy dir

%post
# description: fs ....
chmod +x /usr/local/bin/fs