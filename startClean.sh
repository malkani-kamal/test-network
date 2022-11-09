./network.sh down
./network.sh up -ca -s couchdb
./network.sh createChannel
./network.sh deployCC
./network.sh applicationUp
