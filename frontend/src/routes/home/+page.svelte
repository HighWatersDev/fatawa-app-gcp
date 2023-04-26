<!--<script>-->
<!--    import { signOut } from 'firebase/auth';-->
<!--    import login from '../login/+page.svelte';-->

<!--    let user;-->

<!--    async function logout() {-->
<!--        await signOut(login.signOut());-->
<!--        user = null;-->
<!--    }-->

<!--    $: user = login.currentUser;-->
<!--    console.log(login.currentUser);-->
<!--</script>-->

<!--{#if user}-->
<!--    <p>Welcome, {user.email}!</p>-->
<!--    <button on:click={logout}>Logout</button>-->
<!--{:else}-->
<!--    <p>You are not logged in.</p>-->
<!--    <a href="/login">Login</a>-->
<!--    <a href="/signup">Signup</a>-->
<!--{/if}-->
