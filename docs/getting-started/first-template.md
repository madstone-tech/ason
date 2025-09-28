# ※ Your First Template

> *Learn the sacred art of template creation from foundation to mastery*

This comprehensive guide teaches you how to create powerful, reusable templates that will serve as the foundation for all your future projects.

## Understanding Templates

A template in Ason is a directory structure that serves as a blueprint for generating new projects. Templates use the Pongo2 templating engine (similar to Django/Jinja2) to enable dynamic content generation through variables and filters.

### Template Anatomy

```
my-template/
├── ason.toml              # Template configuration
├── README.md              # Template documentation
├── {{ project_name }}/    # Dynamic directory names
│   ├── src/
│   │   ├── main.{{ language_ext }}  # Dynamic file names
│   │   └── config.json    # Files with variable content
│   ├── tests/
│   └── docs/
├── .gitignore
└── LICENSE
```

## Step 1: Planning Your Template

Before creating a template, consider:

### Template Purpose
- **What type of projects** will this template create?
- **What technologies** does it include?
- **What structure** do projects need?
- **What variables** will users customize?

### Target Audience
- **Skill level**: Beginner, intermediate, or advanced developers?
- **Use case**: Personal projects, team development, enterprise?
- **Customization needs**: High flexibility or opinionated structure?

### Example: Planning a Node.js API Template

```yaml
Purpose: RESTful API with Express.js
Technologies: Node.js, Express, Jest, ESLint
Structure: MVC pattern with controllers, models, routes
Variables:
  - api_name: Name of the API
  - port: Server port
  - database: Database type (mongodb, postgresql, sqlite)
  - auth: Authentication method (jwt, oauth, none)
  - testing: Include tests (true/false)
```

## Step 2: Create the Template Structure

Let's build a comprehensive Node.js API template step by step:

### Basic Directory Structure

```bash
# Create the template directory
mkdir nodejs-api-template
cd nodejs-api-template

# Create the directory structure
mkdir -p {src/{controllers,models,routes,middleware,utils},tests,docs,config}
```

### Configuration File (ason.toml)

```bash
cat > ason.toml << 'EOF'
name = "Node.js API Template"
description = "Professional Node.js REST API with Express, testing, and best practices"
version = "1.0.0"
author = "Your Name"
type = "backend"

# Template variables
[[variables]]
name = "api_name"
description = "Name of the API (e.g., 'User Management API')"
required = true

[[variables]]
name = "project_slug"
description = "Project slug for package.json (auto-generated from api_name)"
required = false
default = "{{ api_name | lower | replace(' ', '-') | replace('_', '-') }}"

[[variables]]
name = "description"
description = "API description"
required = false
default = "A REST API built with Node.js and Express"

[[variables]]
name = "author"
description = "API author"
required = true

[[variables]]
name = "version"
description = "Initial version"
required = false
default = "1.0.0"

[[variables]]
name = "port"
description = "Server port"
required = false
default = "3000"

[[variables]]
name = "database"
description = "Database type"
required = false
default = "mongodb"
options = ["mongodb", "postgresql", "mysql", "sqlite"]

[[variables]]
name = "include_auth"
description = "Include authentication middleware"
required = false
default = "true"
type = "boolean"

[[variables]]
name = "include_tests"
description = "Include test setup"
required = false
default = "true"
type = "boolean"

[[variables]]
name = "license"
description = "License type"
required = false
default = "MIT"

# Files to ignore during processing
ignore = ["node_modules/", "*.log", ".env", "coverage/"]

# Template metadata
tags = ["nodejs", "api", "express", "rest", "backend"]
EOF
```

### Package.json Template

