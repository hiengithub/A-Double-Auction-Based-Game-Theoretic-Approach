/*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/

'use strict';

module.exports.info  = 'startGame';

let bc, contx;
let t = "0"
// let quantity = "0"
// let startKey = "1.25"

function isEmpty(obj) {
  for(var prop in obj) {
    if(obj.hasOwnProperty(prop))
      return false;
  }

  return true;
}

module.exports.init = async function (blockchain, context, args) {

    bc = blockchain;
    contx = context;
   
    return Promise.resolve();
};

module.exports.run = async function() {
    var quantityJson = {};
    while (true) {
        let args = {
            chaincodeFunction: 'startGame',
            chaincodeArguments: ["", JSON.stringify(quantityJson),t]
        };
        console.log(' --- invoke SmartContract transaction --- ');
        console.log('request invokeSmartContract');
        let results = await bc.invokeSmartContract(contx, 'energyTrading2021', 'v0', args, 50);
        // PutState puts the specified `key` and `value` into the transaction's
        // writeset as a data-write proposal. PutState doesn't effect the ledger
        // until the transaction is validated and successfully committed.
        // Returns: Promise will be resolved when the peer has successfully handled the state update request or rejected if any errors.
        Promise.resolve();
        let winners = []
        let price = 0
        
        for (let result of results) {
            // do something with the data
            let shortID = result.GetID().substring(8);
            let executionTime = result.GetTimeFinal() - result.GetTimeCreate();
            let buffer = result.GetResult();
            console.log(`Response: TX [${shortID}] took ${executionTime}ms to execute. Result: {status: ${result.GetStatus()}, 
                value: ""}`);
            
            if (isEmpty(buffer)) {
                console.log("buffer is Empty");
                return results;
            }
            let parseBuffer = JSON.parse(buffer);

            if (parseBuffer[0].winners == null) {  
                return results;
            }

            if (winners.length == 0) {
                //convert string to array
                let temp = parseBuffer[0].winners.replace(/'/g, '"');
                winners = JSON.parse(temp);
            }
            
            price = parseBuffer[0].price;
        }

        console.log("call child_process")
        const spawn = require("child_process").spawnSync;
        var sumQuantity = 0.0

        quantityJson = {};

        for (var i = 0; i < winners.length; i++) {
            // console.log(winners[i])
            let userID = winners[i];
            var pythonProcess = spawn('python',[`/home/doanthanhhien/Documents/energy-trading-project/multiple-end-user/${userID}.py`, price,t]);
            sumQuantity += parseFloat(pythonProcess.stdout);
            let str = pythonProcess.stdout.toString();
            str = str.replace(/\r?\n|\r/g, "");
            // console.log(`optimal quantity ${userID}: ` + str);
            quantityJson[userID] = str;
        }

        console.log("sumQuantity",sumQuantity,'\n')    
    }
    return results;

};

module.exports.end = function() {
    return Promise.resolve();
};
