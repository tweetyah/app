<script lang="ts">
  export let title: string
  export let onClick:Function
  export let icon: string = ""
</script>
  

<button class="bg-slate-800 rounded p-1 text-slate-50 text-sm hover:bg-slate-600" on:click={() => onClick()}>
  {#if icon !== ""}
    <i class="bx {icon}"></i>
  {/if}
  { title }
</button>