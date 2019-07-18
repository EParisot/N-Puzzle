# N-Puzzle
The goal of this project is to solve the N-puzzle ("taquin" in French) game using the A* search algorithm or one of its variants.

We use and dynamically Weighted IDA* to solve the problem with the best solution in a 10s time interval, if time exceed 10s, then Weight is incremented, and so on...

## Usage:
Build :
```
go build
```

Run
```
./N-Puzzle[.exe] [-m mapFile] [-i imageFile] [-d difficulty] [-a heuristic] [-g] [-dg]
                        -m mapFile    = 'map_file.map'
                        -i imageFile  = 'image_file.png'
                        -s size       = map size (int)
                        -h heuristic  = 'heuristic' ('md' (default), 'hd', 'i')
                        -dg (Add numbers to the picture)
                        -g (Active graphic)
```