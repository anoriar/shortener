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

1. В Goland Add Configuration -> go test
2. Run kind = Directory; Directory = к значению, что ide прописало автоматически, надо добавить ```/cmd/e2e```
3. ENVIRONMENT скопировать из ```.env.e2e-example```



## Запуск автотестов локально
Подготовка:
1 Скачать тест и положить в корень проекта
https://github.com/Yandex-Practicum/go-autotests/releases/tag/v0.9.16

2 sudo chmod a+x shortenertestbeta

Запуск:

1 Скомпилировать сервер в папке cmd/shortener
go build -o shortener *.go

В Goland: Edit configurations -> Add -> Go Build
Run kind = Directory
Directory = {your home directory}/shortener/cmd/shortener
Run after build отключить

2 выполнить в корне проекта:
./shortenertestbeta -test.v -test.run=^TestIteration1$ -binary-path=cmd/shortener/shortener

В Goland:
Edit configurations -> Add -> Shell Script
Script text = ./shortenertestbeta -test.v -test.run=^TestIteration1$ -binary-path=cmd/shortener/shortener
Script брать из .github/workflows/shortenertest.yaml и менять для каждой итерации
Пример с переменной окружения SERVER_PORT
SERVER_PORT=$(shuf -i 1024-49151 -n 1); ./shortenertestbeta -test.v -test.run=^TestIteration5$ -binary-path=cmd/shortener/shortener -server-port=$SERVER_PORT

## Как пользоваться линтером Яндекс
1. Скачать бинарник statictest для вашей операционной системы (если у вас apple silicon, дополнительно выполните эту инструкцию)
Ссылка  на скачивание бинарника https://github.com/Yandex-Practicum/go-autotests/releases/tag/v0.9.16
Для ubuntu - statictest
2. Поместить в корень проекта. Дать права 777
3. go vet -vettool=./statictest ./...

## Profiling

Как профилировать:

go tool pprof -http=":9090" -seconds=30 http://127.0.0.1:8081/debug/pprof/heap

За эти 30 секунд создать нагрузку

Смотреть отчет на 9090 порту

С созданием файла:

Создать профиль base.pprof

curl http://127.0.0.1:8081/debug/pprof/heap > ./profiles/base.pprof

Зайти в анализатор через:

go tool pprof -http=":9090" profiles/base.pprof

После оптимизации проделать тоже самое для результируюещего файла result.pprof

Сравнить результаты

pprof -top -diff_base=profiles/base.pprof profiles/result.pprof

## Benchmarks
Создать профили для бенчмарка

В папке cmd/benchmark

go test -bench=. -cpuprofile=cpu.out -memprofile=mem.out

Проанализировать

go tool pprof -http=":9090" benchmark.test mem.out

После оптимизации go test -bench=. -cpuprofile=result-cpu.out -memprofile=result-mem.out

Сравнить 

go tool pprof -top -diff_base=mem.out result-mem.out


# Убрать лишние импорты + gofmt
goimports -local github.com/anoriar/shortener -w ./

## Godoc

godoc -http=:8080

Перейти на 
http://localhost:8080/pkg/github.com/anoriar/shortener/?m=all

# Генерация заглушек godoc (для прохождения 18 теста)

Команда godoc-generate

Либа go install github.com/DimitarPetrov/godoc-generate@latest