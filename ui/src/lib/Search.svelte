<script lang="ts">
	import { Progress } from '@skeletonlabs/skeleton-svelte';

	// Define the shape of the data returned from the API
    interface Result {
        query: string;
        documents: string[];
    }

	// default values for the query
	let query:string = $state("");
	let resultCount:number = $state(5);

	
	let result:Result | null = $state(null); // start as null since no data is loaded
    let loading:boolean = $state(false);
	let progressBarValue:number = $state(10);

	let error:string = $state("");
	
    async function fetchResults(): Promise<void>{
		loading = true;
		
		if (query=="") {
			error = "Invalid query: query must not be empty!"
			loading = false;
			result = null;
			console.error(error);
			return;
		} else {
			error = "";
		}

        try {
			console.log("<Results.svelte> searching: " + query);
			console.log("<Results.svelte> count: " + resultCount);
			
			for (let i = 0; i < 3; i++) {
				progressBarValue += 30;
				await sleep(800); // simulate transfer delay
			}
			progressBarValue = 10;


			const res = await fetch(
                "http://localhost:4000/search?query=" + query + "&count=" + resultCount
            );
            
            if (!res.ok) throw new Error('Failed to fetch results from Search Engine');
            
			const jsonData: Result = await res.json();
			result = jsonData;
        } catch (error) {
            console.error(error);
        } finally {
            loading = false;
        }
    }


	function sleep(ms: number): Promise<void> {
  		return new Promise(resolve => setTimeout(resolve, ms));
	}

</script>

<div class="grid grid-cols-2 gap-4">
	<!-- Search text input -->
	<form class="flex mx-auto w-full space-y-4">
		<input 
			class="input h-8 pl-3 w-64"
			type="search"
			placeholder="Search..."
			bind:value={query} 
		/>
	</form>

	<!-- Search button -->
	<div class="flex mx-auto w-full space-y-4">
		<button type="button" class="btn preset-filled-primary-500 h-8 w-16" onclick={fetchResults} disabled={loading}>Search</button>
	</div>		
</div>

<!-- results -->
{#if error}
	{error}
{/if}
{#if loading}
	<div class="grid grid-cols-1 py-6">
		<Progress value={progressBarValue} />
	</div>
{/if}

{#if result && !loading}
<div class="grid py-6">
	<div class="card p-4 mx-auto justify-end">
		<span class="text-4xl">Results</span>
		<ul>
			<hr class="py-1"/>
			{#each result.documents as d, i}
				<li>
					{i + 1}. <span class="underline">{d}</span>
				</li>
			{/each}
		</ul>
	</div>
</div>
{/if}
