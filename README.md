# 🐍 Snake Game em Go

Um jogo Snake completo desenvolvido em Go como projeto de estudo para aprender os fundamentos e conceitos avançados da linguagem.

![Go Version](https://img.shields.io/badge/Go-1.25.3-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-green)
![Status](https://img.shields.io/badge/status-complete-success)

## 📋 Sobre o Projeto

Este projeto foi desenvolvido como uma jornada de aprendizado em Go, implementando um jogo Snake funcional com múltiplas features avançadas. O objetivo foi aplicar conceitos fundamentais e intermediários da linguagem através de um projeto prático e divertido.

## ✨ Features

- 🎮 **Menu Inicial** - Interface de boas-vindas com instruções
- 🔄 **Reiniciar Jogo** - Pressione R para jogar novamente
- 📊 **Sistema de Níveis** - Velocidade aumenta a cada 50 pontos
- 🏆 **High Score** - Recorde salvo em arquivo persistente
- ⭐ **Power-ups** - Comida especial que vale 50 pontos (20% de chance)
- 🧱 **Obstáculos** - Paredes aleatórias que aumentam com o nível
- 🔊 **Sistema de Sons** - Feedback sonoro para cada ação
- 🎨 **Interface ASCII** - Visual clean com caracteres Unicode

## 🎓 Conceitos de Go Aplicados

### Fundamentos

#### **Variáveis e Tipos**
```go
type Point struct {
    X int
    Y int
}

var score int = 0
direction := "right"
```
Aplicação: Estruturas de dados básicas para posições e estado do jogo.

#### **Structs**
```go
type Game struct {
    Snake      Snake
    Food       Food
    Score      int
    HighScore  int
    Level      int
    Obstacles  []Point
}
```
Aplicação: Organização do estado completo do jogo em estruturas lógicas.

#### **Slices**
```go
g.Snake.Body = append([]Point{newHead}, g.Snake.Body...)
g.Obstacles = []Point{}
```
Aplicação: Lista dinâmica para o corpo da cobra e obstáculos.

#### **Funções e Métodos**
```go
func (g *Game) MoveSnake() { }
func (g *Game) CheckCollision(p Point) bool { }
```
Aplicação: Métodos com receivers para organizar lógica relacionada ao Game.

---

### Intermediário

#### **Ponteiros**
```go
func (g *Game) Reset() {
    g.Score = 0
    g.Level = 1
}
```
Aplicação: Modificar estado do jogo sem criar cópias desnecessárias.

#### **Interfaces**
```go
type beep.Streamer interface {
    Stream(samples [][2]float64) (n int, ok bool)
    Err() error
}
```
Aplicação: Implementação de interface para geração de sons customizados.

#### **Enums (iota)**
```go
const (
    StateMenu GameState = iota
    StatePlaying
    StateGameOver
)
```
Aplicação: State machine para controlar fluxo do jogo.

---

### Avançado

#### **Goroutines**
```go
go game.HandleInput(end)
go playTone(800, 50*time.Millisecond)
```
Aplicação: 
- Leitura de input em paralelo ao game loop
- Sons não-bloqueantes

#### **Channels**
```go
end := make(chan bool)
end <- true
```
Aplicação: Comunicação entre goroutines para sinalizar término do jogo.

#### **Time e Ticker**
```go
ticker := time.NewTicker(game.Speed)
if game.Speed != lastSpeed {
    ticker.Stop()
    ticker = time.NewTicker(game.Speed)
}
```
Aplicação: Game loop com velocidade dinâmica baseada no nível.

#### **Manipulação de Arquivos**
```go
func LoadHighScore() int {
    data, err := os.ReadFile("highscore.txt")
    score, err := strconv.Atoi(strings.TrimSpace(string(data)))
    return score
}
```
Aplicação: Persistência do high score entre sessões.

#### **Packages Externos**
```go
import (
    "github.com/nsf/termbox-go"
    "github.com/faiface/beep"
)
```
Aplicação: 
- `termbox-go`: Renderização no terminal
- `beep`: Sistema de áudio

---

## 🎮 Como Jogar

### Controles
- **↑ ↓ ← →** : Movimentar a cobra
- **ENTER** : Iniciar jogo
- **R** : Reiniciar após game over
- **ESC** : Sair do jogo

### Regras
- **◆** Comida normal: 10 pontos
- **★** Power-up: 50 pontos
- **▓** Obstáculos: Evite!
- A cada 50 pontos você sobe de nível
- Cada nível aumenta velocidade e obstáculos

---

## 🚀 Instalação e Execução

### Pré-requisitos

- Go 1.25.3
- Git

### 1. Clone o Repositório

```bash
git clone https://github.com/seu-usuario/snake-game-go.git
cd snake-game-go
```

### 2. Instale as Dependências

```bash
go get -u github.com/nsf/termbox-go
go get -u github.com/faiface/beep
go get -u github.com/faiface/beep/speaker
```

Ou simplesmente:

```bash
go mod tidy
```

### 3. Execute o Jogo

```bash
go run snake.go
```

### 4. Build (Opcional)

Para gerar um executável:

**Windows:**
```cmd
go build -o snake.exe snake.go
snake.exe
```

**Linux/macOS:**
```bash
go build -o snake snake.go
./snake
```

---

## 📁 Estrutura do Projeto

```
snake-game-go/
├── snake.go    # Código principal
├── go.mod              # Dependências
├── go.sum              # Checksums das dependências
├── highscore.txt       # High score persistente (gerado automaticamente)
└── README.md           # Este arquivo
```

---

## 🏗️ Arquitetura

### State Machine

O jogo utiliza uma máquina de estados para controlar o fluxo:

```
┌─────────┐  ENTER   ┌──────────┐  Game Over  ┌──────────┐
│  Menu   │ ───────> │ Playing  │ ──────────> │GameOver  │
└─────────┘          └──────────┘             └──────────┘
     ^                                              │
     └──────────────────────────────────────────────┘
                          R (Restart)
```

### Game Loop

```go
for {
    select {
    case <-end:
        return
    case <-ticker.C:
        switch game.State {
        case StateMenu:
            game.DrawMenu()
        case StatePlaying:
            game.MoveSnake()
            game.Draw()
        case StateGameOver:
            game.DrawGameOver()
        }
    }
}
```

### Collision Detection

```go
func (g *Game) MoveSnake() {
    if g.CheckWallCollision(newHead) ||
       g.CheckSelfCollision(newHead) ||
       g.CheckObstacleCollision(newHead) {
        g.GameOver = true
    }
}
```

---

## 🎨 Sistema de Renderização

Utilizamos `termbox-go` para desenhar no terminal:

```go
termbox.SetCell(x, y, '█', termbox.ColorGreen, termbox.ColorDefault)
termbox.Flush()
```

**Caracteres usados:**
- `●` - Cabeça da cobra (amarela)
- `█` - Corpo da cobra (verde)
- `◆` - Comida normal (vermelha)
- `★` - Power-up (amarelo/magenta piscando)
- `▓` - Obstáculos (branco)
- `╔═╗║╚╝` - Bordas

---

## 🔊 Sistema de Som

Geração de tons senoidais customizados:

```go
type ToneGenerator struct {
    freq float64
    pos  float64
    sr   beep.SampleRate
}

func (t *ToneGenerator) Stream(samples [][2]float64) (n int, ok bool) {
    for i := range samples {
        v := math.Sin(t.pos * 2 * math.Pi * t.freq / float64(t.sr))
        samples[i][0] = v
        samples[i][1] = v
        t.pos++
    }
    return len(samples), true
}
```

**Sons implementados:**
- Comer comida: 800Hz, 50ms
- Power-up: 600→800→1000Hz (crescente)
- Level up: 1000→1200Hz (duplo)
- Game Over: 400→300→200Hz (descendente)

---

## 📊 Progressão de Dificuldade

| Nível | Velocidade | Obstáculos |
|-------|-----------|-----------|
| 1     | 150ms     | 2         |
| 2     | 140ms     | 4         |
| 3     | 130ms     | 6         |
| 5     | 110ms     | 10        |
| 10+   | 50ms      | 20 (max)  |

---

## 🛠️ Tecnologias Utilizadas

- **Linguagem:** Go 1.25.3
- **Terminal UI:** [termbox-go](https://github.com/nsf/termbox-go)
- **Áudio:** [beep](https://github.com/faiface/beep)
- **Ferramentas:** Go Modules

---

## 📚 Aprendizados

Este projeto cobriu:

✅ **Sintaxe Go:** Variáveis, tipos, structs, funções  
✅ **Controle de Fluxo:** If/else, loops, switch  
✅ **Estruturas de Dados:** Slices, arrays, maps  
✅ **Ponteiros:** Passagem por referência  
✅ **Concorrência:** Goroutines e channels  
✅ **Packages:** Importação e uso de bibliotecas  
✅ **I/O:** Leitura/escrita de arquivos  
✅ **Game Development:** Game loop, collision detection, state machines  
✅ **Audio Programming:** Geração de ondas sonoras  
✅ **Terminal UI:** Renderização de gráficos ASCII  

---

## 🚧 Possíveis Melhorias Futuras

- [ ] Multiplayer local (2 jogadores)
- [ ] Modo sem bordas (cobra atravessa paredes)
- [ ] Top 10 rankings
- [ ] Diferentes temas visuais
- [ ] Configurações de dificuldade
- [ ] Achievements/conquistas
- [ ] Pausa durante o jogo

---

## 📝 Licença

Este projeto está sob a licença MIT. Veja o arquivo `LICENSE` para mais detalhes.

---

## 👨‍💻 Autor

Desenvolvido como projeto de estudo em Go.

**Contato:**
- GitHub: [@mfugissecruz](https://github.com/mfugissecruz)
- Email: mfugissecruz@gmail.com

---

## 🙏 Agradecimentos

- Comunidade Go pela excelente documentação
- [termbox-go](https://github.com/nsf/termbox-go) por facilitar UI no terminal
- [beep](https://github.com/faiface/beep) pelo sistema de áudio simples

---

## 📖 Referências

- [A Tour of Go](https://go.dev/tour/)
- [Go by Example](https://gobyexample.com/)
- [Effective Go](https://go.dev/doc/effective_go)

---

**⭐ Se este projeto foi útil para você, considere dar uma estrela!**