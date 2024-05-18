## Инструкции по запуску

### 1. Склонируйте репозиторий:
`git clone https://github.com/asstrahanec/yadrointernship.git`

### 2. Перейдите в каталог с проектом:
`cd ваш_репозиторий`
### 3. Скомпилируйте программу (если это необходимо):
`go build`
### 4. Запустите программу с аргументами командной строки:
`./ваш_исполняемый_файл test_file.txt`
### Где `test_file.txt` - имя вашего тестового файла.

## Инструкции по запуску с использованием Docker

### 1. Склонируйте репозиторий:
`git clone https://github.com/asstrahanec/yadrointernship.git`
### 2. Перейдите в каталог с проектом:
`cd ваш_репозиторий`
### 3. Постройте Docker образ:
`docker build -t project:latest .`
### 4. Запустите Docker контейнер:
`docker run project:latest`