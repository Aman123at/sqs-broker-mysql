# MySQL-based SQS Broker in Go

This Go application implements a basic message queue system using MySQL. It includes functions to consume messages from the queue and update their status after processing.


## Global Variables
```go
var dbConn *sql.DB
var wg sync.WaitGroup
```

- `dbConn`: Holds the database connection instance.
- `wg`: A `WaitGroup` to manage concurrent consumers.



## Consumer

The `consumer` function:
1. Starts a database transaction.
2. Queries the database for messages with a status of `todo`, locking the row to avoid other consumers from processing the same message.
3. Processes the message and logs it.
4. Updates the message status to `done`.
5. Commits the transaction.
6. Calls `wg.Done()` to signal that the consumer has finished its work.



## Database Table
The application interacts with a table `sbroker`. The schema might look something like:
```sql
CREATE TABLE sbroker (
    id INT AUTO_INCREMENT PRIMARY KEY,
    message TEXT,
    status ENUM('todo', 'done')
);
```


## How it Works

1. Consumers fetch messages with a `todo` status, process them, and update their status to `done`.
2. The `FOR UPDATE SKIP LOCKED` clause ensures that no two consumers process the same message concurrently.





