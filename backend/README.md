# Anisite Go Backend

## Содержание

- [Требования](#требования)
- [Установка и настройка проекта](#установка-и-настройка-проекта)
  - [1. Клонирование репозитория](#1-клонирование-репозитория)
    - [Настройка SSH-ключей (если еще не настроены)](#настройка-ssh-ключей-если-еще-не-настроены)
  - [2. Установка зависимостей и запуск проекта](#2-установка-зависимостей-и-запуск-проекта)
    - [Linux и MacOS:](#linux-и-macos)
    - [Windows (PowerShell):](#windows-powershell)
    - [Windows (CMD):](#windows-cmd)
  - [3. Очистка сборки](#3-очистка-сборки)

## Требования

Перед тем как начать работать с проектом, убедитесь, что на вашем компьютере установлены следующие зависимости:

- **Go** (1.16+): [Установка Go](https://golang.org/doc/install)
- **Go-Swagger**: инструмент для генерации моделей на основе OpenAPI (Swagger).
  Установка:
  ```bash
  go get -u github.com/go-swagger/go-swagger/cmd/swagger
  ```
- **Make**: утилита для автоматизации задач.
  Установка:
  - [Прямая ссылка](https://gnuwin32.sourceforge.net/packages/make.htm)
  - Используя [Chocolatey](https://chocolatey.org/install):
    ```bash
    choco install make
    ```

## Установка и настройка проекта

### 1. Клонирование репозитория

Для работы с репозиторием в этом проекте используется клонирование с помощью SSH. Убедитесь, что у вас настроены
SSH-ключи в вашей учетной записи на GitHub.

#### Настройка SSH-ключей (если еще не настроены)

1. Сгенерируйте новый SSH-ключ (если его у вас еще нет)

   Откройте терминал и выполните команду, заменив `your_email@example.com` на свою почту, привязанную к GitHub:

   ```bash
   ssh-keygen -t ed25519 -C "your_email@example.com"
   ```

   Если вы используете старую систему, которая не поддерживает `ed25519`, используйте RSA:

   ```bash
   ssh-keygen -t rsa -b 4096 -C "your_email@example.com"
   ```

2. Добавьте ваш SSH-ключ в SSH-агент

**Linux**:

```bash
eval "$(ssh-agent -s)"
ssh-add ~/.ssh/id_ed25519  # Или id_rsa, если вы использовали RSA
```

**Windows**:

- Для PowerShell:
  ```shell
  Start-Service ssh-agent
  ssh-add ~\.ssh\id_ed25519  # Или id_rsa, если вы использовали RSA
  ```

- Для CMD:
  ```cmd
  ssh-agent
  ssh-add %USERPROFILE%\.ssh\id_ed25519  # Или id_rsa, если вы использовали RSA
  ```

### 2. Установка зависимостей и запуск проекта

1. Сначала устанавливаются необходимые зависимости Go
2. Потом генерируются необходимые модели из OpenAPI спецификации
3. Билдится и запускается проект

#### Linux и MacOS:

```bash
go mod tidy
make generate-models
make br
```

#### Windows (PowerShell):

```shell
go mod tidy
make generate-models
make br
```

#### Windows (CMD):

```cmd
go mod tidy
make generate-models
make br
```

### 3. Очистка сборки

Для очистки созданных файлов и бинарников используйте команду:

```bash
make clean
```

Эта команда удалит директорию `bin`, содержащую собранные файлы.
