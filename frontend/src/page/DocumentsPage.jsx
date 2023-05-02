import React, {useEffect, useState} from 'react';
import axios from 'axios';
import {auth} from '../firebase';
import { withAuth } from '../components/withAuth';

function DocumentsPage() {
    const [title, setTitle] = useState('');
    const [author, setAuthor] = useState('');
    const [question, setQuestion] = useState('');
    const [answer, setAnswer] = useState('');
    const [searchQuery, setSearchQuery] = useState('');
    const [searchResults, setSearchResults] = useState([]);
    const [successMessage, setSuccessMessage] = useState('');
    const [errorMessage, setErrorMessage] = useState('');
    const handleSubmit = async (event) => {
        event.preventDefault();
        const token = await auth.currentUser.getIdToken();
        const headers = {
            'Authorization': `Bearer ${token}`
        };
        const data = {
            title,
            author,
            question,
            answer
        };
        try {
            const response = await axios.post('http://localhost:8080/v1/documents', data, { headers });
            console.log(response.data);
            setSuccessMessage(`Fatwa with id ${response.data.document_id} created successfully`);
            setErrorMessage('');
        } catch (error) {
            console.error(error);
            setErrorMessage('Failed to create fatwa');
            setSuccessMessage('');
        }
    };

    const handleSearchSubmit = async (event) => {
        event.preventDefault();
        try {
            const token = await auth.currentUser.getIdToken();
            const headers = {
                'Authorization': `Bearer ${token}`
            };
            const response = await axios.get(`http://localhost:8080/v1/documents/search?search=${searchQuery}`, { headers });
            setSearchResults(response.data['documents']);
            setErrorMessage('');
        } catch (error) {
            setSearchResults([]);
            setErrorMessage(error.message);
        }
    };

    return (
        <div>
            <h2>Submit a New Fatwa</h2>
            <form onSubmit={handleSubmit}>
                <label>
                    Title:
                    <input
                        type="text"
                        value={title}
                        onChange={(event) => setTitle(event.target.value)}
                    />
                </label>
                <br />
                <label>
                    Author:
                    <input
                        type="text"
                        value={author}
                        onChange={(event) => setAuthor(event.target.value)}
                    />
                </label>
                <br />
                <label>
                    Question:
                    <input
                        type="text"
                        value={question}
                        onChange={(event) => setQuestion(event.target.value)}
                    />
                </label>
                <br />
                <label>
                    Answer:
                    <input
                        type="text"
                        value={answer}
                        onChange={(event) => setAnswer(event.target.value)}
                    />
                </label>
                <br />
                <button type="submit">Submit</button>
            </form>
            {successMessage && <p>{successMessage}</p>}
            {errorMessage && <p>{errorMessage}</p>}
            <hr />
            <h2>Search for Fatwa</h2>
            <form onSubmit={handleSearchSubmit}>
                <label>
                    Search Query:
                    <input
                        type="text"
                        value={searchQuery}
                        onChange={(event) => setSearchQuery(event.target.value)}
                    />
                </label>
                <br />
                <button type="submit">Search</button>
            </form>
            {errorMessage && <p>{errorMessage}</p>}
            {searchResults.length === 0 ? (
                <p>No search results</p>
            ) : (
                <ul>
                    {searchResults.map((result) => (
                        <li key={result.id}>
                            <h3>{result.title}</h3>
                            <p>Author: {result.author}</p>
                            <p>Question: {result.question}</p>
                            <p>Answer: {result.answer}</p>
                        </li>
                    ))}
                </ul>
            )}
        </div>
    );
}

export default withAuth(DocumentsPage);