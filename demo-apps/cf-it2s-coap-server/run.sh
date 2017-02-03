#!/bin/sh
java -jar ./target/cf-it2s-coap-server-1.0.0-SNAPSHOT.jar >> "debug$(date '+-%d-%m-%y_%H-%M-%S').log" 2>&1 &
