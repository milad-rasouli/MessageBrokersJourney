#!/bin/bash

function encoded_password()
{
    SALT=$(od -A n -t x -N 4 /dev/random)
    PASS=$SALT$(echo -n $1 | xxd -ps | tr -d '\n' | tr -d ' ')
    PASS=$(echo -n $PASS | xxd -r -p | sha256sum |head -c 128)
    PASS=$(echo -n $SALT$PASS | xxd -r -p | base64 | tr -d '\n')
    echo $PASS
}

# put your own password below between the "here"
encoded_password "1234qwer"