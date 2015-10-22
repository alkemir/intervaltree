# intervaltree
A dynamic, thread-safe, self-prunning, self-balancing version of the segment-tree using AVL trees.

This project provides a structure which holds intervals of uint64. Data is structured as a binary search tree,
thus adding intervals and looking them up is efficient. If, at any time, intervals can be merged, this is done, so
memory consumption is kept low. Balancing is automatically done following AVL trees rules, so the structure warrants
scalability on operations.

## Operations
Currently implemented operations are:

* Insert
* Contains

Both have a complexity of O( log n ).

Insert results in either the addition of a node or the expansion of a node's interval.
The latter might involve also the removal of a node. Rebalancing is done afterwards.

Contains is performed as in any ordinary BST.

## Memory
Since this structure is a AVL tree, memory usage is bound to O(n). Take into account that since prunning is performed whenever
possible, you can expect the tree to consume less memory than n, depending on the sparsenes of your intervals.

## Concurrency
Operations are synchronized through a RWLock so the structure is thread-safe and multiple reads (Contains) can be executed at
the same time. This is ideal when you need a synchronized structure that will be updated rarely.

## Errors
Possible error conditions (invalid intervals, overlapping intervals being inserted) are detected and reported. Information
on the value causing the error is returned, so it is possible already to do some rudimentary error handling.

## Tests
Some very basic test cases have been included, but they cover a very narrow range of cases. Since this is a tree in principle
it should be possible to write tests for all possible cases and prove correctness of the implementation, but that is beyond
the scope of my current effort. If you would like to add some cases, please let me know.

## Contribute
This structure lends itself to many more operations (substract intervals, find overlaps, count individual intervals, etc...)
which might be useful. As it is, it works for me, but if you feel like implementing something and sharing it here please be
my guest.

Also, the implementation has not been optimized for performance; if you are feeling adventurous please share your ideas.

Finally, tests are shamingly incomplete. If you enjoy writing tests and see others code burst in flames,
please do so with mine!
