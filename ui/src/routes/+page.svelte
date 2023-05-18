<script lang="ts">
	import { AppBar, AppShell } from "@skeletonlabs/skeleton";
	import ResultsSelector from "../lib/ResultsSelector.svelte";

	import Results from "../lib/Results.svelte";
	import Search from "../lib/Search.svelte";

	let query: string = "";
	let inputSubmitted: boolean = false;
	let rerender: boolean = false;
	let rscount: number = 5;
</script>

<!-- AppShell stores all the page -->
<AppShell>
	<svelte:fragment slot="header">
		<AppBar slotTrail="grid grid-cols-1 justify-items-end text-right">
			<svelte:fragment slot="trail">
				<ResultsSelector bind:rscount />
			</svelte:fragment>
			<h3 class="text-left">Sego</h3>
		</AppBar>
	</svelte:fragment>

	<!-- outer container, align all to center -->
	<div class="container h-full mx-auto flex justify-center items-start py-24">
		<!-- number of columns -->
		<div class="grid grid-cols-1">
			<h1 class="text-center py-2">
				Linux API Documentation Search Engine
			</h1>
			<!-- hr draws a line -->
			<hr />
			<!-- vertical padding from the previous component -->
			<div class="grid py-3">
				<!-- bind variables to propagate child changes to parent -->
				<Search bind:query bind:inputSubmitted bind:rerender />
				<!-- show results once the enter key has been sent -->
				{#if inputSubmitted}
					<!-- rerender when query changes & enter is sent-->
					{#key rerender}
						<Results {query} {rscount} />
					{/key}
				{/if}
			</div>
		</div>
	</div>
</AppShell>
