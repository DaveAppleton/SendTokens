Specifications
==============

0. Created and used at HelloGold for medium sized token transfers

1. Scans a list of <address><ether value> in CSV file and splits into <addresses> & contracts
2. Estimate Gas for contracts
3. Find total gas or if there are any weird addresses
4. Send the tokens recording the transaction ID
5. Use etherscan to check the transaction list for success or failure

NOTE - currently does not estimate the gas needed - please supply default

configuration : JSON FILE
-------------------------

This file was set up to transfer UET 

```
{
    "MAIN_ADDRESS" : "0x27f706edde3aD952EF647Dd67E24e38CD0803DD6",
    "TOKEN_ADDRESS" : "0x27f706edde3aD952EF647Dd67E24e38CD0803DD6",
    "GAS_PRICE" : "1000000000",
    "DEFAULT_GAS" : "145000",
    "SKIP_ROWS" : 1,
    "ADDRESS_COL" : 1,
    "AMOUNT_COL" : 2,
    "AMOUNT_IN_DECIMALS" : true
}
```

Token stuff

* **Main Address** - not currenly used
* **TOKEN_ADDRESS** - address of the token 
* **GAS_PRICE** - the price you wish to pay for the gas
* **DEFAULT_GAS** - currently you need to set the gas limit here 

CSV stuff

* **SKIP_ROWS** - number of header row
* **ADDRESS_COL** - column that holds the address (starting at zero)
* **AMOUNT_COL** - column that holds the amount of tokens to transfer
* **AMOUNT_IN_DECIMALS** - bool value - is the amount in tokens or minimum units [1]

[1] e.g. for 18 dp number

```
true -> 1.0
false -> 1000000000000000000 
```

see example config.json and test.csv