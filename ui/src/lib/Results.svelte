<script lang="ts">
    export let query: string = "";
    export let rscount: number;

    console.log("<Results.svelte> searching: " + query);
    const fetchResult = (async () => {
        const response = await fetch(
            "http://localhost:4000/search?query=" + query + "&count=" + rscount
        );
        return await response.json();
    })();
</script>

{#if query != ""}
    {#await fetchResult then data}
        <div class="grid py-6">
            <div class="card p-4 mx-auto justify-end">
                <span class="text-4xl">Results</span>
                <ul>
                    <hr />
                    {#each data.documents as r, i}
                        <li>
                            {i + 1}. <span class="underline">{r}</span>
                        </li>
                    {/each}
                </ul>
            </div>
        </div>
    {:catch error}
        <p>An error occurred!</p>
    {/await}
{/if}
