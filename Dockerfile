FROM golang:latest

WORKDIR /app

COPY . .

RUN go build -o task.exe

CMD ["./task.exe", "test_file.txt"]
