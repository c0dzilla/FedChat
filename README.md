# FedChat

_End-to-End chat **federation** prototype_

Federation refers to different machines agreeing upon a set of standards to operate in a collective fashion. In the context of communication, chat applications can be federated for decentralization, thereby allowing a large number of users to communicate with each other while being connected to different servers, which themselves interact with each other for universal connectivity. Kind of like how email works with SMTP.

## Architecture


### Registration of newly-joined homeserver

```
	 _____________		   ________________
	|	      |    ip	  |		   | emit to all connected home servers
	|  Homeserver | --------> | Central Server |----------------------------------->
	|_____________|		  |________________|

```

### Messaging

```
   		 ________          ____________		 ____________	       ________
		|        |        |            |        |	     |	      |	       |
   		| Client | <----> | Homeserver | <----> | Homeserver | <----> | Client |
   		|________|        |____________|        |____________|	      |________|
					
```

Whenever a homeserver is created, it emits its IP address to the central server. The job of the central server is to simply relay this information to all connected homeservers.
                             
## Setup

1. Install the package using `go get`.

```  go get github.com/c0dzilla/FedChat/src ```

2. To run a homeserver and connect to the broader federation network:

   ``` cd $GOPATH/bin
./src -address=<IP address of central server ```

5. If the IP of central server is not supplied to the homeserver, it works as a standalone chat. Hence, to run as a simple chat application:

   ``` cd $GOPATH/bin
./chat ```

6.  Homeservers run at port 8080 by default.

## Contributing

Contributions are welcome. Feel free to open an issue or file a pull request.
