#!/bin/sh

echo "Debug mode: Container is running but service not started"
echo "Environment variables:"
env
echo "\nCertificate files:"
ls -la /app/certs/
echo "\nConfiguration file:"
cat $CONFIG_FILE

# Keep container running
tail -f /dev/null 