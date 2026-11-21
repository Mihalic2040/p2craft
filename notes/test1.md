So recently i did some testing with my GF. i run the p2craft server on my laptop behind the NAT on linux.
And she run the client on windows. It looked like the client ddl was hooked up properly to the java side. But we coud not establiss any conncetion. 

But i know that holepunched connection was possible because we have been talking using p2p calls in telegram. So.... What choud have went wrong. 


1. Advertised addreses are bull crap
2. We need more logging 


Component	Purpose
EnableHolePunching()	Enables DCUtR
EnableAutoNATv2()	Helps determine public reachability
EnableAutoRelay()	Discovers relays and obtains reservations
EnableRelay()	Enables relay client functionality
EnableNATPortMap()	Attempts UPnP/NAT-PMP port mapping
Kademlia DHT	Peer discovery and routing
