<script lang="ts">
  import { auth } from './store'
  import { Router } from 'svelte-router-spa'
  import { routes } from './router'
  import { onMount } from 'svelte';

  onMount(() => {
    let authItem = localStorage.getItem("auth")
    if(authItem) {
      let authObj = JSON.parse(authItem)
      auth.set(authObj)
    }
  })
</script>

<Router {routes} />