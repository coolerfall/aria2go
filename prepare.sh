#!/bin/bash

# In this configuration, the following dependent libraries are compiled:
#
# * zlib
# * c-ares
# * expat
# * sqlite3
# * openSSL
# * libssh2

# Compiler and path
PREFIX=$PWD/aria2-lib
C_COMPILER="gcc"
CXX_COMPILER="g++"

# Check tool for download
aria2c --help > /dev/null
if [ "$?" -eq 0 ]; then
    DOWNLOADER="aria2c --check-certificate=false"
else
    DOWNLOADER="wget -c"
fi

echo "Remove old libs..."
rm -rf ${PREFIX}
rm -rf _obj

## Version
ZLIB_V=1.2.11
OPENSSL_V=1.0.2p
EXPAT_V=2.2.6
SQLITE_V=3240000
C_ARES_V=1.14.0
SSH2_V=1.8.0
ARIA2_V=1.34.0

## Dependencies
ZLIB=http://sourceforge.net/projects/libpng/files/zlib/${ZLIB_V}/zlib-${ZLIB_V}.tar.gz
OPENSSL=http://www.openssl.org/source/openssl-${OPENSSL_V}.tar.gz
EXPAT=https://sourceforge.net/projects/expat/files/expat/${EXPAT_V}/expat-${EXPAT_V}.tar.bz2
SQLITE3=http://www.sqlite.org/2018/sqlite-autoconf-${SQLITE_V}.tar.gz
C_ARES=http://c-ares.haxx.se/download/c-ares-${C_ARES_V}.tar.gz
SSH2=https://www.libssh2.org/download/libssh2-${SSH2_V}.tar.gz
ARIA2=https://github.com/aria2/aria2/releases/download/release-${ARIA2_V}/aria2-${ARIA2_V}.tar.bz2

## Config
BUILD_DIRECTORY=/tmp/

## Build
cd ${BUILD_DIRECTORY}

# zlib build
if ! [ -e zlib-${ZLIB_V}.tar.gz ]; then
    ${DOWNLOADER} ${ZLIB}
fi
tar zxvf zlib-${ZLIB_V}.tar.gz
cd zlib-${ZLIB_V}
PKG_CONFIG_PATH=${PREFIX}/lib/pkgconfig/ \
    LD_LIBRARY_PATH=${PREFIX}/lib/ CC="$C_COMPILER" CXX="$CXX_COMPILER" \
    ./configure --prefix=${PREFIX} --static
make -j4
make install

# expat build
cd ..
if ! [ -e expat-${EXPAT_V}.tar.bz2 ]; then
    ${DOWNLOADER} ${EXPAT}
fi
tar jxvf expat-${EXPAT_V}.tar.bz2
cd expat-${EXPAT_V}
PKG_CONFIG_PATH=${PREFIX}/lib/pkgconfig/ \
    LD_LIBRARY_PATH=${PREFIX}/lib/ CC="$C_COMPILER" CXX="$CXX_COMPILER" \
    ./configure --prefix=${PREFIX} --enable-static --enable-shared
make -j4
make install

# c-ares build
cd ..
if ! [ -e c-ares-${C_ARES_V}.tar.gz ]; then
    ${DOWNLOADER} ${C_ARES}
fi
tar zxvf c-ares-${C_ARES_V}.tar.gz
cd c-ares-${C_ARES_V}
PKG_CONFIG_PATH=${PREFIX}/lib/pkgconfig/ \
    LD_LIBRARY_PATH=${PREFIX}/lib/ CC="$C_COMPILER" CXX="$CXX_COMPILER" \
    ./configure --prefix=${PREFIX} --enable-static --disable-shared
make -j4
make install

# openssl build
cd ..
if ! [ -e openssl-${OPENSSL_V}.tar.gz ]; then
    ${DOWNLOADER} ${OPENSSL}
fi
tar zxvf openssl-${OPENSSL_V}.tar.gz
cd openssl-${OPENSSL_V}
PKG_CONFIG_PATH=${PREFIX}/lib/pkgconfig/ \
    LD_LIBRARY_PATH=${PREFIX}/lib/ CC="$C_COMPILER" CXX="$CXX_COMPILER" \
    ./Configure --prefix=${PREFIX} linux-x86_64 shared
make -j4
make install

# sqlite3
cd ..
if ! [ -e sqlite-autoconf-${SQLITE_V}.tar.gz ]; then
    ${DOWNLOADER} ${SQLITE3}
fi
tar zxvf sqlite-autoconf-${SQLITE_V}.tar.gz
cd sqlite-autoconf-${SQLITE_V}
PKG_CONFIG_PATH=${PREFIX}/lib/pkgconfig/ \
    LD_LIBRARY_PATH=${PREFIX}/lib/ CC="$C_COMPILER" CXX="$CXX_COMPILER" \
    ./configure --prefix=${PREFIX} --enable-static --enable-shared
make -j4
make install

# libssh2 build
cd ..
if ! [ -e libssh2-${SSH2_V}.tar.gz ]; then
    ${DOWNLOADER} ${SSH2}
fi
tar zxvf libssh2-${SSH2_V}.tar.gz
cd libssh2-${SSH2_V}
PKG_CONFIG_PATH=${PREFIX}/lib/pkgconfig/ \
    LD_LIBRARY_PATH=${PREFIX}/lib/ CC="$C_COMPILER" CXX="$CXX_COMPILER" \
    ./configure --without-libgcrypt-prefix --with-openssl \
    --without-wincng --prefix=${PREFIX} \
    --enable-static --disable-shared
make -j4
make install

# Build aria2 static library.
cd ..
if ! [ -e aria2-${ARIA2_V}.tar.bz2 ]; then
    ${DOWNLOADER} ${ARIA2}
fi
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
make -j4
make install

# Cleaning
cd ..
rm -rf zlib-${ZLIB_V}
rm -rf expat-${EXPAT_V}
rm -rf c-ares-${C_ARES_V}
rm -rf openssl-${OPENSSL_V}
rm -rf sqlite-autoconf-${SQLITE_V}
rm -rf libssh2-${SSH2_V}
rm -rf aria2-${ARIA2_V}
rm -rf ${PREFIX}/bin

## Generate files for c
cd ${PREFIX}/../
cp ${PREFIX}/include/aria2/aria2.h ../
go tool cgo libaria2.go

echo "Prepare finished!"
