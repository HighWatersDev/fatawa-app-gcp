<div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
        <div>
            <img class="mx-auto h-20 w-auto">
            <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
                Sign in to your account
            </h2>
        </div>
        <form class="mt-8 space-y-6" on:submit|preventDefault={login}>
            <input type="hidden" name="remember" value="true">
            <div class="rounded-md shadow-sm -space-y-px">
                <div class="mb-4">
                    <label for="email-address" class="sr-only">Email address</label>
                    <input id="email-address" name="email" type="email" autocomplete="email" required class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-t-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm" placeholder="Email address" bind:value={email}>
                </div>
                <div class="mb-4">
                    <label for="password" class="sr-only">Password</label>
                    <input id="password" name="password" type="password" autocomplete="current-password" required class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm" placeholder="Password" bind:value={password}>
                </div>
            </div>
            {#if error}
                <p class="text-red-500 text-sm">{error}</p>
            {/if}
            <div class="flex items-center justify-between">
                <div class="flex items-center">
                    <input id="remember-me" name="remember-me" type="checkbox" class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded">
                    <label for="remember-me" class="ml-2 block text-sm text-gray-900">
                        Remember me
                    </label>
                </div>

                <div class="text-sm">
                    <a href="#" class="font-medium text-indigo-600 hover:text-indigo-500">
                        Forgot your password?
                    </a>
                </div>
            </div>

            <div>
                <button type="submit" class="mt-4 w-full inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
                    Sign in
                </button>
            </div>
        </form>
    </div>
</div>


<script>
	import { signInWith } from './../../lib/auth.js';
    import { onMount } from 'svelte';
    import { onAuthStateChanged } from '$lib/auth';
    import { writable } from 'svelte/store';
    import {goto} from "$app/navigation";

    export const user = writable(null);
    const apiUrl = 'http://localhost:8080';

    let email = '';
    let password = '';
    let error = '';
    let endpoint = '/v1/verify';

    async function login() {
        try {
            const userCredential = await signInWith(email, password);
            const idToken = await userCredential.getIdToken(); // user.idToken;
            console.log(idToken);
            const response = await fetch(`${apiUrl}${endpoint}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${idToken}`
                },
                body: JSON.stringify(idToken)
            });
            return await response.text();

        } catch (err) {
            error = err.message;
        }
    }

    onMount(() => {
        onAuthStateChanged(user => {
            if (user) {
                console.log('User is signed in.');
                console.log(user.email);
                // goto("/home");
            } else {
                console.log('User is signed out.');
            }
        });
    });
</script>



<!--    <script>-->
<!--    import { initializeApp } from 'firebase/app'-->
<!--    import { getAuth, signInWithEmailAndPassword } from 'firebase/auth'-->
<!--    import { writable } from 'svelte/store';-->
<!--    import {goto} from "$app/navigation";-->

<!--    export const user = writable(null);-->
<!--    const apiUrl = 'http://localhost:8080';-->

<!--    const firebaseConfig = {-->
<!--        apiKey: import.meta.env.VITE_FIREBASE_API_KEY, //"AIzaSyCpx&#45;&#45;bLjWj1IbARdTboBknNpiXXFftBjQ",-->
<!--        authDomain: import.meta.env.VITE_FIREBASE_AUTH_DOMAIN,-->
<!--        projectId: import.meta.env.VITE_FIREBASE_PROJECT_ID,-->
<!--        storageBucket: import.meta.env.VITE_FIREBASE_STORAGE_BUCKET,-->
<!--        messagingSenderId: import.meta.env.VITE_FIREBASE_MESSAGING_SENDER_ID, //"1055722927890",-->
<!--        appId: import.meta.env.VITE_FIREBASE_APP_ID //"1:1055722927890:web:9cd295b8d79ac32e03a5cb"-->
<!--    };-->

<!--    // Initialize Firebase-->
<!--    const app = initializeApp(firebaseConfig);-->
<!--    const auth = getAuth(app);-->

<!--    let email = '';-->
<!--    let password = '';-->
<!--    let error = '';-->
<!--    let endpoint = '/v1/login';-->

<!--    async function login() {-->
<!--        try {-->
<!--            const userCredential = await signInWithEmailAndPassword(auth, email, password);-->
<!--            const idToken = await userCredential.user.getIdToken();-->
<!--            await goto("/home");-->
<!--            const response = await fetch(`${apiUrl}${endpoint}`, {-->
<!--                method: 'POST',-->
<!--                headers: {-->
<!--                    'Content-Type': 'application/json',-->
<!--                    'Authorization': `Bearer ${idToken}`-->
<!--                },-->
<!--                body: JSON.stringify(idToken)-->
<!--            });-->
<!--            return await response.json(), userCredential.user.email;-->

<!--        } catch (err) {-->
<!--            error = err.message;-->
<!--        }-->
<!--    }-->

<!--</script>-->
