<script>
    import { onMount } from "svelte";
    import { browser } from "$app/environment";

    let Toaster;

    onMount(async () => {
        if (browser) {
            const module = await import("svelte-sonner");
            Toaster = module.Toaster;
        }
    });
</script>

{#if Toaster}
    <svelte:component this={Toaster} position="top-right" richColors />
{/if}
<slot />
