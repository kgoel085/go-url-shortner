```mermaid
sequenceDiagram
    participant Visitor
    participant API
    participant DB

    Visitor->>API: GET /{code}
    API->>DB: Lookup Short URL by code
    API->>DB: Log Analytics
    API-->>Visitor: Redirect to original URL
```