`Orignal repo`
https://github.com/Wave-95/boards

`Using fork`
https://github.com/leegodden/boards

`Access container`
docker exec -it boards-db-1 /bin/bash

`Access db, once in the container`
psql -U postgres -d boards

`Rebuild container`
docker compose up --build after updated

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////

SETTING UP CONFIGURATION WITH DOCKER AND POSTGRES

# Notes for config.go

1. Imports:
   The file imports essential packages:

- `fmt` for formatting strings and handling errors.
- `os` to access environment variables.

- `validator/v10` from go-playground, a package used for validating that certain conditions are met in our
  structures (in this case, ensuring that needed environment variables are set).

2. Constants
   These are keys expected in the environment variables, representing database connection parameters. These
   constants are used to retrieve their corresponding values from the system's environment.

3. Validation Function
   The `Validate` method on `DatabaseConfig` utilizes the validator package to check if all struct fields
   annotated with `validate:"required"` are indeed provided. This helps catch any missing configuration early before
   the application tries to use these in establishing a database connection.

4. Configuration Wrapper Structure - `type Config struct`
   This `Config` struct design is meant to encapsulate all the configuration parts of the application. Right now, it
   only holds the database configuration, but this design allows easy expansion to include other configurations in
   the future without changing the handling of existing configurations.

5. Load Function
   Load is a function that acts as the initializer for config structures. It retrieves the database configuration
   via a private helper function `getDatabaseConfig`. If this function encounters any error - such as missing
   environment variables, it will directly return an error, aborting further execution.

   `(*Config, error):`
   This specifies what the function returns. `*Config` is a pointer to a `Config` structure. This means when
   Load finishes, it will return a `reference` to a Config structure, which contains configuration data
   (like database configuration, in this context).

   A pointer is used here for reasons of efficiency and to allow the function to return a `nil` value if no
   valid configuration can be loaded.

   On successful retrieval and validation, it wraps `DatabaseConfig` in `Config` and returns it, ready for use by
   other parts of the application. The `&` symbol before `Config` takes the address of the newly created Config struct,
   resulting in a pointer.

   This pointer is returned so that the calling function can use the configuration data while preserving memory efficiency
   and allowing for potential modifications at a shared location in memory.

6. Helper Function to Retrieve Database Configuration
   The `getDatabaseConfig()` function constructs a `DatabaseConfig` structure by fetching each required field from
   the environment using the corresponding key. After populating the struct, it calls `Validate` to ensure all
   required fields are present and correctly populated. If any are missing or faulty (not meeting the validator's
   requirements), an error is returned.

   This function centralizes the retrieval and validation of database configuration, encapsulating error handling
   within the configuration loading process, thereby simplifying external usage.

Overall, the config.go files handles configuration setup and ensures robust initialization of your application by
ensuring the database configuration is properly set up before any database connection is attempted.

# Step-by-Step Flow:

1. main.go calls `Load()`

   - The application begins in the `main.go` file, where it calls the `Load()` function defined in `config.go`.
   - The purpose is to retrieve and establish the application configuration, focusing particularly on database
     configuration before any interactions with the database occur.

2. Load() function in config.go:

   - Once invoked, Load() itself calls the `getDatabaseConfig()` function, also located within `config.go`.
   - Load() is responsible for managing the overall fetching and validation of the configuration data. It acts as a
     wrapper to any configuration setup steps and errors management.

3. getDatabaseConfig() function in config.go

   - The purpose of `getDatabaseConfig()` is to actually fetch the environment variables that are required for the
     database configuration.

   - It constructs a `DatabaseConfig` struct with all necessary database parameters fetched from the environment
     (like `database` `host`, `port`, `username`, and `password`).
   - After assembling the DatabaseConfig, it then calls the `Validate()` method on this configuration instance.

4. Validate() function in config.go

   - This method is a crucial safety check to ensure that all required fields in the DatabaseConfig are properly set.
   - It uses tags defined in the DatabaseConfig struct such as validate: "required" to verify that each required
     configuration setting is populated.

   - If validation fails (meaning one or more required fields are missing or improperly filled), it returns an
     error, which is then propagated back through to `getDatabaseConfig()` to Load().

5. Returning to main.go:

   - If `Validate()` finds issues, both `getDatabaseConfig()` and `Load()` funnel that error back to main.go,
     and typically this would cause the program to halt or handle the error accordingly.
   - If there are no errors, `getDatabaseConfig()` returns a valid `DatabaseConfig` to our `Load()` function.
   - Finally, Load() packages this DatabaseConfig into a broader Config struct and returns it back to main.go.
   - In main.go, the program can then use this configuration data to proceed with other tasks like connecting
     to a database.

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////

# Notes for config_test.go

This test file is structured to ensure that the `Load()` function in `config.go` behaves as expected under
varying environments. This includes scenarios where all required environment variables are present and correctly set,
and scenarios where some or all of these variables are missing.

1. setDBEnvVars() Function:

   - This function sets environment variables which simulate a correctly configured environment necessary for the `Load()`
     function to extract and assemble a valid `DatabaseConfig`.

   - By using `os.Setenv`, it mimics the conditions under which your application would run in a properly configured
     real-world setup, providing values like `localhost`, `5432`, `dbname`, `user`, `password` for the associated DB
     environment variables.

