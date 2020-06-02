# N (distributed node architecture)

Application Framework to write multi-node applications.

In "N" there are two kind of nodes:
* Publishers
* Nodes

## Publishers
Publishers are node that can register new nodes on the network.

Each node with "publishing" enabled looks for a Publisher to send its public and reachable endpoint.

Publishers maintain an internal registry of active nodes. 
Active Nodes are pinged each 3 seconds and are removed 
from list if do not respond.

A Publisher synchronize its internal registry with other Publishers in the network

## Nodes
Nodes can be "published" or not.
Published nodes can respond to requests of other nodes.

Unpublished nodes remain isolated from other nodes, but can use remote endpoints when required from load balancing.
  
## Networks
Nodes are grouped inside a network.
Each network has a defined ID (network_id).

Only nodes inside same network can be reached from load balancing service.

## Dependencies
Server uses [lygo_http_server](../lygo_http/lygo_http_server/readme.md)

