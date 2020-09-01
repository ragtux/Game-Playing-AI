/* Konane AI
* file: konane.go
* auth: ragtux
*/

package main

import (
    "fmt"
    "log"
    "io/ioutil"
    "strconv"
    "math"
    "bytes"
    "os"
)

//NUMBER of COLS and ROWS
var COLS, ROWS = 8, 8
var nodecnt = 0

var (
  Red     = Color("\033[1;31m%s\033[0m")
  Green   = Color("\033[1;32m%s\033[0m")
  Yellow  = Color("\033[1;33m%s\033[0m")
  Purple  = Color("\033[1;34m%s\033[0m")
  Magenta = Color("\033[1;35m%s\033[0m")
  Teal    = Color("\033[1;36m%s\033[0m")
  ket     = Color("\033[1;33m\033[41m%s\033[0m")
  far     = Color("\033[1;30m\033[42m%s\033[0m")
)

func Color(colorString string) func(...interface{}) string {
  sprint := func(args ...interface{}) string {
    return fmt.Sprintf(colorString,
      fmt.Sprint(args...))
  }
  return sprint
}

func readInBoard(filename string) ([][]int, error) {
    var m [][]int
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    rows := bytes.Split(data, []byte{'\n'})
    for r := len(rows) - 1; r >= 0; r-- {
        if len(rows[r]) != 0 {
            break
        }
        rows = rows[:len(rows)-1]
    }
    m = make([][]int, len(rows))
    nCols := 0
    for r, row := range rows {
        cols := bytes.Split(row, []byte{' '})
        if r == 0 {
            nCols = len(cols)
        }
        m[r] = make([]int, nCols)
        for c, col := range cols {
            if c < nCols && len(col) > 0 {
                e, err := strconv.ParseInt(string(col), 0, 8)
                if err != nil {
                    return nil, err
                }
                m[r][c] = int(e)
            }
        }
    }
    return m, nil
}

func b2str(board [][]int, mv []int) string {
    var result string = "  "
    for i,_ := range board[0] {
        result += strconv.Itoa(i) + " "
    }
    result += "\n"
    for x,row := range board {
        result += strconv.Itoa(x) + " "
        for y,_ := range row {
            if mv[0] == x && mv[1] == y {
                switch board[x][y] {
                case 1:
                    result += string(ket("B")) + " "
                case 2:
                    result += string(ket("W")) + " "
                case 0:
                    result += string(ket(".")) + " "
                }
            } else if mv[2] == x && mv[3] == y {
                switch board[x][y] {
                case 1:
                    result += string(ket("B")) + " "
                case 2:
                    result += string(ket("W")) + " "
                case 0:
                    result += string(ket(".")) + " "
                }
            } else {
                switch board[x][y] {
                case 1:
                    result += string(Teal("B")) + " "
                case 2:
                    result += string(Magenta("W")) + " "
                case 0:
                    result += string(Green(".")) + " "
                }
            }
        }
        result += "\n"
    }
    return result
}

func contains(board [][]int, row, col, player int) bool {
    onBoard := row >= 0 && col >= 0 && row < ROWS && col < COLS
    return onBoard && board[row][col] == player
}

//returns distance between (r1,c1) and (r2,c2) in a vertical line on the board
func distance(r1,c1,r2,c2 int) int {
    a := float64(r1-r2 + c1-c2)
    return int(math.Abs(a))
}

/*
* "make hypothetical move"
* PRE: A move for a particular player from (r1,c1) to (r2,c2)
* POST: executes the move on a copy of the current konane board. It
* returns the copy of the board, and does not change the given board.
*/
func makeHypoMove(board [][]int, player int, move []int) [][]int {

    r1,c1,r2,c2 := move[0],move[1],move[2],move[3]

    n := ROWS
    m := COLS

    nextBoard := make([][]int, n)
    data := make([]int, n*m)
    for i := range board {
        start := i*m
        end := start + m
        nextBoard[i] = data[start:end]
        copy(nextBoard[i], board[i])
    }

    dist := distance(r1,c1,r2,c2)
    jumps := dist / 2
    //fmt.Println("dist:",dist)
    //fmt.Println("jumps:",jumps)
    dr := (r2-r1)/dist
    dc := (c2-c1)/dist

    for i := 0 ; i < jumps ; i++ {
        nextBoard[r1][c1] = 0
        nextBoard[r1+dr][c1+dc] = 0
        r1 += 2*dr
        c1 += 2*dc
        nextBoard[r1][c1] = player
    }

    return nextBoard
}

