# String Analyzer Service

A powerful RESTful API service built in Go that analyzes strings and stores their computed properties. Features include advanced filtering, natural language query parsing, and comprehensive property analysis.

## ✨ Features

- **Automatic String Analysis** - Computes 6 different properties for each string
- **SHA256 Hashing** - Unique identification and duplicate detection
- **Advanced Filtering** - Filter by palindrome status, length, word count, and characters
- **Natural Language Queries** - Search using plain English (e.g., "single word palindromes")
- **RESTful API** - Full CRUD operations with proper HTTP status codes
- **In-Memory Storage** - Fast O(1) lookups using hash maps
- **Comprehensive Testing** - Unit and integration tests included
- **Production Grade** - Error handling, validation, and health checks

## 🎯 What Gets Analyzed

For each string, the API automatically computes:

| Property | Description |
|----------|-------------|
| `length` | Total number of characters (UTF-8 aware) |
| `is_palindrome` | Boolean: reads same forwards and backwards (case-insensitive) |
| `unique_characters` | Count of distinct characters |
| `word_count` | Number of whitespace-separated words |
| `sha256_hash` | SHA-256 hash for unique identification |
| `character_frequency_map` | Map of each character to its occurrence count |

## 🚀 Quick Start

### Prerequisites

- Go 1.21+ ([download](https://golang.org/dl/))
- Git


### Local Installation

\`\`\`bash
# Clone the repository
git clone https://github.com/yourusername/string-analyzer-api.git
cd string-analyzer-api

# Download dependencies
go mod download

# Run the server
go run main.go
\`\`\`

The API will start on \`http://localhost:8080\`


## 📚 API Endpoints

### 1. Create/Analyze String

**Request:**
\`\`\`bash
POST /strings
Content-Type: application/json

{
  "value": "hello world"
}
\`\`\`

**Response (201 Created):**
\`\`\`json
{
  "id": "7f83b1657ff1fc53b92dc18148a1d65dfc2d4b1fa3d677284addd200126d9069",
  "value": "hello world",
  "properties": {
    "length": 11,
    "is_palindrome": false,
    "unique_characters": 8,
    "word_count": 2,
    "sha256_hash": "7f83b1657ff1fc53b92dc18148a1d65dfc2d4b1fa3d677284addd200126d9069",
    "character_frequency_map": {
      "h": 1,
      "e": 1,
      "l": 3,
      "o": 2,
      "w": 1,
      "r": 1,
      "d": 1
    }
  },
  "created_at": "2025-10-21T14:30:45.123456Z"
}
\`\`\`

### 2. Get All Strings

**Request:**
\`\`\`bash
GET /strings
\`\`\`

**Response (200 OK):**
\`\`\`json
{
  "data": [
    {
      "id": "hash1",
      "value": "radar",
      "properties": { /* ... */ },
      "created_at": "2025-10-21T14:30:45.123456Z"
    }
  ],
  "count": 1,
  "filters_applied": {}
}
\`\`\`

### 3. Get Specific String

**Request:**
\`\`\`bash
GET /strings/{string_value}
\`\`\`

**Example:**
\`\`\`bash
curl http://localhost:8080/strings/hello%20world
\`\`\`

### 4. Filter Strings

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| \`is_palindrome\` | boolean | Filter for palindromes |
| \`min_length\` | integer | Minimum string length |
| \`max_length\` | integer | Maximum string length |
| \`word_count\` | integer | Exact word count |
| \`contains_character\` | string | Single character to search for |

**Examples:**
\`\`\`bash
# Get palindromes only
curl "http://localhost:8080/strings?is_palindrome=true"

# Get strings 5-15 characters long
curl "http://localhost:8080/strings?min_length=5&max_length=15"

# Get single-word strings containing 'a'
curl "http://localhost:8080/strings?word_count=1&contains_character=a"

# Combine multiple filters
curl "http://localhost:8080/strings?is_palindrome=true&word_count=1&min_length=5"
\`\`\`

### 5. Natural Language Filtering

**Request:**
\`\`\`bash
GET /strings/filter-by-natural-language?query={natural_language_query}
\`\`\`

**Supported Queries:**
- \`"single word palindromes"\` → word_count=1, is_palindrome=true
- \`"strings longer than 10 characters"\` → min_length=11
- \`"strings shorter than 20"\` → max_length=19
- \`"palindromic strings"\` → is_palindrome=true
- \`"strings containing letter z"\` → contains_character=z
- \`"strings with the first vowel"\` → contains_character=a

**Examples:**
\`\`\`bash
# Single word palindromes
curl "http://localhost:8080/strings/filter-by-natural-language?query=single%20word%20palindromes"

# Strings longer than 10 characters
curl "http://localhost:8080/strings/filter-by-natural-language?query=strings%20longer%20than%2010"

# Palindromes containing 'a'
curl "http://localhost:8080/strings/filter-by-natural-language?query=palindromic%20strings%20containing%20letter%20a"
\`\`\`

### 6. Delete String

**Request:**
\`\`\`bash
DELETE /strings/{string_value}
\`\`\`

**Response (204 No Content):** Empty response body

**Example:**
\`\`\`bash
curl -X DELETE "http://localhost:8080/strings/hello%20world"
\`\`\`

## 🧪 Testing

### Run All Tests

\`\`\`bash
go test -v
\`\`\`

### Run Specific Tests

\`\`\`bash
go test -run TestPalindromeDetection -v
go test -run TestCharacterFrequency -v
go test -cover
\`\`\`

### Test with cURL

\`\`\`bash
# Health check
curl http://localhost:8080/health

# Create a string
curl -X POST http://localhost:8080/strings \
  -H "Content-Type: application/json" \
  -d '{"value": "radar"}'

# Get all strings
curl http://localhost:8080/strings

# Filter palindromes
curl "http://localhost:8080/strings?is_palindrome=true"

# Natural language search
curl "http://localhost:8080/strings/filter-by-natural-language?query=single%20word%20palindromes"
\`\`\`

## 📊 Usage Examples

### Complete Workflow

\`\`\`bash
# 1. Create multiple strings
curl -X POST http://localhost:8080/strings \
  -H "Content-Type: application/json" \
  -d '{"value": "radar"}'

curl -X POST http://localhost:8080/strings \
  -H "Content-Type: application/json" \
  -d '{"value": "hello world"}'

curl -X POST http://localhost:8080/strings \
  -H "Content-Type: application/json" \
  -d '{"value": "racecar"}'

# 2. Get all strings
curl http://localhost:8080/strings

# 3. Filter for palindromes
curl "http://localhost:8080/strings?is_palindrome=true"

# 4. Filter for single-word strings
curl "http://localhost:8080/strings?word_count=1"

# 5. Natural language search
curl "http://localhost:8080/strings/filter-by-natural-language?query=single%20word%20palindromes"

# 6. Get specific string
curl "http://localhost:8080/strings/radar"

# 7. Delete a string
curl -X DELETE "http://localhost:8080/strings/radar"
\`\`\`

## 🚀 Heroku Deployment

### Prerequisites

- Heroku Account ([create here](https://www.heroku.com))
- Heroku CLI ([install here](https://devcenter.heroku.com/articles/heroku-cli))

### Deploy Steps

\`\`\`bash
# 1. Login to Heroku
heroku login

# 2. Create app
heroku create string-analyzer-api

# 3. Set environment variables
heroku config:set GIN_MODE=release

# 4. Deploy
git push heroku main

# 5. Monitor
heroku logs --tail

# 6. Test
curl https://string-analyzer-api.herokuapp.com/strings
\`\`\`

### Useful Commands

\`\`\`bash
# View logs
heroku logs --tail --app string-analyzer-api

# Open app in browser
heroku open --app string-analyzer-api

# View config variables
heroku config --app string-analyzer-api

# Restart app
heroku restart --app string-analyzer-api

# Check app status
heroku ps --app string-analyzer-api
\`\`\`

## 📁 Project Structure

\`\`\`
string-analyzer-api/
├── main.go                    # Core API implementation
├── go.mod                    # Go module definition
├── go.sum                    # Dependency checksums
├── Procfile                  # Heroku configuration
├── .buildpacks              # Heroku buildpacks
├── .gitignore               # Git ignore rules
├── README.md                # This file
├── LICENSE                  # MIT License
\`\`\`

## 🔧 Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| \`PORT\` | 8080 | Server port |
| \`GIN_MODE\` | debug | Gin framework mode (debug/release) |

### Set Locally

\`\`\`bash
export PORT=9000
export GIN_MODE=release
go run main.go
\`\`\`

## 🛠️ Development

### Building from Source

\`\`\`bash
# Clone repository
git clone https://github.com/yourusername/string-analyzer-api.git
cd string-analyzer-api

# Install dependencies
go mod download

# Run
go run main.go

# Build binary
go build -o string-analyzer

# Run built binary
./string-analyzer
\`\`\`

### Code Structure

- **String Analysis** - \`analyzeString()\` computes all properties
- **Filtering** - \`matchesFilters()\` applies query parameters
- **Natural Language** - \`parseNaturalLanguageQuery()\` interprets English queries
- **HTTP Handlers** - \`createString()\`, \`getString()\`, \`getAllStrings()\`, \`filterByNaturalLanguage()\`, \`deleteString()\`

## 🧠 Algorithm Details

### Palindrome Detection

- Case-insensitive comparison
- Whitespace ignored
- Time complexity: O(n)

### Character Frequency

- Lowercase character mapping
- Whitespace excluded
- Time complexity: O(n)

### Filtering

- Multiple filters with AND logic
- O(m) where m = number of stored strings
- Results sorted chronologically

## ⚠️ Error Handling

The API returns proper HTTP status codes:

| Status | Situation |
|--------|-----------|
| \`201\` | String successfully created |
| \`200\` | GET request successful |
| \`204\` | DELETE successful |
| \`400\` | Bad request (invalid format) |
| \`404\` | String not found |
| \`409\` | String already exists (duplicate) |
| \`422\` | Invalid data type |

## 🔐 Security

### Current Implementation

- No authentication required
- No rate limiting
- In-memory storage (data lost on restart)

### Production Recommendations

- Add JWT/OAuth authentication
- Implement rate limiting
- Use persistent database (PostgreSQL)
- Enable HTTPS/TLS
- Add input validation and sanitization
- Implement API key management
- Set up monitoring and logging

## 📈 Performance

### Time Complexity

| Operation | Complexity |
|-----------|-----------|
| Create String | O(n) |
| Get String | O(1) |
| List All | O(m) |
| Filter | O(m × f) |
| Delete | O(1) |
| Palindrome Check | O(n) |

Where n = string length, m = total strings, f = number of filters

### Space Complexity

- Single string: O(n)
- Total storage: O(m × n)
- Character frequency: O(1) (max 26+ chars)

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Steps

1. Fork the repository
2. Create your feature branch (\`git checkout -b feature/AmazingFeature\`)
3. Commit your changes (\`git commit -m 'Add some AmazingFeature'\`)
4. Push to the branch (\`git push origin feature/AmazingFeature\`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

**Made with ❤️ | Star this repo if you find it useful! ⭐**
