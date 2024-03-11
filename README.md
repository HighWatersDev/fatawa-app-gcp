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
    "parentDocID": "<Parent_Doc_ID>",
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

#### Get All Documents

- URL: /documents/all

- Method: GET

- Description: Get all parent documents or all child documents when parentDocID is provided.

- Content-Type: application/json

- Request Body:

  - For parent docs:

  ```json
  { }
  ```

  - Example cURL Request:

  ```bash
  curl -X POST "http://your-api-hostname.com/documents/all" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <Your_API_Key>" \
  -d '{}'
  ```
  
  - For child docs:
 
  ```json
  {
  "parentDocID": "<Parent_Doc_ID"
  }
  ```
  
  - Example cURL Request:

  ```bash
  curl -X POST "http://your-api-hostname.com/documents/all" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <Your_API_Key>" \
  -d '{"parentDocID": "Parent_Doc_ID"}'
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

# Python FastAPI Fatawa Processor

## Endpoints

### Upload Files

- **Endpoint**: `/upload`
- **Method**: `POST`
- **Auth Required**: Yes
- **Description**: Uploads files to Azure storage.
- **Parameters**:
  - `path`: Path where the file will be uploaded.
  - `author`: Name of the file's author.

**Example**:
```sh
curl -X POST "http://localhost:8000/upload"      -H "Authorization: Bearer YOUR_TOKEN"      -d "path=/path/to/upload&author=JohnDoe"
```

### List Blob

- **Endpoint**: `/list`
- **Method**: `GET`
- **Auth Required**: Yes
- **Description**: Lists files in a specified blob path.
- **Parameters**:
  - `path`: Blob path to list files from.

**Example**:
```sh
curl "http://localhost:8000/list?path=/myblob/path"      -H "Authorization: Bearer YOUR_TOKEN"
```

### Download

- **Endpoint**: `/download`
- **Method**: `GET`
- **Auth Required**: Yes
- **Description**: Downloads a file from storage.
- **Parameters**:
  - `file_path`: Path of the file to download.
  - `author`: Author of the file.

**Example**:
```sh
curl "http://localhost:8000/download?file_path=/path/to/file&author=JohnDoe"      -H "Authorization: Bearer YOUR_TOKEN"
```

### Transcribe

- **Endpoint**: `/transcribe`
- **Method**: `GET`
- **Auth Required**: Yes
- **Description**: Transcribes audio content from a given blob.
- **Parameters**:
  - `blob`: Path to the blob containing audio to transcribe.

**Example**:
```sh
curl "http://localhost:8000/transcribe?blob=/path/to/audio"      -H "Authorization: Bearer YOUR_TOKEN"
```

### Translate

- **Endpoint**: `/translate`
- **Method**: `GET`
- **Auth Required**: Yes
- **Description**: Translates transcribed content.
- **Parameters**:
  - `blob`: Path to the blob containing transcribed content to translate.

**Example**:
```sh
curl "http://localhost:8000/translate?blob=/path/to/transcribed"      -H "Authorization: Bearer YOUR_TOKEN"
```

### Convert to ACC

- **Endpoint**: `/to_acc`
- **Method**: `GET`
- **Auth Required**: Yes
- **Description**: Converts an audio file to ACC format.
- **Parameters**:
  - `blob`: Path to the local WAV file to convert.

**Example**:
```sh
curl "http://localhost:8000/to_acc?blob=/path/to/wavfile"      -H "Authorization: Bearer YOUR_TOKEN"
```

### Convert to WAV

- **Endpoint**: `/to_wav`
- **Method**: `GET`
- **Auth Required**: Yes
- **Description**: Converts an audio file to WAV format.
- **Parameters**:
  - `blob`: Path to the file to convert.

**Example**:
```sh
curl "http://localhost:8000/to_wav?blob=/path/to/accfile"      -H "Authorization: Bearer YOUR_TOKEN"
```
