#!/bin/sh

ROOT=$(cd "$(dirname "$0")" && pwd)

set -uex

cd "$ROOT"
./gradlew build

cd "$ROOT/app/build/outputs/apk/debug/"
rm -rf app-debug
mkdir app-debug
cd app-debug
unzip ../app-debug.apk

cp AndroidManifest.xml "$ROOT"
cp resources.arsc "$ROOT"
