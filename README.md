## Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` — адрес вашего репозитория на GitHub без префикса `https://`) для создания модуля.

## Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m main template https://github.com/Yandex-Practicum/go-musthave-shortener-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/main .github
```

Затем добавьте полученные изменения в свой репозиторий.

## Запуск автотестов

Для успешного запуска автотестов называйте ветки `iter<number>`, где `<number>` — порядковый номер инкремента. Например, в ветке с названием `iter4` запустятся автотесты для инкрементов с первого по четвёртый.

При мёрже ветки с инкрементом в основную ветку `main` будут запускаться все автотесты.

Подробнее про локальный и автоматический запуск читайте в [README автотестов](https://github.com/Yandex-Practicum/go-autotests).


## Паспорт сервиса:

* cmd/shortener/main.go - rest-сервер на Golang (API методы)
* cmd/e2e/shortener_test.go - сквозной тест

## Запуск проекта:

1. В Goland Add Configuration -> go build
2. Run kind = Directory; Directory = к значению, что ide прописало автоматически, надо добавить ```/cmd/shortener```
3. ENVIRONMENT скопировать из ```.env.server-example```


## Запуск e2e теста:

1. В Goland Add Configuration -> go build
2. Run kind = Directory; Directory = к значению, что ide прописало автоматически, надо добавить ```/cmd/e2e```
3. ENVIRONMENT скопировать из ```.env.e2e-example```