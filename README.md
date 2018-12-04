# FedChat

_End-to-End chat **federation** prototype_

Federation refers to different machines agreeing upon a set of standards to operate in a collective fashion. In the context of communication, chat applications can be federated for decentralization, thereby allowing a large number of users to communicate with each other while being connected to different servers, which themselves interact with each other for universal connectivity. Kind of like how email works with SMTP.

## Architecture

```
   		 ________          ____________		 ________________
		|        |        |            |        |		 |		
   		| Client | <----> | Homeserver | <----> | Central Server |
   		|________|        |____________|        |________________|		
								 ^
   							         |
								 |
								 v
							 ________________
							|		 |
							|   Homeserver   |
							|________________|
								 ^
								 |
								 |
								 v
							   ____________
							  |	       |
							  |   Client   |
							  |____________|

```

Whenever a homeserver is created, it emits its IP address to the central server. The job of the central server is to simply relay this information to all connected homeservers.
                             
## Setup

1. Clone the repository under `src` in your `$GOPATH`:

   ``` cd $GOPATH/src && git clone https://github.com/c0dzilla/FedChat.git ```

2. Generate the binary:

   ``` cd src/ && go install chat.go ```

3. To run as central server:

   ``` cd $GOPATH/bin && ./chat -mode=central ```

4. To run as homeserver:

   ``` cd $GOPATH/bin && ./chat -address=<IP address of central server> ```

5. If the IP of central server is not supplied to the homeserver, it works as a standalone chat. Hence, to run as a simple chat application:

   ``` cd $GOPATH/bin && ./chat ```

   When used as a homeserver/standalone chat, the user can access the chat client at port 8080 of the server's address.

## Contributing

Contributions are welcome. Feel free to open an issue or file a pull request.
