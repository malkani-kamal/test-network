
<!-- *** -->
# Hyperledger Fabric Network

## Prerequisites
[Fabric prerequisites](https://hyperledger-fabric.readthedocs.io/en/release-2.4/prereqs.html)

To make the below microservices up uncomment #km in network.sh

cd asset-microservice
make

cd blockchain-microservice
make



Remove all containers
docker rm -f $(docker ps -aq)

Two org setup
Step 1 (Only first time)
`./install-fabric.sh docker`
`./install-fabric.sh binary`

Step 2
#IMP
---Run startclean.sh to clean old network container and start a fresh.

sh startclean.sh

cd ~/asset-contract
./network.sh deployCC -ccv 2 -ccs 2
go mod tidy


If we change asset or blockchain microservice
./network.sh applicationUp

if network down 
sh start.sh

If chaincode changed
./network.sh deployCC -ccv <version> -ccs <sequence>

./network.sh deployCC -ccv 1 -ccs 1
go mod init erc20HTLC


## Steps to start the network
1. Run `./network.sh up -ca -s couchdb` and wait for the peers to boot successfully.
2. Run `./network.sh createChannel` to create the default channel `spydrachannel`.
3. Run `./network.sh deployCC` to install the default chaincode `../chaincode`

## Steps to install a new chaincode
* Run `./network.sh deployCC -ccv <version> -ccs <sequence>`.

    **NOTE:** 
    * You need to change the chaincode version everytime you install a new chaincode
    * You need to increment the chaincode sequence(integer) everytime you install a new chaincode
    * Check `log.txt` for latest chaincode information

## Steps to stop the network
1. Run `./network.sh down` to stop the network. 

<br>

***

## For more info refer to
* Run `./network.sh` for help
* [Fabric Documentation](https://hyperledger-fabric.readthedocs.io/en/release-2.4/index.html)
* [Faric Samples](https://github.com/hyperledger/fabric-samples) 