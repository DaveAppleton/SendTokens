Specifications
==============

0. Created and used at HelloGold for medium sized token transfers

1. Scans a list of &lt;address&gt;&lt;ether value&gt; in CSV file
2. Estimate Gas for contracts
3. Find total gas or if there are any weird addresses
4. Send the tokens recording the transaction ID
5. Use etherscan to check the transaction list for success or failure

NOTE - currently does not estimate the gas needed - please supply default

configuration : JSON FILE
-------------------------

This demo file was set up to transfer UET 

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

USAGE
-----

build it & run once

`gasAndSendTokens --input=test.csv`

```
Token :  Useless Ethereum Token
Payment Address  0xb85D7A869a61C60CDA6AF549Cc29A76476597ddf



Addresses-----------------------------------> 1
0xffC232afBdB712b54DF9291bcA252aB0D8Dc08f6 501320000000000000000 gas :  145000
Contracts-----------------------------------> 1
0x7da82C7AB4771ff031b66538D2fB9b0B047f6CF9 71910000000000000000 gas :  145000
SKIPPED ------------------------------------> 0

Tokens Required 573.230000000000000000, total gas 290000, total gas price 0.000290000000000000

Payment not enabled
```

Now transfer enough ether and tokens to the payment address, then run

`gasAndSendTokens --input=test.csv --pay=yes`

