import styles from './FeatureCard.module.css';

export function FeatureCard({ imgSrc, alt, title, content }) {
  return (
    <div className={styles.feature}>
      <img src={imgSrc} alt={alt} />
      <div>
        <h2>{title}</h2>
        <p>{content}</p>
      </div>
    </div>
  );
}
