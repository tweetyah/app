<script lang="ts">
  import type { Tweet } from "../models";
  import { name, handle, profileImgUrl } from "../store";
  import Button from "./Button.svelte";
  import ComposerTextarea from "./ComposerTextarea.svelte";

  // private fields
  let _name: string;
  let _handle: string;
  let _profileImgUrl: string;

  // Props
  export let tweet: Tweet;
  export let index: number;
  export let total: number;

  // store
  name.subscribe(value => _name = value)
  handle.subscribe(value => _handle = value)
  profileImgUrl.subscribe(value => _profileImgUrl = value)

  // Functions
  function selectImage() {

  }
</script>

<div class="flex">
  <div>
    <img src={ _profileImgUrl } class="w-[50px] h-[50px] rounded-full m-2" alt="user profile image" />
  </div>
  <div class="flex-1 m-2">
    <div>
      <span class="font-bold">{ _name }</span>
      <span class="italic text-slate-600 text-sm">@{ _handle }</span>
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
        <!-- <Button onClick={() => selectImage()} icon="bx-image-add" title="Add image" /> -->
      </div>
    </div>
  </div>
</div>