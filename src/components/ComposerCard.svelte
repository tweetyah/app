<script lang="ts">
  import type { Tweet } from "../models";
  import ComposerTextarea from "./ComposerTextarea.svelte";

  export let tweet: Tweet;
</script>

<div class="bg-slate-100 rounded flex">
  <div>
    <img src="https://pbs.twimg.com/profile_images/1507551817646977024/H-jpBvbo_400x400.jpg" class="w-[50px] h-[50px] rounded-full m-2" />
  </div>
  <div class="flex-1 m-2">
    <div>
      <span class="font-bold">Brian Morrison II</span>
      <span class="italic text-slate-800">@brianmmdev</span>
    </div>
    <ComposerTextarea bind:value={tweet.content} />
  </div>
</div>