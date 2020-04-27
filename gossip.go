package main

import (
  "fmt"
  "time"
  "math/rand"
  "sync"
)

type Node struct{
  id int
  Hbcounter int
  time int
  dead bool
}
var wg sync.WaitGroup
var HB_mutex sync.Mutex
const num_nodes = 8
const num_neighbors = 2
const max_cycles = 20
const cycle_time = 2

func main() {
  member_ch := make(chan map[int]map[int]Node)
  for me := 0; me < num_nodes; me++ {
      my_HB_Table := make(map[int]Node)
      n := chooseNeighbors(me)
      for i := 0; i < num_neighbors; i++{
        my_HB_Table[n[i]] = Node{id: n[i], Hbcounter: 0, time: 0, dead: false}
      }
      spawnNode(Node{id: me, Hbcounter: 0, time: 0, dead: false}, my_HB_Table, member_ch)
  }
  wg.Wait()
}

func chooseNeighbors(me int) [num_neighbors]int {
  var n [num_neighbors]int
  for i := 0; i < num_neighbors; i++{
    var curr = rand.Intn(num_nodes)
    for curr == me || curr == n[0]{
      curr = rand.Intn(num_nodes)
    }
    n[i] = curr
  }
  fmt.Printf("Chose neighbors %v, for node: %d\n", n, me)
  return n
}

func spawnNode(my_node Node, my_HB_Table map[int]Node, member_ch chan map[int]map[int]Node){
  wg.Add(2)
  go updateHeartBeats(my_node, my_HB_Table, member_ch)
  go listenForTraffic(my_node, my_HB_Table, member_ch)
}

func listenForTraffic(my_node Node, my_HB_Table map[int]Node,
                      member_ch chan map[int]map[int]Node){
  defer wg.Done()
  for i := 0; i < max_cycles; i++{//listening on channel
    var mp = <-member_ch
    for k, v := range mp{//should only give us one iteration
      HB_mutex.Lock()
      _, found := my_HB_Table[k]
      HB_mutex.Unlock()
      if found {//if the information coming in is from a neighbor
        updateTable(k, my_node, v, my_HB_Table)
      }
    }
  }
}

func updateTable(sender_node_id int, my_node Node, new_values map[int]Node, my_HB_Table map[int]Node){
  for k, v := range new_values{//for all the information coming in
    value, found := my_HB_Table[k]
    if found && !value.dead{//if the stuff in the incoming table is in the neighborhood
      if v.time > value.time && v.Hbcounter <= value.Hbcounter{
        // fmt.Printf("Node %d, has killed Node %d\n", my_node.id, v.id)
        HB_mutex.Lock()
        my_HB_Table[k] = Node{id: v.id, Hbcounter: v.Hbcounter, time: v.time, dead: true}
        HB_mutex.Unlock()
        fmt.Printf("Node %d, has killed Node %d\n" + "-found %d in table from node %d\n-updating: %+v to: %+v\n"+ "-NEW Node %d TABLE: %+v\n\n", my_node.id, v.id, k, sender_node_id, value, my_HB_Table[k], my_node.id,my_HB_Table)
      }else if v.time > value.time {//if the information is more recent
        HB_mutex.Lock()
        my_HB_Table[k] = v
        HB_mutex.Unlock()
        fmt.Printf("For node : %d\n" + "-found %d in table from node %d\n-updating: %+v to: %+v\n"+ "-NEW Node %d TABLE: %+v\n\n", my_node.id, k, sender_node_id, value, v, my_node.id,my_HB_Table)
      }
    }
  }
}

// mark node as failed when failed
func updateHeartBeats(my_node Node, my_HB_Table map[int]Node,
                      member_ch chan map[int]map[int]Node){
  timer2 := time.NewTimer(time.Second*cycle_time)
  defer wg.Done()
  var sender_map = make(map[int]map[int]Node)
  my_HB_Table[my_node.id] = my_node
  fmt.Printf("Node: %+v, Initial table: %+v\n", my_node, my_HB_Table)
  sender_map[my_node.id] = my_HB_Table
  for i := 0; i < max_cycles; i++{
    <-timer2.C
    if !(i%2 == 0 && my_node.id%3 == 0){//generating failures for certain nodes
      my_node.Hbcounter += 1
    }
    my_node.time += 1
    HB_mutex.Lock()
    my_HB_Table[my_node.id] = my_node
    HB_mutex.Unlock()
    fmt.Printf("Update for Node %d => time: %d, HB: %d\n", my_node.id, my_node.time, my_node.Hbcounter)
    sender_map[my_node.id] = my_HB_Table
    member_ch <- sender_map
    timer2 = time.NewTimer(time.Second)
  }
}
