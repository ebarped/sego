<script lang="ts">
	import { Autocomplete } from "@skeletonlabs/skeleton";
	import Results from "../lib/Results.svelte";

	let inputSubmitted = false;
	let rerender = false;
	let query = "";

	// handle ENTER press
	const sendInput = (e) => {
		if (e.charCode === 13) {
			query = e.target.value;
			console.log("<MAIN> Searching: " + query);
			inputSubmitted = true;
			rerender = !rerender;
		}
	};
</script>

<!-- outer container, align all to center -->
<div class="container h-full mx-auto flex justify-center items-start py-24">
	<!-- number of columns -->
	<div class="grid grid-cols-1">
		<h1 class="text-center py-2">Linux API Documentation Search Engine</h1>
		<!-- hr draws a line -->
		<hr />
		<!-- vertical padding from the previous component -->
		<div class="grid py-3">
			<!-- center and put maximum size -->
			<div class="container mx-auto justify-center max-w-2xl">
				<input
					class="input"
					type="text"
					placeholder="Search"
					on:keypress={sendInput}
				/>
			</div>

			<!-- show results once the enter key has been sent -->
			{#if inputSubmitted}
				<!-- rerender when query changes & enter is sent-->
				{#key rerender}
					<Results {query} />
				{/key}
			{/if}
		</div>
	</div>
</div>