2. clearDBEnvVars() Function:

   - This function clears the environment variables related to the database configuration by using `os.Unsetenv`.

   - The purpose here is to simulate a scenario where these key configuration values are not set, which should lead
     to errors when trying to load the configuration using `Load()`.

3. Test Cases: `t.Run("valid db env vars",...` & `t.Run("missing db env vars",`

   - `valid db env vars`: After setting the environment variables using `setDBEnvVars()`, it runs `Load()` to check
     if it can successfully load the configuration without any errors. It further asserts that the returned configuration
     is not nil, the database configuration specifically is not nil, and that it matches expected values
     (like `localhost` for the host).

   - `Missing db env_vars`: This test ensures that `Load()` correctly handles the absence of necessary environment variables
     by clearing them first and then attempting to load the configuration. It uses `assert.Error(t, err)` to confirm that
     an error is indeed returned, as expected in such cases.

   - Test Assertions
     The use of `assert` from the `github.com/stretchr/testify/assert` package helps in making the tests concise and readable.
     It provides a range of functions like `assert.NoError`, `assert.NotNil`, and `.assert.Error` to check for conditions such
     as the absence of errors, the non-nullity of objects, and the presence of expected errors, respectively.

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////

# Notes for db.go

The `db.go` file is structured to manage database connections using the Go programming language. It interfaces with PostgreSQL
using the `pgxpool` package from `github.com/jackc/pgx/v5/pgxpool`. This package helps efficiently manage multiple database
connections by using a connection pool.

`type DB struct {...`
DB is a struct that embeds a pointer to a `pgxpool.Pool`. This pool manages and maintains a set of database connections that
can be reused for multiple database queries.

Embedding the pool directly in the `DB` struct allows you to use all the methods of `pgxpool.Pool` on `DB instances` directly,
simplifying the interface for database operations.

`func Connect(cfg config.DatabaseConfig) (*DB, error) {...`
Connect is a function that initializes a connection to the database using configuration data provided via `cfg`, which is of
type `config.DatabaseConfig`.

`buildConnectionURL(cfg)`
constructs the PostgreSQL connection string from the given configuration.

`pgxpool.New(context.Background(), url)`
attempts to create a new database connection pool. If this operation fails (err != nil), the function returns the error,
preventing further execution.

If the connection is successful, it logs a message and returns a pointer to a `DB` instance containing the initialized pool.

`func buildConnectionURL(cfg config.DatabaseForeign) string {...`
This function constructs a PostgreSQL connection URL string using the database configuration provided. The URL includes the
`user`, `password`, `host`, `port`, and `database` name, all extracted from `cfg`.

`sslmode=disable`
This is included in the URL, which specifies that SSL will not be used for this connection. This is typical in development
environments, but for production environments, using SSL is recommended to ensure security.

# User entity & respository

1. in `internal` folder create `entity` folder and add `user.go` and code
2. in `internal` folder create `api/user` folders and add `repository.go`

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

# Notes for repository.go

`type repository struct {...`
`*db.DB` is a pointer to db.DB. We use pointers in go foe the reasons below

Efficiency

- If `db.DB` is a struct type that may be large, making its copy every time we pass it to a function or method can
  take considerable time and memory. By using a pointer, we're instead passing the address to this data, which is
  significantly more efficient.

Shared State

- If multiple functions or methods need to manipulate the same db.DB instance (which is the same database connection here),
  it's crucial to pass the reference to that instance (a.k.a, the pointer to it), not its copy. If we were to pass
  the copies, each method would work on its own copy of the db.DB data, not affecting the original database
  connection object.

`func (r *repository) CreateUser(user entity.User) error {..`
Here, (r *repository) is the receiver of the CreateUser method. This means `CreateUser is a method defined on the type *repository` (which is a pointer to a repository type).

The receiver in this case `r` essentially behaves as a variable within our method that holds the instance the method
has been called upon.

# Example

let's say we have a Go struct that represents a car and a method to start the car.
Here is the Car struct:

```go

type Car struct {
    model  string
    make   string
    status string
}
```

Now, let's define a method `StartCar` for our Car:

```go
func (c *Car) StartCar() {
    c.status = "running"
}
```

In this StartCar method, `c` is our receiver (like `r` in our code). It is the instance of type `*Car` on which
the method is called.

Here's how you would use it:

```go

func main() {
    myCar := &Car{model: "Model S", make: "Tesla", status: "stopped"}
    fmt.Println(myCar.status) // prints: stopped

    myCar.StartCar()
    fmt.Println(myCar.status) // prints: running
}
```

In main, we first create an instance of Car,`myCar`. The `&` operator is used to get the address in memory of the
created `Car` and returns a pointer to it. We then print it's status which is initially `stopped`.

Then, we call `myCar.StartCar()`, which changes the status field of myCar via the receiver `c` inside the
StartCar method (c.status = "running").

After calling StartCar, when we print myCar.status again, it prints "running".

In summary, the receiver (`r` or `c` in these examples) inside a method holds the instance of the type that the
method was called upon. It can be used to access or modify fields of this instance.

`ctx := context.Background()`
THis creates a new context, `ctx`. It's common in Go to pass context as the first parameter of a function,
especially for I/O operations like database accesses. Contexts can be useful to pass along things like
cancellation signals, deadlines, or other request-scoped values to all the functions and methods
handling a request.

`_, err := r.db.Exec(...`
The `Exec` method runs the SQL command on the database.
