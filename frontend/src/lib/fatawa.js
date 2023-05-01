import { writable } from 'svelte/store';
import { auth } from './auth.js';

export const createDocument = async (fields) => {
    const token = auth.currentUser.getIdToken();
    const response = await fetch('http://localhost:8080/v1/documents', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(fields),
    });

    if (!response.ok) {
        const { message } = await response.json();
        throw new Error(message);
    }
};

export const useDocumentForm = () => {
    const title = writable('');
    const author = writable('');
    const question = writable('');
    const answer = writable('');

    const handleSubmit = async (event) => {
        event.preventDefault();

        const fields = {
            title,
            author,
            question,
            answer,
        };

        try {
            await createDocument(fields);
            console.log('Document created successfully');
        } catch (error) {
            console.error(error.message);
        }
    };

    return {
        title,
        author,
        question,
        answer,
        handleSubmit,
    };
};
