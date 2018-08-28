#!/bin/bash

# In this configuration, the following dependent libraries are compiled:
#
# * zlib
# * c-ares
# * expat
# * sqlite3
# * openSSL
# * libssh2

# COMPILER AND PATH
PREFIX=$PWD/aria2-lib
C_COMPILER="gcc"
CXX_COMPILER="g++"

# CHECK TOOL FOR DOWNLOAD
 aria2c --help > /dev/null
 if [ "$?" -eq 0 ] ; then
   DOWNLOADER="aria2c --check-certificate=false"
 else
   DOWNLOADER="wget -c"
 fi

## VERSION
ZLIB_V=1.2.11
OPENSSL_V=1.0.2o
EXPAT_V=2.2.6
SQLITE_V=3240000
C_ARES_V=1.14.0
SSH2_V=1.8.0
ARIA2_V=1.34.0

## DEPENDENCES
ZLIB=http://sourceforge.net/projects/libpng/files/zlib/${ZLIB_V}/zlib-${ZLIB_V}.tar.gz
OPENSSL=http://www.openssl.org/source/openssl-${OPENSSL_V}.tar.gz
EXPAT=https://sourceforge.net/projects/expat/files/expat/${EXPAT_V}/expat-${EXPAT_V}.tar.bz2
SQLITE3=http://www.sqlite.org/2018/sqlite-autoconf-${SQLITE_V}.tar.gz
C_ARES=http://c-ares.haxx.se/download/c-ares-${C_ARES_V}.tar.gz
SSH2=https://www.libssh2.org/download/libssh2-${SSH2_V}.tar.gz
ARIA2=https://github.com/aria2/aria2/releases/download/release-${ARIA2_V}/aria2-${ARIA2_V}.tar.bz2

## CONFIG
BUILD_DIRECTORY=/tmp/

## BUILD
cd ${BUILD_DIRECTORY}
#
 # zlib build
  ${DOWNLOADER} ${ZLIB}
  tar zxvf zlib-${ZLIB_V}.tar.gz
  cd zlib-${ZLIB_V}/
  PKG_CONFIG_PATH=${PREFIX}/lib/pkgconfig/ \
  LD_LIBRARY_PATH=${PREFIX}/lib/ CC="$C_COMPILER" CXX="$CXX_COMPILER" \
  ./configure --prefix=${PREFIX} --static
  make
  make install
#
 # expat build
  cd ..
  ${DOWNLOADER} ${EXPAT}
  tar jxvf expat-${EXPAT_V}.tar.bz2
  cd expat-${EXPAT_V}/
  PKG_CONFIG_PATH=${PREFIX}/lib/pkgconfig/ \
  LD_LIBRARY_PATH=${PREFIX}/lib/ CC="$C_COMPILER" CXX="$CXX_COMPILER" \
  ./configure --prefix=${PREFIX} --enable-static --enable-shared
  make
  make install
#
 # c-ares build
  cd ..
  ${DOWNLOADER} ${C_ARES}
  tar zxvf c-ares-${C_ARES_V}.tar.gz
  cd c-ares-${C_ARES_V}/
  PKG_CONFIG_PATH=${PREFIX}/lib/pkgconfig/ \
  LD_LIBRARY_PATH=${PREFIX}/lib/ CC="$C_COMPILER" CXX="$CXX_COMPILER" \
  ./configure --prefix=${PREFIX} --enable-static --disable-shared
  make
  make install
#
 # Openssl build
  cd ..
  ${DOWNLOADER} ${OPENSSL}
  tar zxvf openssl-${OPENSSL_V}.tar.gz
  cd openssl-${OPENSSL_V}/
  PKG_CONFIG_PATH=${PREFIX}/lib/pkgconfig/ \
  LD_LIBRARY_PATH=${PREFIX}/lib/ CC="$C_COMPILER" CXX="$CXX_COMPILER" \
  ./Configure --prefix=${PREFIX} linux-x86_64 shared
  make
  make install
#
 # sqlite3
  cd ..
  ${DOWNLOADER} ${SQLITE3}
  tar zxvf sqlite-autoconf-${SQLITE_V}.tar.gz
  cd sqlite-autoconf-${SQLITE_V}/
  PKG_CONFIG_PATH=${PREFIX}/lib/pkgconfig/ \
  LD_LIBRARY_PATH=${PREFIX}/lib/ CC="$C_COMPILER" CXX="$CXX_COMPILER" \
  ./configure --prefix=${PREFIX} --enable-static --enable-shared
  make
  make install
#
 # libssh2
  cd ..
  ${DOWNLOADER} ${SSH2}
  tar zxvf libssh2-${SSH2_V}.tar.gz
  cd libssh2-${SSH2_V}/
  rm -rf ${PREFIX}/lib/pkgconfig/libssh2.pc
  PKG_CONFIG_PATH=${PREFIX}/lib/pkgconfig/ \
  LD_LIBRARY_PATH=${PREFIX}/lib/ CC="$C_COMPILER" CXX="$CXX_COMPILER" \
  ./configure --with-libgcrypt --with-openssl \
    --without-wincng --prefix=${PREFIX} \
    --enable-static --disable-shared
  make
  make install
#

 #cleaning
  cd ..
  rm -rf c-ares*
  rm -rf sqlite-autoconf*
  rm -rf zlib-*
  rm -rf expat-*
  rm -rf openssl-*
  rm -rf libssh2-*
#

# Build aria2 static library.
${DOWNLOADER} ${ARIA2}
  tar jxvf aria2-${ARIA2_V}.tar.bz2
  cd aria2-${ARIA2_V}/
PKG_CONFIG_PATH=${PREFIX}/lib/pkgconfig/ \
LD_LIBRARY_PATH=${PREFIX}/lib/ \
CC="$C_COMPILER" \
CXX="$CXX_COMPILER" \
./configure \
    --prefix=${PREFIX} \
    --without-libxml2 \
    --without-libgcrypt \
    --with-openssl \
    --without-libnettle \
    --without-gnutls \
    --without-libgmp \
    --with-libssh2 \
    --with-sqlite3 \
    --enable-libaria2 \
    --enable-shared=no \
    --enable-static=yes
make
make install

## Generate files for c.
cd $PWD
go tool cgo libaria2.go

echo "Prepare finished!"
