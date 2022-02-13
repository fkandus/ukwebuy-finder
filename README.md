# ukwebuy-finder
Queries [CeX](https://uk.webuy.com/) to find out availability and prices of games.

## How to run:
1. Execute `build.bat` (Windows) or `./build.sh` (Linux).
2. Rename `config.json.example` to `config.json`.
3. Set the values in `config.json` according to your own setup (Check the [Config file](#config-file) section).
4. Create an input file (Check the [Input file](#input-file) section). Name it, for example, `input-file.txt`.
5. Execute `finder input-file.txt`

This will write to a file `result-input-file-YYYYMMDD-HHMMSS.txt` and to command line.

## Input File

CSV of `GameID,Action` (See [game-list.txt.example](https://github.com/fkandus/ukwebuy-finder/blob/master/game-list.txt.example)).

To get the GameID go to the URL of a game and look for the `id` parameter.
For example, in `https://uk.webuy.com/product-detail?id=5026555423045&categoryName=playstation4-software&superCatName=gaming&title=red-dead-redemption-2-%282-disc%29-%28no-dlc%29` the Game ID is `5026555423045`.

Action is `buy` or `sell`. This is from the user standpoint ("user wants to buy/sell").

## Config File

- `urls`: API Urls to get different type of information. **Should not be changed.**
  - `detail`: API to get details of game.
  - `store`: API to get store availability for a game based on lat and lon.
  - `location`: API to get (lat, lon) from a City in the UK.
- `locations`: Configuration to find (lat, lon).
  - `city`: City names to send to the Location API.
- `stores`: Configuration for post-processing store data.
  - `matchName`: the store must match (contain) any of these strings to be taken into account.

## Cache

The command will try to store the coordinates of the config cities in a text file to use as cache. The user executing must have the permissions to write in the folder where the executable is.

## Output Example

In the following output `Red Dead Redemption 2`, `Infamous: First Light` and `Infamous 2` are examples of the `buy` action. `Uncharted: Golden Abyss` is an example of the `sell` action.

```
--------------------------------------------------------------------------------
Red Dead Redemption 2 (2 Disc) (No DLC) (Playstation4 Games) - (111719174486)
    Buy for: £20.00

    Glasgow - Union Street: 4+
    Glasgow Sauchiehall: 4+
    Glasgow Forge: 4+
--------------------------------------------------------------------------------
Infamous: First Light (Playstation4 Games) - (711719838814)
    Buy for: £8.00

    Not found in any store.
--------------------------------------------------------------------------------
Infamous 2 (Playstation3 Games) - (711719174486)
    Buy for: £2.00

    Glasgow - Union Street: 1
================================================================================
Uncharted: Golden Abyss (PS Vita Games) - (21171922sb02)
    Sell at: £30.00
================================================================================
Available items total Buy Value: £22.00
All items total Buy Value: £30.00
Total Sell Value: £30.00
Buy-Sell difference (available): £8.00
Buy-Sell difference (total): £0.00
================================================================================
Store counter:
Glasgow - Union Street: 2
Glasgow Sauchiehall: 1
Glasgow Forge: 1
```
