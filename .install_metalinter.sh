#!/usr/bin/env sh

version=v2.0.2
name=gometalinter-"$version"-linux-amd64
dest=gometalinter

if test -e "$dest/$version"; then
  exit 0
fi

rm -rf "$dest"

wget "https://github.com/alecthomas/gometalinter/releases/download/$version/$name.tar.bz2"
tar xf "$name.tar.bz2"
mv "$name" "$dest"
touch "$dest/$version"
