# user-balance-api
Микросервис для работы с балансом пользователейAPI методы

## Оглавление

1. [Запуск приложения](#run)
2. [API методы](#api)
3. [База данных](#database)

<a name="run">Запуск приложения</a>
---
```
$ git clone https://github.com/ArkovKonstantin/user-balance-api.git
$ cd user-balance-api
$ make test
$ make run
```
Приложение будет доступно по адресу http://localhost:9000

<a name="api">API методы</a>
---

### Метод начисления средств на баланс
_Request_
```shell script
curl --request POST \
  --url http://localhost:9000/balance/add \
  --header 'content-type: application/json' \
  --data '{
	"user_id": 1,
	"amount": 10
}'
```
_Responses:_
* `HTTP/1.1 200 OK`
```json
{
  "user_id": 1,
  "balance": 10
}
```
### Метод списания средств с баланса
_Request_
```shell script
curl --request POST \
  --url http://localhost:9000/balance/sub \
  --header 'content-type: application/json' \
  --data '{
	"user_id": 1,
	"amount": 5
}'
```

_Responses:_
* `HTTP/1.1 200 OK`
```json
{
  "user_id": 1,
  "balance": 5
}
```

### Метод перевода средств от пользователя к пользователю 
_Request_
```shell script
curl --request POST \
  --url http://localhost:9000/transfer \
  --header 'content-type: application/json' \
  --data '{
	"sender_id": 1,
	"recipient_id": 2,
	"amount": 5
}'
```

_Responses:_
* `HTTP/1.1 200 OK`
```json
{
  "sender": {
    "user_id": 1,
    "balance": 0
  },
  "recipient": {
    "user_id": 2,
    "balance": 5
  }
}
```

### Метод получения текущего баланса пользователя
_Request_
```shell script
curl --request GET \
  --url 'http://localhost:9000/balance/get?user_id=1'
```
_Responses:_
* `HTTP/1.1 200 OK`
```json
{
  "user_id": 1,
  "balance": 0
}
```

<a name="database">База данных</a>
---
### Схема
![Image](https://github.com/ArkovKonstantin/user-balance-api/raw/master/assets/schema.png) <br>

### Сущности

**User** <br>
Пользователь приложения
* **id** - уникальный идентификатор пользователя
* **name** - имя пользователя

**Account** <br>
Счет пользователя
* **id** - идентификатор счета
* **user_id** - идентификатор пользователя, которому принадлежит счет (между таблицей `User` и `Account` установлена связь 1:1)
* **balance** - баланс отображает кол-во денег на счету у пользователя (в рублях)

**Opeartion** <br>
Совершенные операции, касающиеся изменения баланса пользователя
* **id** - уникальный идентификатор операции
* **member** - идентификатор пользователя, чей баланс был изменен
* **created_at** - время совершенной операции
* **meta** - мета информация об операции (ex: тип, назначение(для перевода), описание)

### Листинг скрипта, инициализирующего базу данных
```sql
create table "user"
(
    id   serial primary key,
    name varchar
);

create table "account"
(
    id      serial,
    user_id int primary key references "user" (id),
    balance int
        constraint positive_balance CHECK (balance >= 0)
);

create table "operation"
(
    id         serial primary key,
    member     int references "user" (id),
    created_at timestamp default now(),
    meta       json
);

-- create users
insert into "user" (name)
values ('kolya'),
       ('petya'),
       ('vasya'),
       ('tanya');
``` 