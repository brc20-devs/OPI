# BRC-20 Swap Indexer

**BRC-20 Swap Indexer** is the a module of OPI. BRC-20 Swap Indexer saves all historical balance changes and all BRC-20 events.

**BRC-20 API** exposes activity on block (block events), balance of a wallet at the start of a given height, current balance of a wallet, block hash and cumulative hash at a given block and hash of all current balances.

The following diagram illustrates the architecture and data flow of the BRC-20 Swap Indexer
<img src="https://github.com/brc20-devs/brc20-swap-indexer/assets/3053743/e48e394f-579f-41ac-a06a-4a77a40b811f" width="512">


# Setup

For detailed installation guides:

- Ubuntu: [installation guide](INSTALL.ubuntu.md)
- Windows: [installation guide](INSTALL.windows.md)

OPI uses PostgreSQL as DB. Before running the indexer, setup a PostgreSQL DB (all modules can write into different databases as well as use a single database).

**Build ord:**

```bash
cd ord; cargo build --release;
```

**Install node modules**

```bash
cd modules/main_index; npm install;
cd ../brc20_api; npm install;
```

_Optional:_
Remove the following from `modules/main_index/node_modules/bitcoinjs-lib/src/payments/p2tr.js`

```js
if (pubkey && pubkey.length) {
  if (!(0, ecc_lib_1.getEccLib)().isXOnlyPoint(pubkey))
    throw new TypeError("Invalid pubkey for p2tr");
}
```

Otherwise, it cannot decode some addresses such as `512057cd4cfa03f27f7b18c2fe45fe2c2e0f7b5ccb034af4dec098977c28562be7a2`

**Install python libraries**

```bash
pip3 install python-dotenv;
pip3 install psycopg2-binary;
python3 -m pip install json5 stdiomask;
```

**Setup .env files and DBs**

Run `reset_init.py` in each module folder (preferrably start from main_index) to initialise .env file, databases and set other necessary files.

# (Optional) Restore from an online backup for faster initial sync

1. Install dependencies: (pbzip2 is optional but greatly impoves decompress speed)

```bash
sudo apt update
sudo apt install postgresql-client-common
sudo apt install postgresql-client-14
sudo apt install pbzip2

python3 -m pip install boto3
python3 -m pip install tqdm
```

2. Run `restore.py`

```bash
cd modules/;
python3 restore.py;
```

# Run

**Main Meta-Protocol Indexer**

```bash
cd modules/main_index;
node index.js;
```

**BRC-20 Indexer**

```bash
cd modules/brc20_swap_index;
docker-compose up -d
```

**BRC-20 API**

```bash
cd modules/brc20_api;
node api.js;
```

```

# Update

- Stop all indexers and apis (preferably starting from main indexer but actually the order shouldn't matter)
- Update the repo (`git pull`)
- Recompile ord (`cd ord; cargo build --release;`)
- Re-run all indexers and apis
- If rebuild is needed, you can run `restore.py` for faster initial sync
```