```bash
cat > package.json << 'EOF'
{
  "name": "{{ project_slug }}",
  "version": "{{ version }}",
  "description": "{{ description }}",
  "main": "src/app.js",
  "scripts": {
    "start": "node src/app.js",
    "dev": "nodemon src/app.js",
    "test": "{% if include_tests %}jest{% else %}echo 'No tests configured'{% endif %}",
    "test:watch": "{% if include_tests %}jest --watch{% else %}echo 'No tests configured'{% endif %}",
    "test:coverage": "{% if include_tests %}jest --coverage{% else %}echo 'No tests configured'{% endif %}",
    "lint": "eslint src/",
    "lint:fix": "eslint src/ --fix"
  },
  "keywords": [
    "api",
    "rest",
    "nodejs",
    "express"{% if database == "mongodb" %},
    "mongodb"{% elif database == "postgresql" %},
    "postgresql"{% elif database == "mysql" %},
    "mysql"{% endif %}
  ],
  "author": "{{ author }}",
  "license": "{{ license }}",
  "dependencies": {
    "express": "^4.18.2",
    "cors": "^2.8.5",
    "helmet": "^6.1.5",
    "dotenv": "^16.0.3"{% if database == "mongodb" %},
    "mongoose": "^7.0.3"{% elif database == "postgresql" %},
    "pg": "^8.10.0",
    "sequelize": "^6.31.0"{% elif database == "mysql" %},
    "mysql2": "^3.2.4",
    "sequelize": "^6.31.0"{% elif database == "sqlite" %},
    "sqlite3": "^5.1.6",
    "sequelize": "^6.31.0"{% endif %}{% if include_auth %},
    "jsonwebtoken": "^9.0.0",
    "bcryptjs": "^2.4.3"{% endif %}
  },
  "devDependencies": {
    "nodemon": "^2.0.22",
    "eslint": "^8.39.0"{% if include_tests %},
    "jest": "^29.5.0",
    "supertest": "^6.3.3"{% endif %}
  }{% if include_tests %},
  "jest": {
    "testEnvironment": "node",
    "collectCoverageFrom": [
      "src/**/*.js",
      "!src/config/**"
    ]
  }{% endif %}
}
EOF
```

### Main Application File

```bash
cat > src/app.js << 'EOF'
const express = require('express');
const cors = require('cors');
const helmet = require('helmet');
require('dotenv').config();

{% if database == "mongodb" %}
const mongoose = require('mongoose');
{% elif database in ["postgresql", "mysql", "sqlite"] %}
const { sequelize } = require('./config/database');
{% endif %}

const app = express();
const PORT = process.env.PORT || {{ port }};

// Middleware
app.use(helmet());
app.use(cors());
app.use(express.json());
app.use(express.urlencoded({ extended: true }));

{% if include_auth %}
// Authentication middleware
const authMiddleware = require('./middleware/auth');
{% endif %}

// Routes
app.get('/', (req, res) => {
  res.json({
    message: 'Welcome to {{ api_name }}!',
    version: '{{ version }}',
    status: 'running'
  });
});

// API Routes
const apiRoutes = require('./routes/api');
app.use('/api', apiRoutes);

{% if include_auth %}
// Protected routes example
app.get('/api/protected', authMiddleware, (req, res) => {
  res.json({ message: 'This is a protected route', user: req.user });
});
{% endif %}

// Error handling middleware
app.use((err, req, res, next) => {
  console.error(err.stack);
  res.status(500).json({ error: 'Something went wrong!' });
});

// 404 handler
app.use('*', (req, res) => {
  res.status(404).json({ error: 'Route not found' });
});

// Database connection and server start
async function startServer() {
  try {
    {% if database == "mongodb" %}
    // Connect to MongoDB
    await mongoose.connect(process.env.MONGODB_URI || 'mongodb://localhost:27017/{{ project_slug }}');
    console.log('Connected to MongoDB');
    {% elif database in ["postgresql", "mysql", "sqlite"] %}
    // Connect to database
    await sequelize.authenticate();
    console.log('Database connection established');

    // Sync database (create tables)
    await sequelize.sync();
    console.log('Database synchronized');
    {% endif %}

    app.listen(PORT, () => {
      console.log(`{{ api_name }} server running on port ${PORT}`);
    });
  } catch (error) {
    console.error('Failed to start server:', error);
    process.exit(1);
  }
}

startServer();

module.exports = app;
EOF
```

