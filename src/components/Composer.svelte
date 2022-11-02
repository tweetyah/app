<script lang="ts">
  import Button from "./Button.svelte";
  import ComposerCard from "./ComposerCard.svelte";
  import SendAtScheduler from "./SendAtScheduler.svelte";

  let sendAt: Date = new Date()
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
  {sendAt}
  {JSON.stringify(tweets)}
  <div class="grid grid-cols-2 gap-2">
    <div id="composer-wrapper">
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
    <div id="composer-preview">
      <div>
        <SendAtScheduler bind:value={sendAt} />
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