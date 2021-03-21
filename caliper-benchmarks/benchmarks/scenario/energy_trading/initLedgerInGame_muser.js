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

module.exports.info  = 'initLedgerInGame';

let bc, contx;

module.exports.init = function(blockchain, context, args) {

    bc = blockchain;
    contx = context;

    return Promise.resolve();
};

function padLeadingZeros(num, size) {
    var s = num+"";
    while (s.length < size) s = "0" + s;
    return s;
}

module.exports.run = function() {

	let numberOfSeller = 1000;
	let numberOfBuyer = 1000;
	var minQuantity = 0.1, maxQuantity = 5;
	var minPriceB = 14, maxPriceB = 22;
	var userInfoArr = [];

	// include node fs module
	var fs = require('fs');
	let i = 0;
	while (i < numberOfBuyer) {
		i += 1
		let rdQuantity = (Math.random() * (maxQuantity - minQuantity) + minQuantity).toFixed(3)
		let rdPrice = (Math.random() * (maxPriceB - minPriceB) + minPriceB).toFixed(1)
		let idNumber = padLeadingZeros(i, 4);
		let content = 'import sys \n';
		content += 'lamda = [4] * 13\n';
		content += `xmaxUser${idNumber} = [${rdQuantity},0,0,0,0,0,0,0,0,0,0,0,0]\n`;
		content += `surplusEnergyUser${idNumber} = [0,0,0,0,0,0,0,0,0,0,0,0,0]\n`;
		content += `lamdaUser${idNumber} = lamda\n`;
		content += `pricesUser${idNumber} = [${rdPrice},0,0,0,0,0,0,0,0,0,0,0,0]\n`;
		content += 'p = float(sys.argv[1])\n';
		content += 't = int(sys.argv[2])\n';
		content += `x = (pricesUser${idNumber}[t] - p)/lamdaUser${idNumber}[t]\n`;
		content += 'if x < 0:\n';
		content += '	x = 0\n';
		content += `elif x > xmaxUser${idNumber}:\n`;
		content += `	x = xmaxUser${idNumber}\n`;
		content += 'print(str(x))\n';
		content += 'sys.stdout.flush()\n';
		fs.writeFile(`//home//doanthanhhien//Documents//energy-trading-project//multiple-end-user//user${idNumber}.py`, content, function (err) {
		  if (err) throw err;
		  // console.log('File is created successfully.');
		});

		//generate user info
		var individualUser = {};
		individualUser["xMax"] = [rdQuantity,"0","0","0","0","0","0","0","0","0","0","0","0"];
		individualUser["surplusEnergy"] = ["0","0","0","0","0","0","0","0","0","0","0","0","0"];
		individualUser["prices"] = [rdPrice,"0","0","0","0","0","0","0","0","0","0","0","0"];
		individualUser["identify"] = `user${idNumber}`;
		individualUser["clearPrice"] = ["0","0","0","0","0","0","0","0","0","0","0","0","0"];
		individualUser["clearQuantity"] = ["0","0","0","0","0","0","0","0","0","0","0","0","0"];
		userInfoArr.push(individualUser);
	}

	var minPriceS = 5, maxPriceS = 13;
	let j = i;
	while (j < numberOfSeller + numberOfBuyer) {
		j += 1
		let rdQuantity = (Math.random() * (maxQuantity - minQuantity) + minQuantity).toFixed(3)
		let rdPrice = (Math.random() * (maxPriceS - minPriceS) + minPriceS).toFixed(1)
		let idNumber = padLeadingZeros(j, 4);
		let content = 'import sys \n';
		content += 'lamda = [4] * 13\n';
		content += `xmaxUser${idNumber} = [0,0,0,0,0,0,0,0,0,0,0,0,0]\n`;
		content += `surplusEnergyUser${idNumber} = [${rdQuantity},0,0,0,0,0,0,0,0,0,0,0,0]\n`;
		content += `lamdaUser${idNumber} = lamda\n`;
		content += `pricesUser${idNumber} = [${rdPrice},0,0,0,0,0,0,0,0,0,0,0,0]\n`;
		content += 'p = float(sys.argv[1])\n';
		content += 't = int(sys.argv[2])\n';
		content += `x = (pricesUser${idNumber}[t] - p)/lamdaUser${idNumber}[t]\n`;
		content += 'if x < 0:\n';
		content += '	x = 0\n';
		content += `elif x > xmaxUser${idNumber}:\n`;
		content += `	x = xmaxUser${idNumber}\n`;
		content += 'print(str(x))\n';
		content += 'sys.stdout.flush()\n';
		fs.writeFile(`//home//doanthanhhien//Documents//energy-trading-project//multiple-end-user//user${idNumber}.py`, content, function (err) {
		  if (err) throw err;
		  // console.log('File is created successfully.');
		}); 

		//generate user info
		var individualUser = {};
		individualUser["xMax"] = ["0","0","0","0","0","0","0","0","0","0","0","0","0"];
		individualUser["surplusEnergy"] = [rdQuantity,"0","0","0","0","0","0","0","0","0","0","0","0"];
		individualUser["prices"] = [rdPrice,"0","0","0","0","0","0","0","0","0","0","0","0"];
		individualUser["identify"] = `user${idNumber}`;
		individualUser["clearPrice"] = ["0","0","0","0","0","0","0","0","0","0","0","0","0"];
		individualUser["clearQuantity"] = ["0","0","0","0","0","0","0","0","0","0","0","0","0"];
		userInfoArr.push(individualUser);
	}

	// console.log(userInfoArr)
	let endKey = `user${padLeadingZeros(numberOfSeller + numberOfBuyer + 1, 4)}`;
    let args = {
        chaincodeFunction: 'initLedger',
        chaincodeArguments: [JSON.stringify(userInfoArr),endKey]
    };

    return bc.invokeSmartContract(contx, 'energyTrading2021', 'v0', args, 50);
};

module.exports.end = function() {
    return Promise.resolve();
};
