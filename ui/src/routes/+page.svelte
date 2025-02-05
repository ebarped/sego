<script lang="ts">
    import Header from "$lib/Header.svelte";
    import Search from "$lib/Search.svelte";

	let {query = "test", resultCount = 5 }:{query: string, resultCount:number} = $props();

</script>

<Header />

<!-- outer container, align all to center -->
<div class="container h-full mx-auto flex justify-center items-start py-12">
    <!-- number of columns -->
    <div class="grid grid-cols-1">
        <h1 class="text-center py-2 text-5xl">
            Linux API Documentation Search Engine
        </h1>
        <!-- hr draws a line -->
        <hr />
        <!-- vertical padding from the previous component -->
        <div class="container justify-center grid py-3">
            <Search/>
        </div>
    </div>
</div>
