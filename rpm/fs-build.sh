#!/bin/bash

CUR_DIR=$(dirname $(readlink -f "$0"))
PROJECT_DIR=$1
PROJECT_NAME=$2
VERSION=$3
RELEASE=$4

# build rpm
export PROJECT_NAME=${PROJECT_NAME}
export VERSION=${VERSION}
export RELEASE=${RELEASE}

rpmbuild_dir=$CUR_DIR/rpmbuild
rm -rf $rpmbuild_dir

rpmbuild --define "_topdir $rpmbuild_dir" -bb fs.spec
find $rpmbuild_dir/RPMS/-name "*.rpm" -exec '{}' $CUR_DIR \;