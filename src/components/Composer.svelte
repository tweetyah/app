<script lang="ts">
  import ComposerCard from "./ComposerCard.svelte";

  const tweets = [{
    content: ""
  }]

  function addTweet() {
    tweets.push({
      content: ""
    })
  }
</script>

<div>
  {JSON.stringify(tweets)}
  {#each tweets as t}
    <ComposerCard bind:tweet={t} />
  {/each}
  <button on:click={() => addTweet()}>Add tweet</button>
</div>