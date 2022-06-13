import styles from "./Loader.module.css";
import bgimg from "../images/moilogo.png";
function Loader() {
  return (
    <div className={styles.loaderpage}>
      <div className={styles.catchphrase}>
        <div className={styles.sentence}>바로 시작하는 모의투자 비트모이!</div>
        <img className={styles.bgimg} src={bgimg}></img>
      </div>
    </div>
  );
}

export default Loader;
