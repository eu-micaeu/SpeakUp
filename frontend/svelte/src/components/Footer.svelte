<script>
  const currentYear = new Date().getFullYear();
</script>

<footer class="footer">
  <div class="container">
    <div class="column">
      <h4>Links úteis</h4>
      <ul>
        <li><a href="/login">Login</a></li>
        <li><a href="#features">Funcionalidades</a></li>
        <li><a href="/about">Sobre</a></li>
      </ul>
    </div>
    <div class="column">
      <h4>Contato</h4>
      <p>Email: contato@speakup.com</p>
      <p>Suporte: suporte@speakup.com</p>
    </div>
  </div>
  <div class="bottomBar">
    <p>&copy; {currentYear} SpeakUp. Todos os direitos reservados.</p>
  </div>
</footer>

<style>
  .footer {
    background-color: #000;
    color: #eee;
    padding: 3rem 0;
    font-size: 0.9rem;
    z-index: 1;
    position: relative;
  }

  .container {
    display: flex;
    flex-wrap: wrap;
    justify-content: space-evenly;
    gap: 2rem;
    margin: 0 auto;
  }

  .column {
    display: flex;
    flex-direction: column;
    flex: 1;
    min-width: 200px;
    max-width: 300px;
    padding: 1rem;
  }

  .column h4 {
    margin-bottom: 0.75rem;
    font-size: 1.1rem;
    color: #ccc;
  }

  .column ul {
    list-style: none;
    padding: 0;
  }

  .column li {
    margin-bottom: 0.5rem;
  }

  .column a {
    color: #fff;
    text-decoration: none;
    transition: color 0.2s;
  }

  .column a:hover {
    color: #ccc;
  }

  .bottomBar {
    margin-top: 2rem;
    text-align: center;
    font-size: 0.8rem;
    color: #777;
  }
</style>
