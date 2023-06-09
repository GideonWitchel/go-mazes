package mazesrv

import "go-mazes/maze"

const MAZEHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Go Mazes</title>
    <style>
        #container {
            display: flex;
            justify-content: center;
            align-items: center;
        }
        th {
            border-style: none;
            background-color: white;
        }
        #maze {
            border-collapse: collapse;
            font-size: 4px;
        }
        .b-t {
            border-top-style: solid;
        }
        .b-b {
            border-bottom-style: solid;
        }
        .b-r {
            border-right-style:solid;
        }
        .b-l {
            border-left-style: solid;
        }
        .c-goal {
            background-color: yellow;
        }

        #buttons {
            display: flex;
            flex-direction: row;
            flex-wrap: wrap;
            justify-content: center;
            gap: 10px;
            margin: 30px;
        }
        #hint-best-path {
            text-align: center;
            padding: 20px;
            font-size: 20px;
        }

        #overlay{
            position: fixed;
            top: 0;
            bottom: 0;
            left: 0;
            right: 0;
            background-color: #000;
            opacity: 0.7;
            z-index: 10000;
        }
        #overlay-content{
            opacity: 1;
            color: white;
            margin: 10em;
            display: flex;
            flex-flow: column;
            align-items: center;
        }
    </style>

    <script>
        // TODO Move to separate js files
        let tickSpeed = {{ .TickSpeed }} ;
        let stepsFull = {{ .MPath }} ;
        let bestSteps = {{ .MBestPath }} ;
        let repeats = {{ .PathRepeats }} ;
        let halt = false;
        let formData = {{ .FormData }} ;

        let gen_dfs = "` + maze.GEN_DFS + `"
        let gen_none = "` + maze.GEN_NONE + `"
        let gen_rand = "` + maze.GEN_RAND + `"

        let solve_bfs_multi = "` + maze.SOLVE_BFS_MULTI + `"
        let solve_bfs_single = "` + maze.SOLVE_BFS_SINGLE + `"
        let solve_dfs_multi = "` + maze.SOLVE_DFS_MULTI + `"

        const timer = ms => new Promise(res => setTimeout(res, ms))
        window.addEventListener("load", async function () {
            document.getElementById("overlay").style.display = "none"
            initFormData()
            document.getElementById("maze").addEventListener("click", async function () {
                halt = !halt;
            });
            if (stepsFull.length > 0) {
                await drawAllPathsSimultaneously();
            }
            if (bestSteps.length > 0) {
                await drawBestPath();
            }
        });

        function initFormData() {
            let formStart = document.getElementById("buttons").firstElementChild.nextElementSibling
            let genStart = formStart.firstElementChild
            switch (formData[0]) {
                case gen_dfs:
                    genStart.selected = true
                    break
                case gen_rand:
                    genStart.nextElementSibling.selected = true
                    break
                default:
                    genStart.nextElementSibling.nextElementSibling.selected = true
                    break
            }
            let solveStart = formStart.nextElementSibling.nextElementSibling.nextElementSibling.firstElementChild
            switch (formData[1]) {
                case solve_bfs_single:
                    solveStart.selected = true
                    break
                case solve_bfs_multi:
                    solveStart.nextElementSibling.selected = true
                    break
                case solve_dfs_multi:
                    solveStart.nextElementSibling.nextElementSibling.selected = true
                    break
            }
            let numObj = formStart.nextElementSibling.nextElementSibling.nextElementSibling.nextElementSibling.nextElementSibling.nextElementSibling
            for (let i = 2; i < formData.length; i++) {
                numObj.value = formData[i]
                numObj = numObj.nextElementSibling.nextElementSibling.nextElementSibling
            }
        }

        async function drawAllPathsSimultaneously(){
            const promises = stepsFull.map(async step => {
                await drawMaze(step, getRandomColor(), 0, 0, repeats);
            })
            await Promise.all(promises);
        }

        async function drawAllPathsSequentially() {
            for(let i = 0; i < stepsFull.length; i++) {
                await drawMaze(stepsFull[i], getRandomColor(), 0, 0, repeats)
            }
        }

        async function drawBestPath() {
            // Wait for button press to draw the best path, if button hasn't already been pressed.
            document.getElementById("hint-best-path").style.display = "block"
            while (!halt) {
                await timer(tickSpeed);
            }
            // Put halt back to normal so drawMaze doesn't instantly pause
            halt = !halt
            document.getElementById("hint-best-path").style.display = "none"
            await drawMaze(bestSteps, "red", 1, 0, repeats);
        }

        // drawMaze goes through every index in steps and adds a class based on color.
        // startOffset is the index of steps that drawMaze starts on (to traverse full array, 0).
        // endOffset is how many before the end of steps that drawMaze will stop (to traverse full array, 0).
        // repeats is the number of cells drawMaze will fill for each time it delays with timer.
        // - This is necessary because waiting for timer, even at 0 tickSpeed, is very slow.
        async function drawMaze(steps, color, startOffset, endOffset, repeats){
            for (let i = startOffset; i < steps.length-endOffset; i++) {
                for (let j = 0; j < repeats; j++) {
                    if(i < steps.length-endOffset) {
                        let goalObj = getObjFromCoords(steps, i)
                        goalObj.style.backgroundColor = color;
                        i++;
                    }
                }
                i--;
                while (halt) {
                    await timer(tickSpeed);
                }
                await timer(tickSpeed);
            }
        }

        function getObjFromCoords(steps, index){
            let row = steps[index][0];
            let col = steps[index][1];

            let maze = document.getElementById("maze");
            // Pick the first row
            let r = maze.firstElementChild.firstElementChild;
            // Step to the correct row
            for (let i = 0; i < row; i++){
                r = r.nextElementSibling;
            }
            // Pick the first column
            let node = r.firstElementChild;
            // Step to the correct column
            for (let i = 0; i < col; i++){
                node = node.nextElementSibling;
            }
            return node;
        }

        function getRandomColor() {
            // Not too dark or light
            const firstLetters = '89ABCD';
            const secondLetters = '0123456789ABCDEF';
            let color = '#';
            for (let i = 0; i < 3; i++) {
                color += firstLetters[Math.floor(Math.random() * firstLetters.length)];
                color += secondLetters[Math.floor(Math.random() * secondLetters.length)];
            }
            return color;
        }

    </script>
