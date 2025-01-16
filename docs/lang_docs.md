# Types

|**Type**|**Description**|**Example**|
|---|---|---|
|**Numeric**|Represents floating-point numbers.|`price = 100.5`|
|**Boolean**|Represents logical values: `true` or `false`.|`is_active = true`|
|**String**|Represents sequences of characters enclosed in double quotes.|`message = "Hello, World!"`|

---
<br>

# Operators

### Arithmetic Operators

|**Operator**|**Description**|**Example**|
|---|---|---|
|`+`|Addition|`sum = 5 + 3`|
|`-`|Subtraction|`diff = 10 - 4`|
|`*`|Multiplication|`prod = 6 * 7`|
|`/`|Division|`div = 20 / 5`|

---

### Relational Operators

|**Operator**|**Description**|**Example**|
|---|---|---|
|`==`|Equal to|`5 == 5 // true`|
|`!=`|Not equal to|`5 != 3 // true`|
|`>`|Greater than|`10 > 5 // true`|
|`<`|Less than|`3 < 7 // true`|
|`>=`|Greater than or equal to|`5 >= 5 // true`|
|`<=`|Less than or equal to|`3 <= 4 // true`|

> **Note:** The `==` operator ensures accurate comparisons for floating-point numbers up to the ninth decimal place, avoiding typical precision issues associated with floating-point arithmetic.

---

### Logical Operators

|**Operator**|**Description**|**Example**|
|---|---|---|
|`and`|Logical AND|`true and false // false`|
|`or`|Logical OR|`true or false // true`|

---

<br>

# Variables
Variables in Cryptor Lang are created by assigning a value directly to a variable name:

```plaintext
a = 2
```

### Key Rules for Variables

1. **Implicit Typing**: Variable types are inferred based on their values:
    - Any valid number is treated as a **float**.
    - The literals `true` and `false` are treated as **boolean**.
    - Any characters enclosed in double quotes are treated as a **string**.
2. **Immutability**: Variables are immutable, meaning they cannot be reassigned after being declared.

---

<br>

# Comments

Cryptor Lang supports single-line comments. Comments begin with `//` and extend to the end of the line. Comments are ignored by the interpreter.

### Example

```plaintext
// This is a single-line comment
price = 100.5 // Comments can be added inline
```

---

<br>

# Control Flow

The only control flow structure in Cryptor Lang is the `if - end` block:

```plaintext
if Condition
    // execute code
end
```

### Usage Example

```plaintext
price = 100.5
if price > 50
    // Take action
end
```

---

<br>

# Functions

Cryptor Lang supports two types of functions: those that perform actions and those that return values. All functions, except for `Take_profit` and `Stop_loss`, require arguments to explicitly specify the parameter names they refer to. Abbreviations are available for some commonly used parameters.

## Action Functions

Action functions are used to perform trading actions, such as placing buy or sell orders.

### `Buy(quantity, percentage)`

- **Description**: Places a MARKET order to buy the base asset. The user must specify either the quantity of the base asset to buy or the percentage of the available quote asset balance to spend on the base asset.
- **Parameters**:
    - `quantity` (optional): The amount of the base asset to buy. This is specified if the user wants to buy a fixed quantity of the base asset.
    - `percentage` (optional): The percentage of the available quote asset balance to spend on the base asset. This is specified if the user wants to use a portion of their quote asset balance.
- **Note**: The user must provide one of the two parameters, `quantity` or `percentage`. If neither is provided, an error will occur.
- **Example**:
    
    ```plaintext
    Buy(quantity = 10)              // Buys 10 units of the base asset.
    Buy(percentage = 50)            // Buys 50% of the available balance of the quote asset in the base asset.
    ```
    

### `Sell()`

- **Description**: Closes the trade by selling the entire amount of the base asset that was previously bought with the `Buy` function. No parameters are needed as the function will automatically sell the entire holding.
- **Example**:
    
    ```plaintext
    Sell() // Sells all units of the base asset acquired by Buy().
    ```
    

### `Stop_loss(stop)`

- **Description**: Places a MARKET sell order if the price falls below the percentage specified by `stop`, relative to the purchase price.
- **Parameters**:
    - `stop` (required): The percentage below the purchase price to trigger the stop-loss order. **The parameter is not explicitly named when calling this function**.
