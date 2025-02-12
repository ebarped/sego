<script lang="ts">
    import { Progress } from '@skeletonlabs/skeleton-svelte';

    // Define the shape of the data returned from the API
    interface Result {
        query: string;
        documents: string[];
    }

    // Default values for the query
    let query: string = $state("");
    let resultCount: number = $state(5);

    let result: Result | null = $state(null); // Start as null since no data is loaded
    let loading: boolean = $state(false);
    let progressBarValue: number = $state(10);

    let error: string = $state("");

    async function fetchResults(): Promise<void> {
        loading = true;

        if (query == "") {
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
                await sleep(800); // Simulate transfer delay
            }
            progressBarValue = 10;

            const res = await fetch(
                "http://localhost:4000/search?query=" + query + "&count=" + resultCount
            );

            if (!res.ok) {
				error ='Failed to fetch results from Search Engine'
				throw new Error('Failed to fetch results from Search Engine');
			}

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

<!-- Main form container with flex-column layout -->
<div class="flex flex-col justify-center items-center space-y-4 py-4">
    <!-- Search input and button -->
    <div class="flex flex-row justify-center items-center space-x-2">
        <input 
            class="input h-8 pl-3 my-2 w-64"
            type="search"
            placeholder="Search..."
            bind:value={query} 
        />
        <!-- Search button -->
        <button 
            type="button" 
            class="btn preset-filled-primary-500 h-8 w-16"
            onclick={fetchResults}
            disabled={loading}
        >
            Search
        </button>
    </div>

    <!-- Show progress bar when loading -->
    {#if loading}
        <div class="w-1/2">
            <Progress value={progressBarValue} />
        </div>
    {/if}
	
	<!-- Results section when loaded -->
	{#if result && !loading}
		<div class="card p-6 mx-auto max-w-3xl">
			<h2 class="text-4xl font-semibold text-left mb-2">Results</h2>
			<hr class="mb-2" />
			<ul class="list-decimal pl-6">
				{#each result.documents as d, i}
					<li>
						<span class="underline">{d}</span>
					</li>
				{/each}
			</ul>
		</div>
	{/if}
	
    <!-- Show error message -->
    {#if error}
        <div class="text-red-500">
			{error}
		</div>
    {/if}
</div>

