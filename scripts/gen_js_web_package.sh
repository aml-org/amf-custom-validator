#!/usr/bin/env bash
cp wrappers/js/lib/main.wasm.gz wrappers/js-web/lib/main.wasm.gz
cd wrappers/js-web
npm install
npm run build:dist
