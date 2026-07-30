```
  mmmm             #  mmmmm           m      |    ⠀⢀⣴⣾⣿⣿⣿⣷⣦⡄⠀⣴⣾⣿⣿⣿⣿⣶⣄⠀⠀
 #"   "  mmm    mmm#  #   "#  mmm   mm#mm    |    ⣰⣿⣿⣿⣿⣿⣿⣿⠋⢠⣾⣿⣿⣿⣿⣿⣿⣿⣿⣧⠀
 "#mmm  "   #  #" "#  #mmmm" "   #    #      |    ⣿⣿⣿⣿⣿⣿⣿⣿⣿⣶⣌⠛⣿⣿⣿⣿⣿⣿⣿⣿⡆ 
     "# m"""#  #   #  #   "m m"""#    #      |    ⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⢁⣼⣿⣿⣿⣿⣿⣿⣿⣿⠁
 "mmm#" "mm"#  "#m##  #    " "mm"#    "mm    |    ⠸⣿⣿⣿⣿⣿⣿⣿⡟⢀⣾⣿⣿⣿⣿⣿⣿⣿⣿⠏⠀
                                             |    ⠀⠙⣿⣿⣿⣿⣿⣿⣄⠻⣿⣿⣿⣿⣿⣿⣿⣿⠏⠀⠀
                                             |    ⠀⠀⠈⠻⣿⣿⣿⣿⣿⣧⡈⢿⣿⣿⣿⣿⡟⠁⠀⠀⠀
                                             |    ⠀⠀⠀⠀⠈⠻⣿⣿⣿⣿⡇⢸⣿⣿⠟⠉⠀⠀⠀⠀⠀
                                             |    ⠀⠀⠀⠀⠀⠀⠈⠙⢿⡿⠀⡿⠛⠁⠀⠀⠀     
```

# SadRat

Um Remote Access Toolkit (RAT) somente de estudo feito em Go.

## Conceitos e Arquitetura

O desenvolvimento desse projeto utiliza Domain Driven Design (DDD), Test Driven Development (TDD), Clean Architecture, e influencias da Arquitetura Hexagonal.

## Requisitos

- Go: 1.26.5
- Ubuntu 24.04

## Como Executar o Projeto

## Testes Unitários

```
go test ./tests/
```

## Testes de Integração

```
go test -v ./tests/integration/
```

## Interface Client

> (!) a interface gráfica ainda não está pronta porque o projeto ainda está em desenvolvimento

```
go run cmd/client/main.go <http://ip-do-server:porta>
```

```
╭──────────────────────────────────────────────────╮╭──────────────────────────────────────────────────╮
│                                                  ││ ⠀⢀⣴⣾⣿⣿⣿⣷⣦⡄⠀⣴⣾⣿⣿⣿⣿⣶⣄⠀⠀                       │
│  mmmm             #  mmmmm           m           ││ ⣰⣿⣿⣿⣿⣿⣿⣿⠋⢠⣾⣿⣿⣿⣿⣿⣿⣿⣿⣧⠀                       │
│ #"   "  mmm    mmm#  #   "#  mmm   mm#mm         ││ ⣿⣿⣿⣿⣿⣿⣿⣿⣿⣶⣌⠛⣿⣿⣿⣿⣿⣿⣿⣿⡆                       │
│ "#mmm  "   #  #" "#  #mmmm" "   #    #           ││ ⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⢁⣼⣿⣿⣿⣿⣿⣿⣿⣿⠁                       │
│     "# m"""#  #   #  #   "m m"""#    #           ││ ⠸⣿⣿⣿⣿⣿⣿⣿⡟⢀⣾⣿⣿⣿⣿⣿⣿⣿⣿⠏⠀                       │
│ "mmm#" "mm"#  "#m##  #    " "mm"#    "mm         ││ ⠀⠙⣿⣿⣿⣿⣿⣿⣄⠻⣿⣿⣿⣿⣿⣿⣿⣿⠏⠀⠀                       │
│                                                  ││ ⠀⠀⠈⠻⣿⣿⣿⣿⣿⣧⡈⢿⣿⣿⣿⣿⡟⠁⠀⠀⠀                       │
│> _                                               ││ ⠀⠀⠀⠀⠈⠻⣿⣿⣿⣿⡇⢸⣿⣿⠟⠉⠀⠀⠀⠀⠀                       │
│                                                  ││ ⠀⠀⠀⠀⠀⠀⠈⠙⢿⡿⠀⡿⠛⠁⠀⠀⠀⠀⠀⠀⠀                       │
│                                                  ││                                                  │
│                                                  ││                                                  │
│                                                  ││[+] SadRat C&C Interface iniciada...              │
│                                                  ││                                                  │
│                                                  ││                                                  │
│                                                  ││                                                  │
│                                                  ││                                                  │
│                                                  ││                                                  │
│                                                  ││                                                  │
│                                                  ││                                                  │
│                                                  ││                                                  │
╰──────────────────────────────────────────────────╯╰──────────────────────────────────────────────────╯

```

# Referencias

- Black Hat Rust - Sylvain Kerkour
- Black Hat Go - Tom Steele