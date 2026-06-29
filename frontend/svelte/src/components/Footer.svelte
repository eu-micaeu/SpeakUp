<script>
  const currentYear = new Date().getFullYear();
</script>

<footer class="footer">
  <div class="footerGlow" aria-hidden="true"></div>
  <div class="container">
    <section class="brandColumn">
      <a href="/" class="brand" aria-label="Voltar para a página inicial">
        <img src="/logo.png" alt="Logo SpeakUp" width="34" height="34" />
        <span>SpeakUp</span>
      </a>
      <p class="description">
        Pratique idiomas todos os dias com chat inteligente, correções instantâneas
        e acompanhamento da sua evolução.
      </p>
    </section>
  </div>

  <div class="bottomBar">
    <p>&copy; {currentYear} SpeakUp. Todos os direitos reservados.</p>
    <!-- <a href="/planos">Conheça os planos</a> -->
  </div>
</footer>

<style>
  .footer {
    position: relative;
    overflow: hidden;
    background: #090909;
    border-top: 1px solid rgba(255, 255, 255, 0.08);
    color: #e8e8e8;
    padding: 3.5rem 0 1.25rem;
    font-size: 0.92rem;
    z-index: 1;
  }

  .footerGlow {
    position: absolute;
    top: -120px;
    left: 50%;
    transform: translateX(-50%);
    width: min(780px, 95vw);
    height: 240px;
    pointer-events: none;
  }

  .container {
    width: min(1120px, calc(100% - 3rem));
    margin: 0 auto;
    display: grid;
    grid-template-columns: 1.4fr repeat(2, minmax(180px, 1fr));
    gap: 2.25rem;
    position: relative;
    z-index: 1;
  }

  .brandColumn {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .brand {
    display: inline-flex;
    align-items: center;
    gap: 0.65rem;
    color: #fff;
    text-decoration: none;
    font-size: 1.35rem;
    font-weight: 700;
    width: fit-content;
  }

  .description {
    margin: 0;
    color: #b8b8b8;
    line-height: 1.65;
    max-width: 38ch;
  }

  .primaryLink {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: fit-content;
    text-decoration: none;
    padding: 0.7rem 1rem;
    border-radius: 10px;
    background: #ffffff;
    color: #090909;
    font-weight: 700;
    transition: transform 0.2s, box-shadow 0.2s;
  }

  .primaryLink:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 18px rgba(255, 255, 255, 0.18);
  }

  .column {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .column h4 {
    margin: 0;
    margin-bottom: 0.75rem;
    font-size: 1rem;
    color: #f2f2f2;
    letter-spacing: 0.02em;
  }

  .column ul {
    margin: 0;
    list-style: none;
    padding: 0;
    display: grid;
    gap: 0.65rem;
  }

  .column a {
    color: #c9c9c9;
    text-decoration: none;
    width: fit-content;
    transition: color 0.2s, transform 0.2s;
  }

  .column a:hover {
    color: #ffffff;
    transform: translateX(2px);
  }

  .contactLink {
    font-weight: 600;
  }

  .status {
    margin: 0.4rem 0 0;
    color: #9c9c9c;
    display: inline-flex;
    align-items: center;
    gap: 0.45rem;
  }

  .dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: #53d688;
    box-shadow: 0 0 0 4px rgba(83, 214, 136, 0.16);
  }

  .bottomBar {
    width: min(1120px, calc(100% - 3rem));
    margin: 2.5rem auto 0;
    padding-top: 1rem;
    border-top: 1px solid rgba(255, 255, 255, 0.08);
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
    font-size: 0.8rem;
    color: #8f8f8f;
    position: relative;
    z-index: 1;
  }

  .bottomBar p {
    margin: 0;
  }

  .bottomBar a {
    color: #d8ddff;
    text-decoration: none;
    font-weight: 600;
  }

  .bottomBar a:hover {
    color: #ffffff;
  }

  @media (max-width: 900px) {
    .container {
      grid-template-columns: 1fr 1fr;
    }

    .brandColumn {
      grid-column: 1 / -1;
    }
  }

  @media (max-width: 640px) {
    .footer {
      padding-top: 3rem;
    }

    .container {
      width: min(1120px, calc(100% - 2rem));
      grid-template-columns: 1fr;
      gap: 1.75rem;
    }

    .bottomBar {
      width: min(1120px, calc(100% - 2rem));
      margin-top: 2rem;
      flex-direction: column;
      align-items: flex-start;
      padding-bottom: 0.25rem;
    }
  }
</style>
