<script>
    //import { useState } from 'svelte';
    import { userAuth } from '$lib/user.ts';
    import { auth } from '$lib/auth.js';
    import { writable } from 'svelte/store';

    //const { user } = useUser();
    const title = writable('');
    const author = writable('');
    const question = writable('');
    const answer = writable('');

    async function handleSubmit(event) {
        event.preventDefault();
        const newDocument = {
            title: $title,
            author: $author,
            question: $question,
            answer: $answer
        };

        const token = auth.currentUser.getIdToken();
        const response = await fetch('/v1/documents', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({newDocument})
        });

        if (response.ok) {
            console.log('Document created successfully');
        } else {
            console.error('Error creating document');
        }

        title.set('');
        author.set('');
        question.set('');
        answer.set('');
    }

</script>

{#if $userAuth}
    <form on:submit={handleSubmit}>
        <label>
            Title:
            <input type="text" on:input={e => title.set(e.target.value)} value={$title} />
        </label>
        <label>
            Author:
            <input type="text" on:input={e => author.set(e.target.value)} value={$author} />
        </label>
        <label>
            Question:
            <input type="text" on:input={e => question.set(e.target.value)} value={$question} />
        </label>
        <label>
            Answer:
            <input type="text" on:input={e => answer.set(e.target.value)} value={$answer} />
        </label>
        <button type="submit">Create Document</button>
    </form>
{:else}
    <p>Please sign in to create a document</p>
{/if}
