# ukwebuy-finder
Queries [CeX](https://uk.webuy.com/) to find out availability and prices of games.

## How to run:
1. Execute `go build finder.go ukwebuy.go location.go configuration.go`.
2. Execute `finder short-game-ids.txt`

This will write to a file `trade-games-YYYYMMDD-HHMMSS.txt` and to command line.

## Input File

CSV of Game ID,"buy" or "sell" action. See `short-game-ids.txt`.

## Config File

- `urls`: API Urls to get different type of information. **Should not be changed.**
  - `detail`: API to get details of game.
  - `store`: API to get store availability for a game based on lat and lon.
  - `location`: API to get (lat, lon) from a City in the UK.
- `locations`: Configuration to find (lat, lon).
  - `city`: City name to send to the Location API.
- `stores`: Configuration for post-processing store data.
  - `matchName`: the store must match (contain) this string to be taken into account.

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

