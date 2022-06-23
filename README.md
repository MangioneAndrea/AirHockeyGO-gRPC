# Airhockey multiplayer with gRPC and GO

Simple standalone app of airhockey fully implemented in go and gRPC.

## Phisics

The prototype is provided with a very rudimentary collision system. It allows the disk to be pushed around, collide with the players and the walls

![Network replication](./collisions.gif)

## Network replication

The prototype is capable of multiplayer and replication. If the server is hosted online, multiple devices can join the same session.
The communication is implemented with grpc

![Network replication](./network.gif)