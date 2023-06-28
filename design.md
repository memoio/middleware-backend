# design

## modules

### auth

登录验证模块，用于接入和验证操作的合法性

### pay

包括读写支付和数据交易支付

原则：
每个账户包含一定量的免费读写空间，由管理员申请


### did

账户did和文件did，用于文件nft化和交易（付费获取did的读取权限）

### data

数据写入和读取

### sdk

js调用api以及常用的工具函数

## address as account

### data structure

#### file

```sol
struct file {
    uint64 sizeByte;
    bool active; // deactive if delete; no one can access
} 

mapping(address => mapping(string => file)) public files; //each account has one
```

#### file dns

```sol
struct fileDNS {
    address owner;
    expire uint256;
}

mapping(string => fileDNS) public fdns; // global file space; buy dns; can replace if expire
```

#### space

```sol
struct package {
    uint64 sizeByte;
    uint64 startDay;
    uint64 durationDay; 
}

struct storage {
    []package packages; // add when buy package
}

mapping(address => storage) public space; // key is address
```

#### read pay

```sol
struct payment {
    uint256 nonce;  // plus one when receiver withdraw
    uint256 money;  
    uint256 expire; // sender can withdraw after expire  
}
address receiver; // receive money 
mapping(address => payment) public pay; // key is sender address
```

### operations

1. based on http; (operation, params...)
2. gateway publickey (cold start)
3. space w/o duration
4. token used to limit (request one by one)

#### config

1. client send (config)
2. server send back config (space contract; file contract; read contract; dns contract)

#### login

1. client generate random string as token, and send (login, addr, token, timestamp, sign(token, timestamp)) to server
2. server check timestamp(< 1 minute) and sign, then reset token and nonce;

#### info

public info:
1. client send (info, addr)
2. server send back addr's info(token, nonce)

private info:
1. client send (info, addr, token, nonce, sign(token, nonce))
2. server validate nonce and sign, then send back addr's info()
  
#### buy space

no duration; upload limit; delay deletion

buy space using ERC20 token

##### direct

1. client submit tx to chain

##### by middleware

need receipt for buying

#### write

client:
1. send (write, addr, token, nonce, sign(token, nonce, hash(payload)), payload)

server: 
1. check nonce and sign
2. check write permission
3. check space (size sign)
4. handle data

#### buy payment (channel)

buy payment using ERC20 token

##### direct

##### by middleware

#### read

1. check nonce and sign
2. check read permission (only read self data)
3. check money (channel sign)
4. send back data

## did as account (todo)

#### account DID

```sol
// adid: "did:memo:<account id>"; account id is 64 bytes
struct accountDID {
    mapping(string => bool) controller; // add/del controller by masterkey; add/del other infos by controlelr
    PublicKey[] verificationMethod;     // public key here
    mapping(string => bool) authentication;  // login 
    mapping(string => bool) assertionMethod;
    mapping(string => uint) capabilityDelegation;
}

mapping(string => accountDID) public accounts; // key is did string
```

#### file DID

1. 
2. read control: buy read permission and add to read 

```sol
// fdid: "did:mfile:<mefs id>"
struct fileDID {
    uint64 sizeByte;
    bool active;                 
    string controller;           // change by masterKey of controller 
    mapping(string => bool) read;// change by capabilityDelegation of controller
}

// from account did => file did => fileDID
mapping(string => mapping(string => fileDID)) public files;
```

#### file dns

```sol
struct fileDNS {
    string owner;
    expire uint256;
}

mapping(string => fileDNS) public fdns; // global file space; can replace if expire
```

## address as account

### address manage

+ login: get auth token for writing; need if only one change its
+ auth:


### space manage

+ create: pay and get storage 
+ 


## did as account
 



each address can generate multiple account DID:

+ account DID: did:memo:<specific-id>; specific-id has length 64 byte; example: did:memo:ce5ac89f84530a1cf2cdee5a0643045a8b0a4995b1c765ba289d7859cfb1193e
+ generate account DID: hex(hash32(address, nonce/specific-string))
+ submit account DID to chain or submitted by middleware

### manage

+ register: create account DID and submit to chain
+ login: get auth token for writing; need if only one change its
+ auth: 

### storage space  

+ create

+ manage
+ write
+ read

### pay

+ write
+ read