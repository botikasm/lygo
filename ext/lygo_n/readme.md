# N (distributed node architecture)

Application Framework to write multi-node applications.

In "N" there are two kind of nodes:
* Publishers
* Nodes

## Dependencies
Http uses [lygo_http_server](../lygo_http/lygo_http_server/readme.md) based on FastHTTP

Network messaging is using Nio [lygo_nio](../../base/lygo_nio/readme.md)


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

## What can I do with a Node?
"Why a Node?"

Basically because on a Node you can build a Distributed Application.

Nodes have some interesting features:
* Able to be published on the internet
* Can host a high performance fully customizable Web Server
* Can expose http APIs and Websocket APIs
* Execute "load-balanced" commands looking for free nodes on the internet
* Can be connected in virtual networks (group of nodes communicating each other)

## What cannot I do with a Node?
"What is not a Node?"

A Node is not a Microservice actor.

I didn't think at Nodes for distributed Microservices, but for scalable
monolithic applications.

Microservices are something else.
With Nodes you can deploy many containers as you need, hosting nodes.
Each Node can serve your client application, written in HTML5 or Go. 
