# Notes

## Chapter 4

### Run MySQL in Docker

[Отсюда](http://sysengineering.ru/notes/ispolzovanie-docker-dlya-mysql-servera)

1. Скачать пакет с нужной версией MySQL сервера.
Например, чтобы скачать версию 5.7, нужно в параметрах команды указать `mysql-server:5.7`.
Если не указать версию, то будет поставлен параметр `latest` - последняя версия.

```text
docker pull mysql/mysql-server:5.7
```

2. Запустить контейнер с MySQL сервером, указав в параметрах имя, пароль для root, язык и
пакет с нужным образом.

```text
docker run --name mysql57 -e MYSQL_ROOT_PASSWORD=passw -e LANG=C.UTF-8 -d mysql/mysql-server:5.7
docker logs mysql57
```

3. Подключиться к серверу баз данных и посмотреть имеющиеся базы данных.

```text
docker exec -it mysql57 mysql -uroot -p
SHOW DATABASES;
exit
```

или так:

```text
docker exec -it mysql57 mysql -u root -p
SHOW DATABASES;
exit
```

4. Импортировать свою базу данных (например тестовую) и выполненить нескольких запросов.

*(Примечание: не проверял)*

```text
docker exec -i mysql57 mysql -uroot -ppassw < C:\Temp\MySQL_SampleDB.sql
docker exec -it mysql57 mysql -uroot -ppassw SampleDB
SHOW TABLES;
SELECT * FROM product;
exit
```

5.

- Посмотреть запущенные контейнеры
- Посмотреть все контейнеры

```text
docker container ls
docker container ls -a
```

- Остановить/запустить/удалить контейнеры по его ID или имени

```text
docker stop 61cdba01396f
docker container stop mysql57
docker start mysql57
docker container rm 61cdba01396f
```

6. Заново запустить именованный контейнер

```text
docker start mysql57
```

7 .Чтобы создать контейнер с возможностью подключаться с компьютера-хоста к базе данных
в контейнере, необходимо пробросить порт из контейнера в хост.

A так же создать пользователя в базе данных, которому разрешено подключаться с клиентов.
*(Примечание: сделал это, в дополнение к созданию пользователя по книге).*

Создание контейнера с последней версией MySQL сервера, с поддержкой русского языка,
с заданием пароля администратора и пробросом порта 3306:

```text
docker run --name mysql_instance -e MYSQL_ROOT_PASSWORD=passw -e LANG=C.UTF-8 -p 3306:3306 -d mysql/mysql-server:latest
```

Далее две команды, которые я не делал:

- Импортирование тестовой базы данных, конфигурация которой сохранена в файле `MySQL_SampleDB.sql`
в папке `C:\Temp`:

```text
docker exec -i mysql_instance mysql -uroot -ppassw < C:\Temp\MySQL_SampleDB.sql
```

- Подключение к MySQL серверу в режиме командной строки непосредственно в контейнере

```text
docker exec -it mysql_instance mysql -uroot -ppassw
```

- Создание нового пользователя `dba` с паролем `dbaPass`,
которому разрешены все действия в базе данных, а так же подключение с внешних клиентов.

*(Примечание: сделал наподобие, см. ниже).*

```text
CREATE USER 'dba'@'localhost' IDENTIFIED BY 'dbaPass';
GRANT ALL PRIVILEGES ON *.* TO 'dba'@'localhost' WITH GRANT OPTION;
CREATE USER 'dba'@'%' IDENTIFIED BY 'dbaPass';
GRANT ALL PRIVILEGES ON *.* TO 'dba'@'%' WITH GRANT OPTION;
FLUSH PRIVILEGES;
exit
```

**ВНИМАНИЕ**: для БД в Docker действительно надо создать дополнительного пользователя
наподобие `'dba'@'%'`, чтобы можно было достучаться "снаружи".

Для примера из книги я дополнительно ввел:

```text
docker exec -it mysqllocal mysql -u root -p

create user 'web'@'%';
grant select, insert, update, delete on snippetbox.* to 'web'@'%';
alter user 'web'@'%' identified by 'pass';
flush privilegies;
exit;
```

### Installing a database driver

Install the last version with the major release number 1:

```text
go get github.com/go-sql-driver/mysql@v1
```

Example versions with the major release number 1:

- v1.7.1
- v1.8.0
- v1.8.1
...

Install the lastest version:

```text
go get github.com/go-sql-driver/mysql
```

Install the specific version:

```text
go get github.com/go-sql-driver/mysql@v1.0.3
```

### Modules and reproducible builds

Verify that the checksums of the downloaded packages on your machine match the entries in `go.sum`:

```text
go mod verify
```

Download all the dependencies for the project:

```text
go mod download
```

### Upgrading packages

To upgrade to latest available *minor* or *patch* release of a package:

```text
go get -u github.com/foo/bar
```

Upgrade to a specific version:

```text
go get -u github.com/foo/bar@v2.0.0
```

### Removing unused package

Postfix the package path with `@none`:

```text
go get github.com/foo/bar@none
```

Or automatically remove any unused packages from `go.mod` and `go.sum` files:

```text
go mod tidy -v
```

## Create snippet. POST request

```text
curl -iL -X POST http://localhost:4000/snippet/create
```

## Working with transactions

Basic pattern:

```go
type ExampleModel struct {
    DB *sql.DB
}
func (m *ExampleModel) ExampleTransaction() error {
    // Calling the Begin() method on the connection pool creates a new sql.Tx
    // object, which represents the in-progress database transaction.
    tx, err := m.DB.Begin()
    if err != nil {
        return err
    }

    // Defer a call to tx.Rollback() to ensure it is always called before the
    // function returns. If the transaction succeeds it will be already be
    // committed by the time tx.Rollback() is called, making tx.Rollback() a
    // no-op. Otherwise, in the event of an error, tx.Rollback() will rollback
    // the changes before the function returns.
    defer tx.Rollback()

    // Call Exec() on the transaction, passing in your statement and any
    // parameters. It's important to notice that tx.Exec() is called on the
    // transaction object just created, NOT the connection pool. Although we're
    // using tx.Exec() here you can also use tx.Query() and tx.QueryRow() in
    // exactly the same way.
    _, err = tx.Exec("INSERT INTO ...")
    if err != nil {
        return err
    }

    // Carry out another transaction in exactly the same way.
    _, err = tx.Exec("UPDATE ...")
    if err != nil {
        return err
    }

    // If there are no errors, the statements in the transaction can be committed
    // to the database with the tx.Commit() method.
    err = tx.Commit()
        return err
    }
```

You must **always** call either `Rollback()` or `Commit()` before your function
returns.

## Prepared statements

Use of the `DB.Prepare()` method to create our own prepared statement once,
and reuse that instead.

```go
// We need somewhere to store the prepared statement for the lifetime of our
// web application. A neat way is to embed in the model alongside the connection
// pool.
type ExampleModel struct {
    DB *sql.DB
    InsertStmt *sql.Stmt
}

// Create a constructor for the model, in which we set up the prepared
// statement.
func NewExampleModel(db *sql.DB) (*ExampleModel, error) {
    // Use the Prepare method to create a new prepared statement for the
    // current connection pool. This returns a sql.Stmt object which represents
    // the prepared statement.
    insertStmt, err := db.Prepare("INSERT INTO ...")
    if err != nil {
        return nil, err
    }

    // Store it in our ExampleModel object, alongside the connection pool.
    return &ExampleModel{db, insertStmt}, nil
}

// Any methods implemented against the ExampleModel object will have access to
// the prepared statement.
func (m *ExampleModel) Insert(args...) error {
    // Notice how we call Exec directly against the prepared statement, rather
    // than against the connection pool? Prepared statements also support the
    // Query and QueryRow methods.
    _, err := m.InsertStmt.Exec(args...)
        return err
}

// In the web application's main function we will need to initialize a new
// ExampleModel struct using the constructor function.
func main() {
    db, err := sql.Open(...)
    if err != nil {
        errorLog.Fatal(err)
    }
    defer db.Close()

    // Create a new ExampleModel object, which includes the prepared statement.
    exampleModel, err := NewExampleModel(db)
    if err != nil {
        errorLog.Fatal(err)
    }

    // Defer a call to Close() on the prepared statement to ensure that it is
    // properly closed before our main function terminates.
    defer exampleModel.InsertStmt.Close()
}
```
