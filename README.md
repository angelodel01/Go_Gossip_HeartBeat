# Go_Gossip_HeartBeat

Angelo De Laurentis For Maria Pantoja
Week 3 Assignment Gossip Heartbeat Protocol

To run this project simply type 'go run gossip.go'

You can track the progress of the gossip nodes speaking to each
other by watching the console print outs.

##REQUIREMENTS##

Have 8 computing nodes each with two neighbors
(can select them at random or keep a fix to
whom they need to exchange heartbeat tables with.
Heartbeat tables contain : id neighbor, hbcounter, time

Every node:
Keeps a Hbcounter gets increase every X amount of time (you choose)
Send their HB tables to its neighbor every Y amount of time (you choose)
Simulate one node failing every Z amount of time and how the tables change

##UNDERSTANDING LOGS##
When you first run the program you will see each node print out
the id's of their chosen neighbors, example:

  Chose neighbors [1 7], for node: 0
  Chose neighbors [7 3], for node: 1
  Chose neighbors [1 6], for node: 2
  Chose neighbors [1 4], for node: 3
  etc...

And then Each node will print out their initial values along with their
Initial Heart Beat tables (note each node includes themselves in their tables), example:

  Node: {id:3 Hbcounter:0 time:0 dead:false}, Initial table: map[1:{id:1 Hbcounter:0 time:0 dead:false} 3:{id:3 Hbcounter:0 time:0 dead:false} 4:{id:4 Hbcounter:0 time:0 dead:false}]
  Node: {id:2 Hbcounter:0 time:0 dead:false}, Initial table: map[1:{id:1 Hbcounter:0 time:0 dead:false} 2:{id:2 Hbcounter:0 time:0 dead:false} 6:{id:6 Hbcounter:0 time:0 dead:false}]

Once all the nodes are initialized they will begin updating the console
you will see a node update their timer and heartbeat by printing out their
new HB(Heartbeat) and time on a log that starts with the word "Update"
example:

  Update for Node 5 => time: 1, HB: 1

You will also see logged when a node updates itself based upon information
received from a neighbor it will denote where it got the new information
and how it has updated it's own table, example:

  For node : 0
  -found 7 in table from node 1
  -updating: {id:7 Hbcounter:5 time:5 dead:false} to: {id:7 Hbcounter:11 time:11 dead:false}
  -NEW Node 0 TABLE: map[0:{id:0 Hbcounter:6 time:13 dead:false} 1:{id:1 Hbcounter:8 time:8 dead:false} 7:{id:7 Hbcounter:11 time:11 dead:false}]

When a node has discovered the death of another node in the network
it will show a special log displaying similar information to an update, example:

  Node 2, has killed Node 6
  -found 6 in table from node 6
  -updating: {id:6 Hbcounter:2 time:4 dead:false} to: {id:6 Hbcounter:2 time:5 dead:true}
  -NEW Node 2 TABLE: map[1:{id:1 Hbcounter:3 time:3 dead:false} 2:{id:2 Hbcounter:2 time:5 dead:false} 6:{id:6 Hbcounter:2 time:5 dead:true}]



##MY IMPLEMENTATION##

The amount of nodes is determined by a constant defined on line 18
named "num_nodes"

The amount of nodes that each node will include in their
neighborhood(random group they will communicate with)
is defined as a constant on line 19 named "num_neighbors"

The amount of cycles the algorithm will run through is defined by
a constant on line 20 named "max_cycles"

The frequency that the nodes update their HBCounters and
Heart Beat Tables are determined by a constant on line 21
named "cycle_time"

To simulate Node Death, for Nodes that are multiples of 3(including 0)
on every even cycle they skip a Heartbeat increment.

Because of the randomness introduced with GoRoutines you may not
detect node death on every run.
To remedy this simply increase the amount of cycles or increase
the amount of Heartbeats skipped.

##NODES##
Each node has 2 GoRoutines that are simultaneously running
a loop that updates their HBCounters and sending out their
Heart Beat Tables
