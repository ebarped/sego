<script lang="ts">
    import { Settings } from 'lucide-svelte';

    let {showMenu}:{showMenu: boolean} = $props();
</script>
  
<button onclick={() => showMenu = !showMenu} type="button" class="btn preset-outlined">
    <Settings size={20} />  <!-- settings icon-->
</button>

{#if showMenu}
<form class="mx-auto w-full max-w-md space-y-4">
    <select class="select">
        <option value="5">5 Results</option>
        <option value="10">10 Results</option>
        <option value="15">15 Results</option>
    </select>
</form>
{/if}

<!-- al darle a settings, se abre un modal en medio con un form con todas las opciones , pec -->
