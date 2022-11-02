<script lang="ts">
  import type { Tweet } from "../models";
  import Button from "./Button.svelte";
  import ComposerTextarea from "./ComposerTextarea.svelte";

  // Props
  export let tweet: Tweet;
  export let index: number;
  export let total: number;

  // Functions
  function selectImage() {

  }
</script>

<div class="flex">
  <div>
    <img src="https://pbs.twimg.com/profile_images/1507551817646977024/H-jpBvbo_400x400.jpg" class="w-[50px] h-[50px] rounded-full m-2" />
  </div>
  <div class="flex-1 m-2">
    <div>
      <span class="font-bold">Brian Morrison II</span>
      <span class="italic text-slate-600 text-sm">@brianmmdev</span>
    </div>
    <ComposerTextarea bind:value={tweet.content} />
    <div class="flex text-sm align-center text-slate-600">
      <div class="flex-1">
        <span class="mr-2">{ tweet.content.length }/240</span>
        {#if total > 1}
          <span class="mr-2">#{ index + 1 }/{total}</span>
        {/if}
      </div>
      <div>
        <Button onClick={() => selectImage()} icon="bx-image-add" title="Add image" />
      </div>
    </div>
  </div>
</div>