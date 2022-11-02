<script lang="ts">
  import Button from "./Button.svelte";
import ComposerCard from "./ComposerCard.svelte";

  let tweets = [{
    content: ""
  }]

  function addTweet() {
    tweets = [...tweets, {
      content: ""
    }]
  }
</script>

<div>
  {JSON.stringify(tweets)}
  <div class="flex">
    <div id="composer-wrapper" class="flex-1">
      <div class="bg-slate-100 rounded mb-2">
        {#each tweets as t, idx}
          <ComposerCard 
            bind:tweet={t} 
            index={idx} 
            total={tweets.length} />
        {/each}
      </div>
      <Button onClick={() => addTweet()} icon="bx-list-plus" title="Add tweet" />
    </div>
    <div id="composer-preview" class="flex-1">
      <div>
        Send at scheduler
      </div>
      <div>
        Retweet at scheduler
      </div>
      <div>
        Category selector
      </div>
      <div>
        Save to library opt
      </div>
    </div>
  </div>
</div>