- **Example**:
    
    ```plaintext
    Stop_loss(5) // Sets a stop-loss at 5% below the purchase price.
    ```
    

### `Take_profit(take)`

- **Description**: Places a MARKET sell order if the price rises above the percentage specified by `take`, relative to the purchase price.
- **Parameters**:
    - `take` (required): The percentage above the purchase price to trigger the take-profit order. **The parameter is not explicitly named when calling this function**
- **Example**:
    
    ```plaintext
    Take_profit(10) // Sets a take-profit at 10% above the purchase price.
    ```
    

### Trade Management Rules

- When a trade is open (initiated by the `Buy(quantity = X)` function), subsequent calls to `Buy()` will be ignored until the trade is closed using one of the sell functions (`Sell()`, `Stop_loss(stop = X)`, or `Take_profit(take = X)`).
- If no trade is currently open, all sell functions will be ignored, as there is no active trade to close.

This ensures that:

1. Only one trade can be active at any time.
2. Functions behave predictably within the context of open and closed trades.

---

## Value-Returning Functions

Value-returning functions compute and return data, such as statistical values. These functions require explicit parameter naming, and abbreviations are available for some commonly used parameters.

### `@Mean(window_size, lag)` or `@Mean(ws, lg)`

- **Description**: Returns the average of the last `window_size` price values.
- **Parameters**:
    - `window_size` or `ws` (required): The number of price values to include in the calculation.
    - `lag` or `lg` (optional): Shifts the time window back by the specified number of units (default is 0).
- **Example**:
    
    ```plaintext
    avg = @Mean(window_size = 10, lag = 0)
    avg_abbr = @Mean(ws = 10, lg = 0) // Using abbreviations
    ```
    

### `@Median(window_size, lag)` or `@Median(ws, lg)`

- **Description**: Returns the median of the last `window_size` price values.
- **Parameters**:
    - `window_size` or `ws` (required): The number of price values to include in the calculation.
    - `lag` or `lg` (optional): Shifts the time window back by the specified number of units (default is 0).
- **Example**:
    
    ```plaintext
    median = @Median(window_size = 10, lag = 0)
    ```
    

### `@Std(window_size, lag)` or `@Std(ws, lg)`

- **Description**: Returns the standard deviation of the last `window_size` price values.
- **Parameters**:
    - `window_size` or `ws` (required): The number of price values to include in the calculation.
    - `lag` or `lg` (optional): Shifts the time window back by the specified number of units (default is 0).
- **Example**:
    
    ```plaintext
    std_dev = @Std(ws = 15, lg = 5) // Using abbreviations
    ```
    

### `@Var(window_size, lag)` or `@Var(ws, lg)`

- **Description**: Returns the variance of the last `window_size` price values.
- **Parameters**:
    - `window_size` or `ws` (required): The number of price values to include in the calculation.
    - `lag` or `lg` (optional): Shifts the time window back by the specified number of units (default is 0).
- **Example**:
    
    ```plaintext
    variance = @Var(ws = 10, lg = 2)
    ```
    

### `@Ema(window_size, lag)` or `@Ema(ws, lg)`

- **Description**: Returns the Exponential Moving Average (EMA) of the last `window_size` price values.
- **Parameters**:
    - `window_size` or `ws` (required): The number of price values to include in the calculation.
    - `lag` or `lg` (optional): Shifts the time window back by the specified number of units (default is 0).
- **Example**:
    
    ```plaintext
    ema = @Ema(window_size = 20, lag = 0)
    ```
    

---

## Example Usage of Functions

```plaintext
// Action Functions
Buy(quantity = 5)            // Buys 5 units of the base asset.
Buy(percentage = 50)         // Buys 50% of the available balance of the quote asset in the base asset.
Sell()                       // Sells all units of the base asset acquired by Buy().
Take_profit(take = 10)       // Sets a take-profit at 10% above the purchase price.
Stop_loss(stop = 5)          // Sets a stop-loss at 5% below the purchase price.

// Value-Returning Functions
avg = @Mean(ws = 10, lg = 0)       // Average of the last 10 price values.
median = @Median(ws = 10, lg = 0)  // Median of the last 10 price values.
std_dev = @Std(ws = 15, lg = 5)    // Standard deviation of the last 15 prices, shifted back by 5 units.
ema = @Ema(ws = 20, lg = 0)        // EMA of the last 20 price values.
```
