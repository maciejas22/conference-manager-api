#!/bin/bash

source .env
source .env.dev

SEED_DIR="./migrations/seed"

for file in $SEED_DIR/*.sql; do
  echo "Applying seed: $file"
  psql $DATABASE_URL -f $file
done
