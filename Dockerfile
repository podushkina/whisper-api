FROM golang:1.22-alpine AS go-builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /whisper-api ./cmd/api

FROM python:3.9-slim

RUN apt-get update && apt-get install -y git ffmpeg && rm -rf /var/lib/apt/lists/*

RUN pip install --upgrade pip
RUN pip install torch torchvision torchaudio --extra-index-url https://download.pytorch.org/whl/cpu
RUN pip install git+https://github.com/openai/whisper.git

# Скачиваем русскую языковую модель и large модель Whisper
# Для этого нужно минимум 8гб ОЗУ для работы
# RUN python -c "import whisper; whisper.load_model('large')"

COPY --from=go-builder /whisper-api /usr/local/bin/whisper-api

WORKDIR /root/

EXPOSE 8080

CMD ["whisper-api"]
