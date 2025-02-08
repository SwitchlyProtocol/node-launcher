# Restore Validator Backup

The `make backup` and `make restore` commands are used to backup and restore to the same namespace. If intending to restore a backup into a new node namespace, the following steps should be followed to migrate secrets and ensure the new node keyring is consistent.

1. Create backups of the original node.

```bash
NAME=node1 TYPE=validator NET=mainnet SERVICE=thornode make backup
NAME=node1 TYPE=validator NET=mainnet SERVICE=bifrost make backup
NAME=node1 TYPE=validator NET=mainnet make mnemonic > mnemonic.txt
NAME=node1 TYPE=validator NET=mainnet make password > password.txt
# check the mnemonic and password file contents for correctness
```

2. Destroy the original node.

```bash
NAME=node1 TYPE=validator NET=mainnet make destroy
```

3. Install the new node - enter the saved mnemonic and password when prompted.

```bash
NAME=node2 TYPE=validator NET=mainnet make install
```

4. Restore the backups.

```bash
NAME=node2 TYPE=validator NET=mainnet SERVICE=thornode make restore-backup
NAME=node2 TYPE=validator NET=mainnet SERVICE=bifrost make restore-backup
```
