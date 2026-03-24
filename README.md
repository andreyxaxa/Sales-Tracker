# Sales Tracker - сервис аналитики продаж
Cервис оперирует категориями, продуктами и продажами. Реализован CRUD для этих сущностей и вариативная аналитика продаж.

[Старт]()

## Обзор

- UI - http://localhost:8080/v1
- Документация API - Swagger - http://localhost:8080/swagger

## API

### Category

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| POST | /v1/category | Создать новую категорию |
| GET | /v1/category/{id} | Получить категорию по ID |
| DELETE | /v1/category/{id} | Удалить категорию по ID |

### Product

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| POST | /v1/product | Создать новый продукт |
| GET | /v1/product/{id} | Получить продукт по ID |
| Patch | /v1/product/{id} | Обновить цену продукта по ID |
| DELETE | /v1/product/{id} | Удалить продукт по ID |

### Sale

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| POST | /v1/sale | Создать запись о продаже |
| GET | /v1/sale/{id} | Получить запись о продаже по ID |
| DELETE | /v1/sale/{id} | Удалить записть о продаже по ID |

### Analytics

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| GET | /v1/analytics/revenue-sum | Сумма дохода за период |
| GET | /v1/analytics/revenue-sum | Средний доход за период |
| GET | /v1/analytics/sales-count | Количество продаж за период |
| GET | /v1/analytics/median | Медиана за период |
| GET | /v1/analytics/p90 | 90-й перцентиль за период |

Query-параметры для запросов аналитики:
- `from` - (обязательный): начало диапазона дат(RFC3339);
- `to` - (обязательный): конец диапазона дат(RFC3339);
- `product-id` - (опциональный): аналитика по конкретному продукту;
- `category-id` - (опциональный): аналитика по конкретной категории;
- `payment-method` - (опциональный): аналитика по конкретному способу оплаты;
- `group-by` - (опциональный): группировка по полю(`day`, `week`, `month`, `category`, `product`);
- `sort-by` - (опциональный): сортировка по полю(`group`, `value`). `group` - имя результирующего поля, `value` - значение. Например, если мы попросим количество продаж за период, с группировкой по `product`, в `group` будет имя продукта, в `value` количество продаж этого продукта.
- `sort-order` - (опциональный): порядок сортировки(`asc`, `ASC`, `desc`, `DESC`).

## Запуск

1. Клонируйте репозиторий
2. В корне создайте `.env` файл, скопируйте туда содержимое [env.example](https://github.com/andreyxaxa/Sales-Tracker/blob/main/.env.example):
   ```
   cp .env.example .env
   ```
3. Запустите сервис:
   ```
   make compose-up
   ```
8. Перейдите на http://localhost:8080/v1 и пользуйтесь сервисом.
<img width="1670" height="1010" alt="image" src="https://github.com/user-attachments/assets/d26a4a9a-5341-4738-b101-64efc1a208e7" />


- Перейдите на http://localhost:8080/swagger и ознакомьтесь с API, если хотите взаимодействовать с сервисом вручную или из стороннего сервиса.

## Прочие `make` команды
docker compose down -v:
```
make compose-down
```
Зависимости:
```
make deps
```
