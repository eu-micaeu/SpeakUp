<script>
  export let imgSrc, alt, title, content;
</script>

<div class="feature">
  <img src={imgSrc} alt={alt} />
  <div>
    <h2>{title}</h2>
    <p>{content}</p>
  </div>
</div>

<style>
  .feature {
    display: flex;
    gap: 2rem;
    align-items: center;
    width: 800px;
    margin: 20px auto;
    padding: 20px;
    justify-content: space-between;
  }

  .feature img {
    width: 100px;
    height: auto;
  }

  .feature h2 {
    margin-bottom: 0.5rem;
    text-align: justify;
  }

  .feature p {
    font-size: 1rem;
    max-width: 600px;
    margin-bottom: 0.5rem;
    text-align: justify;
  }

  .feature div {
    width: 600px;
    opacity: 0.8;
  }

  @media screen and (max-width: 768px) {
    .feature {
      flex-direction: column;
      width: 70%;
      text-align: center;
      margin: 0 auto;
    }

    .feature img {
      width: 80px;
    }

    .feature h2 {
      font-size: 1rem;
    }

    .feature div {
      font-size: 0.8rem;
      width: 100%;
    }
  }
</style>