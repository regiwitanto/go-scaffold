// Mock data untuk UI testing
const mockTemplates = [
    {
        id: "api-echo",
        name: "API with Echo",
        type: "api",
        description: "RESTful API using Echo framework"
    },
    {
        id: "api-chi",
        name: "API with Chi",
        type: "api",
        description: "RESTful API using Chi router"
    },
    {
        id: "api-gin",
        name: "API with Gin",
        type: "api",
        description: "RESTful API using Gin framework"
    },
    {
        id: "api-standard",
        name: "API with Standard Library",
        type: "api",
        description: "RESTful API using Go standard library"
    },
    {
        id: "webapp-echo",
        name: "Web App with Echo",
        type: "webapp",
        description: "Web application using Echo framework with HTML templates"
    },
    {
        id: "webapp-chi",
        name: "Web App with Chi",
        type: "webapp",
        description: "Web application using Chi router with HTML templates"
    },
    {
        id: "webapp-gin",
        name: "Web App with Gin",
        type: "webapp",
        description: "Web application using Gin framework with HTML templates"
    },
    {
        id: "webapp-standard",
        name: "Web App with Standard Library",
        type: "webapp",
        description: "Web application using Go standard library with HTML templates"
    }
];

const mockFeatures = {
    features: [
        {
            id: "access-logging",
            name: "Access Logging",
            description: "Middleware for logging all requests and responses"
        },
        {
            id: "admin-makefile",
            name: "Admin Makefile",
            description: "Makefile with common development tasks"
        },
        {
            id: "automatic-versioning",
            name: "Automatic Versioning",
            description: "Use VCS revision as version number"
        },
        {
            id: "basic-auth",
            name: "Basic Authentication",
            description: "HTTP basic authentication middleware"
        },
        {
            id: "email",
            name: "Email Support",
            description: "Helpers for sending emails via SMTP"
        },
        {
            id: "error-notifications",
            name: "Error Notifications",
            description: "Send error alerts to admin email"
        },
        {
            id: "gitignore",
            name: "Gitignore",
            description: "Common .gitignore file for Go projects"
        },
        {
            id: "live-reload",
            name: "Live Reload",
            description: "Auto-rebuild and restart during development"
        },
        {
            id: "secure-cookies",
            name: "Secure Cookies",
            description: "Signed and encrypted cookie support"
        },
        {
            id: "sql-migrations",
            name: "SQL Migrations",
            description: "Database migration tools"
        }
    ],
    premiumFeatures: [
        {
            id: "automatic-https",
            name: "Automatic HTTPS",
            description: "TLS certificate management via Let's Encrypt",
            isPremium: true
        },
        {
            id: "custom-error-pages",
            name: "Custom Error Pages",
            description: "Custom HTML pages for error responses",
            isPremium: true
        },
        {
            id: "user-accounts",
            name: "User Accounts",
            description: "User authentication and management",
            isPremium: true
        }
    ]
};
