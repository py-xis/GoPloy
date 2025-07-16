#!/bin/bash

echo "Cloning $GIT_REPOSITORY__URL..."
git clone "$GIT_REPOSITORY__URL" /home/app/output

echo "Running build-server..."
exec /home/app/build-server