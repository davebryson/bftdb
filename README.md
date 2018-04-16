## BFT-SQLITE

Tendermint + Sqlite3 = BFT Database Replication.

Inspired by rqlite and some others, this is an example of using a blockchain to
replicate a SQL database. All SQL statements are recorded on the blockchain and inserted
into the local sqlite3 database.  Every node has a local copy of the db replicated
by the transactions on the blockchain.  If a node goes offline, it will automatically
re-sync with the network.

This example uses an in-memory sql db.  So it's very difficult for anyone to locally
change/alter the db -  and `DROP` statements are rejected. You can only submit SQL statements
through the blockchain, you can not directly interact with the database.  However, our simple
REST service allows arbitrary queries against the db.

The command line includes 2 options:
- Start a node (with the embedded database)
- And and interactive console to send SQL statements to the blockchain (and db)

1. In one terminal, fire up the blockchain `bftdb start`
2. In another terminal, fire up the interactive console `bftdb console`

For demo purposes, the blockchain creates a table called `sample` with a single field `name` which is a string (TEXT).

Example use of the console:
```
> insert into sample(name) values('dave')
response Status : 200 OK
response Body   : {"check_tx":{"fee":{}},"deliver_tx":{"fee":{}},"hash":"01B60399F645DD59C5CA257C9346D5E96502B1AF","height":39}

> select * from sample
response Status : 200 OK
response Body   : {"columns":["id","name"],"values":[[1,"dave"]]}
```
