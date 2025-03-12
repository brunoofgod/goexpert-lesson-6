# Leilão API

## Descrição
Este projeto é uma API em Go para gestão de leilões, permitindo a criação de leilões, envio de lances e fechamento automático por tempo.

### **Tecnologias e Pacotes Utilizados**
- **Linguagem**: Go
- **Banco de Dados**: MongoDb
- **Goroutines** para fechamento automático dos leilões
- **Docker** para execução e deploy
- **Testcontainers** `testcontainers-go`

## **Funcionamento**
A API permite criar leilões, registrar lances e acompanhar o status de cada leilão. Os leilões são fechados automaticamente após o tempo definido.

Foi utilizado a lib `testcontainers-go` para fazer o teste utilizando o mongo mockado em um container docker.

## **Pré-requisitos para funcionamento**
- Docker e Docker Compose instalados

## **Instruções de Execução**
### **Rodando com Docker Compose**
1. Clone o repositório:
   ```sh
   git clone https://github.com/brunoofgod/goexpert-lesson-6.git
   cd goexpert-lesson-6
   ```

2. Inicie o projeto com Docker Compose:
   ```sh
   docker-compose up -d
   ```

