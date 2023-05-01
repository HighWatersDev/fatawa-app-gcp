<script>
    import { userAuth } from "$lib/user.ts";
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";

    let user = null;

    async function logout() {
        await fetch("/api/auth/logout", { method: "POST" });
        await goto("/");
    }

    onMount(() => {
        return userAuth.subscribe((user) => {
            if (user) {
                goto('/')
            }
        });
    });

</script>

<nav>
    <ul>
        <li><a href="/">Home</a></li>
        {#if $userAuth}
            <li><a href="/fatawa">Submit Fatawa</a></li>
            <li><button on:click={logout}>Logout</button></li>
        {:else}
            <li><a href="/login">Login</a></li>
            <li><a href="/register">Register</a></li>
        {/if}
    </ul>
</nav>

<main>
    <slot></slot>
</main>
