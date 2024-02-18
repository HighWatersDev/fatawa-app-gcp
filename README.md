# Salafi Fatawa App

# API Documentation

### Overview

This document outlines the usage of the API endpoints for creating, searching, and deleting documents within the salafifatawa Firestore collection. These documents can be either parent documents or child documents (split audios) associated with a parent document.

### Authentication

All API requests require the use of a generated API key, which should be included in the request headers. Replace <Your_API_Key> with your actual API key.


```plaintext
Authorization: Bearer <Your_API_Key>
```

### Endpoints

#### Create Document

- URL: /documents

- Method: POST

- Description: Creates a new parent document or a child document under an existing parent, based on the provided JSON payload.

- Content-Type: application/json

- Request Body:

  - For creating a parent document:

    ```json
    {
    "isSplitAudio": false,
    "parentDocument": {
    "audio": "path/to/parent/audio.mp3",
    "author": "Author Name",
    "complete": true
    }
    }
    ```

  - For creating a child document:

    ```json
    {
    "isSplitAudio": true,
    "parentDocID": "<Parent_Document_ID>",
    "document": {
    "audio": "path/to/child/audio.mp3",
    "title": "Child Title",
    "topic": "Child Topic",
    "author": "Author Name",
    "question": "Child Question",
    "answer": "Child Answer"
    }
    }
    ```
    
- Example cURL Request:

  - For a parent document:

    ```bash
    curl -X POST "http://your-api-hostname.com/documents" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer <Your_API_Key>" \
    -d '{"isSplitAudio": false, "parentDocument": {"audio": "path/to/parent/audio.mp3", "author": "Author Name", "complete": true}}'
    ```
    
  - For a child document:

    ```bash
    curl -X POST "http://your-api-hostname.com/documents" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer <Your_API_Key>" \
    -d '{"isSplitAudio": true, "parentDocID": "<Parent_Document_ID>", "document": {"audio": "path/to/child/audio.mp3", "title": "Child Title", "topic": "Child Topic", "author": "Author Name", "question": "Child Question", "answer": "Child Answer"}}'
    ```
    
#### Search Documents

- URL: /documents/search

- Method: POST

- Description: Searches for documents based on the provided search query. Can search within all parent documents or within the child documents of a specified parent.

- Content-Type: application/json

- Request Body:

```json
{
"query": "test",
"isSplitAudio": true,
"parentDocID": "<Parent_Document_ID>"
}
```

- Example cURL Request:

```bash
curl -X POST "http://your-api-hostname.com/documents/search" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer <Your_API_Key>" \
-d '{"query": "test", "isSplitAudio": true, "parentDocID": "<Parent_Document_ID>"}'
```

#### Delete Document

- URL: /documents/{docID}

- Method: DELETE

- Description: Deletes a specified document. To delete a child document, include the parentDocID as a query parameter.

- Query Parameters:

  - parentDocID (optional): The ID of the parent document, required only when deleting a child document.

- Example cURL Request:

  - Deleting a parent document:

    ```bash
    curl -X DELETE "http://your-api-hostname.com/documents/<Doc_ID>" \
    -H "Authorization: Bearer <Your_API_Key>"
    Deleting a child document:
    ```
    
    ```bash
    curl -X DELETE "http://your-api-hostname.com/documents/<Child_Doc_ID>?parentDocID=<Parent_Doc_ID>" \
    -H "Authorization: Bearer <Your_API_Key>"
    ```

### Notes

Replace http://your-api-hostname.com with your actual API hostname.
Ensure to replace placeholder values like <Your_API_Key>, <Parent_Document_ID>, <Doc_ID>, and <Child_Doc_ID> with actual values appropriate for your requests.
