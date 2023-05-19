## Setup "Peer-to-Peer Energy Trading in Smart Grid Through Blockchain: A Double Auction-Based Game Theoretic Approach"

Document Author: Hien Thanh Doan,
Email: hiendoanht@sch.ac.kr

## Configuration
- Hyperledger Fabric version: 1.4.7
- Hyperledger Caliper version: 0.3.2
- Ubuntu version: 16.04.04 LTS

## Hyperledger Caliper Setup
- Clone caliper-benchmarks sample project from Hyperledger Caliper website 
- Create energy-trading folder in caliper-benchmarks/benchmarks/scenario 
- Then put these file into this folder: ``config_energy_trading.yaml``, ``initLedgerInGame_muser.js``, ``startGame_muser.js``. 

**Note** The file ``initLedgerInGame_muser.js`` will automatically generate the optimized code of the user named User.py. In this file, you can also change the range of price, number of participants, etc... To make the simulation work, please update lines 66 and 104 to yours in ``initLedgerInGame_muser.js`` and ``startGame_muser.js``.

## Smart contract Setup
- Two existing methods to deploy smart contracts in the blockchain network:
	- The first method is to initiate a smart contract using Hyperledger Fabric.
	- Another option is to use Caliper, where the smart contract can be found in the caliper-benchmarks\src\fabric\scenario\energyTrading2021_MultipleUser folder. Caliper automatically creates a smart contract if it does not exist. **This method is suggested as a convenience**.

## Start Matching
- First, please update line 13 to yours in ``start-energy-trading.sh`` file. Note that to get an idea of what the network configuration file (caliper_config.yaml) looks like, please look at https://github.com/hyperledger/caliper-benchmarks/tree/v0.3.2/networks/fabric/v1/v1.4.1.
- Then, In caliper-benchmarks folder execute the command: ``./start-energy-trading.sh`` to start the matching process.
