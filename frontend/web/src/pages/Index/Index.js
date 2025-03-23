import './Index.css';

// Components
import Header from '../../components/Header/Header';
import Footer from '../../components/Footer/Footer';

function Index() {
  return (
    <>
    
      <Header />

      <main>
        
        <div className="cards">

          <div className="card">
            <h2>Card 1</h2>
            <p>Content</p>
          </div>

          <div className="card">
            <h2>Card 2</h2>
            <p>Content</p>
          </div>

          <div className="card">
            <h2>Card 3</h2>
            <p>Content</p>
          </div>

        </div>

      </main>

      <Footer />
    
    </>
  );
}

export default Index;