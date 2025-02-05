<script lang="ts">
	import { AppBar } from "@skeletonlabs/skeleton-svelte";
    import SettingsButton from "./SettingsButton.svelte";

</script>

<AppBar>
	{#snippet lead()}
		<h4 class="h3 h-9 flex items-center">Sego</h4>
	{/snippet}
	{#snippet trail()}
		<SettingsButton showMenu={true}/>
	{/snippet}
</AppBar>
