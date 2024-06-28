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

# GO RESOURCES

Go routines - Go Programming Tutorial
`https://www.youtube.com/watch?v=FJo6eQSWruQ`

Pointers in Go Programming Turorial
`https://www.youtube.com/watch?v=v-ttLYKqaO8&list=PL7g1jYj15RUMMCMDYPyZHN3CaWbt3Rl5y&index=3`
`https://www.youtube.com/watch?v=496xSrA6QQ8`

Advanced Golang: Generics Explained
`https://www.youtube.com/watch?v=WpTKqnfp5dY&t=653s`

Closures in Golang
`https://www.youtube.com/watch?v=jHd0FczIjAE&list=PL7g1jYj15RUMMCMDYPyZHN3CaWbt3Rl5y&index=1`

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////

SETTING UP

1. Setup up Go project - Boards
2. Add Server folder
3. Create `docker-compose.yml` and `Dockerfile` and intialize git and .gitignore
4. in server folder create cmd folder and add initial `main.go` file to just print "Hello World"
5. add `Frontend` folder and create `Nextjs` app with TS and Tailwind

CONFIG

1. in server folder add `internal/config` folders
2. add test.env file withvars
3. in config folder add `config.go` and `config_test.go` and add code
4. Run `config_test.go` to test we get a pass

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
   The `Validate` method on `DatabaseConfig` struct utilizes the validator package to check if all struct fields
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

   A pointer is used here so we dont have to copy the Config Object to increatese efficienecy efficiency and to allow
   the function to return a `nil` value if novalid configuration can be loaded.

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

Overall, the `config.go` files handles configuration setup and ensures robust initialization of your application by
ensuring the database configuration is properly set up before any database connection is attempted.

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////

# Step-by-Step Flow:

1. main.go calls `Load()`

   - The application begins in the `main.go` file, where it calls the `Load()` function defined in `config.go`.
   - The purpose is to retrieve and establish the application configuration, focusing particularly on database
     configuration before any interactions with the database occur.

2. Load() function in config.go:

   - Once invoked, Load() itself calls the `getDatabaseConfig()` function, also located within `config.go`.
   - `Load()` is responsible for managing the overall fetching and validation of the configuration data. It acts as a
     wrapper to any configuration setup steps and errors management.

3. getDatabaseConfig() function in config.go

   - The purpose of `getDatabaseConfig()` is to actually fetch the environment variables that are required for the
     database configuration.

   - It constructs a `DatabaseConfig` struct with all necessary database parameters fetched from the environment
     (like `database` `host`, `port`, `username`, and `password`).

   - After assembling the DatabaseConfig, it then calls the `Validate()` method on this configuration instance.

4. Validate() function in config.go

   - This method is a crucial safety check to ensure that all required fields in the `DatabaseConfig` are properly set.
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

   - This function sets test environment variables which simulate a correctly configured environment necessary
     for the `Load()`function to extract and assemble a valid `DatabaseConfig`.

   - By using `os.Setenv`, it mimics the conditions under which your application would run in a properly configured
     real-world setup, providing values like `localhost`, `5432`, `dbname`, `user`, `password` for the associated DB
     environment variables.

2. clearDBEnvVars() Function:

   - This function clears the environment variables related to the database configuration by using `os.Unsetenv`.

   - The purpose here is to simulate a scenario where these key configuration values are not set, which should lead
     to errors when trying to load the configuration using `Load()`.

3. TestLoad Cases:
   `t.Run("valid db env vars",...` & `t.Run("missing db env vars",`

   - `valid db env vars`: After setting the environment variables using `setDBEnvVars()`, it runs `Load()` to check
     if it can successfully load the configuration without any errors. It further asserts that the returned configuration
     is not nil, the database configuration specifically is not nil, and that it matches expected values
     (like `localhost` for the host).

     The load function is called and returns

     - `cfg` which is a pointer to a Config object that contains configuration information, specifically DatabaseConfig
     - `err` an error object that Load() returns in case something goes wrong while getting the database configuration

   - `Missing db env_vars`: This test ensures that `Load()` correctly handles the absence of necessary environment variables
     by clearing them first and then attempting to load the configuration. It uses `assert.Error(t, err)` to confirm that
     an error is indeed returned, as expected in such cases.

   - Test Assertions
     The use of `assert` from the `github.com/stretchr/testify/assert` package helps in making the tests concise and readable.
     It provides a range of functions like `assert.NoError`, `assert.NotNil`, and `.assert.Error` to check for conditions such
     as the absence of errors, the non-nullity of objects, and the presence of expected errors, respectively.

DB CONFIGURATION

1. in server folder add `db` folder and `migrations` folder in that
2. in db folder add `db.go` and `db_test.go` files

# Notes for db.go

