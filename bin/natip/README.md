# natip 

allocates a public ip on a specified server with ports

        Usage of natip:
          -internal string
                (optional) internal ip
          -password string
                clc password
          -ports string
                ports to open
          -server string
                server id
          -silent
                no output
          -username string
                clc username


## implementation

behavior: create or update public ip with ports. 

when `-internal` is passed, the ip associated with the NIC is used, 
otherwise the last NIC is chosen. 

if no public ip is mapped to the NIC, one is generated and associated. 

the requested ports are opened on the targeted NIC. 

## invocation

`./natip -username XYZ -pasword XYZ -server ALIAS-NODE-01 -ports "TCP/80 TCP/22" -silent`
