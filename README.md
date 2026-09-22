# Codebank

Servico de pagamentos da plataforma Codebank. A aplicacao recebe pagamentos por gRPC, consulta o cartao de credito no PostgreSQL, registra transacoes e publica eventos no Kafka.

## Tecnologias

- Go 1.21 no ambiente Docker
- gRPC e Protocol Buffers
- PostgreSQL
- Apache Kafka
- `librdkafka` para o produtor Kafka

## Pre-requisitos

Para executar com Docker:

- Docker e Docker Compose
- PostgreSQL acessivel pela aplicacao
- Kafka acessivel pela aplicacao
- A rede Docker externa `store-api_default`

Crie a rede uma vez, caso ela ainda nao exista:

```bash
docker network create store-api_default
```

O PostgreSQL e o Kafka sao compartilhados com os demais projetos do repositorio. Os containers desses servicos devem estar em uma rede que permita a comunicacao com o `codebank`.

## Configuracao

Copie o arquivo de exemplo para `.env`:

```bash
cp .env.example .env
```

No Windows PowerShell, use:

```powershell
Copy-Item .env.example .env
```

Variaveis esperadas:

| Variavel | Exemplo | Descricao |
| --- | --- | --- |
| `host` | `db` | Host do PostgreSQL |
| `port` | `5432` | Porta do PostgreSQL |
| `user` | `postgres` | Usuario do PostgreSQL |
| `password` | `root` | Senha do PostgreSQL |
| `dbname` | `codebank` | Banco utilizado pela aplicacao |
| `KafkaBootstrapServers` | `host.docker.internal:9094` | Endereco do Kafka |
| `KafkaTransactionsTopic` | `payments` | Topico dos eventos de pagamento |

Nao versionar o arquivo `.env`, pois ele pode conter credenciais.

## Executar com Docker Compose

A imagem instala as ferramentas necessarias e o Compose inicia a aplicacao com `go run main.go`:

```bash
docker compose up --build
```

Para executar em segundo plano:

```bash
docker compose up --build -d
```

Para acompanhar os logs:

```bash
docker compose logs -f app
```

Para parar os containers:

```bash
docker compose down
```

A aplicacao escuta gRPC na porta `50052` dentro do container.

## Migration

Ao iniciar, a aplicacao:

1. abre a conexao com o PostgreSQL;
2. verifica a conectividade com `Ping`;
3. cria a tabela `schema_migrations`, se necessario;
4. executa `db.sql` dentro de uma transacao;
5. registra a migration aplicada.

A migration e idempotente: ela nao e executada novamente depois de registrada. O arquivo [`db.sql`](db.sql) cria as tabelas `credit_cards` e `transactions`.

## API gRPC

O contrato esta em [`infrastructure/grpc/protofile/payment.proto`](infrastructure/grpc/protofile/payment.proto).

Servico:

```text
PaymentService.Payment
```

Mensagem enviada:

- dados do cartao de credito;
- valor da transacao;
- loja;
- descricao opcional.

O servidor registra a reflexao gRPC, portanto clientes como Evans podem consultar os servicos publicados. Exemplo:

```bash
evans --host localhost --port 50052 -r list
```

Para regenerar os arquivos Go a partir do `.proto`:

```bash
make gen
```

## Executar localmente

Com Go, PostgreSQL e Kafka instalados ou acessiveis localmente:

```bash
go mod download
go run main.go
```

O comando deve ser executado na pasta `codebank`, onde ficam o `.env` e o `db.sql`.

## Estrutura principal

```text
codebank/
|- domain/                    Entidades de dominio
|- dto/                       Objetos de transferencia
|- infrastructure/grpc/       Servidor, servico e arquivos protobuf
|- infrastructure/kafka/      Produtor Kafka
|- infrastructure/repository/ Persistencia PostgreSQL
|- migrations/                Executor das migrations
|- usecase/                   Regras de processamento de pagamentos
|- db.sql                     Schema inicial do banco
|- docker-compose.yaml        Ambiente Docker de desenvolvimento
|- main.go                    Ponto de entrada da aplicacao
```

## Testes

Execute os testes dentro do container ou em um ambiente com Go configurado:

```bash
go test ./...
```
