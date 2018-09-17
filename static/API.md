# API

Host: `http://localhost:{controller_port}`

Can change to  `http://c.sipt.top`  if Shuttle was set as a system proxy.

- [General Shutdown](#general-shutdown)
  - [Reload Configuration Mode](#reload-configuration-mode)
  - [Get Mode](#get-mode)
  - [Change Mode](#change-mode)
- [Server](#server)
  - [Server List](#server-list)
  - [Select Server](#select-server)
  - [Refresh RTT](#refresh-rtt)
- [Records](#records)
  - [Records List](#records-list)
  - [Clear Records](#clear-records)
  - [Show Records in Websocket](#show-records-in-websocket)
- [DNS](#dns)
  - [DNS Cache](#dns-cache)
  - [Clear DNS Cache](#clear-dns-cache)
- [Dump](#dump)
  - [Dump Status](#dump-status)
  - [Dump Switch](#dump-switch)
  - [Dump Data](#dump-data)
  - [Dump Body](#dump-body)
- [Certificate](#certificate)
  - [Generate Certificate](#generate-certificate)
  - [Download Certificate](#download-certificate)
- [MitM Rules](#mitm-rules)
  - [Get MitM Rules](#get-mitm-rules)
  - [Add MitM Rules](#add-mitm-rules)
  - [Remove MitM Rules](#remove-mitm-rules)

## General

#### Shutdown

Exit the program.

```
POST /api/shutdown
```

Response Body:

```js
{
  "code": 0, // 0: success, 1: failed
  "message": "" // code==1, error message
}
```

#### Reload Configuration

Reload configuration if the configuration file was edited.

```
POST /api/reload
```

Response Body:

```js
{
  "code": 0, // 0: success, 1: failed
  "message": "" // code==1, error message
}
```

## Mode

#### Get Mode

Get current Shuttle's mode. (Enum: `RULE`, `REMOTE`,` DIRECT`,`REJECT` )

URL: 

```
/api/mode
```

Method:

```
GET
```

Response Body:

```js
{
  "code": 0, // 0: success, 1: failed
  "message": "", // code==1, error message
  "data": "RULE" // enum: RULE, REMOTE, DIRECT, REJECT 
}
```

#### Change Mode

Change Shuttle's mode. (Enum: `RULE`, `REMOTE`,` DIRECT`,`REJECT` )

```
POST /api/mode/:mode
eg. :mode = RULE or rule  # Change mode to "RULE" mode
```

Response Body:

```js
{
  "code": 0, // 0: success, 1: failed
  "message": "", // code==1, error message
  "data": "RULE" // enum: RULE, REMOTE, DIRECT, REJECT 
}
```

## Server

Manage servers and groups in the Shuttle.

#### Server List

Return all servers and groups in Shuttle. (Value of [Porxy] and [Proxy-Group] in configuration file).

`GLOBAL` is defined by Shuttle, that was used by Shuttle in `REMOTE` mode.

```
GET /api/servers
```

Response Body:

```js
{
  "code": 0, // 0: success, 1: failed
  "message": "", // code==1, error message
  "data": [
    {
      "name": "Auto", // group name
      "servers": [ // server or group array
        {
          "name": "Server1", // server name
          "selected": false, // be selected
          "rtt": "936ms" // Round Trip Time
        },
        {
          "name": "Server2",
          "selected": false,
          "rtt": "100ms"
        }
      ],
      "select_type": "rtt" // group type: Round Trip Time. enum: rtt, select
    },
    {
      "name": "Proxy", // group name
      "servers": [// server or group array
        {
          "name": "Auto", // group name
          "selected": true // be selected
        }
      ],
      "select_type": "select" // group type: manual to select. enum: rtt, select
    }
  ]
}
```



#### Select Server

Select a server or group in the group. (This group's type must be "select").

```
POST /api/server/select
```

Request Header:

```
Content-Type: application/x-www-form-urlencoded; charset=utf-8
```

Request Body:

| Key    | Value Type | Desc        |
| ------ | ---------- | ----------- |
| group  | string     | Group Name  |
| server | string     | Server Name |

Response Body:

```js
{
  "code": 0, // 0: success, 1: failed
  "message": "" // code==1, error message
}
```



#### Refresh RTT

Refresh server's rtt in the group. (This group's type must be "rtt").

URL: 

```
POST /api/server/select/refresh
```

Request Header:

```
Content-Type: application/x-www-form-urlencoded; charset=utf-8
```

Request Parameters:

| Key   | Value Type | Desc       |
| ----- | ---------- | ---------- |
| group | string     | Group Name |

Response Body:

```js
{
  "code": 0, // 0: success, 1: failed
  "message": "" // code==1, error message
}
```



## Records

#### Records List

Return request records. (Max count = 500)

```
GET /api/records
```

Response Body:

```js
{
  "code": 0, // 0: success, 1: failed
  "message": "", // code==1, error message
  "data": [
    {
      "ID": 72,
      "Protocol": "HTTPS",
      "Created": "2018-09-14T01:23:52.194709+08:00",
      "Proxy": {
        "Name": "ProxyServerName",
        "Rtt": 44977632, 
        "ProxyProtocol": "socks",
      },
      "Rule": {
        "Type": "DOMAIN-SUFFIX",
        "Value": "example.com",
        "Policy": "Proxy", // group or server or DIRECT or REJECT
        "Options": null,
        "Comment": ""
      },
      "Status": "Completed", // Status: Active, Completed, Reject
      "Up": 10, // upload 10B
      "Down": 1503, // Download 1503B
      "URL": "http://example.com/view",
      "Dumped": false  // data was dumped
    }
  ]
}
```



### Clear Records

Clear records.

```
DELETE /api/records
```

Response Body:

```js
{
  "code": 0, // 0: success, 1: failed
  "message": "", // code==1, error message
}
```



### Show Records in Websocket

Server push data to client in websocket.

URL:

```
ws://{host}:{port}/api/ws/records
```

Push data:

```js
[{
    "ID":0,
    "Op":4, // Op=4, Append Record
    "Value":{
      "ID": 72,
      "Protocol": "HTTPS",
      "Created": "2018-09-14T01:23:52.194709+08:00",
      "Proxy": {
        "Name": "ProxyServerName",
        "Rtt": 44977632, 
        "ProxyProtocol": "socks",
      },
      "Rule": {
        "Type": "DOMAIN-SUFFIX",
        "Value": "example.com",
        "Policy": "Proxy", // group or server or DIRECT or REJECT
        "Options": null,
        "Comment": ""
      },
      "Status": "Completed", // Status: Active, Completed, Reject
      "Up": 10, // upload 10B
      "Down": 1503, // Download 1503B
      "URL": "http://example.com/view",
      "Dumped": false  // data was dumped
    }
  },{
    "ID":4,
    "Op":2, // Op=2, Upload += {Value}Byte
    "Value":63
  },{
    "ID":4,
    "Op":3, // Op=3, Download += {Value}Byte
    "Value":63
  },{
    "ID":4,
    "Op":1, // Op=1, change status
    "Value":"Completed"
  },{
    "ID":0,
    "Op":5, // Op=5, remove the record where ID == {Value} from storage
    "Value":4
  }]
```

| Op   | Desc                                            |
| ---- | ----------------------------------------------- |
| 1    | Change status. Emun: Active, Completed, Reject. |
| 2    | Upload {Value} byte.                            |
| 3    | Download {Value} byte.                          |
| 4    | Append a record.                                |
| 5    | Remove the record where ID == {Value}.          |



## DNS

#### DNS Cache

Return DNS cache.

```
GET /api/dns
```

Response Body:

```js
{
  "code": 0, // 0: success, 1: failed
  "message": "", // code==1, error message
  "data": [
    {
      "MatchType": "DOMAIN", // enum: DOMAIN, 
      "Domain": "localhost",
      "IPs": [
        "127.0.0.1"
      ],
      "DNSs": null,
      "Type": "static",
      "Country": ""
    }
  ]
}
```



#### Clear DNS Cache

```
DELETE /api/dns
```

Response Body:

```js
{
  "code": 0, // 0: success, 1: failed
  "message": "", // code==1, error message
}
```



## Dump

#### Dump Status

```
GET /api/dump/allow
```

Response Body:

```js
{
  "code": 0, // 0: success, 1: failed
  "message": "", // code==1, error message
  "data": {
    "allow_dump": true,
    "allow_mitm": false
  }
}
```



#### Dump Switch

```
POST /api/dump/allow
```

Request Header:

```
Content-Type: application/x-www-form-urlencoded; charset=utf-8
```

Request Parameters:

| Key        | Value Type | Desc        |
| ---------- | ---------- | ----------- |
| allow_dump | `boolean`  | Enable Dump |
| allow_mitm | `boolean`  | Enable MitM |

Response Body:

```js
{
  "code": 0, // 0: success, 1: failed
  "message": "", // code==1, error message
  "data": {
    "allow_dump": true,
    "allow_mitm": false
  }
}
```



#### Dump Data

```
GET /api/dump/data/:conn_id
```

`conn_id`: record's id

Response Body:

```js
{
  "code": 0,
  "message": "success",
  "data": {
    "ReqHeader": "base64.StdEncoding.Encode({Request Headers})",
    // if len({Request Body} > 2MB) { return base64.StdEncoding.Encode("large body") }
    "ReqBody": "base64.StdEncoding.Encode({Request Body})", 
    // if len({Response Body} > 2MB) { return base64.StdEncoding.Encode("large body") }
    "RespBody": "base64.StdEncoding.Encode({Request Headers})",
    "RespHeader": "base64.StdEncoding.Encode({Request Body})"
  }
}
```



#### Dump Body

```
GET /api/dump/large/:conn_id
```

`conn_id`: record's id

Request Query:

| Key       | Type   | Desc                                                       |
| --------- | ------ | ---------------------------------------------------------- |
| file_name | string | File name for download                                     |
| dump_type | string | Enum: `request`(Request Body) or `response`(Response Body) |
| conn_id   | int    | Record's ID                                                |

Response Header:

```
Content-Type: application/octet-stream
content-disposition: attachment; filename={file_name}
```

Response Body:

```
{Byte Stream}
```



## Certificate

Prepare for MitM.

#### Generate Certificate

```
POST /api/cert
```

Response Body:

```js
{
  "code": 0, // 0: success, 1: failed
  "message": "", // code==1, error message
}
```



#### Download Certificate

```
GET /api/cert
```

Response Header:

```
Content-Type: application/octet-stream
content-disposition: attachment; filename=Shuttle.cer
```

Response Body:

```
{Byte Stream}
```



## MitM Rules

Shuttle will only decrypt traffic to hosts which are declared in "mitm.rules".

#### Get MitM Rules

```
GET /mitm/rules
```

Response Body:

```js
{
  "code": 0, // 0: success, 1: failed
  "message": "", // code==1, error message
  "data": [
    "*.google.com" // host list
  ]
}
```



#### Add MitM Rules

```
POST /mitm/rules
```

Request Query:

| Key    | Type   | Desc                                               |
| ------ | ------ | -------------------------------------------------- |
| domain | string | Prefix wildcard `*` is allowed or just single `*`. |

Response Body:

```js
{
  "code": 0, // 0: success, 1: failed
  "message": "", // code==1, error message
  "data": [
    "*.google.com" // host list
  ]
}
```



#### Remove MitM Rules

```
DELETE /mitm/rules
```

Request Query:

| Key    | Type   | Desc                                               |
| ------ | ------ | -------------------------------------------------- |
| domain | string | Prefix wildcard `*` is allowed or just single `*`. |

Response Body:

```js
{
  "code": 0, // 0: success, 1: failed
  "message": "", // code==1, error message
  "data": [
    "*.google.com" // host list
  ]
}
```

