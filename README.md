# Go Todo API 🚀

Uma API REST simples e eficiente para gerenciamento de tarefas (Todo List), desenvolvida em **Go** (Golang) com banco de dados **PostgreSQL** e empacotada com **Docker**.

Este projeto foi construído utilizando a biblioteca padrão do Go (`net/http`), aproveitando os novos recursos de roteamento nativos do Go 1.22+.

---

## 🛠️ Tecnologias Utilizadas

- **Go 1.26** (utilizando apenas a biblioteca padrão `net/http` para roteamento e servidor)
- **PostgreSQL 15** (banco de dados relacional)
- **RabbitMQ 3** (message broker para eventos assíncronos)
- **lib/pq** e **amqp091-go** (drivers para banco de dados e mensageria)
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

2. Crie e configure o arquivo `.env` com base no arquivo `.env.example`:
   ```bash
   cp .env.example .env
   ```
   *(Nota: O arquivo `.env` é ignorado pelo Git e contém as variáveis de conexão com banco de dados, RabbitMQ e e-mail de notificação).*

3. Inicie os containers com Docker Compose:
   ```bash
   docker compose up --build -d
   ```

A API estará rodando em `http://localhost:8080` e o RabbitMQ em `http://localhost:15672` (painel de gerenciamento com usuário `guest` e senha `guest`). O banco de dados e a fila de mensagens serão criados automaticamente na inicialização.

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

## 🏗️ Arquitetura & Recursos Recentes

A API adota práticas de arquitetura modular, focada em testabilidade, segurança e observabilidade:
- **Middleware de Log**: Registra automaticamente no stdout detalhes de cada requisição (método, caminho, IP, status code HTTP, duração e causa raiz detalhada de eventuais falhas).
- **Tratamento de Erros Centralizado (`apierror`)**: Retorna erros padronizados em JSON (`{"error": "message"}`). Utiliza `errors.As` para extrair detalhes estruturados e ocultar erros sensíveis de banco de dados do cliente final, registrando o erro raiz detalhado apenas no log interno.
- **Mensageria com RabbitMQ**: Quando uma tarefa sofre alterações no CRUD (criação, edição ou deleção), um evento é publicado na fila `task_events` gerenciada pelo RabbitMQ. Um consumidor escuta esta fila de maneira assíncrona, simula o processamento e "envia" uma notificação de e-mail de teste nos logs do container.

---

## 🧪 Como Executar os Testes

O projeto possui uma suíte completa de testes unitários que cobre a lógica do worker, tratamento de erros e handlers HTTP (com mocks para evitar a dependência com o Postgres ou com o RabbitMQ ativo).

Para rodar os testes e verificar a cobertura do código localmente:
```powershell
go test -cover ./...
```

---

## 📁 Estrutura do Projeto

* `main.go`: Ponto de inicialização do banco, do cliente do RabbitMQ com lógica de retry, injeção de dependências e inicialização do consumidor de e-mail.
* `.env.example`: Modelo contendo a definição das variáveis de ambiente necessárias.
* `apierror/`: Pacote para formatação JSON de erros e compatibilidade com erros envelopados.
* `database/`: Pacote para inicialização e migração automática do esquema do banco de dados.
* `middleware/`: Contém os middlewares da aplicação (como o Logger).
* `rabbitmq/`: Cliente responsável por gerenciar a conexão física com o RabbitMQ e a declaração de filas.
* `task/`: Pacote de domínio contendo a definição da entidade, a interface e implementação de persistência, e os handlers HTTP.
* `worker/`: Contém o produtor (`worker.go`) que envia eventos para a fila e o consumidor (`consumer.go`) que processa mensagens e simula o e-mail.
* `Dockerfile`: Configura o build em multi-estágio da aplicação para gerar um container Alpine leve.
* `docker-compose.yml`: Orquestra o container da aplicação, o banco PostgreSQL e o broker RabbitMQ.
