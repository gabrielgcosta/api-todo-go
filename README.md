# Go Todo API 🚀

Uma API REST simples e eficiente para gerenciamento de tarefas (Todo List), desenvolvida em **Go** (Golang) com banco de dados **PostgreSQL** e empacotada com **Docker**.

Este projeto foi construído utilizando a biblioteca padrão do Go (`net/http`), aproveitando os novos recursos de roteamento nativos do Go 1.22+.

---

## 🛠️ Tecnologias Utilizadas

- **Go 1.26** (utilizando apenas a biblioteca padrão `net/http` para roteamento e servidor)
- **PostgreSQL 15** (banco de dados relacional)
- **lib/pq** (driver nativo de PostgreSQL para Go)
- **Docker & Docker Compose** (para conteinerização e fácil execução local)

---

## 🚀 Como Executar o Projeto

Você não precisa ter o Go ou PostgreSQL instalados na sua máquina local, apenas o **Docker** e o **Docker Compose**.

### Passos para rodar:

1. Clone este repositório:
   ```bash
   git clone https://github.com/gabrielgcosta/api-todo-go.git
   cd api
   ```

2. Inicie os containers com Docker Compose:
   ```bash
   docker compose up --build -d
   ```

A API estará rodando em `http://localhost:8080` e criará automaticamente a tabela de tarefas no banco de dados na primeira inicialização.

---

## 📌 Rotas da API

### 1. Criar uma Tarefa
* **Rota:** `POST /tasks`
* **JSON de Entrada:**
  ```json
  {
    "title": "Estudar estruturas de dados em Go"
  }
  ```
* **Resposta (201 Created):**
  ```json
  {
    "id": 1,
    "title": "Estudar estruturas de dados em Go",
    "finished": false,
    "created_at": "2026-06-08T22:37:58Z"
  }
  ```

### 2. Listar todas as Tarefas
* **Rota:** `GET /tasks`
* **Resposta (200 OK):**
  ```json
  [
    {
      "id": 1,
      "title": "Estudar estruturas de dados em Go",
      "finished": false,
      "created_at": "2026-06-08T22:37:58Z"
    }
  ]
  ```

### 3. Atualizar uma Tarefa
* **Rota:** `PUT /tasks/{id}`
* **JSON de Entrada:**
  ```json
  {
    "title": "Estudar estruturas de dados em Go (Concluído)",
    "finished": true
  }
  ```
* **Resposta (200 OK):** `Task updated successfully`

### 4. Deletar uma Tarefa
* **Rota:** `DELETE /tasks/{id}`
* **Resposta (204 No Content)**

---

## 📁 Estrutura do Projeto

* `main.go`: Código fonte da aplicação contendo os handlers HTTP e a conexão com a base de dados.
* `Dockerfile`: Configura o container Alpine leve para rodar o binário compilado.
* `docker-compose.yml`: Orquestra o container da aplicação e o container do banco PostgreSQL com volumes persistentes.
* `.gitignore`: Evita o envio de binários locais compilados para o Git.
