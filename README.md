####################################################################
###This document shows how to setup "Peer-to-Peer Energy Trading in Smart Grid Through Blockchain: A Double Auction-Based Game Theoretic Approach"###
Document Author: Hien Thanh Doan
Email: hiendoanht@sch.ac.kr
####################################################################

###Prerequisites###
- Hyperledger Fabric version: 1.4.7
- Hyperledger Caliper version: 0.3.2
- Ubuntu version: 16.04.04 LTS

###Hyperledger Caliper Setup###
- Clone caliper-benchmarks sample project from Hyperledger Caliper website
- Create energy-trading folder in caliper-benchmarks/benchmarks/scenario 
- Then put these file into this folder: config_energy_trading.yaml, initLedgerInGame_muser.js, startGame_muser.js. Note that the file initLedgerInGame_muser.js is used generate the optimize code of user named User.py. In this file you can change range of price, number of participate, etc... The file startGame_muser.js will be start after the init initLedgerInGame_muser.js sucessfully stared. Please make sure you update the folder path in these file where the source located.

###Smart contract Setup###
- Two existing kind of init smart contract
	- One is directly init smart contract through Hyperledger Fabric
	- Another way is using Caliper, which the smart contract is located at caliper-benchmarks\src\fabric\scenario\energyTrading2021_MultipleUser folder. Caliper will automatically init smart contract if it not exist.

###Start Matching###
In caliper-benchmarks folder execute the command: ./start-energy-trading.sh to start the matching process.


======================End===========================
