#!/bin/sh

echo "================"
echo " Build SafeHODL"
echo "================"

rm ~/.safehodl 2>/dev/null
rm ./../dist/safehodl 2>/dev/null

echo "🛠  Building executable binary…"
garble -literals -tiny -seed=random build -ldflags "-extldflags=-static" -o dist/safehodl .
echo "✅ Successfully built! the build is in \"dist/safehodl\" ready to be used."
