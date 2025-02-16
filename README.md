# Merch Shop

Это мое решение тестового задания в Avito на позицию Backend Developer Go.

- [Merch Shop](#merch-shop-avito-trainee-task)
  - [Подход к решению задачи](#подход-к-решению-задачи)
  - [Архитектура сервера](#архитектура-сервера)
  - [Обработка ошибок](#обработка-ошибок)
  - [Проектирование базы данных](#проектирование-базы-данных)
  - [Установка и запуск](#установка-и-запуск)
    - [Запуск в docker](#запуск-в-docker)
    - [Запуск через go run/make](#запуск-через-go-runmake)
  - [Тестирование](#тестирование)
    - [Нагрузочное тестирование](#нагрузочное-тестирование)
      - [Load Test 1](#load-test-1)
      - [Load Test 2](#load-test-2)
      - [Load Test 3](#load-test-3)
      - [Load Test 4](#load-test-4)
    - [Профилирование во время нагрузочных тестов и после](#профилирование-во-время-нагрузочных-тестов-и-после)
    - [Unit-тесты](#unit-тесты)
    - [Интеграционное тестирование](#интеграционное-тестирование)
    - [Ручное тестирование](#ручное-тестирование)
      - [POST /api/auth](#post-apiauth)
      - [POST /api/info](#post-apiinfo)
      - [POST /api/sendCoin](#post-apisendcoin)
      - [GET /api/buy/:item](#get-apibuyitem)
  - [TODO](#todo)
  - [Вопросы и размышления](#вопросы-и-размышления)
    - [Вопросы к API](#вопросы-к-api)
    - [Конфигурация](#конфигурация)
    - [Хранение паролей](#хранение-паролей)
    - [Кэширование](#кэширование)
    - [Индексы](#индексы)
  - [Использованные технологии](#использованные-технологии)

## Подход к решению задачи

Решил не использовать кодогенерацию API, все написал руками.

Для написания выбирал между `echo,gin и chi`, в итоге решил использовать `chi`.

Для работы с базой данных использовал ORM `gorm`.

## Архитектура сервера

Я решил, что архитектура должна быть 3-х уровневая:

1. API - Получает запросы
2. Workers - Обрабатывают запросы
3. libs - Разбивка кода по слоям
   1_domain_methods - Обработчики реализующие основную логику приложения
   2_generated_models - Сгенерированные структуры
   3_infrastructure - Работа с базой данных
   4_common - Содержит вспомогательные компоненты

Запросы поступают в `API`, структурируются и передаются в `Workers`, запрос выполняется и по возвратному каналу возвращает результат пользователю.

Слои `libs` пронумерованы и должны импортировать только нижестоящие слои, чтобы избежать циклических импортов.

Интерфейс `smart_context.ISmartContext` пронизывает приложение и дает возможность взаимодейстовать с логами, бд и тд

## Обработка ошибок

1. Интерфейс `smart_context.ISmartContext` логирует ошибки уровня базы данных и сервисного уровня
2. Обработчики возвращают ошибки и коды ошибок

## Проектирование базы данных

При проектировании я отталкивался от спецификации.

Я выделил, что необходимо создать следующие таблицы:

1. `auth_users`
2. `doc_users`
3. `doc_user_merchs`
4. `doc_merchs`
5. `doc_transactions`

`auth_users` хранит данные для авторизации пользователя
`doc_users` хранит данные пользователя eg. баланс, имя
`doc_user_merchs` хранит покупки пользователей
`doc_merchs` хранит мерч и стоимости
`doc_transactions` хранит транзакции коинов между пользоватями

Ниже приведена диаграмма полученной БД:

![db](/images/db.png)

Использовал UUID чтобы избежать проблем повторяющихся идентификаторов

## Установка и запуск

FYI. Использую ОС Windows, если у вас другая ОС, то некоторые команды могут отличаться.

Изначально необходимо склонировать репозиторий:

```sh
git clone https://github.com/BDVRepo/merch.git
```

### Запуск в docker

В проекте уже лежат `Dockerfile` и `docker-compose.yml`.

Проект можно запускать с помощью команды:

```sh
docker-compose up --build
```

## Тестирование

### Нагрузочное тестирование

Для проведения нагрузочных тестов использовалась утилита `k6`.

Установка на моей системе:

```sh
npm install k6
```

#### Load Test

Был написан файлик `load.test.js`. Его параметры vus 1000, rps 1000, duration: '1m'.
По его сценарию:
1.Пользователи выполняют аутентификацию
2.Пользователи покупают мерч
3.Пользователи передает монету другому пользователю среди тысячи пользователей тысячи
4.Пользователи запрашивают информацию

Так как тест выполняется на 1000 разных пользователей, транзакции на изменение баланса имеют низкий шанс на отмену. Но если отмена тразакции случается то баланс не изменяется и все монеты остаются у пользователей.
**Тестовый файл**: `load_test.js`

_Тест запускается такой командой:_

```sh
k6 run load.test.js
```

_Результаты тестирования:_

```sh
         /\      Grafana   /‾‾/
    /\  /  \     |\  __   /  /
   /  \/    \    | |/ /  /   ‾‾\
  /          \   |   (  |  (‾)  |
 / __________ \  |_|\_\  \_____/

     execution: local
        script: load.test.js
        output: -

     scenarios: (100.00%) 1 scenario, 1000 max VUs, 1m30s max duration (incl. graceful stop):
              * default: 1000 looping VUs for 1m0s (gracefulStop: 30s)


     ✓ Auth success
     ✓ Token exists
     ✓ Buy success
     ✓ SendCoin success
     ✓ Info success

     checks.........................: 100.00% 14360 out of 14360
     data_received..................: 3.4 MB  49 kB/s
     data_sent......................: 3.4 MB  49 kB/s
     http_req_blocked...............: avg=178.5µs  min=0s      med=0s     max=50ms   p(90)=0s       p(95)=571.67µs
     http_req_connecting............: avg=156.22µs min=0s      med=0s     max=50ms   p(90)=0s       p(95)=0s
     http_req_duration..............: avg=5.79s    min=36.99ms med=5.21s  max=23.7s  p(90)=10.67s   p(95)=13.22s
       { expected_response:true }...: avg=5.79s    min=36.99ms med=5.21s  max=23.7s  p(90)=10.67s   p(95)=13.22s
     http_req_failed................: 0.00%   0 out of 11488
     http_req_receiving.............: avg=101.93µs min=0s      med=0s     max=8ms    p(90)=518.45µs p(95)=998.8µs
     http_req_sending...............: avg=51.21µs  min=0s      med=0s     max=25.8ms p(90)=0s       p(95)=24.21µs
     http_req_tls_handshaking.......: avg=0s       min=0s      med=0s     max=0s     p(90)=0s       p(95)=0s
     http_req_waiting...............: avg=5.79s    min=36.99ms med=5.21s  max=23.7s  p(90)=10.67s   p(95)=13.22s
     http_reqs......................: 11488   166.262275/s
     iteration_duration.............: avg=23.42s   min=9.03s   med=24.31s max=47.41s p(90)=30.88s   p(95)=32.23s
     iterations.....................: 2872    41.565569/s
     vus............................: 159     min=159            max=1000
     vus_max........................: 1000    min=1000           max=1000


running (1m09.1s), 0000/1000 VUs, 2872 complete and 0 interrupted iterations
default ✓ [======================================] 1000 VUs  1m0s
```

### Интеграционное тестирование sqlite

Задачей было покрыть основные обработчики тестами основную бизнес логику. Соответственно, тесты писались для `1_domain_methods`, потому что это там и расположены все обработчики. Использовал sqlite бд создавая её в памяти.

_Для проверки процента покрытия тестами выполним поманду:_

```sh
cd golang
go test ./... -cover
```

**Покрытие тестами:**

```sh
        bdv-avito-merch/apps/gen-type           coverage: 0.0% of statements
        bdv-avito-merch/apps/backend-chi                coverage: 0.0% of statements
        bdv-avito-merch/libs/4_common/env_vars          coverage: 0.0% of statements
        bdv-avito-merch/libs/3_infrastructure/db_manager                coverage: 0.0% of statements
        bdv-avito-merch/libs/1_domain_methods/helpers           coverage: 0.0% of statements
        bdv-avito-merch/libs/4_common/auth              coverage: 0.0% of statements
        bdv-avito-merch/libs/4_common/middleware                coverage: 0.0% of statements
        bdv-avito-merch/libs/2_generated_models/model           coverage: 0.0% of statements
        bdv-avito-merch/libs/4_common/safe_go           coverage: 0.0% of statements
        bdv-avito-merch/libs/4_common/types             coverage: 0.0% of statements
        bdv-avito-merch/libs/4_common/smart_context             coverage: 0.0% of statements
ok      bdv-avito-merch/libs/1_domain_methods/handlers  1.170s  coverage: 40.7% of statements
```

Нас интересует строка `ok      bdv-avito-merch/libs/1_domain_methods/handlers  1.170s  coverage: 40.7% of statements`.

В ней сказано, что **покрыто 40.7% сценариев тестами**, что удовлетворяет условию в минимум 40%.

Таким образом, основные бизнес сценарии были покрыты тестами.

_Запустить тесты можно с помощью :_

```sh
cd golang
go test ./...
```

**Вот результат ее выполнения:**

```sh
PS C:\projects\avito-merch-store> cd golang
PS C:\projects\avito-merch-store\golang> go test ./...
...
ok      bdv-avito-merch/libs/1_domain_methods/handlers  1.097s
```

Все тесты успешны, ошибок не было.

### Ручное тестирование

Для ручного тестирования использовалась программа `Postman`.

#### POST /api/auth

Попробуем аутентифицироваться в первый раз в системе:

![Auth1](/images/auth1.jpg)

Получаем в ответ токен, который будем использовать в дальнейших запросах.

Также аутенцифицируем второго пользователя, его мы будем использовать для передачи монет

![Auth2](/images/auth2.jpg)

Получаем в ответ токен, который будем использовать для передачи монет.

#### POST /api/info

Получим информацию о себе.

![Info1](/images/info1.jpg)

Как мы видим, у `test_user1` имеется 1000 монет, пустой инвентарь и не операций передачи монет и получения монет.

#### POST /api/sendCoin

Отправим 100 монет другому юзеру, а именно юзеру `test_user2`.

![send1](/images/send1.jpg)

Затем поменяв токен на токен юзера `test_user2` передадим 10 монет юзеру `test_user1`

![send2](/images/send2.jpg)

И после этого запросим информацию с юзера `test_user1` и `test_user2`.

![info2](/images/info2.jpg)

![info3](/images/info3.jpg)

Как мы видим, количество и история передачи монет отлично записались.

#### GET /api/buy/:item

Купим несколько предметов, например, 1шт - `hoody`, 2шт - `book` и 3шт - `pen`.

![buy1](/images/buy1.jpg)

![buy2](/images/buy2.jpg)

![buy3](/images/buy3.jpg)

Мы купили 1 раз худи, 2 раза книгу и 3 раза ручку, теперь запросим информацию о себе.

![info4](/images/info4.jpg)

Окей, подведем итог, было 1000 монет, 100 мы отправили test_user2 и 10 получили от test_user2 в ответ = на нашем балансе отсалось 910.
Мы совершили покупки, а именно:
худи=300,
книги 2*50=100,
ручки 3*10= 30
подсчитаем 910-430 = 480
На балансе осталось 480 монет как и ожидалось.
Замечательно! Ручное тестирование выполнено успешно!

## Использованные технологии

1. [gorm](https://github.com/go-gorm/gorm)
2. [jwt](https://github.com/golang-jwt/jwt)
3. [chi](https://github.com/go-chi/chi)
4. [godotenv](https://github.com/joho/godotenv)
5. [zap](https://go.uber.org/zap)
