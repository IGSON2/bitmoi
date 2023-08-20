import styles from "./AdDiv.module.css";
import mockup from "../../../component/images/mockup_rank.png";

function AdDiv() {
  return (
    <div className={styles.addiv}>
      <img className={styles.adimage} src={mockup}></img>
    </div>
  );
}

export default AdDiv;
