<script lang="ts">
  export let svg = "",
    tag = "",
    title = "",
    content = "",
    index = 0;

  let isVisible = false;

  function revealOnScroll(node: HTMLElement) {
    if (typeof IntersectionObserver === "undefined") {
      isVisible = true;
      return;
    }

    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            isVisible = true;
            observer.unobserve(node);
          }
        });
      },
      { threshold: 0.15, rootMargin: "0px 0px -40px 0px" }
    );

    observer.observe(node);

    return {
      destroy() {
        observer.disconnect();
      },
    };
  }
</script>

<div
  class="featureCard"
  class:visible={isVisible}
  style="transition-delay: {index * 90}ms"
  use:revealOnScroll
>
  <div class="cardTop">
    <div class="iconBox">
      {@html svg}
    </div>
    {#if tag}
      <span class="cardTag">{tag}</span>
    {/if}
  </div>

  <div class="cardBody">
    <h2>{title}</h2>
    <p>{content}</p>
  </div>
</div>

<style>
  .featureCard {
    background-color: #111111;
    border: 1px solid #222222;
    border-radius: 12px;
    padding: 1.75rem;
    display: flex;
    flex-direction: column;
    justify-content: flex-start;
    text-align: left;
    opacity: 0;
    transition:
      opacity 0.55s ease,
      border-color 0.25s ease,
      background-color 0.25s ease;
  }

  .featureCard.visible {
    opacity: 1;
  }

  .featureCard.visible:hover {
    border-color: #333333;
    background-color: #151515;
  }

  .cardTop {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 1.35rem;
  }

  .iconBox {
    width: 40px;
    height: 40px;
    background-color: #181818;
    border: 1px solid #262626;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #e4e4e7;
    transition: border-color 0.25s ease, color 0.25s ease;
  }

  .featureCard:hover .iconBox {
    border-color: #3a3a3a;
    color: #ffffff;
  }

  .cardTag {
    font-size: 0.68rem;
    font-weight: 600;
    letter-spacing: 0.06em;
    text-transform: uppercase;
    color: #777777;
    background-color: #181818;
    border: 1px solid #242424;
    padding: 0.22rem 0.55rem;
    border-radius: 4px;
  }

  .cardBody h2 {
    margin: 0 0 0.5rem 0;
    font-size: 1.15rem;
    font-weight: 700;
    line-height: 1.35;
    color: #ffffff;
    letter-spacing: -0.01em;
  }

  .cardBody p {
    margin: 0;
    font-size: 0.9rem;
    line-height: 1.6;
    color: #888888;
  }
</style>

