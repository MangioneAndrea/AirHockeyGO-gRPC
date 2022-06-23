# Airhockey multiplayer with gRPC and GO

Simple standalone app of airhockey fully implemented in go and gRPC using ebiten

## Phisics

<div style="display: flex;">
<p style="width: 50%">
The prototype is provided with a very rudimentary collision system. It allows the disk to be pushed around, collide with the players and the walls
</p>
<div style="width: 50%; display: flex;justify-content:space-evenly;">
<img src="./collisions.gif" alt="collisions" height="400"/>
</div>
</div>

## Network replication

<div style="display: flex;">
<p style="width: 50%">
The prototype is capable of multiplayer and replication. If the server is hosted online, multiple devices can join the
same session. The communication is implemented with grpc
</p>
<div style="width: 50%; display: flex;justify-content:space-evenly;">
<img src="./network.gif" alt="network"  height="400"/>
</div>
</div>