/*
* Checks whether a jump is possible starting at (r,c) and going in the
* direction determined by the row delta, rd, and the column delta, cd.
* The factor is used to recursively check for multiple jumps in the same
* direction. Returns all possible jumps in the given direction.
*/
func check(board [][]int, r, c, opponent int) ([][]int) {

    done := false
    moves := [][]int{}

    for d := 0 ; d < 4 ; d++ {
        i:=1
        switch d {
        case 0: // moving north
            //fmt.Println("move north:")
            for !(done) {
                // adjacent opponent
                x := contains(board,r+i*(-1),c,opponent)
                // empty right after
                y := contains(board,r+(i+1)*(-1),c,0)
                //fmt.Println(x,y)
                if x && y {
                    var m = []int{r,c,r+(i+1)*(-1),c}
                    //fmt.Println("move north:",m)
                    moves = append(moves, m)
                } else {
                    done = true
                    //fmt.Println("move north DONE")
                }
                i+=2
            }
        case 1: // moving east
            //fmt.Println("moving east:")
            for !done {
                x := contains(board,r,c+i,opponent)
                y := contains(board,r,c+(i+1),0)
                if x && y {
                    var m = []int{r,c,r,c+(i+1)}
                    //fmt.Println("move east:",m)
                    moves = append(moves,m)
                } else {
                    done = true
                }
                i+=2
            }
        case 2: // moving south
            for !done {
                x := contains(board,r+i,c,opponent)
                y := contains(board,r+(i+1),c,0)
                //fmt.println(x,",",y)
                //fmt.println("adjacent opponent",r+i,c)
                //fmt.println("empty after",r+(i+1),c)
                if x && y {
                    var m = []int{r,c,r+(i+1),c}
                    moves = append(moves,m)
                } else {
                    done = true
                }
                i+=2
            }
        case 3: // moving west
            for !done {
                x := contains(board,r,c+i*(-1),opponent)
                y := contains(board,r,c+(i+1)*(-1),0)
                if x && y {
                    var m = []int{r,c,r,c+(i+1)*(-1)}
                    moves = append(moves, m)
                } else {
                    done = true
                }
                i+=2
            }
        }
        done = false
    }
    //fmt.Println(moves)
    return moves
}

/*
* Generates and returns all legal moves for the given player using the
* current board configuration.
*/
func generateMoves(board [][]int, player int) (int, int, [][]int) {
    moves := [][]int{}
    movablePieces := 0
    counted := false
    oppon := 3 - player
    for x,rows := range(board) {
        for y := range(rows) {
            if board[x][y] == player {
                // check directon for moves
                m := check(board,x,y,oppon)
                for _,mv := range m {
                    moves = append(moves, mv)
                    if !(counted) {
                        counted = true
                        movablePieces++
                    }
                }
            }
            counted = false
        }
    }
    return movablePieces,len(moves),moves
}

/*
* Returns resulting boards for all possible legal moves from
* given board for given player.
*/
func extendPath(board [][]int, player int) (int,int,[][][]int) {
    movablePieces,totalMoves,moves := generateMoves(board, player)
    boards := [][][]int{}
    for _,move := range(moves) {
        boards = append(boards, makeHypoMove(board, player, move))
    }
    return movablePieces,totalMoves,boards
}

// determine heuristic value of a given board
func eval(board [][]int, player int) int {

    plyMvblPieces,plyTotMoves,_ := generateMoves(board, player)
    oppMvblPieces,oppTotMoves,_ := generateMoves(board, 3 - player)

    //fmt.Println("      OG SEF:",plyMvblPieces,plyTotMoves)
    //fmt.Println("   OG OP SEF:",oppMvblPieces,oppTotMoves)
    //fmt.Println("incoming SEF:",mvps,totmvs)
    //fmt.Println("-------")

    a := plyMvblPieces - oppMvblPieces
    b := plyTotMoves - oppTotMoves
    nodecnt++
    //a := mvps - oppMvps
    //b := totmvs - oppMvps

    return int(a + b)
}

func minAndMax(a []int) (min int, max int) {
    min = a[0]
    max = a[0]
    for _, value := range a {
        if value < min {
            min = value
        }
        if value > max {
            max = value
        }
    }
    return min, max
}

/*
* Uses Minimax algorithm to give best possible board heuristic value up to
* end of game/given search depth assuming optimal opponent.
*/
func minimax(board [][]int,depth,limit,alpha,beta,player int) int {

    if depth >= limit {
        return eval(board,player)
    }

    nextBoards := [][][]int{}
    //newMvps := mvps
    //newTotmvs := totmvs
    isMax := depth % 2 == 0

    if isMax {
        //newMvps,newTotmvs,nextBoards = extendPath(board,player)
        _,_,nextBoards = extendPath(board,player)
    } else {
        //newMvps,newTotmvs,nextBoards = extendPath(board,3-player)
        _,_,nextBoards = extendPath(board,3-player)
    }

    if len(nextBoards) == 0 {
        if isMax {
            return -1000000
        } else {
            return 1000000
        }
    }

    values := []int{}
    newAlpha := alpha
    newBeta := beta
    for _,nextBoard := range(nextBoards) {
        if len(values) > 0 {
            if isMax {
                _,max := minAndMax(values)
                if max >= beta {
                    break
                }
                if max > alpha {
                    newAlpha = max
                }
            } else {
                min,_ := minAndMax(values)
                if min <= alpha {
                    break
                }
                if min < beta {
                    newBeta = min
                }
            }
        }
        values = append(values,minimax(nextBoard,depth+1,limit,newAlpha,newBeta,player))
    }
    //fmt.Println(values)
    //os.Exit(2)
    min,max := minAndMax(values)
    if isMax {
        return max
    } else {
        return min
    }
}

