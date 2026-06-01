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

<svelte:head>
    <!-- SEO Básico -->
    <title>SpeakUp | Aprenda Idiomas com Inteligência Artificial</title>
    <meta name="description" content="Melhore sua pronúncia, escrita e gramática em tempo real com a SpeakUp. A plataforma definitiva de aprendizado de idiomas impulsionada por Inteligência Artificial." />
    <meta name="keywords" content="aprender inglês, praticar pronúncia, inteligência artificial, aprender idiomas, falar inglês fluente, correção gramatical, svelte, ia" />
    <meta name="author" content="SpeakUp" />
    <meta name="robots" content="index, follow" />

    <!-- Open Graph / Facebook / LinkedIn -->
    <meta property="og:type" content="website" />
    <meta property="og:title" content="SpeakUp | Aprenda Idiomas com Inteligência Artificial" />
    <meta property="og:description" content="Melhore sua pronúncia, escrita e gramática em tempo real com a SpeakUp. A plataforma definitiva de aprendizado de idiomas impulsionada por Inteligência Artificial." />
    <meta property="og:site_name" content="SpeakUp" />

    <!-- Twitter -->
    <meta name="twitter:card" content="summary_large_image" />
    <meta name="twitter:title" content="SpeakUp | Aprenda Idiomas com Inteligência Artificial" />
    <meta name="twitter:description" content="Melhore sua pronúncia, escrita e gramática em tempo real com a SpeakUp. A plataforma definitiva de aprendizado de idiomas impulsionada por Inteligência Artificial." />
</svelte:head>

<slot />

