<script lang="ts">
    export let query = "";
    export let inputSubmitted = false;
    export let rerender = false;

    // when ENTER press, update state to trigger a new search and rerender of Results
    const updateState = (e) => {
        if (e.charCode === 13) {
            query = e.target.value;
            inputSubmitted = true;
            rerender = !rerender;
            console.log("<Search.svelte> searching: " + query);
        }
    };
</script>

<!-- center and increase size -->
<div class="container mx-auto justify-center max-w-2xl">
    <input
        class="input"
        type="text"
        placeholder="Search"
        on:keypress={updateState}
    />
</div>
