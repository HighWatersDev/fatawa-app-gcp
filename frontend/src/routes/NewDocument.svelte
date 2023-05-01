<script>
    import { useState } from 'svelte';
    import { post } from '$lib/api';

    let title = '';
    let author = '';
    let question = '';
    let answer = '';

    const submitDocument = async () => {
        const document = {
            title,
            author,
            question,
            answer
        };

        try {
            const response = await post('/documents', document);
            console.log(response);
        } catch (error) {
            console.error(error);
        }
    };
</script>

<main>
    <h1>Submit a New Document</h1>

    <form on:submit|preventDefault={submitDocument}>
        <label for="title">Title:</label>
        <input type="text" id="title" bind:value={title} />

        <label for="author">Author:</label>
        <input type="text" id="author" bind:value={author} />

        <label for="question">Question:</label>
        <textarea id="question" bind:value={question}></textarea>

        <label for="answer">Answer:</label>
        <textarea id="answer" bind:value={answer}></textarea>

        <button type="submit">Submit</button>
    </form>
</main>
