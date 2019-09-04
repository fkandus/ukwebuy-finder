# cex-finder
Queries CeX (uk.webuy.com) to find out availability and prices of games.

How to run:
1. `go build cex.go response.go`
2. `cex short-game-ids.txt`

This will write to a file `trade-games-YYYYMMDD-HHMMSS` and to command line.

Output Example

```
-------------------------------------------------------
Red Dead Redemption 2 (2 Disc) (No DLC)
    Sell Price: 20.00

    Glasgow - Union Street: 4+
    Glasgow Sauchiehall: 4+
    Glasgow Forge: 4+
=======================================================
```
