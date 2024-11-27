# Data Vault

## Introduction
Data vault is a distributed file storage system. It is designed to store small to large files, optimized for files between 1KB and 1GB.

## Features
- Distributed storage: files are stored in multiple nodes. nodes can be added or removed at any time.
- Grouping: files are stored in `groups`. A group is a set of files that are stored in the same nodes.
- Consistent hashing: the system uses consistent hashing with equal weights to distributes accross the nodes.
- In-memory Index: the system uses an in-memory index to keep track of the files and their location on each vault. The index is updated at vault level at every action and is reconstructed at start up.
- REST API: the data vault REST API is consistent between gate keeper and vaults.

## Cluster Architecture
The system is composed of a set of nodes (vaults) and a gateway (gate keeper).

```text
+-----------------+   +-----------------+   +-----------------+
| Vault 1         |   | Vault ...       |   | Vault N         |
+-----------------+   +-----------------+   +-----------------+
         |                    |                    |
         |                    |                    |
         +--------------------+--------------------+
                          |
                          |
                   +-------------+
                   | Gate Keeper |
                   +-------------+
```

## Single Vault Architecture
A single vault can be deployed as a standalone service. It is composed of a REST API and a storage service.