The `db.go` file is structured to manage database connections using the Go programming language. It interfaces with PostgreSQL
using the `pgxpool` package from `github.com/jackc/pgx/v5/pgxpool`. This package helps efficiently manage multiple database
connections by using a connection pool.

`type DB struct {...`
DB is a struct that embeds a pointer to a `pgxpool.Pool`. This pool manages and maintains a set of database connections that
can be reused for multiple database queries. In Go, an embedded pointer is a pointer used within a struct which `inherits` the
methods and fields from another type, effectively embedding that type into the struct. This is similar to `inheritance`
in object-oriented programming, but implemented in Go's own unique way.

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

Here's an example of how the connection string might look:
`postgres://myUser:myPassword@localhost:5432/myDatabase?sslmode=disable`

This function uses the `fmt.Sprintf()` function to replace the `%v` placeholders with the values from the cfg parameter,
which is an instance of `config.DatabaseConfig`. Using `%v` is a convenient way to create a string with variables.

`sslmode=disable`
This is included in the URL, which specifies that SSL will not be used for this connection. This is typical in development
environments, but for production environments, using SSL is recommended to ensure security.

NOTE:
We have created `migrate.go` in db folder and added code but this isnt being used at thie stage and might not be neede

3. update `main.go` to call to include `config.Load()` and `db.Connect()`
4. Note... More code will have to be created before we can use our `db_test.go`

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////

USER ENTRY AND REPOSITORY

1. in `internal` folder create `entity` folder and add `user.go` and code
2. in `internal` folder create `api/user` folders and add `repository.go`

# Notes for repository.go

`type repository struct {...`
`*db.DB` is a pointer to db.DB. We use pointers in go for the reasons below

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
Here, (r *repository) is the receiver of the CreateUser method, meaning CreateUser is now a method defined on the
type `*repository` - a pointer to a repository type.

The receiver in this case `r` essentially behaves as a variable within our method that holds the instance the method
has been called upon.

# Example >>>>>

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

# End example >>>>>>>>>

`ctx := context.Background()`
THis creates a new context, `ctx`. It's common in Go to pass context as the first parameter of a function,
especially for I/O operations like database accesses. Contexts can be useful to pass along things like
`cancellation signals`, `deadlines`, or other request-scoped values to all the functions and methods
handling a request.

`_, err := r.db.Exec(...`
The `Exec` method runs the SQL command on the database.

The `pgerrcode` package is used for better readability and easier maintenance of our code when checking for
specific PostgreSQL error codes.

//////////////////////////////////////////////////////////////////////////////////////////////////////////////

# user service

1. in `api/user` add `repository_mock.go` and add initial code
2. in `api/user` folder create `service.go` and add code
3. in `api/user` folder create `service_test.go` and add code

# Notes for service.go

`user := entity.User{...`
Cretes a new User instance filling in its attributes using values from input

`err := s.repo.CreateUser(user)`
Now that we have a user entity ready to go, we call the CreateUser method on s.repo, which is an instance
of our Repository. This method should save the new user to the database

The pattern used here is known as the service pattern. It involves defining an interface (here `Service`),
implementing that interface (with `service`), and then creating instances of this implementation through
a function (like `NewService`). Let's delve into the reasons, why this approach is quite useful
and advantageous:

1. Decoupling and Modularity:
   By using an interface, you're able to decouple your components which can be swapped out with
   different implementations.

2. Code Reusability and Single Responsibility Principle:
   Interfaces and Services in Go enables us to abstract out complex functionalities and encapsulate
   them in their dedicated services.

3. Easier Testing: It's easier to test your components. You can create a mock version of your Repository
   while testing the Service.

4. Code Readability and Maintenance: It increases readability by making it clear what methods are expected
   to be implemented by a service and are handled by the interface.

5. Dependency Injection:
   It enables easier dependency injection. In your case, service can't do anything without Repository.
   Injecting it through NewService allows you to make this dependency explicit and switch out Repository
   with another implementation whenever you like.

# Create User Handler

1. add `github.com/go-chi/chi/v5 v5.0.8 // indirect` in go.mod
2. add `github.com/go-chi/chi/v5 v5.0.8 h1:lD+NLqFcAi1ovnVZpsnObHGW4xb4J8lNmoYVfECH1Y0=` &
   `github.com/go-chi/chi/v5 v5.0.8/go.mod h1:DslCQbL2OYiznFReuXYUmQ2hGd1aDpCnlMNITLSKoi8=`
   in go.sum

3. create handler.go in api/user folder

# Notes for hander.go

handler.go holds the function that maps API endpoints (such as RESTful routes or HTTP routes) to corresponding
handlers in our user service.

`(api *API) RegisterHandlers(r chi.Router)`
This is a method declaration for the API struct. The RegisterHandlers function takes a router from the Chi router
package as an argument. Chi is a lightweight, idiomatic and composable router for building Go HTTP services.

