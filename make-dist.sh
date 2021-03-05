#!/bin/bash

rm -rf dist 
mkdir -p dist

cd images
dirs=$(ls -d *)
for dir in $dirs; do 
    zip -r ../dist/${dir}.zip ${dir}
done 
