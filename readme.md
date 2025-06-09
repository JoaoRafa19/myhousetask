Claro! Abaixo está um **documento simples, direto e descritivo** do projeto, ideal para registrar a visão geral, o propósito, as tecnologias envolvidas e os primeiros passos técnicos já realizados.

---

# 🧹 Projeto: TaskHome – Gerenciador de Tarefas Domésticas Compartilhadas

## 📘 Visão Geral

O **TaskHome** é um aplicativo colaborativo para gerenciamento de tarefas domésticas entre membros de uma família. O objetivo é promover organização e divisão justa das tarefas, de forma clara, acessível e interativa, com sorteio de tarefas rotineiras e possibilidade de redistribuição entre os membros.

---

## 🧾 Funcionalidades principais

* 👨‍👩‍👧‍👦 **Sistema de famílias**

  * Criação de famílias por um usuário
  * Entrada via código ou link de convite
  * Gerenciamento colaborativo de tarefas

* 📋 **Tarefas e rotinas**

  * Criação de tarefas únicas e rotineiras
  * Sorteio diário de tarefas entre os membros
  * Redistribuição de tarefa (somente com aceite do novo membro)
  * Status de tarefas: `pendente`, `em andamento`, `concluída`

* 📅 **Calendário**

  * Visualização de tarefas por dia
  * Edição colaborativa de tarefas não rotineiras
  * Execução automática de rotinas diárias

* 🔔 **Notificações em tempo real**

  * Atualizações de status e responsabilidades via **WebSocket**

---

## 🛠️ Tecnologias Utilizadas

| Tecnologia    | Uso principal                             |
| ------------- | ----------------------------------------- |
| **Go**        | Backend (linguagem principal)             |
| **sqlc**      | Geração de código a partir de SQL         |
| **MySQL**     | Banco de dados relacional                 |
| **WebSocket** | Comunicação em tempo real                 |
| **gRPC**      | Comunicação entre microserviços e Flutter |
| **Templ**     | Frontend SSR opcional em Go               |
| **React**     | Frontend web alternativo (uso leve)       |
| **Flutter**   | Frontend mobile principal                 |

---

## 📂 Estrutura Inicial do Projeto

```
myapp/
├── db/
│   ├── migrations/        # Scripts SQL (schema do banco)
│   ├── query/             # Queries SQL com nomes para o sqlc
│   └── gen/               # Código gerado por sqlc (structs + funções)
├── internal/
│   ├── api/               # Handlers HTTP (REST)
│   ├── ws/                # Lógica de WebSocket
│   └── services/          # Regras de negócio
├── cmd/
│   └── server/            # main.go (servidor principal)
├── sqlc.yaml              # Configuração do sqlc
└── go.mod
```

---

## 📜 Progresso atual

* ✅ Schema SQL modelado com colunas em MAIÚSCULO e tabelas minúsculas
* ✅ Arquivo `sqlc.yaml` configurado para ler `migrations` e `query`
* ✅ Geração de entidades com `sqlc generate` funcionando
* ✅ Configuração para tratar `NULL` como string vazia
* ✅ Servidor Go inicializado com `chi`

---

## 📈 Próximos passos sugeridos

1. [ ] Criar rotas REST para CRUD de usuários, famílias e tarefas
2. [ ] Implementar sorteio automático de tarefas rotineiras
3. [ ] Criar canal WebSocket para atualização em tempo real
4. [ ] Definir padrão de autenticação (token JWT ou similar)
5. [ ] Iniciar o frontend em Flutter (ou Templ para web)
6. [ ] Escrever testes e cobertura básica de endpoints

---

## ✍️ Observações

* O projeto é orientado a evoluir com múltiplos frontends: **Flutter** no mobile e **Templ/React** na web.
* A lógica de sorteio e redistribuição de tarefas será feita no backend e registrada no banco de forma audível.
* O sistema prioriza **clareza**, **transparência** e **colaboração familiar**.

---

Se quiser, posso salvar isso como um PDF ou Markdown para você versionar no projeto. Deseja exportar em algum formato?


### Gerar arquivo go a partir de arquivo .proto
```shell
protoc --go_out=. --go_opt=paths=source_relative \
         --go-grpc_out=. --go-grpc_opt=paths=source_relative \
         proto/category.proto
```


### Criar ou executar o banco mysql no terminal via docker
```shell
docker run -d \
  --name mysql-myhousetask \
  -p 3308:3306 \
  -e MYSQL_DATABASE=myhousetask \
  -e MYSQL_USER=user \
  -e MYSQL_PASSWORD=root \
  -e MYSQL_ROOT_PASSWORD=root \
  mysql:latest
```