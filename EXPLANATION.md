# Explanation

## B07 - Dokumentasi API
You can read the full documentation here [Swagger](https://app.swaggerhub.com/apis/bangkitdc/single-service/1.0.0)

### Endpoint
By default, there was 3 APIs that you can use. Here it is :

#### User
|Method| URL | Explanataion | Need Auth |
|--|--|--|:--:|
| POST | base_url/login | Login & Get Token | |
| GET | base_url/self | Get User Data | &#10004; |

#### Barang
|Method| URL | Explanataion | Need Auth |
|--|--|--|:--:|
| GET | base_url/barang | Get All Barang | &#10004; |
| POST | base_url/barang | Create Barang | &#10004; |
| GET | base_url/barang/:id | Get Barang by ID | &#10004; |
| PUT | base_url/barang/:id | Update Barang by ID | &#10004; |
| DELETE | base_url/barang/:id | Delete Barang by ID | &#10004; |

#### Perusahaan
|Method| URL | Explanataion | Need Auth |
|--|--|--|:--:|
| GET | base_url/perusahaan | Get All Perusahaan | &#10004; |
| POST | base_url/perusahaan | Create Perusahaan | &#10004; |
| GET | base_url/perusahaan/:id | Get Perusahaan by ID | &#10004; |
| PUT | base_url/perusahaan/:id | Update Perusahaan by ID | &#10004; |
| DELETE | base_url/perusahaan/:id | Delete Perusahaan by ID | &#10004; |

### Extra Endpoint
For Monolith Service purpose only:

#### Barang No Auth
|Method| URL | Explanataion | Need Auth |
|--|--|--|:--:|
| GET | base_url/barang-paginate | Get All Barang With Paginations | |
| GET | base_url/barang-noauth/:id | Get Barang by ID | |
| PUT | base_url/barang-noauth/:id | Update Barang by ID | |
| GET | base_url/barang-noauth-recommendation | Get Barang Recommendation | |

#### Perusahaan No Auth
|Method| URL | Explanataion | Need Auth |
|--|--|--|:--:|
| GET | base_url/perusahaan-noauth/:id | Get Perusahaan by ID | |

## B08 - SOLID
1. Single Responsibility Principle (SRP)

Each component or class should have a single responsibility. I'm using Models-Controllers (MC) pattern and then added Middleware and Helper to support and provide great functionality for the application.

- Models
  - Represents the application's data and connect DB.
  - Responsible for data manipulation, validation, and database interactions.
  - Migration and Seeding.
- Controllers
  - Handles user input and updates the Model.
  - Contains application logic related to routing, data flow, and data constraints.
- Middleware
  - Add cross-cutting concerns to your application, such as authentication.
  - Implementation: JWT token checker, for certain API that needs access.
- Helper
  - Helpers are utility functions that provide common functionalities used across different parts of application.
  - Implementation: Response handler.

2. Open/Closed Principle (OCP)

Entities (classes, modules, functions) should be open for extension but closed for modification.

**Implementation:** <br/>
On my application I don't have to make another route for Get Barang by ID for Monolith Service that don't need to be authenticated when accessing the API. Instead I'm using the same Get Barang by ID method but without the middleware, so that anyone can consume the API. Beside of that I'm using function CheckConstraint that can be used by two conditional, I added checkKode (boolean) on the parameter, I extend the used of that funciton instead of making a new one.

3. Liskov Substitution Principle (LSP)

The Liskov Substitution Principle (LSP) is primarily concerned with the behavior of objects in a class hierarchy and how derived types can be substituted for their base types without affecting the correctness of the program. Since my doesn't involve inheritance or class hierarchies, LSP is not directly applicable in this context.

**Implementation:** <br/>
In Golang, LSP is not enforced through class inheritance but rather through the use of interfaces. Interfaces define behavior, and any type that satisfies the interface can be used interchangeably. Therefore, my implementation towards LSP is to use base response the APIResponse with Data inteface{}, which later will be extend by other type Data, such as PaginatedResponse, BarangResponse, PerusahaanResponse, etc. So the inheritance is not directly applicable, but use the same logic.

4. Interface Segregation Principle (ISP)

The Interface Segregation Principle is about creating specific interfaces that are tailored to the needs of the clients that use them, rather than having large, monolithic interfaces.

**Implementation:**

```go
func GetPerusahaans(w http.ResponseWriter, r *http.Request) {
  // ...
}

func GetPerusahaanByID(w http.ResponseWriter, r *http.Request) {
  // ...
}

func CreatePerusahaan(w http.ResponseWriter, r *http.Request) {
  // ...
}

func UpdatePerusahaan(w http.ResponseWriter, r *http.Request) {
  // ...
}

func DeletePerusahaan(w http.ResponseWriter, r *http.Request) {
  // ...
}
```
In this code, each of the functions above represents an HTTP handler that corresponds to a specific HTTP endpoint (GET, POST, PUT, DELETE) related to the resource. ISP is adhered to in this case because each function represents a separate and specific interface that handles one type of request (HTTP method) related to resources. The functions are segregated based on their specific responsibilities.

```go
type BarangRequest struct {
  // ...
}

type BarangResponse struct {
  // ...
}

type PaginatedResponse struct {
  // ...
}

type Pagination struct {
  // ...
}

type BarangRequestQuantity struct {
  // ...
}
```
In the code, there are several data structures (structs) representing different request and response formats for the "barang" (item) resource in a web server application. Each struct serves a specific purpose and has a clear and focused responsibility:

5. Dependency Inversion Principle (DIP)

The DIP states that high-level modules should not depend on low-level modules; both should depend on abstractions. 

**Implementation:** <br/>
There doesn't seem to be explicit Dependency Inversion Principle (DIP) implementation. The controller directly interacts with the models.DB instance and doesn't use any abstraction or dependency injection to decouple itself from the data access layer. However in JWT middleware, I used JWTClaims. By doing this, I've abstracted the JWT claims, which can be helpful if I ever need to change the implementation of the JWT claims in the future without affecting the middleware's functionality. This promotes better separation of concerns and reduces tight coupling between components.