### Routes and Controllers

```bash
# Create API routes
cat > src/routes/api.js << 'EOF'
const express = require('express');
const router = express.Router();

// Import controllers
const healthController = require('../controllers/healthController');
// const userController = require('../controllers/userController');

// Health check route
router.get('/health', healthController.getHealth);

// Add your API routes here
// router.get('/users', userController.getAllUsers);
// router.post('/users', userController.createUser);
// router.get('/users/:id', userController.getUserById);
// router.put('/users/:id', userController.updateUser);
// router.delete('/users/:id', userController.deleteUser);

module.exports = router;
EOF

# Create health controller
cat > src/controllers/healthController.js << 'EOF'
/**
 * Health check controller for {{ api_name }}
 */

const getHealth = (req, res) => {
  res.json({
    status: 'healthy',
    timestamp: new Date().toISOString(),
    uptime: process.uptime(),
    service: '{{ api_name }}',
    version: '{{ version }}'
  });
};

module.exports = {
  getHealth
};
EOF

# Create example user controller
cat > src/controllers/userController.js << 'EOF'
/**
 * User controller for {{ api_name }}
 * TODO: Implement actual user logic based on your needs
 */

const getAllUsers = async (req, res) => {
  try {
    // TODO: Implement user fetching logic
    res.json({
      message: 'Get all users',
      users: []
    });
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
};

const createUser = async (req, res) => {
  try {
    // TODO: Implement user creation logic
    res.status(201).json({
      message: 'User created successfully',
      user: req.body
    });
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
};

const getUserById = async (req, res) => {
  try {
    const { id } = req.params;
    // TODO: Implement user fetching by ID
    res.json({
      message: `Get user ${id}`,
      user: { id }
    });
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
};

const updateUser = async (req, res) => {
  try {
    const { id } = req.params;
    // TODO: Implement user update logic
    res.json({
      message: `User ${id} updated successfully`,
      user: { id, ...req.body }
    });
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
};

const deleteUser = async (req, res) => {
  try {
    const { id } = req.params;
    // TODO: Implement user deletion logic
    res.json({
      message: `User ${id} deleted successfully`
    });
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
};

module.exports = {
  getAllUsers,
  createUser,
  getUserById,
  updateUser,
  deleteUser
};
EOF
```

### Database Configuration

```bash
{% if database == "mongodb" %}
# MongoDB model example
cat > src/models/User.js << 'EOF'
const mongoose = require('mongoose');

const userSchema = new mongoose.Schema({
  name: {
    type: String,
    required: true,
    trim: true
  },
  email: {
    type: String,
    required: true,
    unique: true,
    lowercase: true,
    trim: true
  },
  password: {
    type: String,
    required: true,
    minlength: 6
  },
  role: {
    type: String,
    enum: ['user', 'admin'],
    default: 'user'
  },
  isActive: {
    type: Boolean,
    default: true
  }
}, {
  timestamps: true
});

module.exports = mongoose.model('User', userSchema);
EOF
{% else %}
# Sequelize configuration
cat > src/config/database.js << 'EOF'
const { Sequelize } = require('sequelize');

{% if database == "postgresql" %}
const sequelize = new Sequelize(
  process.env.DB_NAME || '{{ project_slug }}',
  process.env.DB_USER || 'postgres',
  process.env.DB_PASSWORD || '',
  {
    host: process.env.DB_HOST || 'localhost',
    port: process.env.DB_PORT || 5432,
    dialect: 'postgres',
    logging: process.env.NODE_ENV === 'development' ? console.log : false
  }
);
{% elif database == "mysql" %}
const sequelize = new Sequelize(
  process.env.DB_NAME || '{{ project_slug }}',
  process.env.DB_USER || 'root',
  process.env.DB_PASSWORD || '',
  {
    host: process.env.DB_HOST || 'localhost',
    port: process.env.DB_PORT || 3306,
    dialect: 'mysql',
    logging: process.env.NODE_ENV === 'development' ? console.log : false
  }
);
{% elif database == "sqlite" %}
const sequelize = new Sequelize({
  dialect: 'sqlite',
  storage: process.env.DB_PATH || './{{ project_slug }}.sqlite',
  logging: process.env.NODE_ENV === 'development' ? console.log : false
});
{% endif %}

module.exports = { sequelize };
EOF

# Sequelize User model
cat > src/models/User.js << 'EOF'
const { DataTypes } = require('sequelize');
const { sequelize } = require('../config/database');

const User = sequelize.define('User', {
  id: {
    type: DataTypes.INTEGER,
    primaryKey: true,
    autoIncrement: true
  },
  name: {
    type: DataTypes.STRING,
    allowNull: false,
    validate: {
      notEmpty: true,
      len: [1, 100]
    }
  },
  email: {
    type: DataTypes.STRING,
    allowNull: false,
    unique: true,
    validate: {
      isEmail: true
    }
  },
  password: {
    type: DataTypes.STRING,
    allowNull: false,
    validate: {
      len: [6, 255]
    }
  },
  role: {
    type: DataTypes.ENUM('user', 'admin'),
    defaultValue: 'user'
  },
  isActive: {
    type: DataTypes.BOOLEAN,
    defaultValue: true
  }
}, {
  tableName: 'users',
  timestamps: true
});

module.exports = User;
EOF
{% endif %}
```

