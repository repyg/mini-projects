package main

import (
	"encoding/json"
	"log"
	"net/http"
	"text/template"
)

const gridSize = 100

// Grid представляет собой сетку для игры "Жизнь"
type Grid [gridSize][gridSize]bool

// NewGrid создает новую пустую сетку
func NewGrid() *Grid {
	return &Grid{}
}

// NextState вычисляет следующее состояние сетки на основе текущего состояния
func (g *Grid) NextState() {
	var newGrid Grid
	for y := 0; y < gridSize; y++ {
		for x := 0; x < gridSize; x++ {
			livingNeighbors := g.countLivingNeighbors(x, y)
			if g[y][x] && (livingNeighbors == 2 || livingNeighbors == 3) {
				newGrid[y][x] = true
			} else if !g[y][x] && livingNeighbors == 3 {
				newGrid[y][x] = true
			} else {
				newGrid[y][x] = false
			}
		}
	}
	*g = newGrid
}

// countLivingNeighbors подсчитывает количество живых соседей клетки в сетке
func (g *Grid) countLivingNeighbors(x, y int) int {
	count := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			}
			ni, nj := (x+i+gridSize)%gridSize, (y+j+gridSize)%gridSize
			if g[nj][ni] {
				count++
			}
		}
	}
	return count
}

// homeHandler обрабатывает запросы на домашнюю страницу
func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("public/home.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// gameHandler обрабатывает запросы на игровую страницу
func gameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var grid Grid
		err := json.NewDecoder(r.Body).Decode(&grid)
		if err != nil {
			log.Printf("Error decoding JSON: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		grid.NextState()
		err = json.NewEncoder(w).Encode(grid)
		if err != nil {
			log.Printf("Error encoding JSON: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	} else {
		// Отображение страницы с канвасом для игры
		tmpl, err := template.ParseFiles("public/game.html")
		if err != nil {
			log.Printf("Error parsing template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

func main() {
	// Обработка маршрутов
	http.HandleFunc("/", homeHandler)     // Домашняя страница
	http.HandleFunc("/game", gameHandler) // Страница игры

	// Подключение статических файлов (CSS, JS, изображения)
	// Этот обработчик сервирует файлы из папки public
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	// Запуск сервера
	log.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
