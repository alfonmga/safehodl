#!/bin/sh

AES_SECRET_32_BYTES_KEY=$(hexdump -n 16 -e '4/4 "%08X" 1 "\n"' /dev/urandom)

echo "================"
echo " Build SafeHODL"
echo "================"

rm ~/.safehodl 2>/dev/null
rm ./../dist/safehodl 2>/dev/null

echo "Set a PIN code for secure access to SafeHODL:"
read -s pincode

if ! [ ${#pincode} -ge 1 ]; then
    echo "ERROR: PIN code is too short."
    exit 0
fi

echo "Building…"
garble -literals -tiny -seed=random build -ldflags "-extldflags=-static -X main.Secret32BytesKeyAES=$AES_SECRET_32_BYTES_KEY -X main.PinCode=$pincode" -o dist/safehodl .
echo "Successfully built! the build is in \"dist/safehodl\" ready to be used."
