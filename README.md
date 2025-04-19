[![Coverage Status](https://coveralls.io/repos/github/WhaleShip/avito-shop/badge.svg?branch=main)](https://coveralls.io/github/WhaleShip/avito-shop?branch=main) 
[![Quality](https://github.com/WhaleShip/avito-shop/actions/workflows/main.yaml/badge.svg)](https://github.com/WhaleShip/avito-shop/actions/workflows/main.yaml)
# Тестовое задание Avito
> текст задания можно найти [щёлкнув тут](https://github.com/avito-tech/tech-internship/blob/main/Tech%20Internships/Backend/Backend-trainee-assignment-winter-2025/Backend-trainee-assignment-winter-2025.md) 

## Краткая сводка

- ### для реализации использовался go 1.23.4 с fiber v2.
- ### Для регистрации используется JWT токен через алгоритм HS256
- ### в качестве бд используется postgresql с pgBouncer
- ### unit тесты в папках с кодом, покрытие 54.9%
- ### интегарционныые тесты реализованы для всех ключевых сценариев, находятся в папке tests/integreation_tests 
- ### результат нагрузочного теста в [отдельной папке](tests/stress_test) в REAMDE и csv формате
- ### [конфигурация линтера](.golangci.yaml)


> [!IMPORTANT]  
> Дисклеймер: Разработка велась под Linux, у некоторых команд могут возникать проблемы с запуском на Windows, написал решение для всех изввестных мне проблем, но не могу быть уверен.


## Запуск

### 1. Создать .env
> на windows работает только из gitbash

```sh
make env 
```

или переименовать [examplse.env](example.env) в .env


### 2. Запустить через докер
```sh
make run
```

или

```sh
docker compose up
```

у виндовс могут быть проблемы со скриптом баунсера, если такое происходит
```sh
dos2unix docker/scripts/entrypoint.sh
```

## Остальной функционал
> если проблемы с использованием makefile, все соответствующие команды можно найти в нём же
### unit тесты

- Запустить тесты
```sh
make test
```

- Посмотреть покрытие
```sh
make cover // через консоль
make cover-html // через html файл
```


### Интеграционные тесты
> для выполнения нужен докер, так что запускать нужно со среды где есть docker (из .devcontainer не получится) <br>
```sh
make test-int
```


### Линтер

- Установка
```sh
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.5
```

- Запуск
```sh
make lint
```


### Результаты нагрузочного тестирования в [отдельной папке](tests/stress_test)


### связаться со мной можно через [telegram](https://t.me/PanHater)
