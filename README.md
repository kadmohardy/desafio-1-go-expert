# Desafio 1 - Go Expert FullCycle

Projeto client/server responsável por obter cotações USD-BRL

# Executando o projeto
## Passo 1: Executando o servidor

Entrar na pasta /server e executar o comando

```
go run main.go
```

Feito isso, é importante observar nos logs a seguintes mensagens: 
```
2023/10/05 15:56:15 Starting server...
2023/10/05 15:56:15 DB initialized
```

## Passo 2: Executando o cliente
Entrar na pasta /client e executar o comando

```
go run main.go
```
Feito isso, o arquivo `cotacao.txt` é criado no diretório.

## Extras
É possível configurar um banco de dados utilizar o arquivo docker-compose. Para isso, execute o comando:

```
docker-compose up 
```