### Authentication Middleware (Conditional)

```bash
{% if include_auth %}
cat > src/middleware/auth.js << 'EOF'
const jwt = require('jsonwebtoken');

const authMiddleware = (req, res, next) => {
  const token = req.header('Authorization')?.replace('Bearer ', '');

  if (!token) {
    return res.status(401).json({ error: 'Access denied. No token provided.' });
  }

  try {
    const decoded = jwt.verify(token, process.env.JWT_SECRET || 'your-secret-key');
    req.user = decoded;
    next();
  } catch (error) {
    res.status(400).json({ error: 'Invalid token.' });
  }
};

module.exports = authMiddleware;
EOF
{% endif %}
```

### Environment Configuration

```bash
cat > .env.example << 'EOF'
# Server Configuration
PORT={{ port }}
NODE_ENV=development

# Database Configuration
{% if database == "mongodb" %}
MONGODB_URI=mongodb://localhost:27017/{{ project_slug }}
{% elif database == "postgresql" %}
DB_HOST=localhost
DB_PORT=5432
DB_NAME={{ project_slug }}
DB_USER=postgres
DB_PASSWORD=
{% elif database == "mysql" %}
DB_HOST=localhost
DB_PORT=3306
DB_NAME={{ project_slug }}
DB_USER=root
DB_PASSWORD=
{% elif database == "sqlite" %}
DB_PATH=./{{ project_slug }}.sqlite
{% endif %}

{% if include_auth %}
# Authentication
JWT_SECRET=your-very-secure-secret-key-here
JWT_EXPIRES_IN=24h
{% endif %}

# API Configuration
API_VERSION=v1
CORS_ORIGIN=http://localhost:3000
EOF
```

### Tests (Conditional)

