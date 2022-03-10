#!/bin/sh

# Update the dynamic environment to include the values from the container.
node envToJson.js | tee public/dynamicEnv.json dist/dynamicEnv.json
http-server dist
