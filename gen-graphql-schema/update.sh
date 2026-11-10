#!/usr/bin/env bash
# Regenerates ../schema.graphql and the HTML docs in ../doc/schema
# from the live domonda GraphQL API. Run it after every change
# of the API schema and commit the result.
set -euo pipefail
cd "$(dirname "$0")"

npm ci --no-audit --no-fund --ignore-scripts

node print-schema.mjs > ../schema.graphql.tmp
mv ../schema.graphql.tmp ../schema.graphql

# graphdoc --force replaces the whole output directory, logo included
npx --no-install graphdoc --force
mkdir -p ../doc/schema/img
cp ../img/domonda_logo_schriftzug.png ../doc/schema/img/
