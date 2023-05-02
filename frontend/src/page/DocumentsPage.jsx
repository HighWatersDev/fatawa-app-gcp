import React, {useEffect, useState} from 'react';
import axios from 'axios';
import {auth} from '../firebase';

function withAuth(Component) {
    return function AuthenticatedComponent(props) {
        const [user, setUser] = useState(null);

        useEffect(() => {
            return auth.onAuthStateChanged((user) => {
                setUser(user);
            });
        }, []);

        if (!user) {
            return <h1>403 - Access Forbidden</h1>;
        }

        return <Component {...props} />;
    };
}

function DocumentsPage() {
    const [title, setTitle] = useState('');
    const [author, setAuthor] = useState('');
    const [question, setQuestion] = useState('');
    const [answer, setAnswer] = useState('');
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
        } catch (error) {
            console.error(error);
        }
    };

    return (
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
    );
}

export default withAuth(DocumentsPage);