`r.Post("/users", api.HandleCreateUser)`
This line registers a new HTTP POST route to the path `/users`. Any HTTP POST requests to `/users` will be handled by
the `HandleCreateUser` function of the API.

The purpose of defining routes and handlers this way is to maintain a clean and modular structure in our application,
enabling us to handle different paths with different functions. This gives a good separation of concerns where each
function does one thing and does it well.

/////////////////////////////////////////////////////////////////////////////////////////////////

3. In `server/internal/api/user/` folder add `api.go`

# Notes for api.go

`type API struct...`
Holds an instance of `Service` (which would implement the `CreateUser` method) as a field. This allows the API
struct to use or `expose` the `CreateUser` functionality of the Service interface located in service.go

`func NewAPI`
A constructor used to initialize and return a new instance of the `API struct`, using a given service that
implements the `Service` interface

`func (req *CreateUserRequest) Validate() error {...`
A method of the CreateUserRequest struct. It validates the fields of the CreateUserRequest structure. This
function returns an error, which will be nil if the validation is successful or an explicit error if there
is a validation issue.

`func (api *API) HandleCreateUser(w http.ResponseWriter...`
a handler function for a HTTP endpoint in a web server. It processes HTTP requests that aim to create a
new user.

It first declares a variable named `createUserRequest` of type CreateUserRequest. It then reads the HTTP
request body `(r.Body)`, expecting it to contain a JSON-encoded representation of a CreateUserRequest.

Lets break down this line: `json.NewDecoder(r.Body).Decode(&createUserRequest)`

`json.NewDecoder(r.Body)`
This part creates a JSON decoder using the HTTP request's body. This decoder will be able to read and decode JSON
content present in the HTTP request's body.

`.Decode(&createUserRequest)`
This uses the JSON decoder to read the JSON from the request's body and convert it to a Go object. The &createUserRequest
argument provides a memory location where the converted Go object should be stored. This conversion is based on the
structure of CreateUserRequest struct.

Once the request data is decoded into a Go struct, it attempts to validate the request using the Validate() function

If the validation is successful, it will proceed to create the user. It uses CreateUser Input(createUserRequest) to
create a new CreateUserInput instance from the request, then passes this to the api.service.CreateUser() method to actually
create the user.

If the user is created successfully, the function then creates a CreateUserResponse using the details of the newly
created user, encodes this as JSON, and sends it back to the client with an HTTP status 201 Created.

`json.Decode()` is used to decode or parse JSON data into Go data structures.
`json.Encode()` is used to encode or stringify Go data structures into JSON

This is a fairly standard pattern for HTTP handlers in Go: `decode` `request`, `validate` `request`, perform
an action, then send a response.

Create `api_test.go` and add code

///////////////////////////////////////////////////////////////////////////////////////////////////////////

> > # Notes for api_test.go

`payload := strings.NewReader(...):`
This is preparing a payload to simulate the client request body. It uses strings.NewReader to create
a Reader with the JSON request that could be sent from a client.

`res := httptest.NewRecorder()`:
This creates a ResponseRecorder to record the HTTP response.

`req, _ := http.NewRequest(...)`
creates a new HTTP request object. This represents the incoming HTTP request.

`api.HandleCreateUser(res, req)`
This is the function call to HandleCreateUser you're actually testing. It passes the response recorder
and the HTTP request to record what would have been a HTTP response for this request.

`assert.Equal(t, res.Code, http.StatusCreated)`
This is an assertion using the testify package's assert module to check that the HTTP status code
that was written to res in HandleCreateUser is indeed http.StatusCreated (i.e., 201). If it isn't, the
test will fail.

This test ensures that when provided with a valid CreateUser request, the HandleCreateUser function
correctly creates auser and responds with a HTTP 201 (Created) status.

# A NOTE ON TESTING IN GO

In a well-structured application, you will generally have tests in three tiers: `repository tests`,
`service tests`, and `API (or controller)` tests. Each of these tiers is responsible for testing different
aspects of the application.

`Repository (or DAO or model) tests`:
These tests are focused on data access - the layer of your application that interacts directly with your
data store. In these tests, you would ensure that your data access code correctly saves, retrieves, updates,
and deletes data in your database as expected.

`Service (or business logic) tests:`
Here you're testing the layer of your application that contains your core business logic. This is
usually where you enforce business rules and handle interactions between different data models. Service
tests would validate things like manipulating a model based on some event, calling the repository layer to
reflect these changes, handling transactions etc.

`API (or controller) tests:`
These tests focus on the layer of your application that interacts with external clients by sending and
receiving data over the network. API tests validate things like the routing (do requests to /path get handled by
the right function?),request and response formatting (do we properly process incoming request data? do we send
correctly formatted response data?), proper error handling (do we respond with the correct HTTP status codes
and error messages?), etc.

///////////////////////////////////////////////////////////////////////////////////////////////////////////
