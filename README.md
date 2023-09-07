# ApiServers

## kvApiServer

Server that provides a Rest Api for a kv datastore

### Commands

#### add

#### upd

#### del

#### get

#### list

#### info

#### entries

### Api Request Cli

http:[adr]/Get/db/cmd?key=keyval[&val=value]  

## kvLogin

Serve that tests login request. Info is provided in body and cli.  

http:[adr]/Post/db/cmd?user=username&pwd=password  
http:[adr]/Get/db/cmd?user=username&pwd=password  


