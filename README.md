# ukwebuy-finder
Queries CeX (uk.webuy.com) to find out availability and prices of games.

## How to run:
1. Execute `go build cex.go response.go configuration.go`.
2. Execute `cex short-game-ids.txt`

This will write to a file `trade-games-YYYYMMDD-HHMMSS.txt` and to command line.

## Input File

CSV of Game ID,"buy" or "sell" action. See `short-game-ids.txt`. 

## Output Example

```
-------------------------------------------------------
Red Dead Redemption 2 (2 Disc) (No DLC)
    Sell Price: 20.00

    Glasgow - Union Street: 4+
    Glasgow Sauchiehall: 4+
    Glasgow Forge: 4+
=======================================================
Uncharted: Golden Abyss
    Exchange Price: 10.00
=======================================================
```

