# VWAP

![technology Go](https://img.shields.io/badge/technology-go-blue.svg)

## Overview

The goal of this project is to create a real-time `VWAP`  (volume-weighted average price) calculation engine. For this was used the coinbase websocket feed to stream in trade executions and update the VWAP for each trading pair as updates become available. 

## Default parameters: 

- Trading pairs:  `BTC-USD`, `ETH-USD`, `ETH-BTC` .
- Sliding window: 200 . 
- Coinbase URL: `wss://ws-feed.exchange.coinbase.com` .

## Program structure

The service is composed of two main components:

- A `websocket` client that pulls data off on trade executions. By default is used coinbase websocket feed. 
- A `vwapcalculator` to calculate the VWAP. Use a list to push the data points and a map to save the  calculated trading pairs.

## Performance

- Having a list of datapoints with a hash map to store the cumulative values for each trading pair favors efficiency, since this avoids having to loop over all the data points to calculate the vwap. So this strategy reduces the complexity from O(N) to O(1).

- Precision calculation:
  `Decimal library` is used because it's more exact for money, this means, it is less performant but in favors to the correctness. Also, the use of this library is more intuitive than using big.Int, mainly for mathematical operations where errors are sought to be avoided.

## How To Run This Project

- Clone this repo to your workspace. 
- Type `make run` .

## Tests

 - Unit Test: `make test-unit`

 - Integration Test: `make test-intergration`


## Improvements

It is necessary to improve the test cases to cover all border cases. 