//// Given board returns the best move.
func getBestMove(board [][]int, player int, depth int) (int, []int) {

    values := []int{}
    largestIndex := 0
    _,_,moves := generateMoves(board, player)
    if len(moves) == 0 {
        return 0,values
    }
    alpha := -100000
    for _,move := range(moves) {
        values = append(values,minimax(makeHypoMove(board,player,move),1,depth,alpha,100000,player))
        _,max := minAndMax(values)
        if max > alpha {
            alpha = max
        }
    }
    m := 0
    for i, e := range values {
        if i == 0 || e > m {
            m = e
            largestIndex = i
        }
    }
    //os.Exit(3)
    return len(moves), moves[largestIndex]
}

/*
* Given two instances of players, will play out a game
* between them.  Returns 'B' if black wins, or 'W' if
* white wins. When show is true, it will display each move
* in the game.
* def makeHypoMove(board, player, move):
*/

func playAIvsAI(board [][]int, blackLevel, whiteLevel int) {
    bfactor := 0
    move := []int{}
    start := []int{0,0,0,0}
    fmt.Println("Black LV",blackLevel,"vs White LV",whiteLevel)
    turn := 1
    numJumps := 0
    fmt.Println("START:")
    fmt.Println(b2str(board,start))
    for {
        fmt.Println("[",turn,"] PLAYER B's TURN")
        bfactor,move = getBestMove(board,1,blackLevel)
        if len(move) == 0 {
            fmt.Printf("\n\nW WINS!!\n")
            break
        }
        board = makeHypoMove(board, 1, move)
        fmt.Printf("%s",b2str(board,move))
        turn++
        fmt.Println("     ",ket(move))
        numJumps = int(math.Abs(float64(distance(move[0],move[1],move[2],move[3]))/2))
        fmt.Println("         jumps: ", numJumps)
        fmt.Println("     nodecount: ", nodecnt)
        fmt.Println("      b factor: ", bfactor)
        nodecnt = 0
        fmt.Printf("\n")
        fmt.Println("[",turn,"] PLAYER W's TURN")
        bfactor,move = getBestMove(board,2,whiteLevel)
        if len(move) == 0 {
            fmt.Printf("\n\nB WINS!!\n")
            break
        }
        board = makeHypoMove(board, 2, move)
        fmt.Printf("%s",b2str(board,move))
        turn++
        fmt.Println("     ",ket(move))
        numJumps = int(math.Abs(float64(distance(move[0],move[1],move[2],move[3]))/2))
        fmt.Println("         jumps: ", numJumps)
        fmt.Println("     nodecount: ", nodecnt)
        fmt.Println("      b factor: ", bfactor)
        nodecnt = 0
        fmt.Printf("\n")
    }
    fmt.Println("\n\nGAME OVER")
}

func input(x []int, err error) []int {
    if err != nil {
        return x
    }
    var d int
    n, err := fmt.Scanf("%d", &d)
    if n == 1 {
        x = append(x, d)
    }
    return input(x, err)
}

func playHvsAI(board [][]int, whiteLevel int) {
    move := []int{}
    start := []int{0,0,0,0}
    fmt.Println("Black (Human) vs White LV",whiteLevel)
    turn := 1
    numJumps := 0
    fmt.Println("START:")
    fmt.Println(b2str(board,start))
    for {
        fmt.Println("[",turn,"] PLAYER B's TURN")
        //move = getBestMove(board,"B",blackLevel)
        fmt.Printf("Enter Move: ")
        move = input([]int{},nil)

        if len(move) == 0 {
            fmt.Printf("\n\nW WINS!!\n")
            break
        }
        board = makeHypoMove(board, 1, move)
        fmt.Println(b2str(board,move))
        turn++
        fmt.Println(move)
        numJumps = int(math.Abs(float64(distance(move[0],move[1],move[2],move[3]))/2))
        fmt.Println("jumps: ", numJumps)
        fmt.Println("[",turn,"] PLAYER W's TURN")
        _,move = getBestMove(board,2,whiteLevel)
        if len(move) == 0 {
            fmt.Printf("\n\nB WINS!!\n")
            break
        }
        board = makeHypoMove(board, 2, move)
        fmt.Println(b2str(board,move))
        turn++
        fmt.Println(move)
        numJumps = int(math.Abs(float64(distance(move[0],move[1],move[2],move[3]))/2))
        fmt.Println("jumps: ", numJumps)
    }
    fmt.Println("\n\nGAME OVER")
}

func main() {
    board, err := readInBoard("mjt.dat")
    if err != nil {
        log.Fatalf("readLines: %s", err)
    }

    //mv := []int{0,1,0,3}
    //fmt.Println(b2str(board,mv))
    playAIvsAI(board,8,8)
    //playHvsAI(board,8)
    os.Exit(3)
}
