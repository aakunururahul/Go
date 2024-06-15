class Node:
    def __init__(self, key, val):
        self.key = key
        self.val = val
        self.nxt = None
        self.prev = None

class LRUCache:

    def __init__(self, capacity: int):
        self.capacity = capacity
        self.entries = {}
        self.rear = Node(-1, -1)
        self.front = Node(-1, -1)
        self.front.nxt = self.rear
        self.rear.prev = self.front

    def get(self, key: int) -> int:
        if key in self.entries.keys():
            node = self._remove(self.entries[key])
            self._add(node)
            return node.val
        else:
            return -1

    def put(self, key: int, value: int) -> None:
        if key in self.entries.keys():
            self._remove(self.entries[key])

        elif len(self.entries) >= self.capacity:
            rnode = self._remove(self.rear.prev)
            del self.entries[rnode.key]
        
        self.entries[key] = Node(key, value)
        self._add(self.entries[key])

    def _add(self, node) -> None:
        node.prev = self.front
        node.nxt = self.front.nxt
        self.front.nxt.prev = node
        self.front.nxt = node

    def _remove(self, rnode) -> Node:
        n1, n2 = rnode.prev, rnode.nxt
        n1.nxt, n2.prev = n2, n1
        rnode.prev, rnode.nxt = None, None
        return rnode
        
# Your LRUCache object will be instantiated and called as such:
# obj = LRUCache(capacity)
# param_1 = obj.get(key)
# obj.put(key,value)