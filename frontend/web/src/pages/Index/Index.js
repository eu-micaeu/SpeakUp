import './Index.css';

// Components
import Header from '../../components/Header/Header';
import Footer from '../../components/Footer/Footer';

function Index() {
  return (
    <>

      <Header />

      <main>

        <h1>SpeakUp</h1>

        <div className="cards">

          <div className="card">
            <img src="./chat.png"></img>
            <p>
              SpeakUp é uma plataforma inovadora que utiliza inteligência artificial para auxiliar no aprendizado de novos idiomas de maneira eficiente e personalizada.
              Com foco na melhoria contínua, a IA analisa seus textos em tempo real, identificando erros de digitação e gramaticais.
              A cada erro encontrado, a ferramenta destaca as correções necessárias, oferecendo explicações detalhadas sobre o que foi corrigido e o motivo.
              Dessa forma, você não apenas corrige suas falhas, mas também aprende com cada interação, desenvolvendo suas habilidades linguísticas de maneira prática e envolvente.
            </p>
          </div>

          <div className="card">
            <img src="./ia.png"></img>
            <p>O SpeakUp é uma plataforma de ensino de idiomas que utiliza inteligência artificial para oferecer uma experiência de aprendizado personalizada e eficiente.
              A IA analisa seus textos em tempo real, identificando erros de digitação e gramaticais, e destaca as correções necessárias para que você aprenda com cada interação.
              Além disso, a ferramenta oferece explicações detalhadas sobre o que foi corrigido e o motivo, permitindo que você compreenda melhor as regras gramaticais e melhore suas habilidades linguísticas.
              Com o SpeakUp, você pode praticar o idioma de forma prática e envolvente, corrigindo suas falhas e aprimorando suas habilidades com a ajuda da inteligência artificial.
            </p>
          </div>

          <div className="card">
            <img src="./mic.png"></img>
            <p>
              Atualmente em fase de desenvolvimento e implementação, trará uma inovadora funcionalidade de reconhecimento de voz.
              Com essa ferramenta, você poderá utilizar o microfone para captar as palavras que pronunciar enquanto pratica o idioma. A inteligência artificial do SpeakUp analisará sua fala em tempo real, identificando erros de pronúncia e realizando correções instantaneamente.
              Além disso, fornecerá um feedback detalhado sobre as áreas em que você pode melhorar, permitindo um aprendizado mais dinâmico e interativo.
            </p>
          </div>

        </div>

        <hr></hr>

        <h1>Idiomas disponiveis para o aprendizado:</h1>

        <div className="languages">

          <div className="language">
            <img src="./EUA.png"></img>
          </div>

          <div className="language">
            <img src="./japao.png"></img>
          </div>

        </div>

      </main>

      <Footer />

    </>
  );
}

export default Index;