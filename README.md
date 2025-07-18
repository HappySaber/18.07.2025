## Инструкция

### Запуск приложения на вашем устройстве

Должен быть уже установлен Go.

1. Запускаем IDE и создаём новый проект.
2. В консоль прописываем команду:

   ```bash
   git clone https://github.com/HappySaber/18.07.2025
   ```

3. После заходим в рабочий каталог smthToZip

4. Подключаем все зависимости:

   ```bash
   go mod tidy
   ```

5. Запускаем программу:

   ```bash
   go run cmd/tozip/main.go
   ```

Программа запустится по адресу [http://localhost:8080/](http://localhost:8080/).

### Ручки для тестирования через сервисы проверки API

- **POST** [http://localhost:8080/tasks](http://localhost:8080/tasks) - создаст Task  в которую можно добавлять urls для скачивания
- **POST** [http://localhost:8080/tasks/:id/urls](http://localhost:8080/tasks/:id/urls) - создаст .zip и скачает туда картинку или .pdf файл по ссылке которую нужно ввести в raw(пример ввода ссылки 
{
    "urls":["https://i.pinimg.com/originals/b2/dc/9c/b2dc9c2cee44e45672ad6e3994563ac2.jpg"]
}
- **GET** [http://localhost:8080/new/](http://localhost:8080/tasks/:id) - выдаст статус по задаче и ссылку на скачивание архива

