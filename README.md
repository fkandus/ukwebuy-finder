# ukwebuy-finder
Queries [CeX](https://uk.webuy.com/) to find out availability and prices of games.

## How to run:
1. Execute `go build finder.go ukwebuy.go location.go configuration.go`.
2. Rename `config.json.example` to `config.json`.
3. Set the values in `config.json` according to your own setup (Check the [Config file](#config-file) section).
4. Create an input file (Check the [Input file](#input-file) section). Name it, for example, `input-file.txt`.
5. Execute `finder input-file.txt`

This will write to a file `trade-games-YYYYMMDD-HHMMSS.txt` and to command line.

## Input File

CSV of `GameID,Action` (See [short-game-ids.txt](https://github.com/fkandus/ukwebuy-finder/blob/master/short-game-ids.txt)).

To get the GameID go to the URL of a game and look for the `id` parameter.
For example, in `https://uk.webuy.com/product-detail?id=5026555423045&categoryName=playstation4-software&superCatName=gaming&title=red-dead-redemption-2-%282-disc%29-%28no-dlc%29` the Game ID is `5026555423045`.

Action is `buy` or `sell`. This is from the user standpoint ("user wants to buy/sell").

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

The following output shows `Red Dead Redemption 2` as an example of the `buy` action and `Uncharted: Golden Abyss` as an example of the `sell` action.

```
-------------------------------------------------------
Red Dead Redemption 2 (2 Disc) (No DLC)
    Buy for: £20.00

    Glasgow - Union Street: 4+
    Glasgow Sauchiehall: 4+
    Glasgow Forge: 4+
=======================================================
Uncharted: Golden Abyss
    Sell at: £10.00
=======================================================
```

