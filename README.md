# NDID Smart Contract

tendermint ABCI app

## Note

Test this app with command below

```sh
TENDERMINT_ADDRESS=http://localhost:45000 go test -v
```

## Add new validator (For testing)
get PubKey from pub_key.data in priv_validator.json 
```sh
curl -s 'localhost:45000/broadcast_tx_commit?tx="val:PubKey"'
```

## Prerequisites

* Go version >= 1.9.2

  * [Install Go](https://golang.org/dl/) by following [installation instructions.](https://golang.org/doc/install)
  * Set GOPATH environment variable (https://github.com/golang/go/wiki/SettingGOPATH)

* Tendermint 0.19.5

    ```sh
    mkdir -p $GOPATH/src/github.com/tendermint
    cd $GOPATH/src/github.com/tendermint
    git clone https://github.com/tendermint/tendermint.git
    cd tendermint
    git checkout v0.19.5
    make get_tools
    make get_vendor_deps
    make install
    ```

## Setup

1.  Create a directory for the project

    ```sh
    mkdir -p $GOPATH/src/github.com/ndidplatform/smart-contract
    ```

2.  Clone the project
    ```sh
    git clone https://github.com/ndidplatform/smart-contract.git $GOPATH/src/github.com/ndidplatform/smart-contract
    ```

3.  Get dependency (tendermint ABCI)

    ```sh
    cd $GOPATH/src/github.com/ndidplatform/smart-contract/abci
    dep ensure
    ```

### Run IdP node

1.  Run ABCI server

    ```sh
    cd $GOPATH/src/github.com/ndidplatform/smart-contract

    DB_NAME=IdP_DB go run abci/server.go tcp://127.0.0.1:46000
    ```

2.  Run tendermint

    ```sh
    cd $GOPATH/src/github.com/ndidplatform/smart-contract

    tendermint --home ./config/tendermint/IdP unsafe_reset_all && tendermint --home ./config/tendermint/IdP node --consensus.create_empty_blocks=false
    ```

### Run RP node

1.  Run ABCI server

    ```sh
    cd $GOPATH/src/github.com/ndidplatform/smart-contract

    DB_NAME=RP_DB go run abci/server.go tcp://127.0.0.1:46001
    ```

2.  Run tendermint

    ```sh
    cd $GOPATH/src/github.com/ndidplatform/smart-contract

    tendermint --home ./config/tendermint/RP unsafe_reset_all && tendermint --home ./config/tendermint/RP node --consensus.create_empty_blocks=false
    ```
    
### Run AS node

1.  Run ABCI server

    ```sh
    cd $GOPATH/src/github.com/ndidplatform/smart-contract

    DB_NAME=AS_DB go run abci/server.go tcp://127.0.0.1:46001
    ```

2.  Run tendermint

    ```sh
    cd $GOPATH/src/github.com/ndidplatform/smart-contract

    tendermint --home ./config/tendermint/AS unsafe_reset_all && tendermint --home ./config/tendermint/AS node --consensus.create_empty_blocks=false
    ```

## Run in Docker
Required
- Docker CE 17.06+ [Install docker](https://docs.docker.com/install/)
- docker-compose 1.14.0+ [Install docker-compose](https://docs.docker.com/compose/install/)

### Run

```
docker-compose -f docker/docker-compose.yml up
```

### Build

```
./docker/build.sh
```

### Note

* To run docker container without building image, run command in **Run** section (no building required). It will run docker container with image from Dockerhub (https://hub.docker.com/r/ndidplatform/mart-contract/).
* To pull latest image from Dockerhub, run `docker pull ndidplatform/smart-contract`
    
## IMPORTANT NOTE

1.  You must start IDP, RP and AS nodes in order to run the platform.
2.  After starting BOTH nodes, please wait for

    ```
    Commit
    Commit
    ```

    to show in the first terminal (`go run abci ...`) of both processes before starting `api` processes.

3.  When IDP node and RP node run on separate machines, please edit `seeds` in `config/tendermint/{RP or IdP}/config/config.toml` to match address of another machines.

## Technical details to connect with `api`

Interact with `api` in BASE64 format data.

# Broadcast tx format
```sh
functionName|parameter|nonce|base64(sign(param+nonce))|nodeID
```

# Query format
```sh
functionName|parameter
```

# Create transaction function

## InitNDID
### Parameter
```sh
{
  "node_id": "NDID",
  "public_key": "-----BEGIN RSA PUBLIC KEY-----\nMIIBCgKCAQEA30i6deo6vqxPdoxA9pUpuBag/cVwEVWO8dds5QDfu/z957zxXUCY\nRxaiRWGAbOta4K5/7cxlsqI8fCvoSyAa/B7GTSc3vivK/GWUFP+sQ/Mj6C/fgw5p\nxK/+olBzfzLMDEOwFRbnYtPtbWozfvceq77fEReTUdBGRLak7twxLrRPNzIu/Gqv\nn5AR8urXyF4r143CgReGkXTTmOvHpHu98kCQSINFuwBB98RLFuWdVwkrHyzaGnym\nQu+0OR1Z+1MDIQ9WlViD1iaJhYKA6a0G0O4Nns6ISPYSh7W7fI31gWTgHUZN5iTk\nLb9t27DpW9G+DXryq+Pnl5c+z7es/7T34QIDAQAB\n-----END RSA PUBLIC KEY-----\n"
}
```
### Expected Output
```sh
log: "success"
```

## RegisterNode
### Parameter
Posible role is RP,IdP and AS
```sh
{
  "node_id": "IdP1",
  "public_key": "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAwx9oT44DmDRiQJ1K0b9Q\nolEsrQ51hBUDq3oCKTffBikYenSUQNimVCsVBfNpKhZqpW56hH0mtgLbI7QgZGj9\ncNBMzSLMolltw0EerF0Ckz0Svvie1/oFJ1a0Cf4bdKKW6wRzL+aFVvelmNlLoSZX\noCpxUPQq7SMLoYEK1c+e3l3H0bfh6TAVt7APOQEFhXy9MRt83oVSAGW36gdNEksm\nz1WIT/C1XcHHVwCIJGSdZw5F6Y2gBjtiLsiFtpKfxQAPwBvDi7uS0PUdN7YQ/G69\nb0FgoE6qivDTqYfr80Y345Qe/qPGDvfne7oA8DIbRV+Kd5s4tFn/cC0Wd+jvrZJ7\njwIDAQAB\n-----END PUBLIC KEY-----\n",
  "master_public_key": "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEArdcKj/gAetVyg6Nn2lDi\nm/UJYQsQCav60EVbECm5EVT8WgnpzO+GrRyBtxqWUdtGar7d6orLh1RX1ikU7Yx2\nSA8Xlf+ZDaCELba/85Nb+IppLBdPywixgumoto9G9dDGSnPkHAlq5lXXA1eeUS7j\niU1lf37lwTZaO0COAuu8Vt9GcwYPh7SSf4/eXabQGbo/TMUVpXX1w5N1A07Qh5DG\nr/ZKzEE9/5bJJJRS635OA2T4gIY9XRWYiTxtiZz6AFCxP92Cjz/sNvSc/Cuvwi15\nycS4C35tjM8iT5djsRcR+MJeXyvurkaYgMGJTDIWub/A5oavVD3VwusZZNZvpDpD\nPwIDAQAB\n-----END PUBLIC KEY-----\n",
  "node_name": "IdP Number 1 from ...",
  "role": "IdP",
  "max_ial": 3,
  "max_aal": 3
}
```
### Expected Output
```sh
log: "success"
```

## RegisterMsqDestination
### Parameter
```sh
{
  "users": [
    {
      "hash_id": "ece8921066562be07ba4ec44449646fc3b48d6b8a660a2e1e6a4bc7117edebba",
      "ial": 3
    }
  ],
  "node_id": "IdP1"
}
```
### Expected Output
```sh
log: "success"
```

## AddService
### Parameter
```sh
{
  "service_id": "statement",
  "service_name": "Bank statement"
}
```
### Expected Output
```sh
log: "success"
```

## DeleteService
### Parameter
```sh
{
  "service_id": "statement"
}
```
### Expected Output
```sh
log: "success"
```

## RegisterServiceDestination
### Parameter
```sh
{
  "service_id": "statement",
  "node_id": "AS1",
  "min_ial": 1.1,
  "min_aal": 1.2
}
```
### Expected Output
```sh
log: "success"
```

## RegisterMsqAddress
### Parameter
```sh
{
  "node_id": "IdP1",
  "ip": "192.168.3.99",
  "port": 8000
}
```
### Expected Output
```sh
log: "success"
```

## CreateRequest
### Parameter
```sh
{
  "request_id": "ef6f4c9c-818b-42b8-8904-3d97c4c520f6",
  "min_idp": 1,
  "min_aal": 3,
  "min_ial": 3,
  "timeout": 259200,
  "data_request_list": [
    {
      "service_id": "statement",
      "as_id_list": [
        "AS1",
        "AS2"
      ],
      "count": 1,
      "request_params_hash": "hash"
    }
  ],
  "request_message_hash": "hash('Please allow...')",
  "mode": 3
}
```
### Expected Output
```sh
log: "success"
```

## CreateIdpResponse
### Parameter
```sh
{
  "request_id": "ef6f4c9c-818b-42b8-8904-3d97c4c520f6",
  "aal": 3,
  "ial": 3,
  "status": "accept",
  "signature": "signature",
  "identity_proof": "Magic"
}
```
### Expected Output
```sh
log: "success"
```

## SignData
### Parameter
```sh
{
  "service_id": "statement",
  "request_id": "ef6f4c9c-818b-42b8-8904-3d97c4c520f6",
  "signature": "sign(data,asKey)"
}
```
### Expected Output
```sh
log: "success"
```

## SetNodeToken
### Parameter
```sh
{
  "node_id": "RP1",
  "amount": 100
}
```
### Expected Output
```sh
log: "success"
```

## AddNodeToken
### Parameter
```sh
{
  "node_id": "RP1",
  "amount": 111.11
}
```
### Expected Output
```sh
log: "success"
```

## ReduceNodeToken
### Parameter
```sh
{
  "node_id": "RP1",
  "amount": 61.11
}
```
### Expected Output
```sh
log: "success"
```

## SetPriceFunc
### Parameter
```sh
{
  "func": "CreateRequest",
  "price": 99.99
}
```
### Expected Output
```sh
log: "success"
```

## CloseRequest
### Parameter
```sh
{
  "requestId": "ef6f4c9c-818b-42b8-8904-3d97c4c520f6"
}
```
### Expected Output
```sh
log: "success"
```

## TimeOutRequest
### Parameter
```sh
{
  "requestId": "ef6f4c9c-818b-42b8-8904-3d97c4c11111"
}
```
### Expected Output
```sh
log: "success"
```

## AddNamespace
### Parameter
```sh
{
  "namespace": "CID",
  "description": "Citizen ID"
}
```
### Expected Output
```sh
log: "success"
```

## DeleteNamespace
### Parameter
```sh
{
  "namespace": "Tel"
}
```
### Expected Output
```sh
log: "success"
```

## UpdateNode
### Parameter
```sh
{
  "public_key": "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEArdcKj/gAetVyg6Nn2lDi\nm/UJYQsQCav60EVbECm5EVT8WgnpzO+GrRyBtxqWUdtGar7d6orLh1RX1ikU7Yx2\nSA8Xlf+ZDaCELba/85Nb+IppLBdPywixgumoto9G9dDGSnPkHAlq5lXXA1eeUS7j\niU1lf37lwTZaO0COAuu8Vt9GcwYPh7SSf4/eXabQGbo/TMUVpXX1w5N1A07Qh5DG\nr/ZKzEE9/5bJJJRS635OA2T4gIY9XRWYiTxtiZz6AFCxP92Cjz/sNvSc/Cuvwi15\nycS4C35tjM8iT5djsRcR+MJeXyvurkaYgMGJTDIWub/A5oavVD3VwusZZNZvpDpD\nPwIDAQAB\n-----END PUBLIC KEY-----\n",
  "master_public_key": ""
}
```
### Expected Output
```sh
log: "success"
```

## CreateIdentity
### Parameter
```sh
{
  "accessor_id": "accessor_id",
  "accessor_type": "accessor_type",
  "accessor_public_key": "accessor_public_key",
  "accessor_group_id": "accessor_group_id"
}
```
### Expected Output
```sh
log: "success"
```

## AddAccessorMethod
### Parameter
```sh
{
  "accessor_id": "accessor_id_2",
  "accessor_type": "accessor_type_2",
  "accessor_public_key": "accessor_public_key_2",
  "accessor_group_id": "accessor_group_id",
  "request_id": "ef6f4c9c-818b-42b8-8904-3d97c4c520f6"
}
```

## SetValidator
### Parameter
```sh
{
  "public_key": "5/6rEo7aQYq31J32higcxi3i8xp9MG/r5Ho5NemwZ+g=",
  "power": 0
}
```
### Expected Output
```sh
log: "success"
```

## SetDataReceived
### Parameter
```sh
{
  "requestId": "ef6f4c9c-818b-42b8-8904-3d97c4c520f6",
  "service_id": "statement",
  "as_id": "AS1"
}
```
### Expected Output
```sh
log: "success"
```

# Query function

## GetNodeMasterPublicKey
### Parameter
```sh
{
  "node_id": "RP1"
}
```
### Expected Output
```sh
{
  "master_public_key": "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1QXXrV7X1b8uFL1PW7+F\nimlAwxwbEMG5hFru1CN8WsRt8ZVQIkXRpiwNNXh1GS0Qmshnv8pKaNCZ5q5wFdUe\nlYspZHVRbIkHiQAaEU5yG9SyavHsDntUOd50PQ3nC71feW+ff8tvQcJ7+gqf8nZ6\nUAWpG4bvakPtrJ81h4/Qc23vhtbcouP0adgdw6UA0kcdGhTESYMBU0dx/NNysvJh\nNx36z2UU6kbQ3a2/bINEZAgLfJ7/Y+/647+tc7bUYdqj3dNkbnk1xiXh5dTLsiow\n5Xvukpy2uA44M/r2Q5VRfbH2ZrBZlgf/XEOZs7zppySgaTWRB5eDTm+YxxyOyykn\n8wIDAQAB\n-----END PUBLIC KEY-----\n"
}
```

## GetNodePublicKey
### Parameter
```sh
{
  "node_id": "RP1"
}
```
### Expected Output
```sh
{
  "public_key": "-----BEGIN RSA PUBLIC KEY-----\nMIIBCgKCAQEAwCB4UBzQcnd6GAzPgbt9j2idW23qKZrsvldPNifmOPLfLlMusv4E\ncyJf4L42/aQbTn1rVSu1blGkuCK+oRlKWmZEWh3xv9qrwCwov9Jme/KOE98zOMB1\n0/xwnYotPadV0de80wGvKT7OlBlGulQRRhhgENNCPSxdUlozrPhrzGstXDr9zTYQ\noR3UD/7Ntmew3mnXvKj/8+U48hw913Xn6btBP3Uqg2OurXDGdrWciWgIMDEGyk65\nNOc8FOGa4AjYXzyi9TqOIfmysWhzKzU+fLysZQo10DfznnQN3w9+pI+20j2zB6gg\npL75RjZKYgHU49pbvjF/eOSTOg9o5HwX0wIDAQAB\n-----END RSA PUBLIC KEY-----\n"
}
```

## GetIdpNodes
### Parameter
```sh
{
  "hash_id": "ece8921066562be07ba4ec44449646fc3b48d6b8a660a2e1e6a4bc7117edebba",
  "min_ial": 3,
  "min_aal": 3,
}
```
### Expected Output
```sh
{
  "node": [
    {
      "id": "IdP1",
      "name": "IdP Number 1 from ...",
      "max_ial": 3,
      "max_aal": 3
    },
    {
      "id": "IdP2",
      "name": "",
      "max_ial": 3,
      "max_aal": 3
    }
  ]
}
```

## GetAsNodesByServiceId
### Parameter
```sh
{
  "service_id": "statement"
}
```
### Expected Output
```sh
{
  "node": [
    {
      "id": "AS1",
      "name": "AS1",
      "min_ial": 1.1,
      "min_aal": 1.2,
      "service_name": "Bank statement"
    }
  ]
}
```

## GetMsqAddress
### Parameter
```sh
{
  "node_id": "IdP1"
}
```
### Expected Output
```sh
{
  "ip": "192.168.3.99",
  "port": 8000
}
```

## GetRequest
### Parameter
```sh
{
  "requestId": "ef6f4c9c-818b-42b8-8904-3d97c4c520f6"
}
```
### Expected Output
```sh
{
  "closed": false,
  "timed_out": true,
  "request_message_hash": "hash('Please allow...')"
}
```

## GetRequestDetail
### Parameter
```sh
{
  "requestId": "ef6f4c9c-818b-42b8-8904-3d97c4c520f6"
}
```
### Expected Output
```sh
{
  "request_id": "ef6f4c9c-818b-42b8-8904-3d97c4c520f6",
  "min_idp": 1,
  "min_aal": 3,
  "min_ial": 3,
  "request_timeout": 259200,
  "data_request_list": [
    {
      "service_id": "statement",
      "as_id_list": [
        "AS1",
        "AS2"
      ],
      "count": 1,
      "request_params_hash": "hash",
      "answered_as_id_list": [
        "AS1"
      ],
      "received_data_from_list": [
        "AS1"
      ]
    }
  ],
  "request_message_hash": "hash('Please allow...')",
  "responses": [
    {
      "request_id": "ef6f4c9c-818b-42b8-8904-3d97c4c520f6",
      "aal": 3,
      "ial": 3,
      "status": "accept",
      "signature": "signature",
      "identity_proof": "Magic",
      "private_proof_hash": "",
      "idp_id": "IdP1"
    }
  ],
  "closed": false,
  "timed_out": false,
  "special": false,
  "mode": 3
}
```

## GetNodeToken
### Parameter
```sh
{
  "node_id": "RP1"
}
```
### Expected Output
```sh
{
  "amount": 100
}
```

## GetPriceFunc
### Parameter
```sh
{
  "func": "CreateRequest"
}
```
### Expected Output
```sh
{
  "price": 99.99
}
```

## GetUsedTokenReport
### Parameter
```sh
{
  "node_id": "AS1"
}
```
### Expected Output
```sh
[
  {
    "method": "RegisterServiceDestination",
    "price": 1,
    "data": ""
  },
  {
    "method": "SignData",
    "price": 1,
    "data": "ef6f4c9c-818b-42b8-8904-3d97c4c520f6"
  }
]
```

## GetServiceDetail
### Parameter
```sh
{
  "service_id": "statement"
}
```
### Expected Output
```sh
{
  "service_name": "Bank statement"
}
```

## GetNamespaceList
### Parameter
```sh
```
### Expected Output
```sh
[
  {
    "namespace": "CID",
    "description": "Citizen ID"
  }
]
```

## CheckExistingIdentity
### Parameter
```sh
{
   "hash_id": "hash(ns+id)"
}
```
### Expected Output
```sh
{
  "exist": true
}
```

## GetAccessorGroupID
### Parameter
```sh
{
  "accessor_id": "accessor_id_001"
}
```
### Expected Output
```sh
{
  "accessor_group_id": "accessor_group_id"
}
```

## GetAccessorKey
### Parameter
```sh
{
  "accessor_id": "accessor_id_001"
}
```
### Expected Output
```sh
{
  "accessor_public_key": "accessor_public_key"
}
```

## GetServiceList
### Parameter
```sh
```
### Expected Output
```sh
[
  {
    "service_id": "statement",
    "service_name": "Bank statement"
  }
]
```
