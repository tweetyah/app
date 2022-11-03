<script lang="ts">
  import { onMount } from "svelte";
  import { navigateTo } from "svelte-router-spa";
  import { setAuth } from '../store'

  onMount(async () => {
    let query = location.search
    query = query.replace("?", "")
    let spl = query.split("&")
    let map = {}
    spl.forEach(kvp => {
      map[kvp.split("=")[0]] = kvp.split("=")[1]
    })

    let res = await fetch("/.netlify/functions/auth", {
      method: "post",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({
        code: map["code"]
      })
    })
    let json = await res.json()

    localStorage.setItem("auth", JSON.stringify(json))
    setAuth(json)

    navigateTo("/")
  })
</script>

<div>auth!</div>