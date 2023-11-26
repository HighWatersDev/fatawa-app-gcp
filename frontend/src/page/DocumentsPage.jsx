import React, {useEffect, useState} from 'react';
import axios from 'axios';
import {auth} from '../firebase';
import { withAuth } from '../components/withAuth';
import Header from '../components/Header';

function DocumentsPage() {
    const [title, setTitle] = useState('');
    const [author, setAuthor] = useState('');
    const [question, setQuestion] = useState('');
    const [answer, setAnswer] = useState('');
    const [searchQuery, setSearchQuery] = useState('');
    const [searchResults, setSearchResults] = useState([]);
    const [fatwaId, setFatwaId] = useState('');
    const [fatwaResult, setFatwaResult] = useState([]);
    const [fatwas, setFatwas] = useState([]);
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
            const response = await axios.get(`http://localhost:8080/v1/documents/search?query=${searchQuery}`, { headers });
            const documents = response.data && response.data['documents'];
            setSearchResults(documents || []);
            // setSearchResults(response.data['documents']);
            setErrorMessage('');
        } catch (error) {
            setSearchResults([]);
            setErrorMessage(error.message);
        }
    };

    const handleSearchByIdSubmit = async (event) => {
        event.preventDefault();
        try {
            const token = await auth.currentUser.getIdToken();
            const headers = {
                'Authorization': `Bearer ${token}`
            };
            const response = await axios.get(`http://localhost:8080/v1/documents/${fatwaId}`, { headers });
            setFatwaResult([response.data]);
            setErrorMessage('');
        } catch (error) {
            setFatwaResult([]);
            setErrorMessage(error.message);
        }
    };

    const fetchAllDocuments = async () => {
        try {
            const token = await auth.currentUser.getIdToken();
            const headers = {
                'Authorization': `Bearer ${token}`
            };
            const response = await axios.get(`http://localhost:8080/v1/documents/all`, { headers });
            setFatwas(response.data || []);
            setErrorMessage('');
        } catch (error) {
            setFatwas([]);
            setErrorMessage('Error fetching documents. Please try again later.');
        }
    };

    return (
        <div>
            <Header/>
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
                <br/>
                <label>
                    Author:
                    <input
                        type="text"
                        value={author}
                        onChange={(event) => setAuthor(event.target.value)}
                    />
                </label>
                <br/>
                <label>
                    Question:
                    <input
                        type="text"
                        value={question}
                        onChange={(event) => setQuestion(event.target.value)}
                    />
                </label>
                <br/>
                <label>
                    Answer:
                    <input
                        type="text"
                        value={answer}
                        onChange={(event) => setAnswer(event.target.value)}
                    />
                </label>
                <br/>
                <button type="submit">Submit</button>
            </form>
            {successMessage && <p>{successMessage}</p>}
            {errorMessage && <p>{errorMessage}</p>}
            <hr/>
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
                <br/>
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
                            <p>Id: {result.id}</p>
                            <p>Author: {result.author}</p>
                            <p>Question: {result.question}</p>
                            <p>Answer: {result.answer}</p>
                        </li>
                    ))}
                </ul>
            )}
            <hr/>
            <h2>Get Fatwa by Id</h2>
            <form onSubmit={handleSearchByIdSubmit}>
                <label>
                    Fatwa Id:
                    <input
                        type="text"
                        value={fatwaId}
                        onChange={(event) => setFatwaId(event.target.value)}
                    />
                </label>
                <br/>
                <button type="submit">Get Fatwa</button>
            </form>
            {errorMessage && <p>{errorMessage}</p>}
            {fatwaResult.length === 0 ? (
                <p>No fatwa with this id found</p>
            ) : (
                <ul>
                    {fatwaResult.map((result) => (
                        <li key={result.id}>
                            <h3>{result.title}</h3>
                            <p>Author: {result.author}</p>
                            <p>Question: {result.question}</p>
                            <p>Answer: {result.answer}</p>
                        </li>
                    ))}
                </ul>
            )}
                <h1>All Fatawa</h1>
                <button onClick={fetchAllDocuments}>Fetch All Fatawa</button>
                {errorMessage && <p>{errorMessage}</p>}
                <div style={{maxHeight: '200px', overflowY: 'auto', border: '1px solid #ccc', padding: '10px'}}>
                    <ul>
                        {fatwas.map((fatwa) => (
                            <li key={fatwa.id}>
                                <h3>{fatwa.title}</h3>
                                <p>Id: {fatwa.id}</p>
                                <p>Author: {fatwa.author}</p>
                                <p>Question: {fatwa.question}</p>
                                <p>Answer: {fatwa.answer}</p>
                            </li>
                        ))}
                    </ul>
                </div>
        </div>
    );
}

export default withAuth(DocumentsPage);