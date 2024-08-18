const gridSize = 50;
let grid = Array.from({ length: gridSize }, () => Array(gridSize).fill(false));
let intervalId;
let isRunning = false;

const canvas = document.getElementById('gameCanvas');
const ctx = canvas.getContext('2d');
canvas.width = window.innerWidth;
canvas.height = window.innerWidth;
const cellSize = canvas.width / gridSize;

canvas.addEventListener('click', (event) => {
    const rect = canvas.getBoundingClientRect();
    const x = Math.floor((event.clientX - rect.left) / cellSize);
    const y = Math.floor((event.clientY - rect.top) / cellSize);
    grid[y][x] = !grid[y][x];
    drawGrid();
});

document.getElementById('startBtn').addEventListener('click', () => {
    if (!isRunning) {
        intervalId = setInterval(stepGame, 100);
        isRunning = true;
    }
});

document.getElementById('pauseBtn').addEventListener('click', () => {
    clearInterval(intervalId);
    isRunning = false;
});

function drawGrid() {
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    for (let y = 0; y < gridSize; y++) {
        for (let x = 0; x < gridSize; x++) {
            ctx.fillStyle = grid[y][x] ? '#9b59b6' : '#000000';
            ctx.fillRect(x * cellSize, y * cellSize, cellSize, cellSize);
        }
    }
}

function stepGame() {
    fetch('/game', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(grid),
    })
    .then(response => response.json())
    .then(newGrid => {
        grid = newGrid;
        drawGrid();
    })
    .catch(error => console.error('Error:', error));
}

drawGrid();