```bash
{% if include_tests %}
cat > tests/app.test.js << 'EOF'
const request = require('supertest');
const app = require('../src/app');

describe('{{ api_name }} API', () => {
  describe('GET /', () => {
    it('should return welcome message', async () => {
      const response = await request(app)
        .get('/')
        .expect(200);

      expect(response.body).toHaveProperty('message');
      expect(response.body).toHaveProperty('version', '{{ version }}');
      expect(response.body).toHaveProperty('status', 'running');
    });
  });

  describe('GET /api/health', () => {
    it('should return health status', async () => {
      const response = await request(app)
        .get('/api/health')
        .expect(200);

      expect(response.body).toHaveProperty('status', 'healthy');
      expect(response.body).toHaveProperty('service', '{{ api_name }}');
      expect(response.body).toHaveProperty('version', '{{ version }}');
    });
  });

  describe('GET /nonexistent', () => {
    it('should return 404 for unknown routes', async () => {
      const response = await request(app)
        .get('/nonexistent')
        .expect(404);

      expect(response.body).toHaveProperty('error', 'Route not found');
    });
  });
});
EOF

cat > tests/controllers/healthController.test.js << 'EOF'
const request = require('supertest');
const app = require('../../src/app');

describe('Health Controller', () => {
  describe('GET /api/health', () => {
    it('should return health information', async () => {
      const response = await request(app)
        .get('/api/health')
        .expect(200);

      expect(response.body).toMatchObject({
        status: 'healthy',
        service: '{{ api_name }}',
        version: '{{ version }}'
      });

      expect(response.body).toHaveProperty('timestamp');
      expect(response.body).toHaveProperty('uptime');
    });
  });
});
EOF
{% endif %}
```

### Documentation and Configuration Files

```bash
# README.md
cat > README.md << 'EOF'
# {{ api_name }}

{{ description }}

## Features

- RESTful API with Express.js
- {% if database == "mongodb" %}MongoDB{% elif database == "postgresql" %}PostgreSQL{% elif database == "mysql" %}MySQL{% elif database == "sqlite" %}SQLite{% endif %} database integration
{% if include_auth %}- JWT authentication
{% endif %}{% if include_tests %}- Comprehensive test suite
{% endif %}- ESLint for code quality
- Environment-based configuration
- CORS and security headers
- Error handling middleware

## Prerequisites

- Node.js (v16 or higher)
- npm or yarn
{% if database == "mongodb" %}
- MongoDB
{% elif database == "postgresql" %}
- PostgreSQL
{% elif database == "mysql" %}
- MySQL
{% endif %}

## Installation

```bash
# Clone/copy the project
cd {{ project_slug }}

# Install dependencies
npm install

# Copy environment file
cp .env.example .env

# Configure your .env file with appropriate values
```

## Configuration

Edit the `.env` file with your specific configuration:

```env
PORT={{ port }}
{% if database == "mongodb" %}
MONGODB_URI=mongodb://localhost:27017/{{ project_slug }}
{% elif database == "postgresql" %}
DB_HOST=localhost
DB_NAME={{ project_slug }}
DB_USER=your_username
DB_PASSWORD=your_password
{% elif database == "mysql" %}
DB_HOST=localhost
DB_NAME={{ project_slug }}
DB_USER=your_username
DB_PASSWORD=your_password
{% elif database == "sqlite" %}
DB_PATH=./{{ project_slug }}.sqlite
{% endif %}
{% if include_auth %}
JWT_SECRET=your-secret-key
{% endif %}
```

## Running the Application

```bash
# Development mode with auto-reload
npm run dev

# Production mode
npm start
```

The API will be available at `http://localhost:{{ port }}`

## API Endpoints

### Health Check
- `GET /` - Welcome message
- `GET /api/health` - Health status

### Example Routes (TODO: Implement)
- `GET /api/users` - Get all users
- `POST /api/users` - Create user
- `GET /api/users/:id` - Get user by ID
- `PUT /api/users/:id` - Update user
- `DELETE /api/users/:id` - Delete user

{% if include_auth %}
### Authentication
- `POST /api/auth/login` - User login (TODO: Implement)
- `POST /api/auth/register` - User registration (TODO: Implement)
- `GET /api/protected` - Protected route example
{% endif %}

{% if include_tests %}
## Testing

```bash
# Run tests
npm test

# Run tests in watch mode
npm run test:watch

# Run tests with coverage
npm run test:coverage
```
{% endif %}

## Code Quality

```bash
# Lint code
npm run lint

# Fix linting issues
npm run lint:fix
```

## Project Structure

