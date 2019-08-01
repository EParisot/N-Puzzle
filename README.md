# N-Puzzle
The goal of this project is to solve the N-puzzle ("taquin" in French) game using the A* search algorithm or one of its variants.

We use and dynamically Weighted IDA* to solve the problem with the best solution in a 10s time interval, if time exceed 10s, then Weight is incremented, and so on...

![](n_puzzle.gif)

## Usage:
Build :
```
go build
```

Run
```
./N-Puzzle[.exe] [-m mapFile] [-i imageFile] [-d difficulty] [-a heuristic] [-gs] [-g] [-dg]
			-m mapFile    = 'map_file.map'
			-i imageFile  = 'image_file.png'
			-s size       = map size (int)
			-h heuristic  = 'heuristic' ('md' (default), 'hd', 'ed', 'lc')
			-gs           = Greedy Search (cost g(x) = 0)
			-g            = Graphical Interface
			-dg           = Add numbers to the picture
```