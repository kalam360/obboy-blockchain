from os import stat
import sys
from typing import Optional

import uvicorn

import hashlib
import json

from time import time
from uuid import uuid4

from fastapi import FastAPI
from fastapi.exceptions import HTTPException
from pydantic import BaseModel

import requests
from urllib.parse import urlparse


class Blockchain(object):
    difficulty_target = "0000"
    def hash_block(self, block):
        block_encoded = json.dumps(block, sort_keys=True).encode()

        return hashlib.sha256(block_encoded).hexdigest()

    def __init__(self):
        # nodes
        self.nodes = set()
        # Stores all the blocks in this list
        self.chain = []
        # Store current transactions in temporarily which will be transfered to blocks when mined
        self.current_transaction = []
        # Create a genesis block with a specified fixed hash of previous block starts with index 0
        genesis_hash = self.hash_block("genesis_block")
        self.append_block(prev_block_hash = genesis_hash, nonce=self.proof_of_work(0, genesis_hash, []))

    def add_node(self, address):
        parsed_url = urlparse(address)
        self.nodes.add(parsed_url.netloc)
        print(parsed_url.netloc)

    def valid_chain(self, chain):
        last_block = chain[0]
        current_index = 1
        while current_index < len(chain):
            block = chain[current_index]
            if block['prev_block_hash'] != self.hash_block(last_block):
                return False
            
            # check for valid nonce
            if not self.valid_proof(current_index, block['prev_block_hash'], block['transactions'], block['nonce']):
                return False

            # Move to the next block
            last_block = block
            current_index += 1
        ## The chain is valid
        return True

    def update_blockchain(self):

        # get the nodes from the network
        neighbours = self.nodes
        new_chain = None
        
        # find the chains longer than this node
        max_length = len(self.chain)

        ## Grab all the nodes and verify 
        for node in neighbours:
            # get the blockchain from other nodes
            response = requests.get(f'http://{node}/blockchain')

            if response.status_code == 200:
                length = response.json()['length']
                chain = response.json()['chain']
                ## check if the chain is longer and the chain is valid
                if length > max_length and self.valid_chain(chain):
                    max_length = length
                    new_chain = chain
            ## Replace this node chain with valid and longer chain
            if new_chain:
                self.chain = new_chain
            return True
        return False

    def proof_of_work(self, index, prev_block_hash, transactions):
        # start with nonce 0
        nonce = 0

        # hashing the nonce with previous block to see if its valid
        while self.valid_proof(index, prev_block_hash, transactions, nonce) is False:
            nonce += 1
        
        return nonce

    def valid_proof(self, index, prev_block_hash, transactions, nonce):
        # Create a string containing hash of the previous block and content
        content = f'{index}{prev_block_hash}{transactions}{nonce}'.encode()
        # sha256 hash with content
        content_hash = hashlib.sha256(content).hexdigest()

        # Check if hash meets difficulty target
        return content_hash[:len(self.difficulty_target)] == self.difficulty_target

    ## Create new block to add in the blockchain
    def append_block(self, nonce, prev_block_hash):
        block = {
            'index': len(self.chain),
            'timestamp': time(),
            'transactions': self.current_transaction,
            'nonce': nonce,
            'prev_block_hash': prev_block_hash
        }
        self.current_transaction = []
        self.chain.append(block)

        return block

    def add_transactions(self, sender, reciever, value, data):
        self.current_transaction.append({
            'reciever': reciever,
            'sender': sender,
            'value': value,
            'data': data
        })
        return self.last_block['index'] + 1

    @property
    def last_block(self):
        # returns the last block in the blockchain
        return self.chain[-1]

    
app = FastAPI()

## generate a globally unique address for this node

node_indentifier = str(uuid4()).replace('-', '')

## instantiate the blockchain 

blockchain = Blockchain()


# return the entire blockchain

@app.get('/blockchain', status_code= 200)
def full_chain():
    return {
        'chain': blockchain.chain,
        'height': len(blockchain.chain)
    }

class Nodes(BaseModel):
    nodes: list

@app.post('/nodes/add_nodes', status_code=201)
def add_nodes(nodes: Nodes):
    # get the nodes passed
    nodes_list = nodes.nodes
    if len(nodes_list) == 0:
        return HTTPException(status_code=400, detail= "Missing nodes info")

    for node in nodes_list:
        blockchain.add_node(node)

    return {
        "message": "New nodes added",
        "nodes": list(blockchain.nodes)
    }
    
@app.get("/nodes/sync", status_code= 200)
def sync():
    updated = blockchain.update_blockchain()
    if updated:
        response = {
            "message": "The chain is updated to the latest",
            "blockchain": blockchain.chain
        }
    else:
        response = {
            "message": "Current chain is the latest", 
            "blockchain": blockchain.chain
        }
    return response


@app.get('/mine', status_code= 200)
def mine_block():
    blockchain.add_transactions(sender="0", reciever=node_indentifier, value=1, data="mining")

    # obtain hash of the previous block
    last_block_hash = blockchain.hash_block(blockchain.last_block)

    # using PoW, get the nonce for the new block to be added to blockchain
    index = len(blockchain.chain)

    nonce = blockchain.proof_of_work(index, last_block_hash, blockchain.current_transaction)

    # add new block to the blockchain using last block hash and current nonce

    block = blockchain.append_block(nonce, last_block_hash)
  
    return {
        'message': "new block mined",
        'index': block['index'],
        'prev_block_hash': block['prev_block_hash'],
        'nonce': block['nonce'],
        'transactions': block['transactions']
    }

class Tx(BaseModel):
    sender: str
    reciever: str
    value: float
    data: Optional[str] = None



## adding transactions
@app.post("/transactions/new", status_code=201)
def new_transaction(tx: Tx):
    blockchain.add_transactions(sender=tx.sender, reciever=tx.reciever, value=tx.value, data= tx.data)
    return blockchain.current_transaction[-1]
    

if __name__ == "__main__":
    uvicorn.run("blockchain:app", host="127.0.0.1", port = int(sys.argv[1]), log_level="info", reload=True)
