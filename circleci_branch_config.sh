#!/bin/bash
set -euo pipefail

case $CIRCLE_BRANCH in
  'master')
    export STAGE=production
    export NODE_ENV=production
  ;;
  *)
  echo "no configuration for $CIRCLE_BRANCH"
  exit 1
  ;;
esac

echo STAGE is $STAGE
echo NODE_ENV is $NODE_ENV
