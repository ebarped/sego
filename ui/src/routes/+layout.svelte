<script lang="ts">
	// The ordering of these imports is critical to your app working properly
	import "@skeletonlabs/skeleton/themes/theme-crimson.css";
	// If you have source.organizeImports set to true in VSCode, then it will auto change this ordering
	import "@skeletonlabs/skeleton/styles/all.css";
	// Most of your app wide CSS should be put in this file
	import "../app.postcss";
	import { AppBar, AppShell } from "@skeletonlabs/skeleton";
</script>

<AppShell>
	<svelte:fragment slot="header">
		<AppBar>
			<h3 class="text-left">Sego</h3>
		</AppBar>
	</svelte:fragment>

	<!-- slot represents the actual content -->
	<slot />
</AppShell>