```
{{ project_slug }}/
├── src/
│   ├── controllers/     # Route controllers
│   ├── models/         # Database models
│   ├── routes/         # API routes
│   ├── middleware/     # Custom middleware
│   ├── utils/          # Utility functions
│   ├── config/         # Configuration files
│   └── app.js          # Application entry point
{% if include_tests %}├── tests/              # Test files
{% endif %}├── docs/               # Documentation
├── .env.example        # Environment template
├── package.json        # Dependencies and scripts
└── README.md          # This file
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run the test suite
6. Submit a pull request

## License

{{ license }}

## Author

{{ author }}
EOF

# ESLint configuration
cat > .eslintrc.js << 'EOF'
module.exports = {
  env: {
    node: true,
    es2021: true,
{% if include_tests %}
    jest: true
{% endif %}
  },
  extends: [
    'eslint:recommended'
  ],
  parserOptions: {
    ecmaVersion: 12,
    sourceType: 'module'
  },
  rules: {
    'no-console': 'warn',
    'no-unused-vars': 'error',
    'no-undef': 'error',
    'prefer-const': 'error',
    'no-var': 'error'
  }
};
EOF

# .gitignore
cat > .gitignore << 'EOF'
# Dependencies
node_modules/
npm-debug.log*
yarn-debug.log*
yarn-error.log*

# Environment variables
.env
.env.local
.env.development.local
.env.test.local
.env.production.local

# Logs
logs
*.log
loglevel

# Runtime data
pids
*.pid
*.seed
*.pid.lock

# Coverage directory used by tools like istanbul
coverage/
*.lcov

# nyc test coverage
.nyc_output

# Database
{% if database == "sqlite" %}
*.sqlite
*.sqlite3
*.db
{% endif %}

# IDE files
.vscode/
.idea/
*.swp
*.swo

# OS generated files
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db

# Build output
dist/
build/

# Optional npm cache directory
.npm

# Optional eslint cache
.eslintcache
EOF
```

## Step 3: Template Testing

Before adding your template to the registry, it's crucial to test it thoroughly:

### Validate Template Structure

```bash
# Test the template with validation
ason validate ./nodejs-api-template --strict

# Fix any issues found
ason validate ./nodejs-api-template --fix
```

### Test with Different Variable Combinations

```bash
# Test minimal variables
ason new ./nodejs-api-template test-minimal \
  --var api_name="Test API" \
  --var author="Test Author" \
  --dry-run

# Test with all options
ason new ./nodejs-api-template test-full \
  --var api_name="Complete Test API" \
  --var author="Test Author" \
  --var description="A comprehensive test API" \
  --var database="postgresql" \
  --var include_auth="true" \
  --var include_tests="true" \
  --dry-run

# Test different database options
for db in mongodb postgresql mysql sqlite; do
  echo "Testing with $db database..."
  ason new ./nodejs-api-template "test-$db" \
    --var api_name="$db API" \
    --var author="Test" \
    --var database="$db" \
    --dry-run
done
```

### Verify Generated Projects

```bash
# Generate a real project for testing
ason new ./nodejs-api-template test-api \
  --var api_name="Test API" \
  --var author="Your Name" \
  --var database="sqlite" \
  --var include_tests="true"

# Verify the structure
cd test-api
ls -la
cat package.json

# Test if it works
npm install
npm test
npm run lint
```

## Step 4: Adding to Registry

Once your template is thoroughly tested:

```bash
# Add to your registry
ason add nodejs-api ./nodejs-api-template \
  --type backend \
  --description "Professional Node.js REST API with Express, testing, and best practices"

# Verify it was added
ason list

# Test from registry
ason new nodejs-api my-new-api \
  --var api_name="My New API" \
  --var author="Your Name"
```

## Advanced Template Techniques

### Dynamic File Names

```bash
# Create files with dynamic names
echo "// {{ service_name }} configuration" > "config/{{ service_name | lower }}.config.js"

# Dynamic directory names
mkdir -p "src/{{ module_name }}"
```

### Conditional File Generation

Use Pongo2 conditional logic to include/exclude entire files:

```pongo2
{% if include_docker %}
# Dockerfile content here
FROM node:18-alpine
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production
COPY . .
EXPOSE {{ port }}
CMD ["npm", "start"]
{% endif %}
```