</head>
<body>
    <div id="overlay">
        <div id="overlay-content">
            <p id="overlay-title" style="font-size: 5em;">Loading...</p>
            <p id="overlay-body" style="font-size: 1em;">Mazes! Mazes Everywhere!</p>
        </div>
    </div>

    <div id="buttons-container">
        <form action="/" id="buttons">
            <label for="generateAlgorithm">Generation algorithm:</label>
            <select name="generateAlgorithm" id="generateAlgorithm">
                <option value="` + maze.GEN_DFS + `" selected>DFS</option>
                <option value="` + maze.GEN_RAND + `">Random</option>
                <option value="` + maze.GEN_NONE + `">None</option>
            </select>
            <br>
            <label for="solveAlgorithm">Solving algorithm:</label>
            <select name="solveAlgorithm" id="solveAlgorithm">
                <option value="` + maze.SOLVE_BFS_SINGLE + `" selected>BFS</option>
                <option value="` + maze.SOLVE_BFS_MULTI + `">BFS Multithreaded</option>
                <option value="` + maze.SOLVE_DFS_MULTI + `">DFS Multithreaded</option>
            </select>
            <br>
            <label for="width">Width:</label>
            <input type="number" id="width" name="width" min="3" max="500" value="200">
            <br>
            <label for="height">Height:</label>
            <input type="number" id="height" name="height" min="3" max="500" value="112">
            <br>
            <label for="tickSpeed">Milliseconds per tick:</label>
            <input type="number" id="tickSpeed" name="tickSpeed" min="1" max="500" value="1">
            <br>
            <label for="repeats">Number of squares filled every tick: </label>
            <input type="number" id="repeats" name="repeats" min="1" max="1000" value="100">
            <br>
            <label for="density">Density (for randomly generated mazes): </label>
            <input type="number" id="density" name="density" min="1" max="100" value="15">
            <br>
            <input type="submit" value="Submit">
        </form>
    </div>
    <div id="container">
        <table id="maze">
            {{ range .MStyles }}<tr>
                {{ range . }}<th class="{{.}}">     </th>
                {{ end }}
            </tr>
            {{ end }}
        </table>
    </div>
    <h4 style="display: none" id="hint-best-path">Click on the maze to draw the solution!</h4>
</body>
</html>
`
