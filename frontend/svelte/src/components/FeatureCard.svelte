<script>
  export let imgSrc,
    alt,
    title,
    content,
    index = 0;
</script>

<div class="feature" style="animation-delay: {index * 0.2}s">
  <div class="iconWrapper">
    <img src={imgSrc} {alt} />
    <div class="iconGlow"></div>
  </div>
  <div class="content">
    <h2>{title}</h2>
    <p>{content}</p>
  </div>
</div>

<style>
  .feature {
    display: flex;
    gap: 2rem;
    align-items: center;
    width: 100%;
    max-width: 900px;
    margin: 0 auto 2rem;
    padding: 2.5rem;
    background: rgba(255, 255, 255, 0.03);
    border-radius: 8px;
    border: 1px solid rgba(255, 255, 255, 0.08);
    transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
    position: relative;
    overflow: hidden;
    opacity: 0;
    animation: fadeInUp 0.8s ease-out forwards;
  }

  @keyframes fadeInUp {
    from {
      opacity: 0;
      transform: translateY(30px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .feature::before {
    content: "";
    position: absolute;
    top: 0;
    left: -100%;
    width: 100%;
    height: 100%;
    background: linear-gradient(
      90deg,
      transparent,
      rgba(255, 255, 255, 0.05),
      transparent
    );
    transition: left 0.5s;
  }

  .feature:hover::before {
    left: 100%;
  }

  .feature:hover {
    transform: translateY(-8px);
    border-color: rgba(255, 255, 255, 0.2);
    background: rgba(255, 255, 255, 0.06);
  }

  .iconWrapper {
    position: relative;
    flex-shrink: 0;
    width: 120px;
    height: 120px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(
      135deg,
      rgba(255, 255, 255, 0.1),
      rgba(255, 255, 255, 0.05)
    );
    border-radius: 6px;
    border: 1px solid rgba(255, 255, 255, 0.15);
    transition: all 0.4s ease;
  }

  .iconGlow {
    position: absolute;
    width: 100%;
    height: 100%;
    border-radius: 6px;
    opacity: 0;
  }

  .feature:hover .iconGlow {
    opacity: 1;
  }

  .iconWrapper img {
    width: 70px;
    height: 70px;
    object-fit: contain;
    transition: all 0.4s ease;
    position: relative;
    z-index: 1;
  }

  .content {
    flex: 1;
    min-width: 0;
  }

  .content h2 {
    margin-bottom: 1rem;
    font-size: 1.3rem;
    font-weight: 700;
    line-height: 1.4;
    color: #ffffff;
    transition: color 0.3s;
  }

  .feature:hover .content h2 {
    background: linear-gradient(90deg, #ffffff, #e0e0e0);
    background-clip: text;
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
  }

  .content p {
    font-size: 1.05rem;
    line-height: 1.7;
    opacity: 0.85;
    color: #e0e0e0;
    transition: opacity 0.3s;
  }

  .feature:hover .content p {
    opacity: 1;
  }

  @media screen and (max-width: 768px) {
    .feature {
      flex-direction: column;
      text-align: center;
      padding: 2rem 1.5rem;
      margin-bottom: 1.5rem;
    }

    .iconWrapper {
      width: 100px;
      height: 100px;
    }

    .iconWrapper img {
      width: 60px;
      height: 60px;
    }

    .content h2 {
      font-size: 1.1rem;
    }

    .content p {
      font-size: 0.95rem;
    }
  }
</style>
