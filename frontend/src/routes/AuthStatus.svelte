<!--<script lang="ts">-->
<!--  import { userAuth } from "$lib/user";-->
<!--  import { page } from '$app/stores';-->
<!--  import { auth } from '$lib/auth';-->
<!--  console.log($page)-->
<!--</script>-->

<!--<div class="my-4">-->
<!--  {#if $page.data.userSession}-->
<!--    <button>Sign Out</button>-->
<!--    {$page.data.userSession.name} ({$page.data.userSession.email})-->
<!--  {:else}-->
<!--    <button>Sign In with Google</button>-->
<!--    (visitor)-->
<!--  {/if}-->
<!--  {#if $userAuth}-->
<!--    <div class="text-gray-400">-->
<!--      Authenticated as <span class="font-bold text-black">{$userAuth.email}</span>-->
<!--    </div>-->
<!--  {:else}-->
<!--    <div class="text-gray-400">-->
<!--      No one is authenticated-->
<!--    </div>-->
<!--  {/if}-->
<!--</div>-->

<script>
  import { user, signOut, signIn } from '$lib/stores/auth';
  import { page } from '$app/stores';
  console.log($page)
</script>

{#if $user}
  {$user.email}
  <button on:click={() => signOut()}>Log out</button>
{:else}
  Not logged in.
  <button on:click={() => signIn()}>Log in</button>
{/if}