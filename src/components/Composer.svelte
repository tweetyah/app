<script lang="ts">
  import Accordion from "./Accordion.svelte";
  import AccordionNode from "./AccordionNode.svelte";
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
      <div class="bg-white shadow-sm rounded mb-2">
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
      <Accordion>
        <AccordionNode title="Send at">
          <div class="px-3">
            <SendAtScheduler bind:value={sendAt} />
          </div>
        </AccordionNode>
        <AccordionNode title="Retweet at">
          Retweet at
        </AccordionNode>
        <AccordionNode title="Categories">
          Categories
        </AccordionNode>
        <AccordionNode title="Other">
          Save to library opt
        </AccordionNode>
      </Accordion>
    </div>
  </div>
</div>