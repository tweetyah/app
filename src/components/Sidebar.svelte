<script>
  import NavLink from "./NavLink.svelte";

</script>
<div id="sidebar" class="flex flex-col justify-left h-100 w-[250px] m-2 p-2 rounded bg-slate-800 text-slate-100">
  <div class="p-3 text-lg">Tweetyah!</div>
  <hr class="border-slate-700 my-2" />
  <NavLink title="Home" icon="bx-home" to="/" />
  <NavLink title="Library" icon="bx-collection" to="/library" />
</div>