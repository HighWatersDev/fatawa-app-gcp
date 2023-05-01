<!-- src/routes/NewDocument.svelte -->

<script>
    import { onMount } from 'svelte';
    import { goto } from '$app/navigation';
    import { post } from '$lib/api';
    import { useUser } from '$lib/user.js';

    let title = '';
    let author = '';
    let question = '';
    let answer = '';
    let errorMessage = '';
    const { user } = useUser();

    onMount(() => {
        if (!user) {
            goto('/login');
        }
    });

    async function handleSubmit(event) {
        event.preventDefault();
        errorMessage = '';
        if (!title || !author || !question || !answer) {
            errorMessage = 'All fields are required.';
            return;
        }
        try {
            const response = await post('/api/v1/documents', {
                title,
                author,
                question,
                answer
            });
            if (response.ok) {
                goto('/documents');
            } else {
                const data = await response.json();
                errorMessage = data.message || 'Failed to create document.';
            }
        } catch (error) {
            errorMessage = error.message || 'Failed to create document.';
        }
    }
</script>

<main>
    {#if errorMessage}
        <p>{errorMessage}</p>
    {/if}
    <h1>New Document</h1>
    <form on:submit={handleSubmit}>
        <label>
            Title
            <input type="text" bind:value={title} />
        </label>
        <label>
            Author
            <input type="text" bind:value={author} />
        </label>
        <label>
            Question
            <textarea bind:value={question} />
        </label>
        <label>
            Answer
            <textarea bind:value={answer} />
        </label>
        <button type="submit">Submit</button>
    </form>
</main>

<style>
    label {
        display: block;
        margin-bottom: 0.5rem;
    }
    input[type="text"], textarea {
        display: block;
        width: 100%;
        padding: 0.5rem;
        margin-bottom: 1rem;
        border: 1px solid #ccc;
        border-radius: 4px;
    }
    button[type="submit"] {
        padding: 0.5rem 1rem;
        background-color: #4caf50;
        color: white;
        border: none;
        border-radius: 4px;
        cursor: pointer;
    }
</style>
