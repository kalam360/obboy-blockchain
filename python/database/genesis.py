import json

with open("./genesis.json") as f:
    genesis = json.load(f)


class Tx:
    def __init__(self) -> None:
        self.To = ""
        self.From = ""
        self.Value = None
        self.Data = None


class State:
    def __init__(self) -> None:
        self.balances = {}
        self.txMempool = []
    
    def apply(tx: Tx) -> None:
        pass