### Complex Variable Logic

```pongo2
{# Advanced variable manipulation #}
{{ project_name | title | replace(" ", "") }}Service

{# Conditional defaults #}
{% if database == "mongodb" %}
  {% set db_port = "27017" %}
{% elif database == "postgresql" %}
  {% set db_port = "5432" %}
{% elif database == "mysql" %}
  {% set db_port = "3306" %}
{% endif %}

{# Loops for repetitive content #}
{% for route in api_routes %}
router.{{ route.method | lower }}('{{ route.path }}', {{ route.handler }});
{% endfor %}
```

### Template Inheritance

Create base templates and extend them:

```bash
# base-api-template/ason.yaml
extends: "base-template"
variables:
  - name: custom_var
    extends_from: "base"
```

## Best Practices

### 1. Template Organization

```bash
# Good structure
my-template/
├── ason.toml              # Always include configuration
├── README.md              # Document your template
├── .env.example           # Environment template
├── src/                   # Logical code organization
├── tests/                 # Include tests
├── docs/                  # Additional documentation
└── scripts/               # Utility scripts
```

### 2. Variable Design

```toml
# Good variable design
[[variables]]
name = "project_name"
description = "Human-readable project name"
required = true
example = "My Awesome Project"

[[variables]]
name = "project_slug"
description = "URL-safe project identifier"
required = false
default = "{{ project_name | lower | replace(' ', '-') }}"
example = "my-awesome-project"

[[variables]]
name = "database_type"
description = "Database system to use"
required = false
default = "postgresql"
options = ["postgresql", "mysql", "mongodb", "sqlite"]
example = "postgresql"
```

### 3. Documentation

```markdown
# Template README.md should include:

## Overview
What does this template create?

## Variables
List all variables with descriptions and examples

## Prerequisites
What needs to be installed?

## Usage Examples
Show common usage patterns

## Structure
Explain the generated project structure

## Next Steps
What to do after generation
```

### 4. Error Prevention

```toml
# Validate variable formats
[[variables]]
name = "port"
type = "number"
validate = "^[0-9]{2,5}$"

[[variables]]
name = "email"
validate = "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
```

## Troubleshooting Templates

### Common Issues

1. **Syntax Errors**
   ```bash
   # Test template syntax
   ason validate my-template --check templates
   ```

2. **Missing Variables**
   ```bash
   # Check variable usage
   ason validate my-template --check variables
   ```

3. **File Permissions**
   ```bash
   # Check file permissions in template
   find my-template -type f ! -perm 644
   ```

4. **Large Templates**
   ```bash
   # Check template size
   du -sh my-template
   ```

### Testing Strategies

```bash
# Test script for templates
#!/bin/bash
template_dir="$1"
test_output="test-output"

echo "Testing template: $template_dir"

# Validate template
ason validate "$template_dir" --strict || exit 1

# Test generation with minimal variables
ason new "$template_dir" "$test_output" \
  --var project_name="Test Project" \
  --var author="Test Author" \
  --dry-run || exit 1

# Test generation with full variables
ason new "$template_dir" "$test_output" \
  --var project_name="Test Project" \
  --var author="Test Author" \
  --var description="Test Description" \
  --var version="0.1.0" || exit 1

# Verify generated project
cd "$test_output"
if [ -f "package.json" ]; then
  npm install || exit 1
  npm test || exit 1
fi

echo "Template test completed successfully!"
```

## Next Steps

Now that you've created your first comprehensive template:

1. **[Template Registry Guide](../guides/registry.md)** - Learn to manage multiple templates
2. **[Variable Systems Guide](../guides/variables.md)** - Master advanced variable techniques
3. **[Advanced Templating Guide](../guides/advanced-templating.md)** - Explore Pongo2 features
4. **[Examples](../examples/)** - Study real-world template examples

---

*You have learned the sacred art of template creation. Now go forth and scaffold the future! 🪇*