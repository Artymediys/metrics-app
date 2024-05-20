# Сервис метрик

## Описание
Сервис для формирвоания и отправки метрик.  
Бизнес-метрики получаются с БД.

### Есть два вида метрик, предоставляемые сервисом
1) Метрики `Prometheus`, доступные по пути `some.example.ru/metrics-app/metrics`:
  - `policies_purchased_today_gauge` - Количество купленных за сегодня полисов
  - `authentications_today_gauge` - Количество аутентификаций за сегодня

<br>

2) Таблицы с данными, которые отправляются на заданные Email-адреса с почтового адреса `fatal-alert@some.example.ru`:
  - Данные о наличии проблем с сертификатами  
  - Данные об оплаченных, но не выпущенных полисах.

## Сборка и запуск приложения
Проект разрабатывался на версии `Go` - `1.22.1`.  

Приложение необходимо запускать в корневой директории проекта (там где находится `main.go`).  
Запустить приложение можно командой `go run main.go --contour <local/demo/preprod/prod>`.

Также приложение можно собрать в исполняемый файл.  
Это делается командой `go build -o ./metrics-app ./main.go` в случае `Linux/MacOS`
или же `go build -o ./metrics-app.exe ./main.go` в случае `Windows`.  
Приложение соберётся конкретно под ту ОС и архитектуру процессора,
на которой была запущена команда сборки.
Если есть необходимость собрать приложение под другой случай,
тогда надо предварительно установить переменные окружения `GOOS` и `GOARCH`.

Примеры:
- `env GOOS=linux GOARCH=amd64 go build -o metrics-app main.go`
- `env GOOS=darwin GOARCH=arm64 go build -o metrics-app main.go`
- `env GOOS=windows GOARCH=386 go build -o metrics-app.exe main.go`

Полный список доступных вариантов можно узнать командой `go tool dist list`.

## Конфигурация
В директории `./config` должен лежать как минимум один файл с конфигурацией под необходимый контур. 
Название файла оформляется следующим образом: `config.<local/demo/preprod/prod>.json`.   
Пример: `config.demo.json`  

Во время запуска приложения можно указать необходимый контур (как было показано в секции [Сборка и запуск приложения](#сборка-и-запуск-приложения))