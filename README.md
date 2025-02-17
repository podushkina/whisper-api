
# Whisper API

Whisper API — это self-hosted RESTful API для интеграции в проекты, требующие распознавания аудио. Он преобразует аудиофайлы в текст и обеспечивает аутентификацию пользователей по токену. С Whisper API вы можете создавать, получать и удалять задания на распознавание аудио, а также отслеживать состояние API и статистику использования для каждого токена. Поскольку это self-hosted решение, вы сохраняете полный контроль над данными и обеспечиваете их конфиденциальность.

- [Пример использования](#пример-использования)
- [Технологии](#технологии)
- [Установка](#установка)
- [Конфигурация](#конфигурация)
- [API Документация](#api-документация)
- [Будущие улучшения](#будущие-улучшения)
- [Лицензия](#лицензия)


## Пример использования
### Видео демонстрация
В этом видео представлен пример использования API. Показан созданный [фронтенд](https://github.com/podushkina/whisper-frontend), который демонстрирует возможности и функционал проекта. Обратите внимание, что по умолчанию звук в плеере отключен, поэтому не забудьте включить его, чтобы послушать аудио.
### Описание видео
В видео показано:

1. **0:00 - 0:10 Получение токена аутентификации**:
    - Получение токена для отправки запросов.

2. **0:20 - 1:18 Сравнение аудио и результата транскрипции на русском языке**:
    - Создание нового задания на транскрипцию аудиофайла на русском языке.
    - Просмотр результата транскрипции и сопоставление аудио с полученным текстом.


https://github.com/user-attachments/assets/ed83acfb-1c76-4c04-bbbc-788c3190b4c8

[ссылка на полную версию видео](https://drive.google.com/file/d/1esnCw6sJTidS5Fi2Nk8G8ezcUFiKDycR/view)


### Дополнительные способы использования
Этот пример — лишь самый простой и наглядный, но Whisper API можно использовать гораздо шире. Вот несколько конкретных примеров проектов, которые можно реализовать:

#### Реальные примеры из моей жизни

1. **Подкасты и видео**: По ссылке на YouTube автоматически скачивается видео, которое затем через мое API преобразуется из аудиоконтента подкастов и видеороликов в текст. Это полезно в те моменты, когда нельзя посмотреть видео, а субтитры YouTube читать неудобно. Сплошной текст видео намного лучше.
2. **Запись лекций**: Использую API для распознавания записей лекций в университете. Во время лекции делаю фотографии, и в итоге получается удобный конспект (если удалось записать нормальный звук лекции).
3. **Транскрибация видео на иностранном языке**: Часто на YouTube нет субтитров для видео на иностранном языке. Используя API, можно получить файл формата VTT и легко добавить субтитры.

#### Возможные варианты использования бизнесом

4. **Протоколирование встреч**: Если нужно вести записи деловых встреч или интервью приема на работу, можно разработать микросервис, который с помощью API автоматически преобразует аудиозаписи в текст. Затем можно использовать локальные LLM размером 7-8b параметров, чтобы суммаризировать эти записи.
5. **Анализ звонков в колл-центре**: Можно создать сервис, который транскрибирует записи разговоров с клиентами для анализа качества обслуживания. Это поможет выявить частые проблемы, когда сотрудник говорит не по скрипту, повышая удовлетворенность клиентов. В целом можно выявить отдельный класс вопросов, которые можно добавить на сайт в FAQ или включить в роботизированную версию ответов по телефону.

## Технологии
### Go
- [Go](https://golang.org/) : Используется для разработки основного приложения.
    - **Веб-фреймворк**: Используется [chi](https://github.com/go-chi/chi) для маршрутизации HTTP-запросов.
    - **JSON Web Tokens (JWT)**: Для аутентификации и авторизации пользователей с помощью [jwt-go](https://github.com/dgrijalva/jwt-go).
    - **Redis клиент**: Для взаимодействия с Redis используется [go-redis](https://github.com/go-redis/redis).

### Redis
- [Redis](https://redis.io/) : Используется для хранения данных и управления сессиями.

### Docker
- [Docker](https://www.docker.com/) : Используется для контейнеризации приложения, что облегчает его деплой и масштабирование. Docker позволяет легко поднять API и взаимодействовать с ним.

### OpenAI Whisper
- [OpenAI Whisper](https://github.com/openai/whisper) : Библиотека для транскрипции аудиофайлов, обеспечивает высокое качество распознавания речи. Whisper может:
    - Транскрибировать аудиофайлы в текст.
    - Поддерживает различные языки.
    - Настраивается с помощью параметров, таких как модель, язык, формат вывода, температура и количество потоков.

### Python
- [Python](https://www.python.org/) : Используется для упрощения установки OpenAI Whisper. Также, используя Python в Dockerfile, можно настроить, какая модель будет загружена.


## Установка

1. Клонируйте репозиторий:
   ```sh
   git clone https://github.com/podushkina/whisper-api.git
   cd whisper-api
   ```

2. Скопируйте `.env.example` в `.env` в корневой директории проекта
   ```sh
   cp .env.example .env
   ```
   данный раздел объяснен в [Конфигурация](#конфигурация).


3. Инициализируйте Go модуль:
   ```sh
   go mod init 
   ```

4. Установите зависимости Go:
   ```sh
   go mod tidy
   ```

5. Если у вас не установлен Docker и Docker Compose, выполните следующие команды для их установки:
   ```sh
    curl -sSL https://get.docker.com | sh
    sudo usermod -aG docker $(whoami)
    sudo apt update
    sudo apt install docker-compose
   ```

6. Запустите Docker Compose для сборки и запуска контейнеров:
   ```sh
   docker-compose up --build
   ```

После успешного запуска вы увидите логи от API и Redis в консоли. Приложение будет доступно по адресу [http://localhost:8080](http://localhost:8080).

7. Чтобы остановите контейнеры, напишите:
   ```sh
   docker-compose down
   ```
   
## Конфигурация

Настройки сервера задаются через переменные окружения:

- `SERVER_PORT`: Порт для сервера API (по умолчанию: `8080`)
- `REDIS_ADDR`: Адрес Redis сервера
- `REDIS_PASSWORD`: Пароль для Redis сервера
- `REDIS_DB`: База данных Redis
- `JWT_SECRET`: Секретный ключ для JWT
- `WHISPER_PATH`: Путь к бинарному файлу Whisper

Создайте файл `.env` в корне проекта и добавьте в него ваши настройки:

```env
SERVER_PORT=8080
REDIS_ADDR=redis:6379
REDIS_PASSWORD=
REDIS_DB=0
JWT_SECRET=your_very_secret_jwt_key_here
WHISPER_PATH=/usr/local/bin/whisper
```

## API Документация

API документировано с использованием OpenAPI (Swagger). Описание API находится в файле `swagger.yaml`. Вы можете просмотреть его, перейдя по следующей ссылке: [swagger.yaml](https://github.com/podushkina/whisper-api/blob/main/api/swagger.yaml%20)

<img width="743" alt="image" src="https://github.com/user-attachments/assets/5f6238b1-8504-4f77-a4bd-ae6dfb5bbc7a">


### Примеры запросов

- Получение токена:
  ```sh
  curl -X POST "http://localhost:8080/auth/token" -H "Content-Type: application/json" -d '{"api_key": "your-api-key"}'
  ```

- Создание задания на транскрипцию:
  ```sh
  curl -X POST "http://localhost:8080/transcribe" -H "Authorization: Bearer <token>" -F "audio=@/path/to/audio.wav" -F "language=en" -F "output_format=json" -F "model=base" -F "temperature=0.7" -F "threads=4"
  ```

- Получение списка заданий:
  ```sh
  curl -X GET "http://localhost:8080/transcribe" -H "Authorization: Bearer <token>"
  ```

## Будущие улучшения

- Реализовать поддержку Cross-Origin Resource Sharing (CORS), чтобы API могло быть доступно с разных доменов.
- Увеличить количество параметров, поддерживаемых Whisper, для более настраиваемых вариантов транскрипции.
- Внедрить возможности многопоточности для обработки нескольких задач транскрипции одновременно.
- Добавить больше функциональности в API, включая:
  - Реализовать поддержку загрузки больших аудиофайлов путем их разделения на части и последовательной обработки.
  - Включить возможность транскрипции в реальном времени для потокового аудио.
  
## Лицензия

Этот проект лицензирован под лицензией MIT. См. файл [LICENSE](LICENSE) для получения дополнительной информации.
